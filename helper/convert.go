package helper

import (
	"fmt"
	"log"
	"math"
	"net"
	"regexp"
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
	return net.ParseIP(ipv4Addr)
}

func decodeIPv4(encodedIpv4Addr net.IP) string {
	return encodedIpv4Addr.String()
}

// bit to byte length conversion
// セグフォしそうなので一応注意
func bitwidthToBytes(bitwidth int32) int {
	return int(math.Ceil(float64(bitwidth) / 8.0))
}

func encodeNum(number, bitwidth int32) ([]byte, error) {
	//byteLen := bitwidthToBytes(bitwidth)
	/* origNumber := number

	if number < 0 {
		if number <- (int(math.Pow(2, float64(bitwidth)-1))) {
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
	}*/
	return nil, nil
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
	case int32:
		encodedbytes, err = encodeNum(v, bitwidth)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Encoding objects of %v is not supported", v)

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
