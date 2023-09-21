package helper

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"regexp"
	"strings"
)

const (
	macPattern = `^([\da-fA-F]{2}:){5}([\da-fA-F]{2})$`
	ipPattern  = `^(\d{1,3}\.){3}(\d{1,3})$`
)

func matchesMac(macAddr string) bool {
	macPattern, err := regexp.Compile(macPattern)
	if err != nil {
		log.Fatalf("matchesMac function faital: %v\n", err)

	}
	return macPattern.MatchString(macAddr)
}

func encodeMac(macAddr string) ([]byte, error) {
	return net.ParseMAC(macAddr)
}

func decodeMac(encodedMacAddr net.HardwareAddr) string {
	return encodedMacAddr.String()
}

func matchesIPv4(ipv4Addr string) bool {
	ipPattern, err := regexp.Compile(ipPattern)
	if err != nil {
		log.Fatalf("matchesIPv4 function faital: %v\n", err)
	}
	return ipPattern.MatchString(ipv4Addr)
}

func encodeIPv4(ipv4Addr string) []byte {
	return net.ParseIP(ipv4Addr).To4()
}

func decodeIPv4(encodedIpv4Addr net.IP) string {
	return encodedIpv4Addr.String()
}

// bit to byte length conversion
func bitwidthToBytes(bitwidth int32) int {
	return int(math.Ceil(float64(bitwidth) / 8.0))
}

func encodeNum(number, bitwidth int32) ([]byte, error) {
	byteLen := bitwidthToBytes(bitwidth)
	origNumber := number

	if number < 0 {
		if number < -(1 << (bitwidth - 1)) {
			return nil, fmt.Errorf("nagative namuber %d has 2's complete representation that does nao fit in %d bits", number, bitwidth)
		}
		number = (1 << bitwidth) + number
	}
	numStr := fmt.Sprintf("%x", number)
	if origNumber < 0 {
		fmt.Printf("CONVERT_NEGATIVE_NUMBER debug: origNumber=%d number=%d bitwidth=%d numStr=%s\n",
			origNumber, number, bitwidth, numStr)
	}
	if number >= 1<<bitwidth {
		return nil, fmt.Errorf("number, %d, does not fit in %d bits", number, bitwidth)
	}

	// Create a hex string with leading zeros to fill the required byte length.
	paddedNumStr := strings.Repeat("0", byteLen*2-len(numStr)) + numStr
	decodedBytes, err := hexStringToBytes(paddedNumStr)
	if err != nil {
		return nil, err
	}
	return decodedBytes, nil
}

func hexStringToBytes(hexStr string) ([]byte, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	bytesLen := len(hexStr) / 2
	result := make([]byte, bytesLen)

	for i := 0; i < bytesLen; i++ {
		n, err := fmt.Sscanf(hexStr[2*i:2*i+2], "%02x", &result[i])
		if err != nil || n != 1 {
			return nil, errors.New("failed to decode hex string to bytes")
		}
	}

	return result, nil
}

func decodeNum() int {
	return 0
}

func encode(data interface{}, bitwidth int32) ([]byte, error) {
	bytelen := bitwidthToBytes(bitwidth)
	var encodedbytes []byte
	var err error
	switch v := data.(type) {
	case string:
		if matchesMac(v) {
			encodedbytes, err = encodeMac(v)
			if err != nil {
				return nil, err
			}
		} else if matchesIPv4(v) {
			encodedbytes = encodeIPv4(v)
		} else {
			encodedbytes = []byte(v)
		}
	case int:
		encodedbytes, err = encodeNum(int32(v), bitwidth)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("encoding objects of %v is not supported", v)

	}
	if len(encodedbytes) != bytelen {
		return nil, fmt.Errorf("encodedbytes and bytelen length is not equal")
	}
	return encodedbytes, nil
}

func encodedDst(dst string) ([]byte, error) {
	var encodedbytes []byte
	var err error
	if matchesMac(dst) {
		encodedbytes, err = encodeMac(dst)
		if err != nil {
			return nil, err
		}
	} else if matchesIPv4(dst) {
		encodedbytes = encodeIPv4(dst)
	} else {
		//imm
		encodedbytes = []byte(dst)
	}
	return encodedbytes, err
}
