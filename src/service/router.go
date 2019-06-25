package service

import "src/action"

// 设置路由
func setRouter() {
	{
		index := new(action.Index)
		App.Router.GET("/", index.Demo)
	}
}
