package utils

import (
	"net"
	"strconv"
	"strings"
)

// SplitHostPort 将主机名与端口号分离。
func SplitHostPort(host string) (string, uint64) {
	p0 := strings.LastIndexByte(host, ':')
	if p0 < 0 {
		return strings.TrimSpace(host), 0
	}
	if p0 > 0 && host[p0-1] == ':' {
		return strings.TrimSpace(host), 0
	}

	hap0 := host[:p0]
	hap1 := host[p0+1:]
	if pv, err := strconv.ParseUint(hap1, 10, 64); err != nil {
		return strings.TrimSpace(hap0), 0
	} else {
		return strings.TrimSpace(hap0), pv
	}
}

// GetHostAddress 获取指定地址对应的IP地址。
// address 可以是非回环IP地址，或者表示本机的回环地址，或者是主机名。
func GetHostAddress(address string) (net.IP, error) {
	if ipNet, err := net.ResolveIPAddr("ip", address); err != nil {
		return nil, err
	} else if ipNet.IP.IsLoopback() {
		// 如果指定的地址是个回环地址（即表示本机），那么尝试获取本机的非回环地址。
		if addrs, err := net.InterfaceAddrs(); err != nil {
			return nil, err
		} else {
			// 遍历所有绑定到本地的地址，查找非回环地址。
			var nLocalIp net.IP
			for _, addr := range addrs {
				if localIpNet, ok := addr.(*net.IPNet); ok {
					localIp := localIpNet.IP
					if !localIp.IsLoopback() &&
						!localIp.IsLinkLocalMulticast() &&
						!localIp.IsLinkLocalUnicast() {
						if localIpNetV4 := localIp.To4(); localIpNetV4 != nil {
							return localIpNetV4, nil
						} else if len(nLocalIp) == 0 {
							nLocalIp = localIp
						}
					}
				}
			}
			return nLocalIp, nil
		}
	} else {
		// 非回环地址，直接返回。
		return ipNet.IP, nil
	}
}

// GetLocalIpAddress 获取本地IP地址，此方法会尽可能返回IPv4地址。
func GetLocalIpAddress() string {
	if ipAddr, err := GetHostAddress("localhost"); err != nil {
		return ""
	} else {
		return ipAddr.String()
	}
}
