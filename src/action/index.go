package action

import (
	"github.com/dxvgef/tsing"
)

type Index struct{}

func (Index) Json(ctx tsing.Context) error {
	respData := makeRespMapData()
	respData.Data["username"] = "test"
	respData.Data["password"] = "123456"
	return JSON(ctx, 200, respData)
}
