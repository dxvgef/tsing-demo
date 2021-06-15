package global

import (
	"time"

	"github.com/rs/zerolog/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var EtcdCli *clientv3.Client

func SetEtcdCli() (err error) {
	EtcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:            RuntimeConfig.Etcd.Endpoints,
		DialTimeout:          time.Duration(RuntimeConfig.Etcd.DialTimeout) * time.Second,
		Username:             RuntimeConfig.Etcd.Username,
		Password:             RuntimeConfig.Etcd.Password,
		AutoSyncInterval:     time.Duration(RuntimeConfig.Etcd.AutoSyncInterval) * time.Second,
		DialKeepAliveTime:    time.Duration(RuntimeConfig.Etcd.DialKeepAliveTime) * time.Second,
		DialKeepAliveTimeout: time.Duration(RuntimeConfig.Etcd.DialKeepAliveTimeout) * time.Second,
		MaxCallSendMsgSize:   int(RuntimeConfig.Etcd.MaxCallSendMsgSize),
		MaxCallRecvMsgSize:   int(RuntimeConfig.Etcd.MaxCallRecvMsgSize),
		RejectOldCluster:     RuntimeConfig.Etcd.RejectOldCluster,
		PermitWithoutStream:  RuntimeConfig.Etcd.PermitWithoutStream,
	})
	if err != nil {
		log.Err(err).Caller().Send()
	}
	return
}
