package service

import (
	"net/http/pprof"

	"github.com/dxvgef/tsing"
)

// pprof路由
func pprofRouter() {
	router := app.Group("/debug/pprof")
	router.GET("/", indexHandler)
	router.GET("/heap", heapHandler)
	router.GET("/block", blockHandler)
	router.GET("/threadcreate", threadCreateHandler)
	router.GET("/cmdline", cmdlineHandler)
	router.GET("/profile", profileHandler)
	router.GET("/symbol", symbolHandler)
	router.POST("/symbol", symbolHandler)
	router.GET("/trace", traceHandler)
	router.GET("/mutex", mutexHandler)
	router.GET("/goroutine", goroutineHandler)
	router.GET("/allocs", allocsHandler)
}

func indexHandler(ctx *tsing.Context) error {
	pprof.Index(ctx.ResponseWriter, ctx.Request)
	return nil
}

func heapHandler(ctx *tsing.Context) error {
	pprof.Handler("heap").ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

func goroutineHandler(ctx *tsing.Context) error {
	pprof.Handler("goroutine").ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

func allocsHandler(ctx *tsing.Context) error {
	pprof.Handler("allocs").ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

func blockHandler(ctx *tsing.Context) error {
	pprof.Handler("block").ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

func threadCreateHandler(ctx *tsing.Context) error {
	pprof.Handler("threadcreate").ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

func cmdlineHandler(ctx *tsing.Context) error {
	pprof.Cmdline(ctx.ResponseWriter, ctx.Request)
	return nil
}

func profileHandler(ctx *tsing.Context) error {
	pprof.Profile(ctx.ResponseWriter, ctx.Request)
	return nil
}

func symbolHandler(ctx *tsing.Context) error {
	pprof.Symbol(ctx.ResponseWriter, ctx.Request)
	return nil
}

func traceHandler(ctx *tsing.Context) error {
	pprof.Trace(ctx.ResponseWriter, ctx.Request)
	return nil
}

func mutexHandler(ctx *tsing.Context) error {
	pprof.Handler("mutex").ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}
