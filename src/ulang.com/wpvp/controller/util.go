package controller

import (
	"mime/multipart"
	"fmt"
	"io"
	"time"
	"bytes"
)

// CopyFile copy file
func ConvertFileToBytes(header *multipart.FileHeader) ([]byte, error) {
	buffer := new(bytes.Buffer)
	src, err := header.Open()
	if err != nil {
		fmt.Println("Open file error." + header.Filename)
		return nil, err
	}
	defer src.Close()
	io.Copy(buffer, src)
	return buffer.Bytes(), err
}


type TimeStamp int64

//func (t TimeStamp) String() string {
//	tm := time.Unix(t, 0)
//	return tm.Format("2006-01-02")
//}
//
//func foo() {
//	ts := TimeStamp(time.Now().Unix()).String()
//}
func GetDateTime(timestamp int64) string {
	//格式化为字符串,tm为Time类型
time.Now().Unix()
	tm := time.Unix(timestamp, 0)

	//fmt.Println(tm.Format("2006-01-02 03:04:05 PM"))

	//fmt.Println(tm.Format("02/01/2006 15:04:05 PM"))

	return tm.Format("2006-01-02")
}
