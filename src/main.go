package main

import (
	"local/global"
	"local/service"
)

func main() {
	// 设置默认logger
	global.SetDefaultLogger()

	// 加载配置
	if err := global.LoadConfig(); err != nil {
		return
	}

	// 根据配置设置logger
	if err := global.SetLogger(); err != nil {
		return
	}

	// 设置sessions
	if err := global.SetSessions(); err != nil {
		return
	}

	// 设置数据库
	if err := global.SetDatabase(); err != nil {
		return
	}

	// 启动服务
	if err := service.Start(); err != nil {
		return
	}
}
