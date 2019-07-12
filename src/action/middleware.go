package action

import (
	"time"

	"github.com/dxvgef/tsing"
	"github.com/gbrlsnchs/jwt/v3"
)

func CheckJWT(ctx tsing.Context) (tsing.Context, error) {
	respData := makeRespMapData()
	tokenStr, err := ctx.FormValue("token")
	if err != nil {
		respData.Error = err.Error()
		err := JSON(ctx, 401, respData)
		return ctx.Break(err)
	}

	var accessToken AccessToken

	alg := jwt.NewHS256([]byte("secret"))
	_, err = jwt.Verify([]byte(tokenStr), alg, &accessToken,
		jwt.ValidatePayload(&accessToken.Payload, jwt.ExpirationTimeValidator(time.Now())),
	)
	if err != nil {
		respData.Error = err.Error()
		err := JSON(ctx, 401, respData)
		return ctx.Break(err)
	}
	respData.Data["id"] = accessToken.Data.ID
	respData.Data["username"] = accessToken.Data.Username
	return ctx.Continue()
}
