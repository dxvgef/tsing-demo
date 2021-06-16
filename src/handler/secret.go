package handler

import (
	"local/global"

	"github.com/dxvgef/tsing"
)

// 验证连接密码
func Auth(ctx *tsing.Context) error {
	authStr := ctx.Request.Header.Get("Authorization")
	if authStr != global.RuntimeConfig.Service.Secret {
		ctx.Abort()
		return ctx.Status(401)
	}
	return nil
}
