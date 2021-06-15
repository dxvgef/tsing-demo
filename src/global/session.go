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
		Addr:     RuntimeConfig.Session.RedisAddr,
		Username: RuntimeConfig.Session.RedisUsername,
		Password: RuntimeConfig.Session.RedisPassword,
		Prefix:   RuntimeConfig.Session.RedisKeyPrefix,
		DB:       RuntimeConfig.Session.RedisDB,
	}); err != nil {
		return
	}

	// 创建引擎
	if Sessions, err = sessions.New(&sessions.Config{
		Key:         RuntimeConfig.Session.CookieName,
		HTTPOnly:    RuntimeConfig.Session.HTTPOnly,
		Secure:      RuntimeConfig.Session.Secure,
		Path:        "/",
		IdleTimeout: RuntimeConfig.Session.IdleTimeout, // 空闲超时时间(秒)
	}, storage); err != nil {
		return
	}
	return
}
