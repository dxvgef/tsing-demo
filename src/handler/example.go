package handler

import (
	"time"

	"github.com/dxvgef/filter/v2"
	"github.com/dxvgef/tsing"
	"github.com/gbrlsnchs/jwt/v3"
)

type AccessToken struct {
	Data struct {
		ID       int64  `json:"id,omitempty"`
		Username string `json:"username,omitempty"`
	} `json:"data,omitempty"`
	jwt.Payload
}

type Example struct{}

// 登录，签发token
func (*Example) Login(ctx *tsing.Context) error {
	var respData RespData
	var accessToken AccessToken
	accessToken.Data.ID = 123
	accessToken.Data.Username = "dxvgef"
	accessToken.ExpirationTime = jwt.NumericDate(time.Now().Add(5 * time.Minute))

	alg := jwt.NewHS256([]byte("secret"))
	token, err := jwt.Sign(accessToken, alg)
	if err != nil {
		return err
	}
	respData.Data = string(token)

	return JSON(ctx, 200, respData)
}

/*
// 读写session的处理器
func (*Example) Session(ctx *tsing.Context) error {
	var respData RespData

	session, err := global.Session.Use(ctx.Request, ctx.ResponseWriter)
	if err != nil {
		return err
	}
	if err = session.Set("test", "tsing"); err != nil {
		return err
	}
	session.Get("test")
	sessValue, err := session.Get("test").String()
	if err != nil {
		return err
	}
	respData.Data = sessValue

	return JSON(ctx, 200, respData)
}
*/

func (*Example) Index(ctx *tsing.Context) error {
	var reqData struct {
		username string
		password string
	}
	err := filter.Batch(
		filter.String(ctx.Post("username"), "账号").
			Require().RemoveSpace().MinLength(3).MaxLength(16).
			Set(&reqData.username),
		filter.String(ctx.Post("password"), "密码").
			Require().MinLength(6).MaxLength(32).HasDigit().HasUpper().HasLower().HasSymbol().
			Set(&reqData.username),
	)

	if err != nil {
		return err
	}
	return String(ctx, 200, "身份验证通过，欢迎使用Tsing")
}
