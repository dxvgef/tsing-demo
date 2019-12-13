package main

import (
	"errors"
	stdLog "log"

	"github.com/rs/zerolog/log"

	"local/global"
	"local/service"
)

func main() {
	stdLog.SetFlags(stdLog.Lshortfile)

	var err error

	// 加载TOML配置文件
	// if err = global.LoadTOMLConfig(); err != nil {
	// 	log.Fatal(err.Error())
	// 	return
	// }
	// 加载YAML配置文件
	if err = global.LoadYAMLConfig(); err != nil {
		stdLog.Fatal(err.Error())
		return
	}

	// 监视并热更新配置
	// if err = global.WatchConfig(); err != nil {
	// 	stdLog.Fatal(err.Error())
	// 	return
	// }

	// 设置日志记录器
	if err = global.SetLogger(); err != nil {
		stdLog.Fatal(err.Error())
		return
	}

	log.Fatal().Err(errors.New("测试一下"))

	// 设置ID节点
	if err = global.SetIDnode(); err != nil {
		log.Fatal().Err(err)
		return
	}

	// 设置数据库
	if err = global.SetDatabase(); err != nil {
		log.Fatal().Err(err)
		return
	}

	// 设置Session
	if err = global.SetSessions(); err != nil {
		log.Fatal().Err(err)
		return
	}

	// 设置etcd client
	// if err = global.SetETCDClient(); err != nil {
	// 	global.Logger.Caller.Fatal(err.Error())
	// 	return
	// }

	// 从etcd加载远程配置
	// if err = global.LoadRemoteConfigFromETCD(); err != nil {
	// 	global.Logger.Caller.Fatal(err.Error())
	// 	return
	// }

	// 启动服务
	service.Start()
}
