package controller

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"ulang.com/wpvp/model"
	"encoding/json"
	"ulang.com/wpvp/conf"
	"github.com/donnie4w/go-logger/logger"
)

func AdminGetRetrievalTasks(c echo.Context) error {
	cond := bson.M{}
	sort := "-CreateTime"
	tasks, err := model.GetRetrievalTask(&cond, &sort)

	if err != nil {
		return c.JSON(http.StatusOK, ResponseMsg(S_DATABASE_ERROR))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"Code": S_OK,
		"Data": tasks,
	})
}
func AdminGetRetrievalTasksList(c echo.Context) error {

	pageStart := 0
	pageSize := 30
	start := c.Param("page_start")
	if len(start) != 0 {
		pageStart, _ = strconv.Atoi(start)
	}

	size := c.Param("page_size")
	if len(size) != 0 {
		pageSize, _ = strconv.Atoi(size)
	}

	cond := bson.M{}
	sort := "-CreateTime"
	tasks, cnt, err := model.ListRetrievalTask(&cond, &sort, pageStart, pageSize)
	if err != nil {
		return c.JSON(http.StatusOK, ResponseMsg(S_DATABASE_ERROR))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"Code":      S_OK,
		"Data":      tasks,
		"Count":     cnt,
		"PageCount": (cnt + pageSize - 1) / pageSize,
	})
}

func AdminInsertRetrievalTask(c echo.Context) error {
	timeStart := time.Now() // get current time
	s_code := S_OK
	retrieval, _ := NewRetrievalCondition(c)
	condition, _ := model.NewQueryCondition(retrieval)
	logger.Info("condition=", condition)
	var task *model.RetrievalTask
	var err error
	task = &model.RetrievalTask{
		Lang:        retrieval.Lang,
		Gender:      retrieval.Gender,
		City:        retrieval.City,
		MaxCount:    retrieval.MaxCount,
		Threshold:   retrieval.Threshold,
		SpeakerName: retrieval.Speaker,
		VoiceFile:   retrieval.VoiceFile,
		CreateTime:  time.Now().Unix(),
	}

	task.TaskId, _ = model.SysCfgIncRetrievalTask()
	var result *ForwardResult
	if len(retrieval.VoiceFile) == 0 {
		logger.Info("*************no file retrieve...")
		var jsonbytes []byte
		task.TaskType = Retrieval_Mobile
		fr := &ForwardRequest{condition, task.MaxCount}
		addr := conf.Config.Proxy_Addr + "/api/objects/query"
		jsonbytes, err = json.Marshal(fr)
		if err != nil {
			s_code = S_JSON_ERROR
		}
		task.Query = string(jsonbytes)
		result, err = RemoteGetSpeakers(addr, task.Query)
		if err != nil {
			s_code = S_CONVERT_ERROR
		} else {
			task.Total = result.count
		}

	} else {
		logger.Info("*************file retrieve...")
		task.TaskType = Retrieval_Voice
		addr := conf.Config.Proxy_Addr + "/api/objects/search"
		result, err = RemoteFileRetrieveSpeaker(addr, task)
		if err != nil {
			s_code = S_CONVERT_ERROR
		} else {
			task.Total = result.count
		}
	}
	if s_code == S_OK {
		task.State = true
	}
	task.Duration = fmt.Sprintf("%s", time.Since(timeStart))
	if task.VoiceFile != nil {
		task.VoiceFile = nil
	}
	err = model.NewRetrieveTask(task)
	if err != nil {
		s_code = S_DATABASE_ERROR
		return c.JSON(http.StatusOK, ResponseError(s_code, err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"Code": s_code,
		"Data": task,
	})
}

func NewRetrievalCondition(c echo.Context) (*model.Condition, error) {
	condition := &model.Condition{}
	if err := c.Bind(condition); err != nil {
		return condition, err
	}
	VoiceFile, err := c.FormFile("VoiceFile")
	if err != nil {
		condition.VoiceFile = nil
	} else {
		bytes, _ := ConvertFileToBytes(VoiceFile)
		condition.VoiceFile = bytes
	}

	if condition.MaxCount == 0 {
		condition.MaxCount = 60
	}

	if condition.Threshold == 0 {
		condition.Threshold = 75
	}
	return condition, nil
}
