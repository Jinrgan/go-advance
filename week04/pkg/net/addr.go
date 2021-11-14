package net

import (
	"fmt"
	"net"
)

func GetPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0")
	if err != nil {
		return 0, fmt.Errorf("cannot resolve tcp addr: %v", err)
	}

	return addr.Port, nil
}
