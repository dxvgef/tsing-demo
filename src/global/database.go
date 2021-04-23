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
	// 连接数据库
	if err := SetDatabase(); err != nil {
		log.Error().Msg(err.Error())
	}
	return ctx, nil
}

// AfterQuery 查询后钩子
func (QueryHook) AfterQuery(ctx context.Context, qe *pg.QueryEvent) error {
	if !Config.Database.StmtLog {
		return nil
	}
	// 记录SQL语句
	stmt, err := qe.FormattedQuery()
	if err != nil {
		log.Error().Msg(err.Error())
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
	if Config.Database.Addr == "" {
		return errors.New("配置文件[database]节点的addr参数不正确")
	}
	if Config.Database.User == "" {
		return errors.New("配置文件[database]节点的user参数不正确")
	}
	if Config.Database.Name == "" {
		return errors.New("配置文件[database]节点的name参数不正确")
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

	// 注册查询钩子
	if Config.Database.StmtLog {
		DB.AddQueryHook(QueryHook{})
	}

	return nil
}
