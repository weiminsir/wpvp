package main

import (
	"os"
	"log"
	"ulang.com/wpvp/conf"
	"ulang.com/wpvp/router"
	"flag"
	"ulang.com/wpvp/model"
	"github.com/donnie4w/go-logger/logger"
)

func main() {

	mode := flag.String("mode", "dev", "mode")
	flag.Parse()
	if len(os.Args) == 2 {
		runmodel := os.Args[1]
		log.Println("run model:", runmodel)
	}
	conf.InitConf(*mode)
	initLog()
	model.Start()
	router.Start()

}

/**
 * 初始化日志设置
 */
func initLog() {
	logger.SetConsole(conf.Config.Log.PrintConsole)
	logger.SetRollingFile(conf.Config.Log.LogPath, conf.Config.Log.LogFileName,
		conf.Config.Log.MaxNumber, conf.Config.Log.MaxSize, logger.MB)

	switch conf.Config.Log.Level {
	case 0:
		logger.SetLevel(logger.ALL)
	case 1:
		logger.SetLevel(logger.DEBUG)
	case 2:
		logger.SetLevel(logger.INFO)
	case 3:
		logger.SetLevel(logger.WARN)
	case 4:
		logger.SetLevel(logger.ERROR)
	case 5:
		logger.SetLevel(logger.FATAL)
	case 6:
		logger.SetLevel(logger.OFF)
	default:
		logger.SetLevel(logger.DEBUG)
	}
}
