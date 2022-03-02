package util

import (
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// 获取本机内网IP
func GetInternalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		// 非环回地址
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ``
}

// 获取当前程序运行的文件夹
// 因为使用相对路径的话在不同目录下调用可执行文件会找不到
func GetExecDir() string {
	dir, err := filepath.Abs(path.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	// windows下分隔符不一致
	return strings.ReplaceAll(dir, "\\", "/")
}
