// 路由

package service

import (
	"local/global"
	"local/handler"
)

// 设置路由
func setRouter() {
	if global.RuntimeConfig.Service.Debug {
		setDebugRouter()
	}
	engine.GET("/set", handler.Set)
	engine.GET("/get", handler.Get)
}
