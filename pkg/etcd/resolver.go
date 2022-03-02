package etcd

import (
	"context"
	"fmt"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

type etcdBuilder struct {
	c *clientv3.Client
}

func NewBuilder(client *clientv3.Client) (resolver.Builder, error) {
	return &etcdBuilder{
		c: client,
	}, nil
}

// grpc resolver.Target Scheme、EndPoint已弃用,改为URL
// eg: "dns://some_authority/foo.bar"
// URL.Schema: dns
// URL.Host: some_authority		=> Authority
// URP.PATH: foo.bar			=> Endpoint
// 此处约定 URL.Host为空, URL.Scheme="etcd", URL.PATH=服务名称+本地[ip:port]
// 即客户端调用时Dial etcd:///server127.0.0.1:8001
// (ip地址不应该用环回地址,不然可能会产生冲突,此处仅为示例)

const scheme = "etcd"

// 构造一个Resolver
// 核心要完成cc.UpdateState更新真实的IP地址
func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &etcdResolver{
		c:      b.c,
		cc:     cc,
		ctx:    ctx,
		cancel: cancel,
		target: target.URL.Path,
	}

	// 服务发现
	go r.watch()

	return r, nil
}

func (b *etcdBuilder) Scheme() string {
	return scheme
}

type etcdResolver struct {
	c  *clientv3.Client
	cc resolver.ClientConn

	ctx    context.Context
	cancel context.CancelFunc
	target string
}

func (r *etcdResolver) watch() {
	resp, err := r.c.Get(r.ctx, r.target, clientv3.WithPrefix())
	if err != nil {
		return
	}

	// map形式,方便delete
	addrs := make(map[string]resolver.Address)
	for _, kv := range resp.Kvs {
		addrs[string(kv.Key)] = resolver.Address{
			Addr: string(kv.Value),
		}
	}

	// 初次更新
	r.cc.UpdateState(resolver.State{
		Addresses: convertToGRPCAddress(addrs),
	})

	// 开启watch
	wChan := r.c.Watch(r.ctx, r.target, clientv3.WithPrefix())
	for {
		select {
		case <-r.ctx.Done():
			return
		case wRes, ok := <-wChan:
			if !ok {
				fmt.Println(`resolver: watch closed`)
				return
			}
			if wRes.Err() != nil {
				fmt.Println(`resolver: watch failed `, err.Error())
				return
			}
			update := false
			for _, e := range wRes.Events {
				update = true
				switch e.Type {
				case mvccpb.PUT:
					addrs[string(e.Kv.Key)] = resolver.Address{
						Addr: string(e.Kv.Value),
					}
				case mvccpb.DELETE:
					delete(addrs, string(e.Kv.Key))
				}
			}

			if update {
				r.cc.UpdateState(resolver.State{
					Addresses: convertToGRPCAddress(addrs),
				})
			}
		}
	}
}

func (r *etcdResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *etcdResolver) Close() {
	r.cancel()
}

func convertToGRPCAddress(addrs map[string]resolver.Address) []resolver.Address {
	var result []resolver.Address
	for _, addr := range addrs {
		result = append(result, addr)
	}
	return result
}
