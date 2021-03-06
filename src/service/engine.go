package service

import (
	"context"
	"net"
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
func setService() (err error) {
	var (
		config   tsing.Config
		rootPath string
	)
	if global.Config.ServiceCenter.Addr != "" {
		if _, err = net.ResolveTCPAddr("tcp", global.Config.ServiceCenter.Addr); err != nil {
			log.Err(err).Caller().Str("addr", global.Config.ServiceCenter.Addr).Msg("服务中心参数值异常")
			return
		}
	}
	config.EventHandler = handler.EventHandler
	config.Recover = true
	config.EventShortPath = true
	if global.Config.Debug {
		config.EventSource = true
		config.EventTrace = true
	}
	config.EventHandlerError = true // 一定要处理handler返回的错误
	rootPath, err = os.Getwd()
	if err == nil {
		config.RootPath = rootPath + "/src/"
	}

	engine = tsing.New(&config)

	// 设置路由
	setRouter()

	// 如果是调试模式，注册pprof路由
	if global.Config.Debug {
		setDebugRouter()
	}

	// 设置HTTP服务
	if global.Config.Service.HTTPPort > 0 {
		httpServer = &http.Server{
			Addr:              global.Config.Service.IP + ":" + strconv.FormatUint(uint64(global.Config.Service.HTTPPort), 10),
			Handler:           engine,                                                               // 调度器
			ReadTimeout:       time.Duration(global.Config.Service.ReadTimeout) * time.Second,       // 读取超时
			WriteTimeout:      time.Duration(global.Config.Service.WriteTimeout) * time.Second,      // 响应超时
			IdleTimeout:       time.Duration(global.Config.Service.IdleTimeout) * time.Second,       // 连接空闲超时
			ReadHeaderTimeout: time.Duration(global.Config.Service.ReadHeaderTimeout) * time.Second, // http header读取超时
		}
	}
	return
}

func Start() (err error) {
	// 设置服务
	if err = setService(); err != nil {
		return
	}

	if global.Config.Service.HTTPPort > 0 {
		// 启动http服务
		go func() {
			log.Info().Str("监听地址", httpServer.Addr).Msg("启动HTTP服务")
			if err = httpServer.ListenAndServe(); err != nil {
				if err.Error() == http.ErrServerClosed.Error() {
					log.Info().Msg("HTTP服务已关闭")
				} else {
					log.Fatal().Err(err).Str("addr", httpServer.Addr).Caller().Msg("启动HTTP服务失败")
				}
				return
			}
		}()
	} else {
		log.Warn().Uint16("HTTP Port", global.Config.Service.HTTPPort).Msg("HTTP服务未启用")
	}

	// 设置服务中心
	if global.Config.ServiceCenter.Addr != "" {
		go SetCenter()
	}

	// 监听进程退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 退出HTTP服务
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Config.Service.QuitWaitTimeout)*time.Second)
	defer cancel()
	if httpServer != nil {
		if err = httpServer.Shutdown(ctx); err != nil {
			log.Err(err).Caller().Msg("HTTP服务关闭时出现异常")
		}
	}
	return
}
