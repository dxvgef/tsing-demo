// 路由

package service

import (
	"local/handler"
)

// 设置路由
func setRouter() {
	engine.GET("/set", handler.Set)
	engine.GET("/get", handler.Get)
}
