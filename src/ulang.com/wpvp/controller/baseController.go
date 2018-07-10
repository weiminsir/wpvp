package controller

import (
	"github.com/labstack/echo"
	"ulang.com/wpvp/model"
	"io/ioutil"
	"net/http"
	"log"
	"strings"
	"os"
	"path/filepath"
	"net/url"
)

type BaseController struct {
	echo.Context
	jsonData map[string]interface{}
	userme   model.SysUser
}

func ResponseMsg(code int) interface{} {
	return map[string]interface{}{
		"Code":    code,
		"Message": GetMessage(code),
	}
}
func ResponseError(code int, err error) interface{} {
	return map[string]interface{}{
		"Code":    code,
		"Message": err.Error(),
	}
}
func ResponseOK(data interface{}) interface{} {
	return map[string]interface{}{
		"Code":    S_OK,
		"Message": GetMessage(S_OK),
		"Data":    data,
	}
}

var Gender = []int{0, 1, 2}
var Lang = []string{"英语", "日语", "韩语", "汉语", "闽南语", "沪语", "维语"}
//var Lang = []string{"英语", "法语",
//	"德语", "哈萨克语",
//	"日语", "韩语",
//	"汉语", "闽南语",
//	"俄语", "沪语",
//	"土耳其语", "维语",
//	"粤语", "安多藏语",
//	"康巴藏语", "拉萨藏语"}

type StatResult struct {
	Date   string `json:"Date,omitempty"`
	Count  int64  `json:"Count"`
	Lang   string `json:"Lang,omitempty"`
	Gender int    `json:"Gender,omitempty"`
}

func ForwardPostRequest(addr string, values map[string][]string) ([]byte) {
	resp, err := http.PostForm(addr, values)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	return body
}
func ForwardRawPostRequest(url string, jsonString string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonString))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func ForwardGetRequest(addr string) ([]byte, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func Download(output echo.Context, file string, filename ...string) {
	// check get file error, file not found or other error.
	if _, err := os.Stat(file); err != nil {
		http.ServeFile(output.Response().Writer, output.Request(), file)
		return
	}
	var fName string
	if len(filename) > 0 && filename[0] != "" {
		fName = filename[0]
	} else {
		fName = filepath.Base(file)
	}
	output.Response().Header().Set("Content-Disposition", "attachment; filename="+url.QueryEscape(fName))
	output.Response().Header().Set("Content-Description", "File Transfer")
	output.Response().Header().Set("Content-Type", "application/octet-stream")
	output.Response().Header().Set("Content-Transfer-Encoding", "binary")
	output.Response().Header().Set("Expires", "0")
	output.Response().Header().Set("Cache-Control", "must-revalidate")
	output.Response().Header().Set("Pragma", "public")
	http.ServeFile(output.Response().Writer, output.Request(), file)
}
