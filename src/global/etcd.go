package global

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"
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
		return
	}
	return
}
