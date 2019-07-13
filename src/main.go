package main

import (
	"log"

	"src/global"
	"src/service"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// 设置全局配置
	if err := global.SetConfig(); err != nil {
		log.Println(err.Error())
		return
	}

	// 设置日志记录器
	if err := global.SetLogger(); err != nil {
		log.Println(err.Error())
		return
	}

	// 设置ID节点
	if err := global.SetIDnode(); err != nil {
		global.ServiceLogger.Error(err.Error())
		return
	}

	// 设置数据库
	if err := global.SetDatabase(); err != nil {
		global.ServiceLogger.Error(err.Error())
		return
	}

	// 设置Session
	if err := global.SetSessions(); err != nil {
		global.ServiceLogger.Error(err.Error())
		return
	}

	// 启动服务
	service.Start()
}
