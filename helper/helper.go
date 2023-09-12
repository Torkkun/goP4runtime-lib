package helper

import (
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	v1 "github.com/p4lang/p4runtime/go/p4/config/v1"
)

type P4InfoHelper struct {
	P4Info *v1.P4Info
}

func NewP4InfoHelper(p4InfoFilePath string) *P4InfoHelper {
	p4info := new(v1.P4Info)
	file, err := os.ReadFile(p4InfoFilePath)
	if err != nil {
		log.Fatalln("OS Read File Faital: %v", err)
	}

	if err := proto.UnmarshalText(string(file), p4info); err != nil {
		log.Fatalf("UnmarshalText error: %v", err)
	}
	return &P4InfoHelper{
		P4Info: p4info,
	}
}

func (P4InfoHelper) get() {

}

func (p4h *P4InfoHelper) getID() {

}

func (p4h *P4InfoHelper) getName() {

}

func (p4h *P4InfoHelper) getAlias() {

}
