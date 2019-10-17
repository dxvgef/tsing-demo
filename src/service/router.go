package service

import (
	"local/action"
)

// 设置路由
func setRouter() {

	handler := new(action.Example)

	app.Router.GET("/sign", handler.SignJWT)
	app.Router.GET("/session", handler.Session)
	{
		adminRouter := app.Router.GROUP("/admin", action.CheckJWT)
		adminRouter.GET("", handler.Admin)
	}
}
