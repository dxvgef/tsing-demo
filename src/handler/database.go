package handler

import (
	"errors"

	"local/global"

	"github.com/dxvgef/tsing"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/rs/xid"
)

type DatabaseHandler struct {
	tableName struct{} `pg:"tsing"`
	ID        int64    `pg:"id,pk"`
	XID       string   `pg:"xid,notnull,unique"`
}

// 创建数据表
func (*DatabaseHandler) CreateTable(ctx *tsing.Context) (err error) {
	if err = global.DB.Model(&DatabaseHandler{}).CreateTable(&orm.CreateTableOptions{}); err != nil {
		return ctx.Caller(err)
	}
	return ctx.Status(204)
}

// 写数据
func (model DatabaseHandler) Add(ctx *tsing.Context) (err error) {
	var (
		result pg.Result
	)
	model.ID = global.SnowflakeNode.Generate().Int64()
	model.XID = xid.New().String()
	result, err = global.DB.Model(&model).Insert()
	if err != nil {
		return ctx.Caller(err)
	}
	if result.RowsAffected() == 0 {
		return ctx.Caller(errors.New("写入数据失败"))
	}
	return ctx.Status(204)
}

// 读数据
func (*DatabaseHandler) Get(ctx *tsing.Context) (err error) {
	var (
		result   []DatabaseHandler
		total    int
		respData = make(map[string]interface{})
	)
	total, err = global.DB.Model(&DatabaseHandler{}).SelectAndCount(&result)
	if err != nil {
		return ctx.Caller(err)
	}
	respData["total"] = total
	respData["rows"] = respData
	return ctx.JSON(200, &respData)
}

// 删除数据
func (model DatabaseHandler) Delete(ctx *tsing.Context) (err error) {
	if err = global.DB.Model(&model).DropTable(&orm.DropTableOptions{}); err != nil {
		return ctx.Caller(err)
	}
	return ctx.Status(204)
}
