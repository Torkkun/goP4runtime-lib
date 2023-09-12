package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	serverAddr := flag.String("a", "", "address and port of the switch\n P4 Runtime server (e.g. 192.168.0.1:50051)")
	deviceId := flag.Int("d", 0, "Internal device ID to use in P4Runtime messages")
	protoDumpFile := flag.String("p", "", "path to file where to dump protobuf message sent to the switch")
	runtimeConfFile := flag.String("c", "", "path to input runtime configuration file (JSON)")
	flag.Parse()
	newsw := newSwitch(runtimeConfFile)
	newsw.addr = *serverAddr
	newsw.deviceId = *deviceId
	newsw.protoDumpFilePath = *protoDumpFile

}

type swConf struct {
}

func checkSwitchConf() {

}

type Switch struct {
	addr              string
	deviceId          int
	swConfFile        string
	workDir           string
	protoDumpFilePath string
}

func newSwitch(runtimeConfFile *string) *Switch {
	if _, err := os.Stat(*runtimeConfFile); os.IsNotExist(err) {
		log.Fatalf("File %s does not exist!: %v", *runtimeConfFile, err)
	}
	fileAbsPath, err := filepath.Abs(*runtimeConfFile)
	if err != nil {
		log.Fatalf("Filepath Absolute Path missing: %v", err)
	}
	// working dir
	workdir := filepath.Dir(fileAbsPath)

	//load sw conf
	swConfFile, err := os.Open(*runtimeConfFile)
	if err != nil {
		log.Fatal(err)
	}
	defer swConfFile.Close()

	return &Switch{
		workDir: workdir,
	}
}

func programSwitch(sw Switch) {

}

func validateTableEntry() {

}

func insertTableEntry() {

}

func jsonLoadByteified(*os.File) {

}

func byteify() {

}

func tableEntryToString() {

}

func groupEntryToString() {

}

func cloneEntryToString() {

}

func insertMulticastGroupEntry() {

}

func insertCloneGroupEntry() {

}
