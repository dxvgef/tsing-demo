package service

import (
	"local/handler"
)

// 设置路由
func setRouter() {

	example := new(handler.Example)

	app.Router.GET("/login", example.SignToken)
	// app.Router.GET("/session", example.Session)

	adminRouter := app.Group("/", handler.CheckToken)
	adminRouter.GET("", example.Admin)

}
