package action

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
		if global.LocalConfig.Service.NotFoundEvent {
			log.Error().Int("status", event.Status).
				Str("method", event.Request.Method).
				Str("uri", event.Request.RequestURI).Msg(http.StatusText(404))
		}
	case 405:
		if global.LocalConfig.Service.MethodNotAllowedEvent {
			log.Error().Int("status", event.Status).
				Str("method", event.Request.Method).
				Str("uri", event.Request.RequestURI).Msg(http.StatusText(405))
		}
	case 500:
		e := log.Error()
		if global.LocalConfig.Service.Trigger {
			e.Str("caller", " "+event.Trigger.File+":"+strconv.Itoa(event.Trigger.Line)+" ").
				Str("func", event.Trigger.Func)
		}

		if global.LocalConfig.Service.Trace {
			var trace []string
			for k := range event.Trace {
				trace = append(trace, event.Trace[k])
			}
			e.Strs("trace", trace)
		}

		e.Err(event.Message)
	}

	if global.LocalConfig.Service.Debug {
		if _, err := event.ResponseWriter.Write([]byte(event.Message.Error())); err != nil {
			log.Err(err)
		}
	} else {
		if _, err := event.ResponseWriter.Write([]byte(http.StatusText(event.Status))); err != nil {
			log.Err(err)
		}
	}
}
