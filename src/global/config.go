package global

import (
	"flag"
	"log"

	"github.com/fsnotify/fsnotify"
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

// 设置配置文件
func SetConfig(watch bool) error {
	configFilePath := flag.String("c", "./config.toml", "配置文件路径")
	flag.Parse()

	viper.SetConfigFile(*configFilePath) // 配置文件路径
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := LoadConfig(); err != nil {
		return err
	}

	if watch {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Println("配置文件发生变更，重载配置")
			if err := LoadConfig(); err != nil {
				log.Fatalln(err.Error())
				return
			}
		})
	}
	return nil
}

// 加载配置
func LoadConfig() (err error) {
	if err = viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		return
	}
	if err = viper.Unmarshal(&Config); err != nil {
		log.Println(err.Error())
		return
	}
	return nil
}
