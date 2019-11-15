package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"local/action"
	"local/global"

	"github.com/dxvgef/tsing"
)

var app *tsing.App
var appServer *http.Server

func Config() {
	app = tsing.New()
	app.Config.HandleOPTIONS = global.LocalConfig.Service.HandleOPTIONS
	app.Config.FixPath = global.LocalConfig.Service.FixPath
	app.Config.EventTrace = global.LocalConfig.Service.EventTrace
	app.Config.EventTrigger = global.LocalConfig.Service.EventTrigger
	app.Config.RedirectTrailingSlash = global.LocalConfig.Service.RedirectTrailingSlash
	app.Config.EventHandler = action.EventHandler

	// 禁用panic处理器，提升性能

	// 设置路由
	setRouter()

	// 如果是调试模式，注册pprof路由
	if global.LocalConfig.Service.Debug {
		pprofRouter()
	}

	// 定义HTTP服务
	appServer = &http.Server{
		Addr:              global.LocalConfig.Service.IP + ":" + strconv.Itoa(global.LocalConfig.Service.Port),
		Handler:           app,                                                                       // 调度器
		ErrorLog:          global.Logger.StdError,                                                    // 日志记录器
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
		global.Logger.Default.Info("HTTP服务 " + appServer.Addr + " 启动成功")
		if err := appServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				global.Logger.Default.Info("服务已退出")
			} else {
				global.Logger.Caller.Fatal(err.Error())
			}
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
		global.Logger.Caller.Fatal(err.Error())
	}
}
