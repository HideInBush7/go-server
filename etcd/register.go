package etcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdRegistration struct {
	c       *clientv3.Client
	leaseID clientv3.LeaseID
}

func NewRegistration(cli *clientv3.Client) (*etcdRegistration, error) {
	return &etcdRegistration{
		c: cli,
	}, nil
}

// 服务注册
func (r *etcdRegistration) Register(serverName string, addr string, ttl int64) error {
	// 创建租约
	lease, err := r.c.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	r.leaseID = lease.ID

	// 服务携带租约注册
	_, err = r.c.Put(context.Background(), serverName+addr, addr, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}

	// 自动续租
	keepAlive, err := r.c.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		return err
	}

	go func() {
		for {
			if res := <-keepAlive; res == nil {
				fmt.Println(`keepAlive closed...`)
				return
			}
		}
	}()

	return nil
}

func (r *etcdRegistration) Unregister() error {
	_, err := r.c.Revoke(context.TODO(), r.leaseID)
	return err
}
