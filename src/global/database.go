package global

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
)

var DB *pg.DB

// QueryHook 查询钩子
type QueryHook struct{}

// BeforeQuery 查询前钩子
func (QueryHook) BeforeQuery(ctx context.Context, qe *pg.QueryEvent) (context.Context, error) {
	// 自动重连数据库
	if DB == nil {
		if err := SetDatabase(); err != nil {

		}
	}
	return ctx, nil
}

// AfterQuery 查询后钩子
func (QueryHook) AfterQuery(ctx context.Context, qe *pg.QueryEvent) error {
	if !Config.Debug {
		return nil
	}
	// 记录SQL语句
	stmt, err := qe.FormattedQuery()
	if err != nil {
		log.Err(err).Caller().Send()
	}

	log.Debug().Msg(BytesToStr(stmt))

	return nil
}

// 设置数据库
func SetDatabase() (err error) {
	// 判断是否需要重启赋值
	if DB != nil {
		return
	}
	// 读取配置文件
	if _, err = net.ResolveTCPAddr("tcp", Config.Database.Addr); err != nil {
		log.Err(err).Caller().Str("addr", Config.Database.Addr).Msg("数据库配置失败")
		return
	}
	if Config.Database.User == "" {
		err = errors.New("user参数不正确")
		log.Err(err).Caller().Str("user", Config.Database.User).Msg("数据库配置失败")
		return
	}
	if Config.Database.Name == "" {
		err = errors.New("name参数不正确")
		log.Err(err).Caller().Str("name", Config.Database.Name).Msg("数据库配置失败")
		return
	}

	// 连接数据库
	DB = pg.Connect(&pg.Options{
		Addr:         Config.Database.Addr,
		User:         Config.Database.User,
		Password:     Config.Database.Password,
		Database:     Config.Database.Name,
		DialTimeout:  time.Duration(Config.Database.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(Config.Database.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(Config.Database.WriteTimeout) * time.Second,
		PoolSize:     int(Config.Database.PoolSize),
	})

	// 注册钩子
	DB.AddQueryHook(QueryHook{})

	return
}
