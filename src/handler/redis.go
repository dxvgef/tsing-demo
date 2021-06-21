package handler

import (
	"context"
	"time"

	"local/global"

	"github.com/dxvgef/filter/v2"
	"github.com/dxvgef/tsing"
)

type RedisHandler struct{}

// 写数据
func (*RedisHandler) Set(ctx *tsing.Context) error {
	var (
		err   error
		key   string
		value string
	)
	key, err = filter.String(ctx.Post("key"), "key").Require("").String()
	if err != nil {
		return ctx.String(400, "key参数不正确")
	}
	value, err = filter.String(ctx.Post("value"), "value").Require("").String()
	if err != nil {
		return ctx.String(400, "value参数不正确")
	}

	redisCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = global.RedisCli.Set(redisCtx, key, value, 0).Err(); err != nil {
		return ctx.Caller(err)
	}

	return ctx.Status(204)
}

// 读数据
func (*RedisHandler) Get(ctx *tsing.Context) (err error) {
	var (
		key   string
		value string
	)
	key, err = filter.String(ctx.Path("key"), "key").Require("").String()
	if err != nil {
		return ctx.String(400, "key参数不正确")
	}
	redisCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	value, err = global.RedisCli.Get(redisCtx, key).Result()
	if err != nil {
		return ctx.Caller(err)
	}
	return ctx.String(200, value)
}
