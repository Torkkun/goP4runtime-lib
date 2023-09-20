package architecture

import (
	simpleswitch "gop4runtimelib/shimpleswitch"
	"log"
	"os"
)

type Bmv2SwitchConnection struct {
	*simpleswitch.SwitchConnection
}

func (Bmv2SwitchConnection) BuildDeviceConfig(bmv2jpath string) []byte {
	// Build the device config for BMv2
	var deviceConfig []byte
	f, err := os.ReadFile(bmv2jpath)
	if err != nil {
		log.Fatalf("build Device Config failed: %v", err)
	}
	deviceConfig = f
	return deviceConfig
}
