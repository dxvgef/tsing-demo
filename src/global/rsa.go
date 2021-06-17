package global

import (
	"github.com/dxvgef/rsalib"
	"github.com/rs/zerolog/log"
)

var PrivateKey rsalib.PrivateKey

func MakeRSAKey() (err error) {
	if err = PrivateKey.New(2048); err != nil {
		log.Err(err).Caller().Msg("生成RSA私钥失败")
	}
	return
}
