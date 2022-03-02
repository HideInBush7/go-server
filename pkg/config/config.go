package config

import (
	"path/filepath"

	"github.com/HideInBush7/go-server/pkg/util"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// 配置文件和输入参数初始化 默认读取可执行文件同目录下的config.yaml
func Init() {
	pflag.String(`config`, filepath.Join(util.GetExecDir(), `config.yaml`), `config file`)
	pflag.String(`ip`, util.GetInternalIp(), `ip address`)
	pflag.String(`port`, `8000`, `port`)

	pflag.String(`logpath`, util.GetExecDir(), `logpath`)
	pflag.Bool(`debug`, true, `open console log`)
	pflag.String(`loglevel`, `debug`, `log level: Panic/Fatal/Error/Warn/Info/Debug/Trace`)
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)
	viper.SetConfigFile(viper.GetString(`config`))
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
