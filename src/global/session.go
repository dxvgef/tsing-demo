package global

import (
	"errors"
	"net"

	"github.com/dxvgef/sessions"
	"github.com/dxvgef/sessions/storage/redis"
	"github.com/rs/zerolog/log"
)

// Sessions 引擎
var Sessions *sessions.Engine

// SetSessions 设置session引擎
func SetSessions() (err error) {
	if _, err = net.ResolveTCPAddr("tcp", Config.Session.RedisAddr); err != nil {
		log.Err(err).Caller().Str("addr", Config.Session.RedisAddr).Msg("Session存储器配置失败")
		return
	}
	if Config.Session.RedisDB > 16 {
		err = errors.New("库索引号不能大于16")
		log.Err(err).Caller().Uint8("db", Config.Session.RedisDB).Msg("Session存储器配置失败")
		return
	}
	// 创建存储器
	var storage sessions.Storage
	if storage, err = redis.New(&redis.Config{
		Addr:     Config.Session.RedisAddr,
		Username: Config.Session.RedisUsername,
		Password: Config.Session.RedisPassword,
		Prefix:   Config.Session.RedisKeyPrefix,
		DB:       Config.Session.RedisDB,
	}); err != nil {
		log.Err(err).Caller().Msg("Session存储器配置失败")
		return
	}

	// 创建引擎
	if Sessions, err = sessions.New(&sessions.Config{
		Key:         Config.Session.CookieKey,
		HTTPOnly:    Config.Session.HTTPOnly,
		Secure:      Config.Session.Secure,
		Path:        "/",
		IdleTimeout: Config.Session.IdleTimeout, // 空闲超时时间(秒)
	}, storage); err != nil {
		log.Err(err).Caller().Msg("Session引擎配置失败")
		return
	}
	return
}
