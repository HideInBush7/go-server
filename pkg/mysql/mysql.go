package mysql

import (
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var dbs = sqlxManager{
	DBs: make(map[string]*sqlx.DB),
}

type sqlxManager struct {
	DBs map[string]*sqlx.DB
	mu  sync.RWMutex
}

type mysqlConfig struct {
	Alias        string `json:"alias"`
	Dsn          string `json:"dsn"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxIdleConns int    `json:"maxIdleConns"`
}

func Init() {
	var conf = []*mysqlConfig{}
	viper.UnmarshalKey(`mysql`, &conf)

	for _, v := range conf {
		register(v)
	}
}

func register(conf *mysqlConfig) {
	db, err := sqlx.Connect(`mysql`, conf.Dsn)
	if err != nil {
		logrus.WithField(`mysql_config`, conf).Error(err)
		return
	}

	if conf.MaxIdleConns <= 0 {
		conf.MaxIdleConns = 10
	}
	if conf.MaxOpenConns <= 0 {
		conf.MaxIdleConns = 20
	}
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxOpenConns)

	dbs.mu.Lock()
	defer dbs.mu.Unlock()
	dbs.DBs[conf.Alias] = db
}

// 通过别名 'alias' 获取当前*sqlx.DB连接 默认alias=default
func GetInstance(alias ...string) *sqlx.DB {
	var dbName string
	if len(alias) == 0 || alias == nil {
		dbName = `default`
	} else {
		dbName = alias[0]
	}

	dbs.mu.RLock()
	defer dbs.mu.RUnlock()
	db, ok := dbs.DBs[dbName]
	if !ok {
		logrus.WithField(`alias`, dbName).Error(`db alias does't register`)
		return nil
	}
	return db
}
