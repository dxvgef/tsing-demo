package global

import (
	"time"

	"github.com/rs/zerolog/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var EtcdCli *clientv3.Client

func SetEtcdCli() (err error) {
	EtcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:            Config.Etcd.Endpoints,
		DialTimeout:          time.Duration(Config.Etcd.DialTimeout) * time.Second,
		Username:             Config.Etcd.Username,
		Password:             Config.Etcd.Password,
		AutoSyncInterval:     time.Duration(Config.Etcd.AutoSyncInterval) * time.Second,
		DialKeepAliveTime:    time.Duration(Config.Etcd.DialKeepAliveTime) * time.Second,
		DialKeepAliveTimeout: time.Duration(Config.Etcd.DialKeepAliveTimeout) * time.Second,
		MaxCallSendMsgSize:   int(Config.Etcd.MaxCallSendMsgSize),
		MaxCallRecvMsgSize:   int(Config.Etcd.MaxCallRecvMsgSize),
		RejectOldCluster:     Config.Etcd.RejectOldCluster,
		PermitWithoutStream:  Config.Etcd.PermitWithoutStream,
	})
	if err != nil {
		log.Err(err).Caller().Send()
	}
	return
}
