package action

import (
	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"
)

type Index struct{}

func (Index) Json(ctx tsing.Context) error {
	respData := makeRespMapData()
	respData.Data["username"] = "test"
	respData.Data["password"] = "123456"
	return JSON(ctx, 200, respData)
}

func (Index) String(ctx tsing.Context) error {
	var reqData struct {
		username string
		password string
	}
	// 过滤多个元素
	err := filter.MSet(
		// 要过滤的元素
		filter.El(
			&reqData.username, // 要接收过滤结果的变量
			// 数据来源于get参数username
			filter.FromString(ctx.RawQueryValue("username")).
				RemoveSpace().  // 移除所有空格
				MinLength(3).   // 要求最小长度
				MaxLength(16)), // 要求最大长度
		filter.El(&reqData.password,
			filter.FromString(ctx.RawQueryValue("password")).
				MinLength(6).MaxLength(32).HasDigit().HasUpper().HasLower().HasSymbol(),
		),
	)
	if err != nil {
		return ctx.Event(err)
	}
	return String(ctx, 200, err.Error())
}
