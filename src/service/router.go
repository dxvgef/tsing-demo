package service

import (
	"local/handler"
)

// 设置路由
func setRouter() {

	example := new(handler.Example)

	app.GET("/login", example.SignToken)
	// app.GET("/session", example.Session)

	app.Dir("/dir/", "./")
	app.File("/config.yaml", "./config.yaml")

	adminRouter := app.Group("/", handler.CheckToken)
	adminRouter.GET("", example.Admin)

}
