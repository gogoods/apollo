package apollo

import (
	"net"
	"os"
	"strings"
)

var ip string

const (
	prefixIP = "ip="
)

func init() {

	for _, v := range os.Args {
		if strings.HasPrefix(v, prefixIP) {
			ip = strings.Replace(v, prefixIP, "", 1)
			return
		}
	}

	conn, err := net.Dial("udp", "1.2.3.4:80")
	if err != nil {
		ip = "127.0.0.0"
		return
	}
	defer conn.Close()
	local := conn.LocalAddr().(*net.UDPAddr)
	ip = local.IP.String()
}

func LocalIP() string {
	return ip
}
