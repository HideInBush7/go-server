package util

import "net"

func GetInternalIp() string {
	addrs, _ := net.InterfaceAddrs()

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ``
}
