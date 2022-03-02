
## 服务发现用例

配置config.yaml中etcd地址

本机有etcd就直接127.0.0.1:6379,此处为docker机的运行环境

可以把config.yaml放在main.go的同级目录,省略参数 --config

可以在不同的机器上开启服务端,即为单服务的分布式部署,通过etcd进行服务注册与发现

```            
# 服务端 -> 注册服务
cd example/order/server
# 开启第一个服务 默认跑在本机的8000端口
go run server.go --config=../config.yaml
# 开启同一个服务的另一端口 8001
go run server.go --config=../config.yaml -port=8001

# 客户端 -> 发现服务
cd example/order/client
# 发现服务,以轮询的方式请求服务
go run client.go --config=../config.yaml

```