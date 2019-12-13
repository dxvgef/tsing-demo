package global

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
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
	// 记录SQL语句
	stmt, err := qe.FormattedQuery()
	if err != nil {
		log.Error().Caller(1).Msg(err.Error())
	}

	skip := 5
	// nolint:gocritic
	if stmt == "BEGIN" {
		skip = 6
	} else if strings.HasPrefix(stmt, "UPDATE") {
		skip = 7
	} else if strings.HasPrefix(stmt, "INSERT") {
		skip = 6
	} else if strings.HasPrefix(stmt, "DELETE") {
		skip = 7
	}

	log.Debug().Caller(skip).Msg(stmt)

	return nil
}

// SetDatabase 设置数据库
func SetDatabase() error {
	// 判断是否需要重启赋值
	if DB != nil {
		return nil
	}
	// 读取配置文件
	if LocalConfig.Database.Addr == "" {
		return errors.New("配置文件[database]节点的addr参数不正确")
	}
	if LocalConfig.Database.User == "" {
		return errors.New("配置文件[database]节点的user参数不正确")
	}
	if LocalConfig.Database.Name == "" {
		return errors.New("配置文件[database]节点的name参数不正确")
	}

	// 连接数据库
	DB = pg.Connect(&pg.Options{
		Addr:         LocalConfig.Database.Addr,
		User:         LocalConfig.Database.User,
		Password:     LocalConfig.Database.Password,
		Database:     LocalConfig.Database.Name,
		DialTimeout:  time.Duration(LocalConfig.Database.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(LocalConfig.Database.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(LocalConfig.Database.WriteTimeout) * time.Second,
		PoolSize:     LocalConfig.Database.PoolSize,
	})

	// 注册查询钩子
	if LocalConfig.Database.StmtLog {
		DB.AddQueryHook(QueryHook{})
	}

	return nil
}
