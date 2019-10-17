package service

import (
	"local/action"
)

// 设置路由
func setRouter() {

	handler := new(action.Example)

	App.Router.GET("/sign", handler.SignJWT)
	App.Router.GET("/session", handler.Session)
	{
		adminRouter := App.Router.GROUP("/admin", action.CheckJWT)
		adminRouter.GET("", handler.Admin)
	}
}
