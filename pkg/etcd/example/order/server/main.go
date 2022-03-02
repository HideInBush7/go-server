package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/HideInBush7/go-server/pkg/config"
	"github.com/HideInBush7/go-server/pkg/etcd"
	"github.com/HideInBush7/go-server/pkg/etcd/example/order/rpcserver"
	"github.com/spf13/viper"
)

func main() {
	config.Init()
	// 服务地址
	ServiceAddr := net.JoinHostPort(viper.GetString(`ip`), viper.GetString(`port`))

	client, err := etcd.DefaultClient()
	if err != nil {
		panic(err)
	}
	reg, err := etcd.NewRegistration(client)
	if err != nil {
		panic(err)
	}

	err = reg.Register(`/order.rpc`, ServiceAddr, 3)
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

	go rpcserver.NewServer()

	s := <-sig
	err = reg.Unregister()
	fmt.Println("服务注销: ", `err: `, err)
	fmt.Println("receive signal ", s.String())
}
