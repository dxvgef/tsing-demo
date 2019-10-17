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

var app *tsing.App
var appServer *http.Server

func Config() {
	app = tsing.New()
	app.Config.HandleOPTIONS = global.Config.Service.HandleOPTIONS
	app.Config.FixPath = global.Config.Service.FixPath
	app.Config.EventTrace = global.Config.Service.EventTrace
	app.Config.EventTrigger = global.Config.Service.EventTrigger
	app.Config.RedirectTrailingSlash = global.Config.Service.RedirectTrailingSlash
	app.Config.EventHandler = action.EventHandler

	// 设置路由
	setRouter()

	// 定义HTTP服务
	appServer = &http.Server{
		Addr:              global.Config.Service.IP + ":" + strconv.Itoa(global.Config.Service.Port),
		Handler:           app,                                                                  // 调度器
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
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// 指定退出超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Config.Service.QuitWaitTimeout)*time.Second)
	defer cancel()
	if err := appServer.Shutdown(ctx); err != nil {
		global.Logger.Caller.Fatal(err.Error())
	}
}
