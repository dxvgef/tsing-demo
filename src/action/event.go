package action

import (
	"net/http"

	"github.com/dxvgef/tsing"

	"src/global"
)

func EventHandler(event tsing.Event) {
	// 状态码
	event.ResponseWriter.WriteHeader(event.Status)

	switch event.Status {
	case 404:
		// 如果配置文件禁用404记录
		if global.Config.Logger.EnableNotFound == true {
			global.ServiceLogger.Error(event.Request.Method + " " + event.Request.RequestURI + " " + http.StatusText(event.Status))
		}
	case 405:
		if global.Config.Logger.EnableMethodNotAllowed == true {
			global.ServiceLogger.Error(event.Request.Method + " " + event.Request.RequestURI + " " + http.StatusText(event.Status))
		}
	case 500:
		var trace string
		if global.Config.Logger.EnableTrace == true {
			l := len(event.Trace)
			if l == 1 {
				trace = event.Trace[0]
			} else {
				for i := 0; i < l; i++ {
					if event.Trace[i] != ":0" {
						if i > 0 {
							trace += "\r\n"
						}
						trace += event.Trace[i]
					}
				}
			}
		}
		if len(event.Trace) == 1 {
			global.ServiceLogger.Error(trace + " " + event.Message.Error())
		} else {
			global.ServiceLogger.Error(event.Message.Error() + "\r\n" + trace)
		}
	}

	if global.Config.Service.Debug == true {
		event.ResponseWriter.Write([]byte(event.Message.Error()))
	} else {
		event.ResponseWriter.Write([]byte(http.StatusText(event.Status)))
	}
}
