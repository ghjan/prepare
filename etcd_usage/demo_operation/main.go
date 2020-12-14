package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"prepare/etcd_usage"
)

func main() {
	var (
		client       *clientv3.Client
		err          error
		getOp, putOp clientv3.Op
		opResp       clientv3.OpResponse
	)
	job_key := "/cron/jobs/job4"

	//建立连接
	endpoints := []string{"192.168.1.133:2379",}
	if client, err = etcd_usage.GetEtcdClient(endpoints); err != nil {
		fmt.Println("main, etcd_usage.GetEtcdClient error:", err)
		return
	}

	//Op:operation
	putOp = clientv3.OpPut(job_key, "i am job4")

	//执行OP
	if opResp, err = client.Do(context.TODO(), putOp); err != nil {
		fmt.Println("= client.Do error:", err)
		return
	}
	//打印
	fmt.Println("写入Revision:", opResp.Put().Header.Revision)

	//Op:operation
	getOp = clientv3.OpGet(job_key)

	//执行OP
	if opResp, err = client.Do(context.TODO(), getOp); err != nil {
		fmt.Println("= client.Do error:", err)
		return
	}
	//打印
	fmt.Println("数据Revision:", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据Value:", string(opResp.Get().Kvs[0].Value))

	fmt.Println("end of main")
}
