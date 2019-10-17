package action

import (
	"net/http"

	"github.com/dxvgef/tsing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"local/global"
)

// tsing的事件处理器
func EventHandler(event *tsing.Event) {
	// 先输出状态码
	event.ResponseWriter.WriteHeader(event.Status)

	// 根据状态码做不同的日志处理
	switch event.Status {
	case 404:
		if global.Config.Service.NotFoundEvent == true {
			global.Logger.Default.Error(
				event.Message.Error(),
				zap.Int("status", event.Status),
				zap.String("method", event.Request.Method),
				zap.String("uri", event.Request.RequestURI),
			)
		}
	case 405:
		if global.Config.Service.MethodNotAllowedEvent == true {
			global.Logger.Default.Warn(
				event.Message.Error(),
				zap.Int("status", event.Status),
				zap.String("method", event.Request.Method),
				zap.String("uri", event.Request.RequestURI),
			)
		}
	case 500:
		var fields []zapcore.Field
		fields = append(fields, zap.Int("status", event.Status))
		fields = append(fields, zap.String("file", event.Trigger.File))
		fields = append(fields, zap.Int("line", event.Trigger.Line))
		fields = append(fields, zap.String("func", event.Trigger.Func))
		if global.Config.Service.EventTrigger == true {
			fields = append(fields, zap.String("file", event.Trigger.File))
			fields = append(fields, zap.Int("line", event.Trigger.Line))
			fields = append(fields, zap.String("func", event.Trigger.Func))
		}

		if global.Config.Service.EventTrace == true {
			var trace []string
			for k := range event.Trace {
				trace = append(trace, event.Trace[k])
			}
			fields = append(fields, zap.Strings("trace", trace))
		}

		global.Logger.Default.Error(event.Message.Error(), fields...)
	}

	if global.Config.Service.Debug == true {
		if _, err := event.ResponseWriter.Write([]byte(event.Message.Error())); err != nil {
			global.Logger.Default.Error(err.Error())
		}
	} else {
		if _, err := event.ResponseWriter.Write([]byte(http.StatusText(event.Status))); err != nil {
			global.Logger.Default.Error(err.Error())
		}
	}
}
