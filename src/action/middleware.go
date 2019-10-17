package action

import (
	"time"

	"github.com/dxvgef/tsing"
	"github.com/gbrlsnchs/jwt/v3"

	"local/global"
)

func CheckJWT(ctx *tsing.Context) error {
	global.Logger.Caller.Debug("执行了CheckJWT中间件")
	respData := makeRespMapData()
	tokenStr, exist := ctx.Post("token")
	if !exist {
		respData.Error = "token字段不存在"
		return JSON(ctx, 401, respData)
	}

	var accessToken AccessToken

	alg := jwt.NewHS256([]byte("secret"))
	_, err := jwt.Verify([]byte(tokenStr), alg, &accessToken,
		jwt.ValidatePayload(&accessToken.Payload, jwt.ExpirationTimeValidator(time.Now())),
	)
	if err != nil {
		respData.Error = err.Error()
		err := JSON(ctx, 401, respData)
		return err
	}
	respData.Data["id"] = accessToken.Data.ID
	respData.Data["username"] = accessToken.Data.Username

	// 执行后面的处理器
	ctx.Next()

	return nil
}
