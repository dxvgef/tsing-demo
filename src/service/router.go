// 路由

package service

import (
	"local/handler"
)

// 设置路由
func setRouter() {
	// 在路由组中注册路由
	var sessionHandler handler.SessionHandler
	sessionRouter := engine.Group("/session")
	sessionRouter.PUT("/", sessionHandler.Put)
	sessionRouter.GET("/:key", sessionHandler.Get)
	sessionRouter.DELETE("/", sessionHandler.Destroy)

	var databaseHandler handler.DatabaseHandler
	databaseRouter := engine.Group("/database")
	databaseRouter.PUT("/", databaseHandler.CreateTable)
	databaseRouter.POST("/", databaseHandler.Add)
	databaseRouter.GET("/", databaseHandler.Get)
	databaseRouter.DELETE("/", databaseHandler.Delete)

	var jwtHandler handler.JWTHandler
	jwtRouter := engine.Group("/jwt")
	jwtRouter.POST("/", jwtHandler.Sign)
	jwtRouter.GET("/", jwtHandler.Auth)

	var redisHandler handler.RedisHandler
	redisRouter := engine.Group("/redis")
	redisRouter.POST("/", redisHandler.Set)
	redisRouter.GET("/:key", redisHandler.Get)

	// 定义路由组，同时注册一个处理器做为中间件
	// secretRouter := engine.Group("/secret", handler.Auth)
}
