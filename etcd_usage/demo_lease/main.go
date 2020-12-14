package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"prepare/etcd_usage"
	"time"
)

func main() {
	var (
		client   *clientv3.Client
		err      error
		finished chan bool
	)

	//建立连接
	endpoints := []string{"192.168.1.133:2379",}
	if client, err = etcd_usage.GetEtcdClient(endpoints); err != nil {
		fmt.Println("main, etcd_usage.GetEtcdClient error:", err)
		return
	}

	fmt.Println("-----UsageLease1--------")
	UsageLease1(client, 10)
	fmt.Println("-----UsageLease2--------")
	finished = UsageLease2(client, 30*time.Second)
	<-finished
	fmt.Println("end of main")
}

func UsageLease1(client *clientv3.Client, ttl int64) {
	var (
		err error
		ctx context.Context
		//lease              clientv3.Lease
		leaseId            clientv3.LeaseID
		leaseGrantResponse *clientv3.LeaseGrantResponse
		presponse          *clientv3.PutResponse
		greponse           *clientv3.GetResponse
	)
	ctx = context.TODO()
	if leaseGrantResponse, err = etcd_usage.GetNewLease(client, ctx, ttl); err != nil {
		fmt.Println("UsageLease1, etcd_usage.GetNewLease error:", err)
		return
	}
	leaseId = leaseGrantResponse.ID
	if presponse, err = client.Put(ctx, "/cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println("UsageLease1, client.Put error:", err)
		return
	}

	fmt.Println("UsageLease1,写入成功", presponse.Header.Revision)

	//定时看一下key过期没有
	for {
		if greponse, err = client.Get(ctx, "/cron/lock/job1"); err != nil {
			fmt.Println("UsageLease1, client.Get error:", err)
			return
		}
		if greponse.Count == 0 {
			fmt.Println("UsageLease1,kv过期了")
			break

		} else {
			fmt.Println("UsageLease1,还没过期", greponse.Kvs)
			time.Sleep(1 * time.Second)
		}
	}
}

func UsageLease2(client *clientv3.Client, timeout time.Duration) (finished chan bool) {
	var (
		err                error
		ctx2               context.Context
		cancelFunc         context.CancelFunc
		lease2             clientv3.Lease
		leaseId2           clientv3.LeaseID
		leaseGrantResponse *clientv3.LeaseGrantResponse
		keepRespChan       <-chan *clientv3.LeaseKeepAliveResponse
		keepResp           *clientv3.LeaseKeepAliveResponse
	)

	lease2 = clientv3.NewLease(client)

	finished = make(chan bool, 0)
	ctx2, cancelFunc = context.WithTimeout(context.TODO(), timeout*2)
	defer cancelFunc()
	ttl := int64(10)
	if leaseGrantResponse, err = etcd_usage.GetNewLease(client, ctx2, ttl); err != nil {
		fmt.Println("UsageLease2, etcd_usage.GetNewLease error:", err)
		return
	}
	leaseId2 = leaseGrantResponse.ID

	//自动续租
	if keepRespChan, err = lease2.KeepAlive(ctx2, leaseId2); err != nil {
		fmt.Println("UsageLease2, lease2.KeepAlive error:", err)
		return
	}

	if err = putWithLease(client, ctx2, leaseId2, "/cron/lock/job2", "UsageLease2"); err != nil {
		time.Sleep(1 * time.Second)
		return
	}

	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil || keepResp == nil {
					fmt.Println("UsageLease2,租约已经失效了")
					fmt.Printf("UsageLease2,keepRespChan:%v,keepResp:%v\n", keepRespChan, keepResp)
					time.Sleep(1 * time.Second)
					finished <- true
					goto END
				} else { //每秒会续租一次，所以就会收到一次应答
					fmt.Println("UsageLease2,收到自动续租应答", keepResp.ID)
				}
				//case <-time.After(timeout):
				//	fmt.Printf("UsageLease2,after timeout:%v,to revoke lease2(%d)\n", timeout, leaseId2)
				//	time.Sleep(1 * time.Second)
				//	cancelFunc()
				//	//lease2.Revoke(ctx2, leaseId2)
			}
		}
	END:
	}()
	return
}

func putWithLease(client *clientv3.Client, ctx context.Context, leaseId clientv3.LeaseID, key string, val string) (err error) {
	var (
		presponse *clientv3.PutResponse
		greponse  *clientv3.GetResponse
	)
	if presponse, err = client.Put(ctx, key, val, clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("putWithLease, 写入成功", presponse.Header.Revision)
	//定时看一下key过期没有
	go func() {
		for {
			if greponse, err = client.Get(ctx, key); err != nil {
				fmt.Println("putWithLease, client.Get error:", err)
				return
			}
			if greponse.Count == 0 {
				fmt.Println("putWithLease,kv过期了")
				break

			} else {
				fmt.Println("putWithLease,还没过期", greponse.Kvs)
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return
}
