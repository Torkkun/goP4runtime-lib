package architecture

import simpleswitch "gop4runtimelib/shimpleswitch"

// DIしてもいい
type Bmv2SwitchConnection struct {
	*simpleswitch.SwitchConnection
}

func (Bmv2SwitchConnection) buildDeviceConfig() {
	// Build the device config for BMv2

}
