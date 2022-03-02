## 大部分功能基于config.Init()的公共库

### 主要有日志(logrus)、mysql(sqlx)、redis(redigo)
### 还有之前实现的etcd服务注册与发现
``` shell
例子: 
cd /pkg/example
修改config.yaml对应文件

# 因为此处使用的是文件执行路径
# 直接go run main.go会在一个临时运行目录,所以先编译后运行
go build . 

./example

```