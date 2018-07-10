package model

import (
	"gopkg.in/mgo.v2"
	"ulang.com/wpvp/conf"
	"log"
	"github.com/donnie4w/go-logger/logger"
	"fmt"
)

var sess *mgo.Session
var mgo_srv string
var db_name string

func Start() {
	mgo_srv = fmt.Sprintf("%s:%s", conf.Config.Database.Db_host, conf.Config.Database.Db_port)
	db_name = conf.Config.Database.Db_name
	session, err := mgo.Dial(mgo_srv)
	if err != nil {
		logger.Fatal("connecting error:", err)
		log.Panicf("Mgo(%s): %v", mgo_srv, err)
	}
	sess = session
	sess.SetMode(mgo.Monotonic, true)
	InitSysCfg()
	InitSysAdmin()
}

func InitSysCfg() {
	_, err := NewSysCfg()
	if err != nil {
		logger.Fatal("Mgo(%s): %v", mgo_srv, err)
	}

}

func InitSysAdmin() {

	user, found := FindSysUser("admin")
	if !found {
		_, err := NewSysUser("admin", "123456")
		if err != nil {
			log.Panicf("new admin error:", err)
		}
	} else {
		logger.Info("login account=", user.Username)
	}
}
