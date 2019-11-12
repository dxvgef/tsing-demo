package global

import (
	"context"
	"log"
)

// 远程配置
var RemoteConfigFromETCD struct {
	SiteName   string
	SiteDomain string
}

// 从ETCD获得远程配置
func LoadRemoteConfigFromETCD() error {
	eResp, err := ETCDClient.Set(context.Background(), "test", "haha", nil)
	if err != nil {
		return err
	}
	log.Println(eResp.Action)
	log.Println(eResp.Node)
	log.Println(eResp.ClusterID)

	eResp, err = ETCDClient.Get(context.Background(), "test", nil)
	if err != nil {
		return err
	}
	log.Println(eResp.Action)
	log.Println(eResp.Node)
	log.Println(eResp.Node.Value)
	log.Println(eResp.ClusterID)
	return nil
}
