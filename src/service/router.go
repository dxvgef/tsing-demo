// 路由

package service

import (
	"local/handler"
)

// 设置路由
func setRouter() {
	// 定义路由组，同时注册一个处理器做为中间件
	secretRouter := engine.Group("", handler.Auth)

	// 在路由组中注册路由
	secretRouter.GET("/set", handler.Set)
	secretRouter.GET("/get", handler.Get)
}
