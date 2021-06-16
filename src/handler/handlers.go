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

// 写session
func Set(ctx *tsing.Context) error {
	sess, err := global.Sessions.Use(ctx.Request, ctx.ResponseWriter)
	if err != nil {
		return err
	}
	if err = sess.Put("test", "haha"); err != nil {
		return err
	}
	return String(ctx, 200, "ok")
}

// 读session
func Get(ctx *tsing.Context) error {
	sess, err := global.Sessions.Use(ctx.Request, ctx.ResponseWriter)
	if err != nil {
		return err
	}
	value, err2 := sess.Get("test").String()
	if err2 != nil {
		return err2
	}
	return String(ctx, 200, value)
}
