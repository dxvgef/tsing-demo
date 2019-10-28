package global

import (
	"errors"
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
)

// 配置文件路径(仅在成功加载配置文件并且解析成功后有值)
var ConfigPath string

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
func SetConfig() error {
	if ConfigPath == "" {
		ConfigPath = *flag.String("c", "./config.toml", "配置文件路径")
		flag.Parse()
	}

	_, err := toml.DecodeFile(ConfigPath, &Config)
	if err != nil {
		ConfigPath = ""
		return err
	}

	return nil
}

// 监视配置文件更新
func WatchConfig() error {
	if ConfigPath == "" {
		return errors.New("配置文件没有成功解析，无法启动监视")
	}
	// 创建监视器
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		done := make(chan bool)
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						log.Println("重载配置文件")
						if err := SetConfig(); err != nil {
							panic(err.Error())
						}
					}
				case err := <-watcher.Errors:
					panic(err.Error())
				}
			}
		}()

		// 添加要监视的文件
		if err = watcher.Add(ConfigPath); err != nil {
			panic(err.Error())
		}
		<-done
	}()
	return nil
}
