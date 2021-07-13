package ip

import (
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
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

func PackIPv4(ip string, byteOrder binary.ByteOrder) []byte {
	host, port, err := net.SplitHostPort(ip)
	if err != nil {
		return nil
	}

	parsedHost := net.ParseIP(host)
	if parsedHost == nil {
		return nil
	}
	parsedHost = parsedHost.To4()

	i, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil
	}

	parsedPort := make([]byte, 2)
	byteOrder.PutUint16(parsedPort, uint16(i))

	return []byte{parsedHost[0], parsedHost[1], parsedHost[2], parsedHost[3], parsedPort[0], parsedPort[1]}
}

func ParseOnion(addr []byte, byteOrder binary.ByteOrder) string {
	if len(addr) < 5 {
		return ""
	}
	host, port := addr[:len(addr)-2], addr[len(addr)-2:]
	onion := strings.ToLower(base32.HexEncoding.EncodeToString(host))
	return fmt.Sprintf("%s.onion:%d", onion, byteOrder.Uint16(port))
}
