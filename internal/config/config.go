package config

import (
	"net"
	"strconv"
)

type Config struct {
	MyIP   string
	MyPort int
}

func (c Config) GetMyIPPort() string {
	return c.MyIP + ":" + strconv.FormatInt(int64(c.MyPort), 10)
}

var c Config

func SetupConfig(appPort int) {
	c.MyIP = getLocalIP()
	//TODO: get from env
	c.MyPort = appPort
}

func GetConfig() Config {
	return c
}

// GetLocalIP returns the non loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
