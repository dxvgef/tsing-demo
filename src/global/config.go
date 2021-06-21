package global

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog/log"
)

const LOCAL = "local"

// 启动参数
var LaunchFlag struct {
	ConfigSource string // 配置来源(local或者服务中心地址'127.0.0.1:10000')
	Env          string // 环境变量
}

// 运行时配置
var Config struct {
	Debug     bool   `json:"-" toml:"-"`
	ServiceID string `toml:"-" json:"-"` // 服务ID

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
		DB        uint8  `json:"db" toml:"db"`
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
	LaunchFlag.ConfigSource = LOCAL
	// LaunchFlag.Env = LOCAL
	// LaunchFlag.ServiceID = ""

	// 服务默认配置
	Config.Debug = true
	Config.Service.ReadTimeout = 10
	Config.Service.ReadHeaderTimeout = 10
	Config.Service.WriteTimeout = 10
	Config.Service.IdleTimeout = 10
	Config.Service.QuitWaitTimeout = 5
	Config.Service.HTTPPort = 80
	Config.Service.EventNotFound = true
	Config.Service.EventMethodNotAllowed = true

	// 日志默认配置
	Config.Logger.Level = "debug"
	Config.Logger.FileMode = 0600
	Config.Logger.Encode = "console"
	Config.Logger.TimeFormat = "y-m-d h:i:s"

	// etcd默认配置
	Config.Etcd.Endpoints = []string{"http://127.0.0.1:2379"}
	Config.Etcd.DialTimeout = 5

	// 数据库默认配置
	Config.Database.Addr = "127.0.0.1:5432"
	Config.Database.User = "postgres"
	Config.Database.DialTimeout = 5
	Config.Database.ReadTimeout = 10
	Config.Database.WriteTimeout = 10
	Config.Database.PoolSize = 200

	// session默认配置
	Config.Session.CookieKey = "sessionid"
	Config.Session.IdleTimeout = 40 * 60
	Config.Session.RedisKeyPrefix = "sess_"

	// redis默认配置
	Config.Redis.Addr = "127.0.0.1:6379"
}

// 加载配置
func LoadConfig() (err error) {
	// 加载本地默认配置
	defaultConfig()

	// 解析启动参数
	flag.StringVar(&LaunchFlag.ConfigSource, "cfg", LaunchFlag.ConfigSource, "配置来源，可以是'local'表示本地或者配置中心地址'ip:port'，默认'local'")
	flag.StringVar(&LaunchFlag.Env, "env", LaunchFlag.Env, "环境变量，默认为空")
	flag.Parse()

	LaunchFlag.Env = strings.ToLower(LaunchFlag.Env)

	log.Info().Str("服务ID", Config.ServiceID).Msg("内置变量")
	log.Info().Str("配置来源(cfg)", LaunchFlag.ConfigSource).Str("环境变量(env)", LaunchFlag.Env).Msg("启动参数")

	// 加载本地配置文件
	if LaunchFlag.ConfigSource == LOCAL {
		// 加载本地配置文件
		return loadConfigFile()
	}
	return
}

// 加载本地配置文件
func loadConfigFile() (err error) {
	var (
		filePath string
		file     *os.File
	)
	if LaunchFlag.Env == "" {
		filePath = "./config.toml"
	} else {
		filePath = filepath.Clean("./config." + LaunchFlag.Env + ".toml")
	}
	file, err = os.Open(filePath) // nolint:gosec
	if err != nil {
		log.Err(err).Caller().Str("path", filePath).Send()
		return
	}

	// 解析配置文件到Config
	err = toml.NewDecoder(file).Decode(&Config)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}
	log.Info().Str("路径", filePath).Msg("加载本地配置文件")
	return
}

// 设置服务ID
func SetServiceID(id string) (err error) {
	if id == "" {
		err = errors.New("服务ID不能为空")
		log.Err(err).Caller().Str("ID", id).Msg("设置服务ID失败")
		return
	}
	Config.ServiceID = id
	return
}
