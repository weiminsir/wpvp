package model

import (
	"gopkg.in/mgo.v2/bson"
)

type SysCfg struct {
	Id         bson.ObjectId `json:"_id,omitempty"bson:"_id,omitempty"`
	NextTaskId int           "NextTaskId"
	NextUserId int           "NextUserId"
}
type SysUser struct {
	UserId      string `json:"UserId,omitempty"bson:"UserId,omitempty"`
	PhoneNumber string `json:"PhoneNumber,omitempty"bson:"PhoneNumber,omitempty"`
	Email       string `json:"Email,omitempty"bson:"Email,omitempty"`
	Username    string `json:"Username,omitempty"bson:"Username,omitempty"`
	Avatar      string `json:"Avatar,omitempty"bson:"Avatar,omitempty"`
	Password    string `json:"Password,omitempty"bson:"Password,omitempty"`
	Role        int    `json:"Role,omitempty"bson:"Role,omitempty"`
}
type Authority struct {
	UserId string "UserId"
	Role   int    "Role" //super 0/ 一级管理员1,以此类推
}

type RetrievalTask struct {
	Id          bson.ObjectId `json:"_id,omitempty"bson:"_id,omitempty"`
	TaskId      string        "TaskId"
	Launcher    string        "Launcher"
	CreateTime  int64         "CreateTime"
	Duration    string        "Duration"
	Total       int           "Total"
	State       bool          "State"
	Query       string        "Query"
	Lang        string        "Lang"
	Gender      string        "Gender"
	City        string        "City"
	MaxCount    int64         "MaxCount"
	Threshold   float32       "Threshold"
	SpeakerName string        "Speaker"
	TaskType    int           "TaskType"
	VoiceFile   []byte        `json:"VoiceFile,omitempty"bson:"VoiceFile,omitempty"`
}
type Condition struct {
	PhoneNum  string  `PhoneNum`
	Lang      string  `Lang`
	Gender    string  `Gender`
	City      string  `City`
	StartTime int64   `StartTime`
	EndTime   int64   `EndTime`
	MaxCount  int64   `MaxCount`
	Threshold float32 `Threshold`
	Speaker   string  `Speaker`
	TaskType  int     `TaskType`
	VoiceFile []byte  `VoiceFile`
	FileDir   string  `FileDir`
}
