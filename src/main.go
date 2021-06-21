package main

import (
	"local/global"
	"local/service"

	"github.com/rs/zerolog/log"
)

func main() {
	var err error

	// 设置默认logger
	global.SetDefaultLogger()

	// 设置服务ID
	if err = global.SetServiceID("tsing-demo"); err != nil {
		log.Fatal().Msg("设置服务ID失败")
		return
	}

	// 加载配置
	if err = global.LoadConfig(); err != nil {
		log.Fatal().Msg("加载配置失败")
		return
	}

	// 设置logger
	if err = global.SetLogger(); err != nil {
		log.Fatal().Msg("设置Logger失败")
		return
	}

	// 设置snowflake
	if err = global.SetSnowflake(); err != nil {
		log.Fatal().Msg("设置snowflake节点失败")
		return
	}

	// 创建RSA密钥
	if err = global.MakeRSAKey(); err != nil {
		log.Fatal().Msg("生成RSA私钥失败")
		return
	}

	// 设置sessions
	if err = global.SetSessions(); err != nil {
		log.Fatal().Msg("设置Session失败")
		return
	}

	// 设置数据库
	if err = global.SetDatabase(); err != nil {
		log.Fatal().Msg("设置数据库失败")
		return
	}

	// 设置Redis
	if err = global.SetRedis(); err != nil {
		log.Fatal().Msg("设置Redis失败")
		return
	}

	// 启动服务，此函数会阻塞
	if err = service.Start(); err != nil {
		log.Fatal().Msg("服务启动异常")
	}
}
