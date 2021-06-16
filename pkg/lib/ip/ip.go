package ip

import (
	"encoding/binary"
	"fmt"
	"net"
)

func ParseIPv4(addr [6]byte) string {
	ip, port := addr[:4], addr[4:6]

	parsedIP := net.IPv4(ip[0], ip[1], ip[2], ip[3])
	parsedPort := binary.BigEndian.Uint16(port)

	return fmt.Sprintf("%s:%d", parsedIP, parsedPort)
}
