package controller

import (
	"strconv"
)

const (
	S_OK             = 0
	S_ERROR          = 1
	S_PAGE_ERROR     = 2
	S_FORM_ERROR     = 3
	S_DATABASE_ERROR = 4
	S_INVALID_PARAM  = 5
	S_PHONE_ERROR    = 6
	S_JSON_ERROR     = 7
	S_GRPC_ERROR     = 8
	S_CONVERT_ERROR  = 9

	S_FILE_ERROR        = 20
	S_FILE_EMPTY        = 21
	S_FILE_UPLOAD_ERROR = 22

	S_INVALID_USERID    = 100
	S_DUP_USERID        = 101
	S_INVALID_PASSWORD  = 102
	S_SHORT_PASSWORD    = 104
	S_PASSWORD_MISMATCH = 105
	S_LOGIN_FAILED      = 106
	S_LOGIN_ELSEWHERE   = 107
	S_ACCOUNT_DISABLED  = 108
	S_LOGIN_REQUIRED    = 109

	S_UNKNOW_ERROR = -1
)

var status_text = map[int]string{
	//常规类错误 0开头
	S_OK:             "操作成功",
	S_ERROR:          "内部错误",
	S_PAGE_ERROR:     "页面错误",
	S_FORM_ERROR:     "表单错误",
	S_DATABASE_ERROR: "数据库操作错误",
	S_INVALID_PARAM:  "无效的参数或参数为空",
	S_PHONE_ERROR:    "手机格式错误",
	S_JSON_ERROR:     "Json错误",
	S_GRPC_ERROR:     "GRPC错误",
	S_CONVERT_ERROR:  "http转发错误",

	// 文件类错误 20开头
	S_FILE_ERROR:        "文件错误",
	S_FILE_EMPTY:        "文件为空",
	S_FILE_UPLOAD_ERROR: "文件上传失败",

	//用户模块错误 100 开头
	S_INVALID_USERID:    "用户Id不存在",
	S_DUP_USERID:        "用户已经存在",
	S_INVALID_PASSWORD:  "密码错误",
	S_SHORT_PASSWORD:    "密码太短(<6)",
	S_PASSWORD_MISMATCH: "重复密码不一致",
	S_LOGIN_FAILED:      "用户不存在或密码错误",
	S_LOGIN_ELSEWHERE:   "用户已在别处登录",
	S_LOGIN_REQUIRED:    "请先登入",
	S_ACCOUNT_DISABLED:  "用户状态未激活",

	S_UNKNOW_ERROR: "未知错误",
	//业务类错误  200 开头
}

func ParseStatus(status string) (string, int) {
	if status == "" {
		return "", -1
	}

	sc, err := strconv.Atoi(status)
	if err != nil {
		return "", -1
	}

	return status_text[sc], sc
}
func GetMessage(status int) string {
	return status_text[status]
}

