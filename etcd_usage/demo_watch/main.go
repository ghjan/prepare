package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"prepare/etcd_usage"
	"time"
)

func main() {
	var (
		client             *clientv3.Client
		err                error
		getResp            *clientv3.GetResponse
		watchStartRevision int64
		watcher            clientv3.Watcher
		watchRespChan      <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
	)
	job_key := "/cron/jobs/job3"

	//建立连接
	endpoints := []string{"192.168.1.133:2379",}
	if client, err = etcd_usage.GetEtcdClient(endpoints); err != nil {
		fmt.Println("main, etcd_usage.GetEtcdClient error:", err)
		return
	}
	ctx, cancelFun := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("---取消")
		cancelFun()
	})
	go func() {
		for {
			client.Put(ctx, job_key, "i am job3")

			client.Delete(ctx, job_key)
			time.Sleep(1 * time.Second)
		}

	}()

	//先GET到当前的值，并监听后续变化
	if getResp, err = client.Get(context.TODO(), job_key); err != nil {
		fmt.Println("main, client.Get error:", err)
	}

	//如果现在的值是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值", string(getResp.Kvs[0].Value))
	}

	//当前etcd集群事务ID，单调递增的
	watchStartRevision = getResp.Header.Revision + 1

	//创建一个watcher
	watcher = clientv3.NewWatcher(client)

	//启动监听
	fmt.Println("从该版本开始向后监听:", watchStartRevision)

	watchRespChan = watcher.Watch(ctx, job_key, clientv3.WithRev(watchStartRevision))

	//处理kv变化事件
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision,
					event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision,
					event.Kv.ModRevision)
			}
		}
	}

	fmt.Println("end of main")
}
