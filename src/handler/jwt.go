package handler

import (
	"time"

	"local/global"

	"github.com/dxvgef/tsing"
	"github.com/pascaldekloe/jwt"
)

type JWTHandler struct {
	UserID   int64  `json:"userID"`
	Nickname string `json:"nickname"`
	jwt.Claims
}

// 验证JWT
func (*JWTHandler) Auth(ctx *tsing.Context) error {
	tokenStr := ctx.Query("token")
	if tokenStr == "" {
		return ctx.Status(400)
	}
	claims, err := jwt.RSACheck([]byte(tokenStr), global.PrivateKey.GetPublicKey().ToRaw())
	if err != nil {
		return ctx.Caller(err)
	}
	if !claims.Valid(time.Now()) {
		return ctx.Status(401)
	}
	return ctx.JSON(200, &claims)
}

// 签发JWT
func (*JWTHandler) Sign(ctx *tsing.Context) (err error) {
	var (
		tokenBytes []byte
		claims     JWTHandler
	)

	claims.UserID = global.SnowflakeNode.Generate().Int64()
	claims.Nickname = "dxvgef"
	claims.Expires = jwt.NewNumericTime(time.Now().Add(10 * time.Minute))

	tokenBytes, err = claims.RSASign(jwt.RS256, global.PrivateKey.ToRaw())
	if err != nil {
		return ctx.Caller(err)
	}
	return ctx.String(200, string(tokenBytes))
}
