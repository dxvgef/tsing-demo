// 路由

package service

import (
	"local/handler"
)

// 设置路由
func setRouter() {
	// 在路由组中注册路由
	sessionRouter := engine.Group("/session")
	var sessionHandler handler.SessionHandler
	sessionRouter.PUT("/", sessionHandler.Put)
	sessionRouter.GET("/:key", sessionHandler.Get)
	sessionRouter.DELETE("/", sessionHandler.Destroy)

	databaseRouter := engine.Group("/database")
	var databaseHandler handler.DatabaseHandler
	databaseRouter.PUT("/", databaseHandler.CreateTable)
	databaseRouter.POST("/", databaseHandler.Add)
	databaseRouter.GET("/", databaseHandler.Get)
	databaseRouter.DELETE("/", databaseHandler.Delete)

	// 定义路由组，同时注册一个处理器做为中间件
	// secretRouter := engine.Group("/secret", handler.Auth)
}
