package utility

import (
	"github.com/d2r2/go-logger"
	"net"
	"regexp"
	"strings"
)

var lg = logger.NewPackageLogger("net", logger.InfoLevel)

// LogNetworkInterfacesAndGetIpAdr logs the network interfaces and returns the first non-loopback
// IPv4 address as a string.
func LogNetworkInterfacesAndGetIpAdr() string {
	var ipAddress string
	interfaces, err := net.Interfaces()
	if err != nil {
		lg.Error(err.Error())
		return "Error"
	}
	reg := regexp.MustCompilePOSIX("^((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.){3}(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])")
	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			lg.Warn(err.Error())
		}
		err = nil
		addresses, _ := byName.Addrs()
		for _, v := range addresses {
			ipv4 := v.String()
			if reg.MatchString(ipv4) {
				lg.Info(ipv4)
				if strings.Index(ipv4, "127.0.") != 0 {
					idx := strings.Index(ipv4, "/")
					if idx > 0 {
						ipAddress = ipv4[0:idx]
					} else {
						ipAddress = ipv4
					}
				}
			}
		}
	}
	return ipAddress
}
