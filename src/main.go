package main

import (
	"log"

	"src/global"
	"src/service"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// 设置全局配置
	err := global.SetConfig()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 设置日志记录器
	err = global.SetLogger()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 启动服务
	service.Start()
}
