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
	stmt, _ := qe.FormattedQuery()
	Logger.Default.Debug(stmt)
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
		ReadTimeout:  time.Duration(Config.Database.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(Config.Database.WriteTimeout) * time.Second,
		PoolSize:     Config.Database.PoolSize,
	})

	// 启用查询日志，将记录SQL语句
	if Config.Database.StmtLog == true {
		DB.AddQueryHook(QueryHook{})
	}

	return nil
}
