package global

import (
	"time"

	"go.etcd.io/etcd/client"
)

// etcd客户端
var ETCDClient client.KeysAPI

// 设置ETCD Client
func SetETCDClient() error {
	config, err := client.New(client.Config{
		Endpoints:               LocalConfig.ETCD.Endpoints,
		Transport:               client.DefaultTransport,
		Username:                LocalConfig.ETCD.Username,
		Password:                LocalConfig.ETCD.Password,
		HeaderTimeoutPerRequest: time.Duration(LocalConfig.ETCD.HeaderTimeoutPerRequest) * time.Second,
	})
	if err != nil {
		return err
	}
	ETCDClient = client.NewKeysAPI(config)
	return nil
}
