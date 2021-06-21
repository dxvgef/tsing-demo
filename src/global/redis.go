package global

import (
	"errors"
	"net"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

var RedisCli *redis.Client

func SetRedis() (err error) {
	_, err = net.ResolveTCPAddr("tcp", Config.Redis.Addr)
	if err != nil {
		log.Err(err).Caller().Str("addr", Config.Redis.Addr).Msg("Redis配置失败")
		return
	}

	if Config.Redis.DB > 16 {
		err = errors.New("库索引号不能大于16")
		log.Fatal().Err(err).Caller().Uint8("db", Config.Redis.DB).Msg("Redis配置失败")
		return
	}
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Addr,
		Username: Config.Redis.Username,
		Password: Config.Redis.Password,
		DB:       int(Config.Redis.DB),
	})

	return
}
