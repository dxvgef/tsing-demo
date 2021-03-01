// 路由

package service

import (
	"local/handler"
)

// 设置路由
func setRouter() {
	engine.GET("/", handler.Demo)
	engine.GET("/demo2", handler.Demo2)
}
