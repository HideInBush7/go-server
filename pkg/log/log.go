package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 日志初始化配置 依赖于config.Init()
func Init() {
	// 日志等级
	level, err := logrus.ParseLevel(viper.GetString(`loglevel`))
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(level)

	// 设置软链和日志分割
	writer, err := rotatelogs.New(
		viper.GetString(`logfile`)+"%Y-%m-%d"+".log",
		rotatelogs.WithLinkName("log.log"),
	)
	if err != nil {
		logrus.Panic(err)
	}

	// 打开控制台输出 通过设置hook
	if viper.GetBool(`debug`) {
		hook := lfshook.NewHook(writer, &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			PrettyPrint:     false, //是否格式化json格式
		})
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			DisableColors:   false,
		})
		logrus.AddHook(hook)
		return
	}

	// 关闭控制台输出,直接SetOutPut
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     false,
	})
	logrus.SetOutput(writer)
}
