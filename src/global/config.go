package global

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const LOCAL = "local"

// 启动配置
var LaunchConfig struct {
	ConfigSource string // 配置来源(local或者服务中心地址'127.0.0.1:10000')
	Env          string // 环境变量
	ServiceID    string // 服务ID
}

// 运行时配置
var RuntimeConfig struct {
	Debug bool `json:"debug" toml:"debug"`

	Service struct {
		Secret                string `json:"secret" toml:"secret"`
		IP                    string `json:"-" toml:"-"`
		ReadTimeout           uint   `json:"read_timeout" toml:"read_timeout"`
		ReadHeaderTimeout     uint   `json:"read_header_timeout" toml:"read_header_timeout"`
		WriteTimeout          uint   `json:"write_timeout" toml:"write_timeout"`
		IdleTimeout           uint   `json:"idle_timeout" toml:"idle_timeout"`
		QuitWaitTimeout       uint   `json:"quit_wait_timeout" toml:"quit_wait_timeout"`
		HTTPPort              uint16 `json:"http_port" toml:"http_port"`
		HTTPSPort             uint16 `json:"https_port" toml:"https_port"`
		HTTP2                 bool   `json:"http2" toml:"http2"`
		EventNotFound         bool   `json:"event_not_found" toml:"event_not_found"`
		EventMethodNotAllowed bool   `json:"event_method_not_allowed" toml:"event_method_not_allowed"`
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
		Username  string `json:"username" toml:"username"`
		Password  string `json:"password" toml:"password"`
		KeyPrefix string `json:"key_prefix" toml:"key_prefix"`
	} `json:"redis" toml:"redis"`

	Database struct {
		Addr         string `json:"addr" toml:"addr"`
		User         string `json:"user" toml:"user"`
		Password     string `json:"password" toml:"password"`
		Name         string `json:"name" toml:"name"`
		DialTimeout  uint   `json:"dial_timeout" toml:"dial_timeout"`
		ReadTimeout  uint   `json:"read_timeout" toml:"read_timeout"`
		WriteTimeout uint   `json:"write_timeout" toml:"write_timeout"`
		PoolSize     uint   `json:"pool_size" toml:"pool_size"`
	} `json:"database" toml:"database"`

	Session struct {
		CookieKey      string `json:"cookie_key" toml:"cookie_key"`
		RedisAddr      string `json:"redis_addr" toml:"redis_addr"`
		RedisUsername  string `json:"redis_username" toml:"redis_username"`
		RedisPassword  string `json:"redis_password" toml:"redis_password"`
		RedisKeyPrefix string `json:"redis_key_prefix" toml:"redis_key_prefix"`
		IdleTimeout    uint   `json:"idle_timeout" toml:"idle_timeout"`
		HTTPOnly       bool   `json:"http_only" toml:"http_only"`
		Secure         bool   `json:"secure" toml:"secure"`
		RedisDB        uint8  `json:"redis_db" toml:"redis_db"`
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

// 设置本地默认配置
func defaultConfig() {
	// 环境变量
	LaunchConfig.ConfigSource = LOCAL
	// LaunchConfig.Env = LOCAL
	LaunchConfig.ServiceID = "tsing-demo"

	// 服务默认配置
	RuntimeConfig.Debug = true
	RuntimeConfig.Service.ReadTimeout = 10
	RuntimeConfig.Service.ReadHeaderTimeout = 10
	RuntimeConfig.Service.WriteTimeout = 10
	RuntimeConfig.Service.IdleTimeout = 10
	RuntimeConfig.Service.QuitWaitTimeout = 5
	RuntimeConfig.Service.HTTPPort = 80
	RuntimeConfig.Service.EventNotFound = true
	RuntimeConfig.Service.EventMethodNotAllowed = true

	// 日志默认配置
	RuntimeConfig.Logger.Level = "debug"
	RuntimeConfig.Logger.FileMode = 0600
	RuntimeConfig.Logger.Encode = "console"
	RuntimeConfig.Logger.TimeFormat = "y-m-d h:i:s"

	// etcd默认配置
	RuntimeConfig.Etcd.Endpoints = []string{"http://127.0.0.1:2379"}
	RuntimeConfig.Etcd.DialTimeout = 5

	// 数据库默认配置
	RuntimeConfig.Database.Addr = "127.0.0.1:5432"
	RuntimeConfig.Database.User = "postgres"
	RuntimeConfig.Database.DialTimeout = 5
	RuntimeConfig.Database.ReadTimeout = 10
	RuntimeConfig.Database.WriteTimeout = 10
	RuntimeConfig.Database.PoolSize = 200

	// session默认配置
	RuntimeConfig.Session.CookieKey = "sessionid"
	RuntimeConfig.Session.IdleTimeout = 40 * 60

	// redis默认配置
	RuntimeConfig.Redis.Addr = "127.0.0.1:6379"
	RuntimeConfig.Redis.KeyPrefix = "sess_"
}

// 加载配置
func LoadConfig() (err error) {
	// 加载本地默认配置
	defaultConfig()

	// 解析启动参数
	flag.StringVar(&LaunchConfig.ConfigSource, "cfg", LaunchConfig.ConfigSource, "配置来源，可以是'local'表示本地或者配置中心地址'ip:port'，默认'local'")
	flag.StringVar(&LaunchConfig.Env, "env", LaunchConfig.Env, "环境变量，默认为空")
	flag.Parse()

	LaunchConfig.Env = strings.ToLower(LaunchConfig.Env)

	log.Info().Str("服务ID", LaunchConfig.ServiceID).Msg("内置变量")
	log.Info().Str("配置来源(cfg)", LaunchConfig.ConfigSource).Str("环境变量(env)", LaunchConfig.Env).Msg("启动参数")

	// 加载本地配置文件
	if LaunchConfig.ConfigSource == LOCAL {
		// 加载本地配置文件
		if err = loadConfigFile(); err != nil {
			return err
		}
		return nil
	}

	// 加载远程配置
	if err = loadRemoteConfig(); err != nil {
		log.Err(err).Caller().Send()
		return err
	}
	return nil
}

// 加载本地配置文件
func loadConfigFile() error {
	var filePath string
	if LaunchConfig.Env == "" {
		filePath = "./config.toml"
	} else {
		filePath = filepath.Clean("./config." + LaunchConfig.Env + ".toml")
	}
	file, err := os.Open(filePath) // nolint:gosec
	if err != nil {
		log.Err(err).Caller().Str("path", filePath).Msg("无法读取本地配置文件")
		return err
	}

	// 解析配置文件到Config
	err = toml.NewDecoder(file).Decode(&RuntimeConfig)
	if err != nil {
		log.Err(err).Caller().Msg("解析本地配置文件失败")
		return err
	}
	log.Info().Str("路径", filePath).Msg("加载本地配置文件")
	return nil
}

// 加载远程配置
func loadRemoteConfig() (err error) {
	var (
		resp *clientv3.GetResponse
		key  strings.Builder
	)
	// 实例化etcd客户端
	if EtcdCli == nil {
		if err = SetEtcdCli(); err != nil {
			log.Err(err).Caller().Msg("设置etcd客户端失败")
			return err
		}
	}
	key.WriteString("/")
	key.WriteString(LaunchConfig.ServiceID)
	key.WriteString("/")
	key.WriteString(LaunchConfig.Env)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(RuntimeConfig.Etcd.DialTimeout)*time.Second)
	defer cancel()
	// 取出前缀下的所有的key
	if resp, err = EtcdCli.Get(ctx, key.String(), clientv3.WithPrefix()); err != nil {
		log.Err(err).Caller().Send()
		return
	}
	// 遍历key来加载配置
	for k := range resp.Kvs {
		remoteKey := string(resp.Kvs[k].Key)
		// 加载公用配置
		if remoteKey == key.String() {
			if err = json.Unmarshal(resp.Kvs[k].Value, &RuntimeConfig); err != nil {
				log.Err(err).Caller().Send()
				return err
			}
			log.Info().Str("Key", remoteKey).Msg("加载远程公用配置成功")
		}
		// 加载定制配置
		if strings.HasPrefix(remoteKey, key.String()+"/"+LaunchConfig.ServiceID) {
			if err = json.Unmarshal(resp.Kvs[k].Value, &RuntimeConfig); err != nil {
				log.Err(err).Caller().Send()
				return err
			}
			log.Info().Str("Key", remoteKey).Msg("加载远程定制配置成功")
		}
	}
	log.Info().Strs("etcd", RuntimeConfig.Etcd.Endpoints).Msg("远程配置加载")
	return nil
}
