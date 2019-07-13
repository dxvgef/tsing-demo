package service

import (
	"src/action"
)

// 设置路由
func setRouter() {
	{
		index := new(action.Example)
		App.Router.GET("/", index.SignJWT)
		App.Router.GET("/session", index.Session)
		{
			admin := App.Router.GROUP("/admin", action.CheckJWT)
			admin.GET("", index.Admin)
		}
	}
}
