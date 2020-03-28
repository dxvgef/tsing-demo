package service

import (
	"local/handler"
)

// 设置路由
func setRouter() {

	example := new(handler.Example)

	// 注册一个静态目录路由
	app.Dir("/dir/", "./")
	// 注册一个静态文件路由
	app.File("/config", "./config.yaml")

	// 登录，获得token
	app.GET("/login", example.Login)

	// 读写session的演示
	// app.GET("/session", example.Session)

	// 注册一个根路径的路由组，并添加检查token的处理器
	secretRouter := app.Group("/", handler.CheckToken)
	// 组内的首页路由
	secretRouter.GET("", example.Index)

}
