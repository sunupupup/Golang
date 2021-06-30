package main

import "zap-demo/log"

//简单用zap和lumberjack实现了记录日志，并且进行日志切割的功能
//后面完善：将不同等级的日志放到不同文件夹
func main() {
	log.CreatDir()
	log.InitLogger("info")
	for {
		log.Logger.Info("info")
		log.Logger.Warn("warning")
		log.Logger.Error("error")
	}
}
