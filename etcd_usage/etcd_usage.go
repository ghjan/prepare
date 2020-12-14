package etcd_usage

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

var (
	config clientv3.Config
)

func GetEtcdClient(endpoints []string) (client *clientv3.Client, err error) {
	//客户端配置
	config = clientv3.Config{
		Endpoints:   endpoints, //集群列表
		DialTimeout: 5 * time.Second,
	}
	//建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
	}
	return
}

func GetNewLease(client *clientv3.Client, ctx context.Context, ttl int64) (leaseGrantResponse *clientv3.LeaseGrantResponse, err error) {
	if ctx == nil {
		ctx = context.TODO()
	}
	lease := clientv3.NewLease(client)
	if leaseGrantResponse, err = lease.Grant(ctx, ttl); err != nil {
		fmt.Println("UsageLease1, lease.Grant error:", err)
		return nil, nil
	}
	return leaseGrantResponse, err
}

func GetNewLeaseKeepAlive(client *clientv3.Client, ctx context.Context, ttl int64) (lease clientv3.Lease,
	keepRespChan <-chan *clientv3.LeaseKeepAliveResponse, leaseGrantResponse *clientv3.LeaseGrantResponse, err error) {
	if ctx == nil {
		ctx = context.TODO()
	}
	leaseGrantResponse, err = GetNewLease(client, ctx, ttl)
	lease = clientv3.NewLease(client)

	if leaseGrantResponse, err = GetNewLease(client, nil, ttl); err != nil {
		fmt.Println("UsageLease2, etcd_usage.GetNewLease error:", err)
		return
	}

	//自动续租
	if keepRespChan, err = lease.KeepAlive(ctx, leaseGrantResponse.ID); err != nil {
		fmt.Println("UsageLease2, lease2.KeepAlive error:", err)
		return
	}
	return
}
