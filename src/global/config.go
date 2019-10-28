package global

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

// 配置文件路径(仅在成功加载配置文件并且解析成功后有值)
var ConfigPath string

// Config 全局配置
var Config struct {
	Service struct {
		IP                    string `yaml:"ip",toml:"ip"`
		Port                  int    `yaml:"port",toml:"port"`
		ReadTimeout           int    `yaml:"readTimeout",toml:"readTimeout"`
		ReadHeaderTimeout     int    `yaml:"readHeaderTimeout",toml:"readHeaderTimeout"`
		WriteTimeout          int    `yaml:"writeTimeout",toml:"writeTimeout"`
		IdleTimeout           int    `yaml:"idleTimeout",toml:"idleTimeout"`
		QuitWaitTimeout       int    `yaml:"quitWaitTimeout",toml:"quitWaitTimeout"`
		Debug                 bool   `yaml:"debug",toml:"debug"`
		EventTrigger          bool   `yaml:"eventTrigger",toml:"eventTrigger"`
		EventTrace            bool   `yaml:"eventTrace",toml:"eventTrace"`
		NotFoundEvent         bool   `yaml:"notFoundEvent",toml:"notFoundEvent"`
		MethodNotAllowedEvent bool   `yaml:"methodNotAllowedEvent",toml:"methodNotAllowedEvent"`
		HandleOPTIONS         bool   `yaml:"handleOPTIONS",toml:"handleOPTIONS"`
		FixPath               bool   `yaml:"fixPath",toml:"fixPath"`
		RedirectTrailingSlash bool   `yaml:"redirectTrailingSlash",toml:"redirectTrailingSlash"`
	} `yaml:"service",toml:"service"`
	Logger struct {
		Level      string `yaml:"level",toml:"level"`
		Outputs    string `yaml:"outputs",toml:"outputs"`
		Encode     string `yaml:"encode",toml:"encode"`
		ColorLevel bool   `yaml:"colorLevel",toml:"colorLevel"`
	} `yaml:"logger",toml:"logger"`
	Snowflake struct {
		Epoch int64 `yaml:"epoch",toml:"epoch"`
		Node  int64 `yaml:"node",toml:"node"`
	} `yaml:"snowflake",toml:"snowflake"`
	Database struct {
		Addr         string `yaml:"addr",toml:"addr"`
		User         string `yaml:"user",toml:"user"`
		Password     string `yaml:"password",toml:"password"`
		Name         string `yaml:"name",toml:"name"`
		StmtLog      bool   `yaml:"stmtLog",toml:"stmtLog"`
		ReadTimeout  int    `yaml:"readTimeout",toml:"readTimeout"`
		WriteTimeout int    `yaml:"writeTimeout",toml:"writeTimeout"`
		PoolSize     int    `yaml:"poolSize",toml:"poolSize"`
	} `yaml:"database",toml:"database"`
	Session struct {
		Key            string `yaml:"key",toml:"key"`
		CookieName     string `yaml:"cookieName",toml:"cookieName"`
		HTTPOnly       bool   `yaml:"httpOnly",toml:"httpOnly"`
		Secure         bool   `yaml:"secure",toml:"secure"`
		MaxAge         int    `yaml:"maxAge",toml:"maxAge"`
		IdleTime       int    `yaml:"idleTime",toml:"idleTime"`
		RedisAddr      string `yaml:"redisAddr",toml:"redisAddr"`
		RedisDB        int    `yaml:"redisDB",toml:"redisDB"`
		RedisKeyPrefix string `yaml:"redisKeyPrefix",toml:"redisKeyPrefix"`
	} `yaml:"session",toml:"session"`
}

// 设置TOML配置文件
func SetTOMLConfig() error {
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

// 设置YAML配置文件
func SetYAMLConfig() error {
	if ConfigPath == "" {
		ConfigPath = *flag.String("c", "./config.yaml", "配置文件路径")
		flag.Parse()
	}
	file, err := os.Open(ConfigPath)
	if err != nil {
		return err
	}
	err = yaml.NewDecoder(file).Decode(&Config)
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
						// if err := SetTOMLConfig(); err != nil {
						// 	panic(err.Error())
						// }
						if err := SetYAMLConfig(); err != nil {
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
