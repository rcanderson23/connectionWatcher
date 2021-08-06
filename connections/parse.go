package connections

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"
)

// parseEndpoint splits the endpoint hex string delimited by ':' into IP and port
// returns the IP, port, and errors
func parseEndpoint(ep string) (net.IP, uint16, error) {
	if len(ep) != 13 {
		err := fmt.Errorf("length of string does not equal 13: %d", len(ep))
		return nil, 0, err
	}

	split := strings.Split(ep, ":")
	if len(split) != 2 {
		err := fmt.Errorf("2 strings expected from split, got %d", len(split))
		return nil, 0, err
	}

	ip, err := parseIPv4(split[0])
	if err != nil {
		return nil, 0, err
	}

	port, err := parsePort(split[1])
	if err != nil {
		return nil, 0, err
	}

	return ip, port, nil
}

// parseIPv4 expects an IP in hex little endian and returns an IP struct
func parseIPv4(s string) (net.IP, error) {
	ipBytes, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	ip := net.IPv4(ipBytes[3], ipBytes[2], ipBytes[1], ipBytes[0])
	if ip == nil {
		return nil, errors.New("failed to parse ip")
	}

	return ip, nil
}

// parsePort expects a port in hex and returns the uint16
func parsePort(s string) (uint16, error) {
	portBytes, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(portBytes), nil
}
