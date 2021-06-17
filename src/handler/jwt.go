package handler

import (
	"crypto/rand"
	"crypto/rsa"
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

var privateKey *rsa.PrivateKey

// 验证JWT
func (*JWTHandler) Auth(ctx *tsing.Context) error {
	tokenStr := ctx.Query("token")
	if tokenStr == "" {
		return ctx.Status(400)
	}
	if privateKey == nil {
		return ctx.String(500, "私钥没有创建")
	}
	claims, err := jwt.RSACheck([]byte(tokenStr), &privateKey.PublicKey)
	if err != nil {
		return ctx.Caller(err)
	}
	return ctx.JSON(200, &claims)
}

// 签发JWT
func (*JWTHandler) Sign(ctx *tsing.Context) (err error) {
	var (
		tokenBytes []byte
		claims     JWTHandler
	)
	// 生成私钥文件
	if privateKey == nil {
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return ctx.Caller(err)
		}
	}

	claims.UserID = global.SnowflakeNode.Generate().Int64()
	claims.Nickname = "dxvgef"
	claims.Expires = jwt.NewNumericTime(time.Now().Add(10 * time.Minute))

	tokenBytes, err = claims.RSASign(jwt.RS256, privateKey)
	if err != nil {
		return ctx.Caller(err)
	}
	return ctx.String(200, string(tokenBytes))
}
