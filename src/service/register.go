// 服务注册

package service

import (
	tsingCenter "github.com/dxvgef/tsing-center-go"
	"github.com/rs/zerolog/log"

	"local/global"
)

var tc *tsingCenter.Client

// 设置服务中心
func SetCenter() {
	var (
		err  error
		ip   string
		port uint16
	)

	// 新建tsing center客户端实例
	tc, err = tsingCenter.New(tsingCenter.Config{
		Addr:          global.Config.ServiceCenter.Addr,          // tsing center api地址
		Secret:        global.Config.ServiceCenter.Secret,        // tsing center api请求密钥
		TouchInterval: global.Config.ServiceCenter.TouchInterval, // 自动触活的间隔时间(秒)
		Timeout:       global.Config.ServiceCenter.Timeout,       // api请求超时时间(秒)
	})
	if err != nil {
		log.Fatal().Err(err).Caller().Send()
	}

	// 获取IP地址
	if ip, _, err = tc.GetIP(); err != nil {
		log.Fatal().Err(err).Caller().Send()
	}
	// 获得端口
	// nolint:gocritic
	if global.Config.Service.HTTPSPort > 0 {
		port = global.Config.Service.HTTPSPort
	} else if global.Config.Service.HTTPPort > 0 {
		port = global.Config.Service.HTTPPort
	} else {
		log.Fatal().Err(err).Caller().Msg("服务端口号无效")
	}

	// 注册服务
	if _, err = tc.SetService(tsingCenter.Service{
		ID:          global.Config.ServiceID,
		LoadBalance: "SWRR",
	}); err != nil {
		log.Fatal().Err(err).Caller().Send()
	}

	// 注册节点
	if _, err = tc.SetNode(global.Config.ServiceID, tsingCenter.Node{
		IP:     ip,
		Port:   port,
		TTL:    global.Config.ServiceCenter.TTL,
		Weight: global.Config.ServiceCenter.Weight,
	}); err != nil {
		log.Fatal().Err(err).Caller().Send()
	}

	log.Info().Str("ServiceID", global.Config.ServiceID).Uint("AutoTouchInterval", global.Config.ServiceCenter.TouchInterval).Str("IP", ip).Uint16("Port", port).Uint("TTL", global.Config.ServiceCenter.TTL).Uint("Weight", global.Config.ServiceCenter.Weight).Msg("服务注册成功")

	// 服务发现
	// var node tsingCenter.Node
	// node, _, err = tc.DiscoverService(global.setService.ServiceID)
	// if err != nil {
	// 	log.Fatal().Err(err).Caller().Msg("服务发现失败")
	// }
	// if node.IP == "" {
	// 	log.Debug().Caller().Msg("该服务没有有效的节点，应该访问501状态码")
	// } else {
	// 	log.Debug().Str("ip", node.IP).Uint16("port", node.Port).Caller().Send()
	// }

	tc.AutoTouchNode(global.Config.ServiceID, ip, port, func(status int, err error) {
		log.Err(err).Caller().Msg("自动触活失败")
	})
}
