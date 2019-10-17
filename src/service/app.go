package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dxvgef/tsing"

	"local/action"
	"local/global"
)

var App *tsing.App
var AppServer *http.Server

func Config() {
	App = tsing.New()
	App.Config.HandleOPTIONS = global.Config.Service.HandleOPTIONS
	App.Config.FixPath = global.Config.Service.FixPath
	App.Config.EventTrace = global.Config.Service.EventTrace
	App.Config.EventTrigger = global.Config.Service.EventTrace
	App.Config.RedirectTrailingSlash = global.Config.Service.RedirectTrailingSlash
	App.Config.EventHandler = action.EventHandler

	// 设置路由
	setRouter()

	// 定义HTTP服务
	AppServer = &http.Server{
		Addr:              global.Config.Service.IP + ":" + strconv.Itoa(global.Config.Service.Port),
		Handler:           App,                                                                  // 调度器
		ErrorLog:          global.Logger.StdError,                                               // 日志记录器
		ReadTimeout:       time.Duration(global.Config.Service.ReadTimeout) * time.Second,       // 读取超时
		WriteTimeout:      time.Duration(global.Config.Service.WriteTimeout) * time.Second,      // 响应超时
		IdleTimeout:       time.Duration(global.Config.Service.IdleTimeout) * time.Second,       // 连接空闲超时
		ReadHeaderTimeout: time.Duration(global.Config.Service.ReadHeaderTimeout) * time.Second, // http header读取超时
	}
}

func Start() {
	Config()
	// 在新协程中启动服务，方便实现退出等待
	go func() {
		global.Logger.Default.Info("HTTP服务 " + AppServer.Addr + " 启动成功")
		if err := AppServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				global.Logger.Default.Info("服务已退出")
			} else {
				global.Logger.Caller.Fatal(err.Error())
			}
		}
	}()

	// 退出进程时等待
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// 指定退出超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Config.Service.QuitWaitTimeout)*time.Second)
	defer cancel()
	if err := AppServer.Shutdown(ctx); err != nil {
		global.Logger.Caller.Fatal(err.Error())
	}
}
