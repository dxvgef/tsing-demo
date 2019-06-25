package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dxvgef/tsing"

	"src/action"
	"src/global"
)

var App *tsing.App

func Start() {
	App = tsing.New()
	App.Event.EnableTrace = global.Config.Logger.EnableTrace
	App.Event.ShortCaller = true
	App.Event.Handler = action.EventHandler

	// 设置路由
	setRouter()

	// 定义HTTP服务
	svr := &http.Server{
		Addr:    global.Config.Service.IP + ":" + strconv.Itoa(global.Config.Service.Port),
		Handler: App, // 调度器
		// ErrorLog:          global.ServiceLogger, // 错误日志记录器
		ReadTimeout:       10 * time.Second, // 读取超时
		WriteTimeout:      10 * time.Second, // 响应超时
		IdleTimeout:       10 * time.Second, // 连接空闲超时
		ReadHeaderTimeout: 10 * time.Second, // http header读取超时
	}

	// 在新协程中启动服务，方便实现退出等待
	go func() {
		global.ServiceLogger.Info("HTTP服务 " + svr.Addr + " 启动成功")
		if err := svr.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				global.ServiceLogger.Info("服务已退出")
			} else {
				global.ServiceLogger.Error(err.Error())
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
	if err := svr.Shutdown(ctx); err != nil {
		global.ServiceLogger.Error(err.Error())
	}
}
