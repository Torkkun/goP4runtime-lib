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

type EntityType int

const (
	TableEntity EntityType = iota
	CounterEntity
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

func (pih *P4InfoHelper) GetID(et EntityType, name string) (uint32, error) {
	switch et {
	case TableEntity:
		for _, v := range pih.P4Info.Tables {
			if v.Preamble.Name == name {
				return v.Preamble.Id, nil
			}
		}
	case CounterEntity:
		for _, v := range pih.P4Info.Counters {
			if v.Preamble.Name == name {
				return v.Preamble.Id, nil
			}
		}
	}
	return 0, fmt.Errorf("While Get ID, But not exist name in Entity")
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
			continue
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

type ExactValue string
type LpmValue struct {
	Dst string
	Id  int32
}

// valueの与え方はもうちょっと考える
// 処理がちょっと不明
func (pih *P4InfoHelper) GetMatchFieldPb(tablename string, matchfiledname string, value interface{}) (*v1.FieldMatch, error) {
	p4infomatch, err := pih.GetMatchField(tablename, isSameName(matchfiledname))
	if err != nil {
		return nil, err
	}
	//bitwidth := p4infomatch.Bitwidth
	p4rtm := new(v1.FieldMatch)
	p4rtm.FieldId = p4infomatch.Id
	mtype := p4infomatch.GetMatchType()
	switch mtype {
	case v1conf.MatchField_EXACT:
		exact := new(v1.FieldMatch_Exact)
		//exact.Value, err = encode(value.(ExactValue), bitwidth)
		p4rtm.FieldMatchType = &v1.FieldMatch_Exact_{Exact: exact}
	case v1conf.MatchField_LPM:
		//lpme := p4rtm.GetLpm()
		//lpme.Value, err = encode(value.(LpmValue).string, bitwidth)
		//lpme.PrefixLen = value.(LpmValue).int32

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
func MatchFields(match isMatchFields) TableEntryOptions {
	return func(te *v1.TableEntry) {
		// get match field pbに	table nameとmatch field nameとvalueを引数に取らせる
		// pending

		te.Match = append(te.Match)
	}
}

type isMatchFields interface {
	isMatchFields()
}

type MatchField_EXACT struct {
	ExactField map[string]uint32
}

type MatchField_LPM struct {
	LpmField map[string]map[string]uint32
}

func (*MatchField_EXACT) isMatchFields() {}

func (*MatchField_LPM) isMatchFields() {}

func DeFaultAction() TableEntryOptions {
	return func(te *v1.TableEntry) {}
}

// TODO:
func ActionName(aname string) TableEntryOptions {
	return func(te *v1.TableEntry) {

	}
}

// TODO:
func ActionParam(aprm map[string]interface{}) TableEntryOptions {
	return func(te *v1.TableEntry) {

	}
}

// 引数の値の与え方は要検討
// とりあえずオプションにしているけど仕様書を見ながら変える
// もし、Fieldがある程度理解度が上がれば、それぞれExactやLPMごとにTableEntry構造体を作成するかも
func (pih *P4InfoHelper) BuildTableEntry(tablename string, options ...TableEntryOptions) (*v1.TableEntry, error) {
	table := new(v1.TableEntry)
	id, err := pih.GetID(TableEntity, tablename)
	if err != nil {
		return nil, err
	}
	table.TableId = id
	//execute setting Option argument when exist Option Parameter in table entry struct
	for _, o := range options {
		o(table)
	}
	return table, nil
}

func (pih *P4InfoHelper) BuidlMulticastGroupEntry() {
}

func (pih *P4InfoHelper) BuildCloneSessionEntry() {

}
