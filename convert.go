package main

import (
	"fmt"
	"net"
	"regexp"
)

const (
	macPattern = `^([\da-fA-F]{2}:){5}([\da-fA-F]{2})$`
	ipPattern  = `^(\d{1,3}\.){3}(\d{1,3})$`
)

func matchesMac(macAddr string) (bool, error) {
	macPattern, err := regexp.Compile(macPattern)
	if err != nil {
		err = fmt.Errorf("matchesMac function err: %v", err)
		return false, err
	}
	return macPattern.MatchString(macAddr), nil
}

func encodeMac(macAddr string) (net.HardwareAddr, error) {
	return net.ParseMAC(macAddr)
}

func decodeMac(encodedMacAddr net.HardwareAddr) string {
	return encodedMacAddr.String()
}

func matchIPv4(ipv4Addr string) (bool, error) {
	ipPattern, err := regexp.Compile(ipPattern)
	if err != nil {
		return false, err
	}
	return ipPattern.MatchString(ipv4Addr), nil
}
