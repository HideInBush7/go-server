package etcd

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var clients = make(map[string]*clientv3.Client)
var mu sync.RWMutex

func NewClient(cfg clientv3.Config, alias string) (*clientv3.Client, error) {
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	mu.Lock()
	defer mu.Unlock()
	clients[alias] = client
	return client, nil
}

func GetInstance(alias string) (*clientv3.Client, error) {
	if alias == `` {
		alias = `default`
	}

	var c *clientv3.Client
	var ok bool

	mu.RLock()
	defer mu.RUnlock()
	if c, ok = clients[alias]; !ok || c == nil {
		return nil, errors.Errorf(`etcd client alias '%s' does't exist`, alias)
	}

	return c, nil
}

func DefaultClient() (client *clientv3.Client, err error) {
	client, err = GetInstance(``)
	if err != nil {
		client, err = NewClient(clientv3.Config{
			Endpoints: viper.GetStringSlice(`etcd.endpoints`),
		}, `default`)
	}
	return
}
