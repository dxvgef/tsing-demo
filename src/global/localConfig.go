package global

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

// 配置文件路径(仅在成功加载配置文件并且解析成功后有值)
var LocalConfigPath string

// LocalConfig 全局配置
var LocalConfig struct {
	Service struct {
		IP                    string `yaml:"ip" toml:"ip"`
		Port                  int    `yaml:"port" toml:"port"`
		ReadTimeout           int    `yaml:"readTimeout" toml:"readTimeout"`
		ReadHeaderTimeout     int    `yaml:"readHeaderTimeout" toml:"readHeaderTimeout"`
		WriteTimeout          int    `yaml:"writeTimeout" toml:"writeTimeout"`
		IdleTimeout           int    `yaml:"idleTimeout" toml:"idleTimeout"`
		QuitWaitTimeout       int    `yaml:"quitWaitTimeout" toml:"quitWaitTimeout"`
		Debug                 bool   `yaml:"debug" toml:"debug"`
		EventTrigger          bool   `yaml:"eventTrigger" toml:"eventTrigger"`
		EventTrace            bool   `yaml:"eventTrace" toml:"eventTrace"`
		NotFoundEvent         bool   `yaml:"notFoundEvent" toml:"notFoundEvent"`
		MethodNotAllowedEvent bool   `yaml:"methodNotAllowedEvent" toml:"methodNotAllowedEvent"`
		HandleOPTIONS         bool   `yaml:"handleOPTIONS" toml:"handleOPTIONS"`
		FixPath               bool   `yaml:"fixPath" toml:"fixPath"`
		RedirectTrailingSlash bool   `yaml:"redirectTrailingSlash" toml:"redirectTrailingSlash"`
	} `yaml:"service" toml:"service"`
	Logger struct {
		Level      string `yaml:"level" toml:"level"`
		Outputs    string `yaml:"outputs" toml:"outputs"`
		Encode     string `yaml:"encode" toml:"encode"`
		ColorLevel bool   `yaml:"colorLevel" toml:"colorLevel"`
	} `yaml:"logger" toml:"logger"`
	Snowflake struct {
		Epoch int64 `yaml:"epoch" toml:"epoch"`
		Node  int64 `yaml:"node" toml:"node"`
	} `yaml:"snowflake" toml:"snowflake"`
	Database struct {
		Addr         string `yaml:"addr" toml:"addr"`
		User         string `yaml:"user" toml:"user"`
		Password     string `yaml:"password" toml:"password"`
		Name         string `yaml:"name" toml:"name"`
		StmtLog      bool   `yaml:"stmtLog" toml:"stmtLog"`
		ReadTimeout  int    `yaml:"readTimeout" toml:"readTimeout"`
		WriteTimeout int    `yaml:"writeTimeout" toml:"writeTimeout"`
		PoolSize     int    `yaml:"poolSize" toml:"poolSize"`
	} `yaml:"database" toml:"database"`
	Session struct {
		Key            string `yaml:"key" toml:"key"`
		CookieName     string `yaml:"cookieName" toml:"cookieName"`
		HTTPOnly       bool   `yaml:"httpOnly" toml:"httpOnly"`
		Secure         bool   `yaml:"secure" toml:"secure"`
		MaxAge         int    `yaml:"maxAge" toml:"maxAge"`
		IdleTime       int    `yaml:"idleTime" toml:"idleTime"`
		RedisAddr      string `yaml:"redisAddr" toml:"redisAddr"`
		RedisDB        int    `yaml:"redisDB" toml:"redisDB"`
		RedisKeyPrefix string `yaml:"redisKeyPrefix" toml:"redisKeyPrefix"`
	} `yaml:"session" toml:"session"`
	ETCD struct {
		Endpoints               []string `yaml:"endpoints" toml:"endpoints"`
		Username                string   `yaml:"username" toml:"username"`
		Password                string   `yaml:"password" toml:"password"`
		HeaderTimeoutPerRequest int      `yaml:"headerTimeoutPerRequest" toml:"headerTimeoutPerRequest"`
		KeyPrefix               string   `yaml:"keyPrefix" toml:"keyPrefix"`
	} `yaml:"etcd" toml:"etcd"`
}

// 加载TOML配置文件
func LoadTOMLConfig() error {
	if LocalConfigPath == "" {
		LocalConfigPath = *flag.String("c", "./config.toml", "配置文件路径")
		flag.Parse()
	}

	_, err := toml.DecodeFile(LocalConfigPath, &LocalConfig)
	if err != nil {
		LocalConfigPath = ""
		return err
	}

	return nil
}

// 加载YAML配置文件
func LoadYAMLConfig() error {
	if LocalConfigPath == "" {
		LocalConfigPath = *flag.String("c", "./config.yaml", "配置文件路径")
		flag.Parse()
	}
	file, err := os.Open(filepath.Clean(LocalConfigPath))
	if err != nil {
		return err
	}
	err = yaml.NewDecoder(file).Decode(&LocalConfig)
	if err != nil {
		LocalConfigPath = ""
		return err
	}
	return nil
}

// 监视配置文件更新
func WatchConfig() error {
	if LocalConfigPath == "" {
		return errors.New("配置文件没有成功解析，无法启动监视")
	}
	// 创建监视器
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if e := watcher.Close(); e != nil {
			}
		}()

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
						if e := LoadYAMLConfig(); e != nil {
							panic(err.Error())
						}
					}
				case e := <-watcher.Errors:
					panic(e.Error())
				}
			}
		}()

		// 添加要监视的文件
		if err = watcher.Add(LocalConfigPath); err != nil {
			panic(err.Error())
		}
		<-done
	}()
	return nil
}
