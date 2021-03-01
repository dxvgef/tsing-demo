package global

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog/log"
)

var Config struct {
	// 启动参数
	Env string `toml:"-" json:"-"`

	// 本地参数
	ProjectID string `json:"-" toml:"project_id"`
	Version   string `json:"-" toml:"version"`

	Common struct {
	} `json:"common" toml:"common"`

	Service struct {
		ID                    string `json:"-" toml:"id"` // 服务ID
		Secret                string `json:"secret" toml:"secret"`
		IP                    string `json:"-" toml:"-"`
		HTTPPort              uint16 `json:"http_port" toml:"http_port"`
		HTTPSPort             uint16 `json:"https_port" toml:"https_port"`
		ReadTimeout           uint   `json:"read_timeout" toml:"read_timeout"`
		ReadHeaderTimeout     uint   `json:"read_header_timeout" toml:"read_header_timeout"`
		WriteTimeout          uint   `json:"write_timeout" toml:"write_timeout"`
		IdleTimeout           uint   `json:"idle_timeout" toml:"idle_timeout"`
		QuitWaitTimeout       uint   `json:"quit_wait_timeout" toml:"quit_wait_timeout"`
		HTTP2                 bool   `json:"http2" toml:"http2"`
		Debug                 bool   `json:"-" toml:"debug"`
		EventSource           bool   `json:"event_source" toml:"event_source"`
		EventTrace            bool   `json:"event_trace" toml:"event_trace"`
		EventNotFound         bool   `json:"event_not_found" toml:"event_not_found"`
		EventMethodNotAllowed bool   `json:"event_method_not_allowed" toml:"event_method_not_allowed"`
		Recover               bool   `json:"-" toml:"-"`
		EventShortPath        bool   `json:"event_short_path" toml:"event_short_path"`
	} `json:"service" toml:"service"`

	Logger struct {
		Level      string      `json:"level" toml:"level"`
		FilePath   string      `json:"file_path" toml:"file_path"`
		Encode     string      `json:"encode" toml:"encode"`
		TimeFormat string      `json:"time_format" toml:"time_format"`
		FileMode   os.FileMode `json:"file_mode" toml:"file_mode"`
		NoColor    bool        `json:"no_color" toml:"no_color"`
	} `json:"logger" toml:"logger"`

	Etcd struct {
		Endpoints            []string `toml:"endpoints"`
		DialTimeout          uint     `toml:"dial_timeout"`
		Username             string   `toml:"username"`
		Password             string   `toml:"password"`
		AutoSyncInterval     uint     `toml:"auto_sync_interval"`
		DialKeepAliveTime    uint     `toml:"dial_keep_alive_time"`
		DialKeepAliveTimeout uint     `toml:"dial_keep_alive_timeout"`
		MaxCallSendMsgSize   uint     `toml:"max_call_send_msg_size"`
		MaxCallRecvMsgSize   uint     `toml:"max_call_recv_msg_size"`
		RejectOldCluster     bool     `toml:"reject_old_cluster"`
		PermitWithoutStream  bool     `toml:"permit_without_stream"`
	} `toml:"etcd"`

	Redis struct {
		Addr      string `json:"addr" toml:"addr"`
		DB        int    `json:"db" toml:"db"`
		Password  string `json:"password" toml:"password"`
		KeyPrefix string `json:"key_prefix" toml:"key_prefix"`
	} `json:"redis" toml:"redis"`

	Database struct {
		Addr         string `json:"addr" toml:"addr"`
		User         string `json:"user" toml:"user"`
		Password     string `json:"password" toml:"password"`
		Name         string `json:"name" toml:"name"`
		StmtLog      bool   `json:"stmt_log" toml:"stmt_log"`
		DialTimeout  uint   `json:"dial_timeout" toml:"dial_timeout"`
		ReadTimeout  uint   `json:"read_timeout" toml:"read_timeout"`
		WriteTimeout uint   `json:"write_timeout" toml:"write_timeout"`
		PoolSize     uint   `json:"pool_size" toml:"pool_size"`
	} `json:"database" toml:"database"`

	Session struct {
		AESKey         string `json:"aes_key" toml:"aes_key"`
		CookieName     string `json:"cookie_name" toml:"cookie_name"`
		MaxAge         uint   `json:"max_age" toml:"max_age"`
		IdleTime       uint   `json:"idle_time" toml:"idle_time"`
		HTTPOnly       bool   `json:"http_only" toml:"http_only"`
		Secure         bool   `json:"secure" toml:"secure"`
		RedisAddr      string `json:"redis_addr" toml:"redis_addr"`
		RedisPassword  string `json:"redis_password" toml:"redis_password"`
		RedisDB        int    `json:"redis_db" toml:"redis_db"`
		RedisKeyPrefix string `json:"redis_key_prefix" toml:"redis_key_prefix"`
	} `json:"session" toml:"session"`

	ServiceCenter struct {
		Addr          string `json:"addr" toml:"addr"`                     // 服务中心api地址
		Secret        string `json:"secret" toml:"secret"`                 // 服务中心api请求密钥
		Timeout       uint   `json:"timeout" toml:"timeout"`               // 服务中心api请求超时时间(秒)
		TTL           uint   `json:"ttl" toml:"ttl"`                       // 节点的生命周期(秒)
		TouchInterval uint   `json:"touch_interval" toml:"touch_interval"` // 节点自动触活的间隔时间(秒)
		Weight        uint   `json:"weight" toml:"weight"`                 // 节点权重值
	} `json:"service_center" toml:"service_center"`
}

