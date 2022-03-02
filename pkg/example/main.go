package main

import (
	"github.com/HideInBush7/go-server/pkg/config"
	"github.com/HideInBush7/go-server/pkg/log"
	"github.com/HideInBush7/go-server/pkg/mysql"
	"github.com/HideInBush7/go-server/pkg/redis"
	"github.com/sirupsen/logrus"
)

func main() {
	config.Init()
	log.Init()
	mysql.Init()
	redis.Init()

	db := mysql.GetInstance()
	var now interface{}
	err := db.Get(&now, `SELECT now()`)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info(now)
	}

	rds := redis.GetInstance()
	reply, err := rds.Do(`PING`)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Infof("%s", reply)
	}
}
