package handler

import (
	"time"

	"github.com/dxvgef/tsing"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/rs/zerolog/log"
)

func CheckToken(ctx *tsing.Context) error {
	log.Debug().Caller().Msg("执行了CheckToken中间件")
	respData := makeRespMapData()
	tokenStr, exist := ctx.QueryParam("token")
	if !exist {
		ctx.Abort()
		respData.Error = "token参数不存在，需要先访问/login获得token，然后带上[GET]token参数发起请求"
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

	return nil
}
