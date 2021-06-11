package global

import (
	"github.com/dxvgef/sessions"
	"github.com/dxvgef/sessions/storage/redis"
)

// Sessions 引擎
var Sessions *sessions.Engine

// SetSessions 设置session引擎
func SetSessions() (err error) {
	// 创建存储器
	var storage sessions.Storage
	if storage, err = redis.New(&redis.Config{
		Addr:     Config.Session.RedisAddr,
		Username: Config.Session.RedisUsername,
		Password: Config.Session.RedisPassword,
		Prefix:   Config.Session.RedisKeyPrefix,
		DB:       Config.Session.RedisDB,
	}); err != nil {
		return
	}

	// 创建引擎
	if Sessions, err = sessions.New(&sessions.Config{
		Key:         Config.Session.CookieName,
		HTTPOnly:    Config.Session.HTTPOnly,
		Secure:      Config.Session.Secure,
		Path:        "/",
		IdleTimeout: Config.Session.IdleTimeout, // 空闲超时时间(秒)
	}, storage); err != nil {
		return
	}
	return
}
