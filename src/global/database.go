package global

import (
	"errors"
	"time"

	"github.com/go-pg/pg"
)

var DB *pg.DB

// QueryHook 查询钩子
type QueryHook struct{}

// BeforeQuery 查询前钩子
func (QueryHook) BeforeQuery(qe *pg.QueryEvent) {
	// 连接数据库
	if err := SetDatabase(); err != nil {
		Logger.Caller.Error(err.Error())
	}
}

// AfterQuery 查询后钩子
func (QueryHook) AfterQuery(qe *pg.QueryEvent) {
	// 记录SQL语句
	stmt, err := qe.FormattedQuery()
	if err != nil {
		Logger.Caller.Warn(err.Error())
	}
	Logger.Default.Debug(stmt)
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
		ReadTimeout:  time.Duration(LocalConfig.Database.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(LocalConfig.Database.WriteTimeout) * time.Second,
		PoolSize:     LocalConfig.Database.PoolSize,
	})

	// 启用查询日志，将记录SQL语句
	if LocalConfig.Database.StmtLog {
		DB.AddQueryHook(QueryHook{})
	}

	return nil
}
