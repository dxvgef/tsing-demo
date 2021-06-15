package main

import (
	"local/global"
	"local/service"

	"github.com/rs/zerolog/log"
)

func main() {
	// 设置默认logger
	global.SetDefaultLogger()

	// 加载配置
	if err := global.LoadConfig(); err != nil {
		log.Fatal().Err(err).Msg("加载配置失败")
	}

	// 根据配置设置logger
	if err := global.SetLogger(); err != nil {
		log.Fatal().Err(err).Msg("设置Logger失败")
	}

	// 设置sessions
	if err := global.SetSessions(); err != nil {
		log.Fatal().Err(err).Msg("设置Sessions失败")
	}

	// 设置数据库
	if err := global.SetDatabase(); err != nil {
		log.Fatal().Err(err).Msg("设置数据库失败")
	}

	// 启动服务
	if err := service.Start(); err != nil {
		log.Fatal().Err(err).Msg("启动服务失败")
	}
}
