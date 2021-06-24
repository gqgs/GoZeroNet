package ip

import (
	"encoding/binary"
	"fmt"
	"net"
)

func ParseIPv4(addr []byte, byteOrder binary.ByteOrder) string {
	if len(addr) != 6 {
		return ""
	}

	ip, port := addr[:4], addr[4:6]

	parsedIP := net.IPv4(ip[0], ip[1], ip[2], ip[3])
	parsedPort := byteOrder.Uint16(port)

	return fmt.Sprintf("%s:%d", parsedIP, parsedPort)
}
