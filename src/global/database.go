package global

import (
	"context"
	"errors"
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
	if err := SetDatabase(); err != nil {
		log.Err(err).Caller().Send()
	}
	return ctx, nil
}

// AfterQuery 查询后钩子
func (QueryHook) AfterQuery(ctx context.Context, qe *pg.QueryEvent) error {
	if !RuntimeConfig.Service.Debug {
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

// SetDatabase 设置数据库
func SetDatabase() error {
	// 判断是否需要重启赋值
	if DB != nil {
		return nil
	}
	// 读取配置文件
	if RuntimeConfig.Database.Addr == "" {
		return errors.New("配置文件[database]节点的addr参数不正确")
	}
	if RuntimeConfig.Database.User == "" {
		return errors.New("配置文件[database]节点的user参数不正确")
	}
	if RuntimeConfig.Database.Name == "" {
		return errors.New("配置文件[database]节点的name参数不正确")
	}

	// 连接数据库
	DB = pg.Connect(&pg.Options{
		Addr:         RuntimeConfig.Database.Addr,
		User:         RuntimeConfig.Database.User,
		Password:     RuntimeConfig.Database.Password,
		Database:     RuntimeConfig.Database.Name,
		DialTimeout:  time.Duration(RuntimeConfig.Database.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(RuntimeConfig.Database.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(RuntimeConfig.Database.WriteTimeout) * time.Second,
		PoolSize:     int(RuntimeConfig.Database.PoolSize),
	})

	// 注册钩子
	DB.AddQueryHook(QueryHook{})

	return nil
}
