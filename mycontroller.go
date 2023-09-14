package main

import (
	"flag"
	"gop4runtimelib/architecture"
	simpleswitch "gop4runtimelib/shimpleswitch"
	"log"
	"os"
)

var SWITCH_TO_HOST_PORT = 1
var SWITCH_TO_SWITCH_PORT = 2

func main() {
	p4info := flag.String(
		"p4info",
		"./build/advanced_tunnel.p4.p4info.txt",
		"p4info proto in text format from p4c")

	bmv2json := flag.String(
		"bmv2json",
		"./build/advanced_tunnel.json",
		"BMv2 JSON file from p4c")
	flag.Parse()

	if _, err := os.Stat(*p4info); err != nil {
		log.Fatalf("File %s does not exist!: %v", *p4info, err)
	}

	if _, err := os.Stat(*bmv2json); err != nil {
		log.Fatalf("File %s does not exist!: %v", *bmv2json, err)
	}

	s1 := architecture.Bmv2SwitchConnection{
		SwitchConnection: simpleswitch.SwitchConnection{
			Name:          "s1",
			Address:       "127.0.0.1:50051",
			DeviceId:      0,
			ProtoDumpFile: "logs/s1-p4runtime-requests.txt",
		},
	}
	s2 := architecture.Bmv2SwitchConnection{
		SwitchConnection: simpleswitch.SwitchConnection{
			Name:          "s2",
			Address:       "127.0.0.1:50052",
			DeviceId:      1,
			ProtoDumpFile: "logs/s2-p4runtime-requests.txt",
		},
	}

	s1.MasterArbitrationUpdate()
	s2.MasterArbitrationUpdate()
}

func writeTunnelRules(id uint32, dstethaddr, dstipaddr string) {

}