// 加载本地默认配置
func loadDefault() {
	// 公用默认配置
	Config.Env = "local" // 环境变量(默认)

	// 服务默认配置
	Config.Service.ReadTimeout = 10
	Config.Service.ReadHeaderTimeout = 10
	Config.Service.WriteTimeout = 10
	Config.Service.IdleTimeout = 10
	Config.Service.QuitWaitTimeout = 5
	Config.Service.Debug = true
	Config.Service.HTTPPort = 80
	Config.Service.Recover = true
	Config.Service.EventSource = true
	Config.Service.EventTrace = true
	Config.Service.EventNotFound = true
	Config.Service.EventMethodNotAllowed = true
	Config.Service.EventShortPath = true

	// 日志默认配置
	Config.Logger.Level = "debug"
	Config.Logger.FileMode = 600
	Config.Logger.Encode = "console"
	Config.Logger.TimeFormat = "y-m-d h:i:s"

	// etcd默认配置
	Config.Etcd.Endpoints = []string{"http://127.0.0.1:2379"}
	Config.Etcd.DialTimeout = 5

	// 数据库默认配置
	Config.Database.Addr = "127.0.0.1:5432"
	Config.Database.User = "postgres"
	Config.Database.StmtLog = true
	Config.Database.DialTimeout = 5
	Config.Database.ReadTimeout = 10
	Config.Database.WriteTimeout = 10
	Config.Database.PoolSize = 200

	// session默认配置
	Config.Session.CookieName = "sessionid"
	Config.Session.MaxAge = 60
	Config.Session.IdleTime = 40

	// redis默认配置
	Config.Redis.Addr = "127.0.0.1:6379"
	Config.Redis.KeyPrefix = "sess_"
}

// 加载配置
func LoadConfig() error {
	// 加载本地默认配置
	loadDefault()

	// 解析启动参数
	flag.StringVar(&Config.Env, "env", Config.Env, "环境变量，默认值:"+Config.Env)
	flag.Parse()

	Config.Env = strings.ToLower(Config.Env)

	// 加载本地配置文件
	err := loadFile()
	if err != nil {
		log.Fatal().Err(err).Caller().Send()
		return err
	}
	log.Info().Msg("本地配置加载成功")

	// 如果不是本地模式启动
	if Config.Env == "local" {
		return nil
	}
	// 实例化etcd客户端
	if EtcdCli == nil {
		if err = SetEtcdCli(); err != nil {
			log.Fatal().Err(err).Caller().Send()
			return err
		}
	}
	// 加载远程配置
	if err = loadRemoteConfig(); err != nil {
		log.Fatal().Err(err).Caller().Send()
		return err
	}
	log.Info().Msg("远程配置加载成功")
	return nil
}

// 加载本地yaml配置文件
func loadFile() error {
	file, err := os.Open(filepath.Clean("./config." + Config.Env + ".toml"))
	if err != nil {
		log.Fatal().Err(err).Caller().Send()
		return err
	}

	// 解析yaml到Config
	err = toml.NewDecoder(file).Decode(&Config)
	if err != nil {
		log.Fatal().Err(err).Caller().Send()
		return err
	}
	return nil
}

// 加载远程配置
func loadRemoteConfig() (err error) {
	var (
		resp *clientv3.GetResponse
		key  strings.Builder
	)
	key.WriteString("/" + Config.ProjectID)
	key.WriteString("/" + Config.Version)
	key.WriteString("/" + Config.Env)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Config.Etcd.DialTimeout)*time.Second)
	defer cancel()
	// 取出前缀下的所有的key
	if resp, err = EtcdCli.Get(ctx, key.String(), clientv3.WithPrefix()); err != nil {
		log.Fatal().Err(err).Caller().Send()
		return
	}
	// 遍历key来加载配置
	for k := range resp.Kvs {
		remoteKey := string(resp.Kvs[k].Key)
		// 加载公用配置
		if remoteKey == key.String() {
			if err = json.Unmarshal(resp.Kvs[k].Value, &Config); err != nil {
				log.Fatal().Err(err).Caller().Send()
				return err
			}
			log.Info().Str("Key", remoteKey).Msg("加载远程公用配置成功")
		}
		// 加载定制配置
		if strings.HasPrefix(remoteKey, key.String()+"/"+Config.Service.ID) {
			if err = json.Unmarshal(resp.Kvs[k].Value, &Config); err != nil {
				log.Fatal().Err(err).Caller().Send()
				return err
			}
			log.Info().Str("Key", remoteKey).Msg("加载远程定制配置成功")
		}
	}
	return nil
}
