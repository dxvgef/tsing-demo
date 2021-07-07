package handler

import (
	"net/http"
	"strconv"

	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"local/global"
)

// tsing的事件处理器
func EventHandler(event *tsing.Event) {
	// 响应状态码
	event.ResponseWriter.WriteHeader(event.Status)

	// 根据状态码做不同的日志记录
	switch event.Status {
	case 404:
		if global.Config.Service.NotFoundEvent {
			log.Error().Int("status", event.Status).
				Str("method", event.Request.Method).
				Str("uri", event.Request.RequestURI).Msg(http.StatusText(http.StatusNotFound))
		}
	case 405:
		if global.Config.Service.MethodNotAllowedEvent {
			log.Error().Int("status", event.Status).
				Str("method", event.Request.Method).
				Str("uri", event.Request.RequestURI).Msg(http.StatusText(http.StatusMethodNotAllowed))
		}
	case 500:
		e := log.Error()
		if global.Config.Debug {
			e.Str("caller", " "+event.Source.File+":"+strconv.Itoa(event.Source.Line)+" ").
				Str("func", event.Source.Func)
		}
		// 如果启用了调试
		if global.Config.Debug {
			// 输出trace
			var trace []string
			for k := range event.Trace {
				trace = append(trace, event.Trace[k])
			}
			e.Strs("trace", trace)
		}

		e.Msg(event.Message.Error())
	}

	// 响应正文
	responseMsg := ""
	if global.Config.Debug {
		responseMsg = event.Message.Error()
	} else {
		responseMsg = http.StatusText(event.Status)
	}
	if _, err := event.ResponseWriter.Write(global.StrToBytes(responseMsg)); err != nil {
		log.Error().Msg(err.Error())
	}
}
