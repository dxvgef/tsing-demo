package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"local/global"
	"local/handler"

	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"
)

var engine *tsing.Engine
var httpServer *http.Server

// 配置引擎
func Config() {
	var config tsing.Config
	config.EventHandler = handler.EventHandler
	config.Recover = global.RuntimeConfig.Service.Recover
	config.EventShortPath = global.RuntimeConfig.Service.EventShortPath
	if global.RuntimeConfig.Service.Debug {
		config.EventSource = true
		config.EventTrace = true
	}
	config.EventHandlerError = true // 一定要处理handler返回的错误
	rootPath, err := os.Getwd()
	if err == nil {
		config.RootPath = rootPath + "/src/"
	}

	engine = tsing.New(&config)

	// 设置路由
	setRouter()

	// 如果是调试模式，注册pprof路由
	if global.RuntimeConfig.Service.Debug {
		setDebugRouter()
	}

	// 设置HTTP服务
	if global.RuntimeConfig.Service.HTTPPort > 0 {
		httpServer = &http.Server{
			Addr:              global.RuntimeConfig.Service.IP + ":" + strconv.FormatUint(uint64(global.RuntimeConfig.Service.HTTPPort), 10),
			Handler:           engine,                                                                      // 调度器
			ReadTimeout:       time.Duration(global.RuntimeConfig.Service.ReadTimeout) * time.Second,       // 读取超时
			WriteTimeout:      time.Duration(global.RuntimeConfig.Service.WriteTimeout) * time.Second,      // 响应超时
			IdleTimeout:       time.Duration(global.RuntimeConfig.Service.IdleTimeout) * time.Second,       // 连接空闲超时
			ReadHeaderTimeout: time.Duration(global.RuntimeConfig.Service.ReadHeaderTimeout) * time.Second, // http header读取超时
		}
	}
}

func Start() (err error) {
	// 配置服务
	Config()

	// 启动http服务
	if global.RuntimeConfig.Service.HTTPPort > 0 {
		go func() {
			log.Info().Str("addr", httpServer.Addr).Msg("启动HTTP服务")
			if err = httpServer.ListenAndServe(); err != nil {
				if err.Error() == http.ErrServerClosed.Error() {
					log.Info().Msg("HTTP服务已关闭")
				} else {
					log.Err(err).Caller().Msg("启动HTTPS服务失败")
				}
				return
			}
		}()
	}

	// 设置服务中心
	if global.RuntimeConfig.ServiceCenter.Addr != "" {
		go SetCenter()
	}

	// 监听进程退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 退出http服务
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.RuntimeConfig.Service.QuitWaitTimeout)*time.Second)
	defer cancel()
	if httpServer != nil {
		if err = httpServer.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error()) // nolint:gocritic
		}
	}
	return
}
