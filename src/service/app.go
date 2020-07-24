package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"local/global"
	"local/handler"
)

var app *tsing.Engine
var appServer *http.Server

func Config() {
	var config tsing.Config
	config.EventHandler = handler.EventHandler
	config.Recover = global.LocalConfig.Service.Recover
	log.Debug().Bool("recover", config.Recover)
	config.EventShortPath = global.LocalConfig.Service.EventShortPath
	config.EventSource = global.LocalConfig.Service.EventSource
	config.EventTrace = global.LocalConfig.Service.EventTrace
	config.EventHandlerError = global.LocalConfig.Service.EventHandlerError
	rootPath, err := os.Getwd()
	if err == nil {
		config.RootPath = rootPath + "/src/"
	}

	app = tsing.New(config)

	// 设置路由
	setRouter()

	// 如果是调试模式，注册pprof路由
	if global.LocalConfig.Service.Debug {
		pprofRouter()
	}

	// 定义HTTP服务
	appServer = &http.Server{
		Addr:    global.LocalConfig.Service.IP + ":" + strconv.Itoa(global.LocalConfig.Service.Port),
		Handler: app, // 调度器
		// ErrorLog:          global.Logger.StdError,                                                    // 日志记录器
		ReadTimeout:       time.Duration(global.LocalConfig.Service.ReadTimeout) * time.Second,       // 读取超时
		WriteTimeout:      time.Duration(global.LocalConfig.Service.WriteTimeout) * time.Second,      // 响应超时
		IdleTimeout:       time.Duration(global.LocalConfig.Service.IdleTimeout) * time.Second,       // 连接空闲超时
		ReadHeaderTimeout: time.Duration(global.LocalConfig.Service.ReadHeaderTimeout) * time.Second, // http header读取超时
	}
}

func Start() {
	Config()
	// 在新协程中启动服务，方便实现退出等待
	go func() {
		log.Info().Msg("HTTP服务 " + appServer.Addr + " 启动成功")
		if err := appServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Info().Msg("服务已关闭")
				return
			}
			log.Fatal().Msg(err.Error())
		}
	}()

	// 退出进程时等待
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 指定退出超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.LocalConfig.Service.QuitWaitTimeout)*time.Second)
	defer cancel()
	if err := appServer.Shutdown(ctx); err != nil {
		log.Fatal().Caller().Msg(err.Error())
	}

	log.Info().Msg("程序已退出")
}
