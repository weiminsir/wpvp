package controller

import (
	"strconv"
	"github.com/labstack/echo"
	"net/http"
	"ulang.com/wpvp/model"
	"ulang.com/wpvp/conf"
	"encoding/json"
	"github.com/tidwall/gjson"
	"errors"
	"log"
	"os"
	"io/ioutil"
	"mime/multipart"
	"io"
	"bytes"
	"github.com/donnie4w/go-logger/logger"
	"path/filepath"
	"mime"
	"net/url"
)

const (
	Retrieval_Mobile   = 0
	Retrieval_Launcher
	Retrieval_Voice    = 2
)

func AdminGetSpkRecords(c echo.Context) error {

	taskId := c.Param("task_id")
	if len(taskId) == 0 {
		return c.JSON(http.StatusOK, ResponseMsg(S_INVALID_PARAM))
	}

	task, found := model.FindRetrievalTask(taskId)
	if !found {
		return c.JSON(http.StatusOK, ResponseMsg(S_INVALID_PARAM))
	}

	var result *ForwardResult
	var err error
	if task.TaskType == Retrieval_Mobile {
		url := conf.Config.Proxy_Addr + "/api/objects/query"
		result, err = RemoteGetSpeakers(url, task.Query)
		if err != nil {
			return c.JSON(http.StatusOK, echo.Map{
				"Code":    S_CONVERT_ERROR,
				"Message": err.Error(),
			})
		}
	}
	if task.TaskType == Retrieval_Voice {
		url := conf.Config.Proxy_Addr + "/api/objects/search"
		result, err = RemoteFileRetrieveSpeaker(url, task)
		if err != nil {
			return c.JSON(http.StatusOK, echo.Map{
				"Code":    S_CONVERT_ERROR,
				"Message": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"Code": S_OK,
		"Data": result.msg,
	})
}
func AdminGetSpeakerStatistics(c echo.Context) error {
	statType := c.QueryParam("type")
	duration := c.QueryParam("duration")
	if len(statType) == 0 {
		statType = "0" //默认type=0
	}
	if len(duration) == 0 {
		duration = "7" //默认一周
	}
	url := conf.Config.Proxy_Addr + "/api/sidb/stat/" + duration + "/" + statType

	result, err := RemoteGetSpeakersStatistic(url)
	if err != nil {
		return c.JSON(http.StatusOK, ResponseMsg(S_GRPC_ERROR))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"Code": S_OK,
		"Data": result.msg,
	})
}

func AdminGetVoice(c echo.Context) error {
	url := conf.Config.Proxy_Addr + c.Request().RequestURI
	name := c.Param("name")
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	var b bytes.Buffer
	io.Copy(&b, res.Body)
	iRead := bytes.NewReader(b.Bytes())
	log.Printf("Size:%d, Len:%d\n", iRead.Size(), iRead.Len())
	return SetVoiceOutput(c, iRead, name)

}

// SetVoiceOutput send image to client
func SetVoiceOutput(c echo.Context, r io.Reader, name string) (err error) {
	w := c.Response().Writer
	ctype := mime.TypeByExtension(filepath.Ext(name))
	c.Response().Header().Set(echo.HeaderContentType, ctype)
	log.Println("ctype=", ctype, "filename=", url.QueryEscape(name))

	c.Response().Header().Set("Content-Disposition", "attachment; filename="+url.QueryEscape(name))
	c.Response().Header().Set(echo.HeaderContentType, "application/octet-stream")
	c.Response().WriteHeader(http.StatusOK)
	_, err = io.Copy(w, r)
	return
}

func RemoteGetSpeakers(url string, jsonStr string) (*ForwardResult, error) {
	logger.Info("RemoteGetSpeakers Query=", jsonStr)
	body, err := ForwardRawPostRequest(url, jsonStr)
	if err != nil {
		return nil, err
	}

	code := gjson.Get(string(body), "code")
	value := gjson.Get(string(body), "msg")
	if !value.Exists() {
		r := &ForwardResult{code: code.String(), msg: nil}
		return r, nil
	}
	var speakers []Speaker
	for _, v := range value.Array() {
		spk := Speaker{}
		json.Unmarshal([]byte(v.String()), &spk)
		speakers = append(speakers, spk)
	}
	r := ForwardResult{msg: speakers, count: len(speakers)}
	return &r, nil
}
func RemoteGetSpeakersStatistic(url string) (*ForwardResult, error) {
	body, err := ForwardGetRequest(url)
	if err != nil {
		return nil, err
	}
	value := gjson.Get(string(body), "msg")
	if !value.Exists() {
		//r := &ForwardResult{code: "200", msg: nil}
		return nil, errors.New("no data")
	}
	var stat []StatResult
	for _, v := range value.Array() {
		s := StatResult{}
		json.Unmarshal([]byte(v.String()), &s)
		stat = append(stat, s)
	}

	r := ForwardResult{msg: stat, count: len(stat)}
	return &r, nil
}
func RemoteFileRetrieveSpeaker(url string, task *model.RetrievalTask) (*ForwardResult, error) {
	paramName := "Voc"
	dir := conf.Config.Retrieve_Dir
	filePath := dir + "/" + task.TaskId + ".wav"
	if task.VoiceFile != nil {
		os.MkdirAll(dir, 0777)
		ioutil.WriteFile(filePath, task.VoiceFile, 0755)
	}
	threshold := strconv.Itoa(int(task.Threshold))
	params := map[string]string{
		"Type":      "0",
		"SpkName":   task.SpeakerName,
		"Lang":      task.Lang,
		"Gender":    task.Gender,
		"City":      task.City,
		"Threshold": threshold,
	}
	logger.Info("params=", params)
	return ForwardFileRetrieveRequest(url, params, paramName, filePath)

}

// start retrieve result from speaker
func ForwardFileRetrieveRequest(uri string, params map[string]string, paramName, path string) (*ForwardResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		log.Panic(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", uri, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	//生成要访问的url
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	logger.Info("file retrieve result=", string(b))

	code := gjson.Get(string(b), "code")
	value := gjson.Get(string(b), "msg")
	if !value.Exists() {
		r := &ForwardResult{code: code.String(), msg: nil}
		return r, nil
	}
	var speakers []Speaker
	for _, v := range value.Array() {
		spk := Speaker{}
		json.Unmarshal([]byte(v.String()), &spk)
		speakers = append(speakers, spk)
	}
	r := ForwardResult{msg: speakers, count: len(speakers)}
	return &r, nil
}
