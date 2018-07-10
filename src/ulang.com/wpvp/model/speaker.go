package model

import (
	"os"
	"io/ioutil"
	"ulang.com/wpvp/conf"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

const DOC_SPEAKER = "speaker"

func SaveRetrieveFile(taskId string, voiceFile [] byte) {
	dir := conf.Config.Retrieve_Dir
	os.MkdirAll(dir, 0777)
	ioutil.WriteFile(dir+"/"+taskId+".wav", voiceFile, 0755)
}

//genQueryCmd 生成查询语句
func NewQueryCondition(result *Condition) (map[string]interface{}, error) {
	if result == nil {
		return nil, nil
	}
	cmd := make(map[string]interface{}) //map要初始化 才能使用
	if result.PhoneNum != "" {
		nums := strings.Split(result.PhoneNum, ",")
		cmd["tel"] = bson.M{"$in": nums}
	}
	if result.Lang != "" {
		cmd["lang"] = result.Lang
	}
	if result.Gender != "" {
		cmd["gender"] = result.Gender
	}
	if result.Speaker != "" {
		cmd["spkname"] = result.Speaker
	}
	if result.City != "" {
		cmd["city"] = result.City
	}
	if result.StartTime != 0 && result.EndTime != 0 {
		cmd["tmTrain"] = map[string]interface{}{
			"Start": result.StartTime,
			"End":   result.EndTime,
		}
	}

	return cmd, nil
}
