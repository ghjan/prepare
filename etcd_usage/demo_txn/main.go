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
		client             *clientv3.Client
		err                error
		lease              clientv3.Lease
		leaseId            clientv3.LeaseID
		leaseGrantResponse *clientv3.LeaseGrantResponse
		keepRespChan       <-chan *clientv3.LeaseKeepAliveResponse
		keepResp           *clientv3.LeaseKeepAliveResponse
		ctx2               context.Context
		cancelFunc         context.CancelFunc
		txn                clientv3.Txn
		txnResp            *clientv3.TxnResponse
	)
	job_key := "/cron/jobs/job5"

	//建立连接
	endpoints := []string{"192.168.1.133:2379",}
	if client, err = etcd_usage.GetEtcdClient(endpoints); err != nil {
		fmt.Println("main, etcd_usage.GetEtcdClient error:", err)
		return
	}
	//lease实现锁自动过期
	//op操作
	//txn事务if else then
	//1.上锁（创建租约，自动续租，拿着租约去抢占一个key)

	//准备一个用于自动取消自动续租的context
	ctx2, cancelFunc = context.WithCancel(context.TODO())

	countOfSecond := int64(10)
	if lease, keepRespChan, leaseGrantResponse, err = etcd_usage.GetNewLeaseKeepAlive(client, ctx2, countOfSecond+1);
		err != nil {
		fmt.Println(" etcd_usage.GetNewLease error:", err)
		return
	}
	leaseId = leaseGrantResponse.ID
	//确保函数退出之前，自动续租会停止
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	//处理续约应答的goroutine
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil || keepResp == nil {
					fmt.Println("租约已经失效了")
					fmt.Printf("keepRespChan:%v,keepResp:%v\n", keepRespChan, keepResp)
					goto END
				} else { //每秒会续租一次，所以就会收到一次应答
					fmt.Println("收到自动续租应答", keepResp.ID)
				}
			}
		}
	END:
	}()
	//if不存在key，then设置它，else抢锁失败

	//创建事务
	txn = client.Txn(context.TODO())
	//定义事务
	txn.If(clientv3.Compare(clientv3.CreateRevision(job_key), "=", 0)).Then(
		clientv3.OpPut(job_key, "", clientv3.WithLease(leaseId))).Else(
		clientv3.OpGet(job_key))
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(" txn.Commit error:", err)
		return

	}

	//判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Printf("锁被占用:%s\n", txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		return
	}

	//2.处理业务
	//在锁内，很安全
	fmt.Println("处理任务")
	time.Sleep(time.Duration(countOfSecond) * time.Second)

	//3.释放锁（取消自动续租，释放租约)
	//defer 会把租约释放

	fmt.Println("end of main")
}
