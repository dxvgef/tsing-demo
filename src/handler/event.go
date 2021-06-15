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
	// 先输出状态码
	event.ResponseWriter.WriteHeader(event.Status)

	// 根据状态码做不同的日志处理
	switch event.Status {
	case 404:
		if global.RuntimeConfig.Service.EventNotFound {
			log.Error().Int("status", event.Status).
				Str("method", event.Request.Method).
				Str("uri", event.Request.RequestURI).Msg(http.StatusText(http.StatusNotFound))
		}
	case 405:
		if global.RuntimeConfig.Service.EventMethodNotAllowed {
			log.Error().Int("status", event.Status).
				Str("method", event.Request.Method).
				Str("uri", event.Request.RequestURI).Msg(http.StatusText(http.StatusMethodNotAllowed))
		}
	case 500:
		e := log.Error()
		if global.RuntimeConfig.Service.Debug {
			e.Str("caller", " "+event.Source.File+":"+strconv.Itoa(event.Source.Line)+" ").
				Str("func", event.Source.Func)
		}

		if global.RuntimeConfig.Service.Debug {
			var trace []string
			for k := range event.Trace {
				trace = append(trace, event.Trace[k])
			}
			e.Strs("trace", trace)
		}

		e.Msg(event.Message.Error())
	}

	responseMsg := ""
	if global.RuntimeConfig.Service.Debug {
		responseMsg = event.Message.Error()
	} else {
		responseMsg = http.StatusText(event.Status)
	}
	if _, err := event.ResponseWriter.Write([]byte(responseMsg)); err != nil {
		log.Error().Msg(err.Error())
	}
}
