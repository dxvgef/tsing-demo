package global

import (
	"errors"
	"flag"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config 全局配置
var Config struct {
	Service struct {
		IP                    string
		Port                  int
		ReadTimeout           int
		ReadHeaderTimeout     int
		WriteTimeout          int
		IdleTimeout           int
		QuitWaitTimeout       int
		Debug                 bool
		EventTrigger          bool
		EventTrace            bool
		NotFoundEvent         bool
		MethodNotAllowedEvent bool
		HandleOPTIONS         bool
		FixPath               bool
		RedirectTrailingSlash bool
	}
	Logger struct {
		Level      string
		Outputs    string
		Encode     string
		ColorLevel bool
	}
	Snowflake struct {
		Epoch int64
		Node  int64
	}
	Database struct {
		Addr         string
		User         string
		Password     string
		Name         string
		StmtLog      bool
		ReadTimeout  int
		WriteTimeout int
		PoolSize     int
	}
	Session struct {
		Key            string
		CookieName     string
		HTTPOnly       bool
		Secure         bool
		MaxAge         int
		IdleTime       int
		RedisAddr      string
		RedisDB        int
		RedisKeyPrefix string
	}
}

// 解析配置文件路径
func parseFilePath(configFilePath string) (configName string, configType string) {
	if configFilePath == "" {
		return
	}
	fullExt := filepath.Ext(configFilePath)
	if fullExt == "" {
		return
	}
	extArr := strings.Split(fullExt, ".")
	if len(extArr) > 1 {
		configType = extArr[1]
	}
	extPos := strings.LastIndex(configFilePath, fullExt)
	configName = configFilePath[:extPos]
	return
}

// 设置全局配置变量
func SetConfig() error {
	configFilePath := flag.String("c", "./config.toml", "配置文件路径")
	flag.Parse()

	configName, configType := parseFilePath(*configFilePath)
	if configName == "" || configType == "" {
		return errors.New("无法读取配置文件")
	}
	viper.SetConfigName(configName) // 配置文件名，不需要扩展名
	viper.SetConfigType(configType) // 文件类型
	viper.AddConfigPath(".")        // 文件路径

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return viper.Unmarshal(&Config)
}
