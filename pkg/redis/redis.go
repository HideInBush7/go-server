package redis

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

var redisPools = redisPoolManager{
	pool: make(map[string]*redis.Pool),
}

type redisPoolManager struct {
	pool map[string]*redis.Pool
	mu   sync.RWMutex
}

// 此处redis配置没有加上密码,有需要可以添加
type redisConfig struct {
	Alias       string `json:"alias"`
	Address     string `json:"address"`
	MaxIdle     int    `json:"maxIdleConns"`
	IdleTimeout int    `json:"idleTimeout"`
}

func Init() {
	redisConf := []*redisConfig{}
	viper.UnmarshalKey(`redis`, &redisConf)
	for _, v := range redisConf {
		register(v)
	}

}

func register(conf *redisConfig) {
	if conf.MaxIdle <= 0 {
		conf.MaxIdle = 3
	}
	if conf.IdleTimeout <= 0 {
		conf.IdleTimeout = 240
	}

	pool := &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(`tcp`, conf.Address)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do(`ping`)
			return err
		},
	}

	redisPools.mu.Lock()
	defer redisPools.mu.Unlock()
	redisPools.pool[conf.Alias] = pool
}

// 通过别名 'alias' 获取当前 redis.Conn 连接 默认alias=default
func GetInstance(alias ...string) redis.Conn {
	var name string
	if len(alias) == 0 || alias == nil {
		name = `default`
	}

	redisPools.mu.RLock()
	defer redisPools.mu.RUnlock()

	if pool, ok := redisPools.pool[name]; ok {
		return pool.Get()
	}
	return nil
}
