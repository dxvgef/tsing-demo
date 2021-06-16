// 路由

package service

import (
	"local/handler"
)

// 设置路由
func setRouter() {
	// 在路由组中注册路由
	testSession := engine.Group("/session")
	var ts handler.TestSession
	testSession.PUT("/", ts.Put)
	testSession.GET("/:key", ts.Get)
	testSession.DELETE("/", ts.Destroy)

	// 定义路由组，同时注册一个处理器做为中间件
	// secretRouter := engine.Group("/secret", handler.Auth)
}
