package config

import (
	"github.com/HideInBush7/go-server/util"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Init() {
	pflag.String(`config`, `./config.yaml`, `config file`)
	pflag.String(`ip`, util.GetInternalIp(), `ip address`)
	pflag.String(`port`, `8000`, `port`)
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigFile(viper.GetString(`config`))
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
