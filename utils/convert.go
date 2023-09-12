package utils

import (
	"fmt"
	"math"
	"net"
	"regexp"
	"strconv"
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

func encodeIPv4(ipv4Addr string) net.IP {
	return net.ParseIP(ipv4Addr)
}

func decodeIPv4(encodedIpv4Addr net.IP) string {
	return encodedIpv4Addr.String()
}

func bitwidthToBytes(bitwidth int) int {
	return int(math.Ceil(float64(bitwidth) / 8.0))
}

// fix:
// そもそも何したい部分だ？
func encodeNum(number, bitwidth int) error {
	//byteLen := bitwidthToBytes(bitwidth)
	origNumber := number

	if number < 0 {
		if number < -(int(math.Pow(2, float64(bitwidth)-1))) {
			return fmt.Errorf("Nagative namuber, %d, has 2's complete representation that does nao fit in %d bits\n", number, bitwidth)
		}
		number = int(math.Pow(2, float64(bitwidth))) + number
	}
	numStr := strconv.Itoa(number)
	if origNumber < 0 {
		fmt.Printf("CONVERT_NEGATIVE_NUMBER debug: origNumber=%d number=%d bitwidth=%d numStr=%s\n",
			origNumber, number, bitwidth, numStr)
	}
	if number >= int(math.Pow(2, float64(bitwidth))) {
		return fmt.Errorf("Number, %d, does not fit in %d bits\n", number, bitwidth)
	}
	return nil
}

func decodeNum() int {
	return 0
}

func encode() {

}
