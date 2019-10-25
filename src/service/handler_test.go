package service

import (
	"github.com/dxvgef/tsing"
	"io/ioutil"
	"local/action"
	"local/global"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var healthData = `
{
"description": "Golang Eureka Discovery Client",
"status": "UP"
}
`

func init() {
	log.SetFlags(log.Lshortfile)

	os.Args = append(os.Args, []string{"-c", "./../../config.toml"}...)

	// 设置全局配置
	if err := global.SetConfig(true); err != nil {
		log.Fatal(err.Error())
		return
	}

	// 设置日志记录器
	if err := global.SetLogger(); err != nil {
		log.Fatal(err.Error())
		return
	}
	Config()
}

//test tsing handlers
func TestTsingHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	app.Router.GET("/health", func(context *tsing.Context) error {
		action.String(context, http.StatusOK, healthData)
		return nil
	})
	app.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Error(err)
		return
	}
	dataset := string(data)
	if healthData != dataset {
		t.Errorf("test health handler err expect: %s ,actual:%s ", healthData, dataset)
	}
}
