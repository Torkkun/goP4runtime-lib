package helper

import (
	"fmt"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	v1conf "github.com/p4lang/p4runtime/go/p4/config/v1"
	v1 "github.com/p4lang/p4runtime/go/p4/v1"
)

type P4InfoHelper struct {
	P4Info *v1conf.P4Info
}

type EntityType string

const (
	entityTypeTables  EntityType = "tables"
	entityTypeActions EntityType = "actions"
)

func NewP4InfoHelper(p4InfoFilePath string) *P4InfoHelper {
	p4info := new(v1conf.P4Info)
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

func (P4InfoHelper) Get() {

}

func (pih *P4InfoHelper) GetID(entitytype EntityType, name string) uint32 {

}

func (pih *P4InfoHelper) GetName() {

}

func (pih *P4InfoHelper) GetAlias() {

}

type MatchFieldOption func(*v1conf.MatchField) bool

func isSameName(name string) MatchFieldOption {
	return func(mf *v1conf.MatchField) bool {
		return mf.Name == name
	}
}

func isSameId(id uint32) MatchFieldOption {
	return func(mf *v1conf.MatchField) bool {
		return mf.Id == id
	}
}

func (pih *P4InfoHelper) GetMatchField(tablename string, option ...MatchFieldOption) (*v1conf.MatchField, error) {
	for _, t := range pih.P4Info.Tables {
		pre := t.Preamble
		if pre.Name != tablename {
			return nil, fmt.Errorf("Table Name not Equal Error: PreambleName=%s, tablename=%s", pre.Name, tablename)
		}
		for _, mf := range t.MatchFields {
			for _, o := range option {
				if o(mf) {
					return mf, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("Not match Fieleds Error")
}

func (pih *P4InfoHelper) GetMatchFieldId() {

}

func (pih *P4InfoHelper) GetMatchFieldName() {

}

// valueの与え方はもうちょっと考える
func (pih *P4InfoHelper) GetMatchFieldPb(tablename string, matchfiledname string, value []interface{}) (*v1.FieldMatch, error) {
	p4infomatch, err := pih.GetMatchField(tablename, isSameName(matchfiledname))
	if err != nil {
		return nil, err
	}
	bitwidth := p4infomatch.Bitwidth
	p4rtm := new(v1.FieldMatch)
	p4rtm.FieldId = p4infomatch.Id
	mtype := p4infomatch.GetMatchType()
	switch mtype {
	case v1conf.MatchField_EXACT:
		exact := p4rtm.GetExact()
		exact.Value, err = encode(value, bitwidth)
	case v1conf.MatchField_LPM:
		lpme := p4rtm.GetLpm()
		lpme.Value, err = encode(value[0], bitwidth)
		lpme.PrefixLen = value[1].(int32)
	case v1conf.MatchField_TERNARY:

	case v1conf.MatchField_RANGE:

	case v1conf.MatchField_OPTIONAL:

	default:
		return nil, fmt.Errorf("MatchFieldType is %s", v1conf.MatchField_UNSPECIFIED.String())
	}
	return p4rtm, nil
}

func (pih *P4InfoHelper) GetMatchFieldValue() {

}

func (pih *P4InfoHelper) GetActionParam() {

}

func (pih *P4InfoHelper) GetActionParamId() {

}

func (pih *P4InfoHelper) GetActionParamName() {

}

func (pih *P4InfoHelper) GetActionParamPb() {

}

type TableEntryOptions func(*v1.TableEntry)

func Priority(priority int32) TableEntryOptions {
	return func(te *v1.TableEntry) {
		te.Priority = priority
	}
}

// mapで受け取る key is matchfield name , value is any id
func MatchFields(match map[string]uint32) TableEntryOptions {
	return func(te *v1.TableEntry) {
		// get match field pbに	table nameとmatch field nameとvalueを引数に取らせる
		// pending
		te.Match = append(te.Match)
	}
}

func DeFaultAction() TableEntryOptions {
	return func(te *v1.TableEntry) {}
}

func ActionName() TableEntryOptions {
	return func(te *v1.TableEntry) {}
}

// 引数の値の与え方は要検討
func (pih *P4InfoHelper) BuildTableEntry(tablename string, TableEntryoptions ...TableEntryOptions) *v1.TableEntry {
	table := new(v1.TableEntry)
	table.TableId = pih.GetID(entityTypeTables, tablename)
	//execute setting Option argument when exist Option Parameter in table entry struct
	for _, o := range TableEntryoptions {
		o(table)
	}
	return table
}

func (pih *P4InfoHelper) BuidlMulticastGroupEntry() {
}

func (pih *P4InfoHelper) BuildCloneSessionEntry() {

}
