package service

import "src/action"

// 设置路由
func setRouter() {
	{
		manager := new(action.Manager)
		App.Router.GET("/auth", manager.Auth)
	}
}
