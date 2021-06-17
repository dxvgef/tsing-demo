package handler

import (
	"local/global"

	"github.com/dxvgef/filter/v2"
	"github.com/dxvgef/sessions"
	"github.com/dxvgef/tsing"
)

type SessionHandler struct{}

// 写session
func (*SessionHandler) Put(ctx *tsing.Context) error {
	var (
		err   error
		key   string
		value string
		sess  *sessions.Session
	)
	key, err = filter.String(ctx.Post("key"), "key").Require("").String()
	if err != nil {
		return ctx.String(400, "key参数不正确")
	}
	value, err = filter.String(ctx.Post("value"), "value").Require("").String()
	if err != nil {
		return ctx.String(400, "value参数不正确")
	}

	sess, err = global.Sessions.Use(ctx.Request, ctx.ResponseWriter)
	if err != nil {
		return err
	}
	if err = sess.Put(key, value); err != nil {
		return ctx.Caller(err)
	}
	return ctx.Status(204)
}

// 读session
func (*SessionHandler) Get(ctx *tsing.Context) (err error) {
	var (
		sess  *sessions.Session
		key   string
		value string
	)
	key, err = filter.String(ctx.Path("key"), "key").Require("").String()
	if err != nil {
		return ctx.String(400, "key参数不正确")
	}
	sess, err = global.Sessions.Use(ctx.Request, ctx.ResponseWriter)
	if err != nil {
		return ctx.Caller(err)
	}
	value, err = sess.Get(key).String()
	if err != nil {
		return ctx.Caller(err)
	}
	return ctx.String(200, value)
}

// 销毁session
func (*SessionHandler) Destroy(ctx *tsing.Context) (err error) {
	var sess *sessions.Session
	sess, err = global.Sessions.Use(ctx.Request, ctx.ResponseWriter)
	if err != nil {
		return ctx.Caller(err)
	}
	if err = sess.Destroy(); err != nil {
		return ctx.Caller(err)
	}
	return ctx.Status(204)
}
