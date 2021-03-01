package handler

import (
	"local/global"

	"github.com/dxvgef/tsing"
)

func Demo(ctx *tsing.Context) error {
	sess, err := global.Session.Use(ctx.Request, ctx.ResponseWriter)
	if err != nil {
		return err
	}
	if err = sess.Set("test", "haha"); err != nil {
		return err
	}
	return String(ctx, 200, "ok")
}

func Demo2(ctx *tsing.Context) error {
	sess, err := global.Session.Use(ctx.Request, ctx.ResponseWriter)
	if err != nil {
		return err
	}
	value, err2 := sess.Get("test").String()
	if err2 != nil {
		return err2
	}
	return String(ctx, 200, value)
}
