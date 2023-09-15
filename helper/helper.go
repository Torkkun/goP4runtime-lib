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
	MatchField
)

func NewP4InfoHelper(p4InfoFilePath string) *P4InfoHelper {
	p4info := new(v1conf.P4Info)
	file, err := os.ReadFile(p4InfoFilePath)
	if err != nil {
		log.Fatalf("os Read File Faital: %v", err)
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
	return 0, fmt.Errorf("while Get ID, But not exist name in Entity")
}

func (pih *P4InfoHelper) GetName(et EntityType, id uint32) (string, error) {
	switch et {
	case TableEntity:
		for _, v := range pih.P4Info.Tables {
			if v.Preamble.Id == id {
				return v.Preamble.Name, nil
			}
		}
	case CounterEntity:
		for _, v := range pih.P4Info.Counters {
			if v.Preamble.Id == id {
				return v.Preamble.Name, nil
			}
		}
	}
	return "", fmt.Errorf("while Get Name, But not exist id in Entity")

}

func (pih *P4InfoHelper) GetMatchFieldId(tname, name string) (*v1conf.MatchField, error) {
	for _, t := range pih.P4Info.Tables {
		pre := t.Preamble
		if pre.Name == tname {
			for _, mf := range t.MatchFields {
				if mf.Name == name {
					return mf, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("can not Find Match Filed Name is %s in %s Table", name, tname)
}

func (pih *P4InfoHelper) GetMatchFieldName(tname string, id uint32) (*v1conf.MatchField, error) {
	for _, t := range pih.P4Info.Tables {
		pre := t.Preamble
		if pre.Name == tname {
			for _, mf := range t.MatchFields {
				if mf.Id == id {
					return mf, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("can not Find Match Filed Name is %d in %s Table", id, tname)
}

// valueの与え方はもうちょっと考える
// 処理がちょっと不明
func (pih *P4InfoHelper) GetMatchFieldPb(tablename string, iemf *IsEntryMatchField) (*v1.FieldMatch, error) {
	mf, err := pih.GetMatchFieldId(tablename, iemf.Name)
	if err != nil {
		return nil, err
	}
	p4rtm := new(v1.FieldMatch)
	p4rtm.FieldId = mf.Id
	bitwidth := mf.Bitwidth
	switch mf.GetMatchType() {
	case v1conf.MatchField_EXACT:
		values := iemf.EntryMatchField.(*EntryMatchFieldExact)
		exact := new(v1.FieldMatch_Exact)
		exact.Value, err = encode(values.Value0, bitwidth)
		if err != nil {
			return nil, err
		}
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

func (pih *P4InfoHelper) GetMatchFieldValue(mf *v1conf.MatchField) {

	switch mf.GetMatchType() {
	case v1conf.MatchField_EXACT:

	case v1conf.MatchField_LPM:
	case v1conf.MatchField_TERNARY:
	case v1conf.MatchField_RANGE:
	default:

	}

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

/* // mapで受け取る key is matchfield name , value is any id
func MatchFields(match isMatchFields, mfmt v1conf.MatchField_MatchType) TableEntryOptions {
	return func(te *v1.TableEntry) {
		// get match field pbに	table nameとmatch field nameとvalueを引数に取らせる
		// 引数からKeyとValueを引き出しGetMatchFieldPBへテーブル名と一緒にぶち込む
		switch mfmt {
		case v1conf.MatchField_EXACT:

			mexact := new(v1.FieldMatch_Exact_) // このかたで返すPbこれは消す

			m.FieldMatchType = mexact
			te.Match = append(te.Match, m)
		case v1conf.MatchField_TERNARY:
		case v1conf.MatchField_LPM:
		case v1conf.MatchField_RANGE:
		case v1conf.MatchField_OPTIONAL:

		}

	}

} */

/* type isMatchFields interface {
	isMatchFields()
}

type MatchField_EXACT struct {
	ExactField map[string]uint32
}

type MatchField_LPM struct {
	LpmField map[string]map[string]uint32
} */

/* func (*MatchField_EXACT) isMatchFields() {}

func (*MatchField_LPM) isMatchFields() {} */

type IsEntryMatchField struct {
	Name            string
	EntryMatchField isEntryMatchFieldType
}

type isEntryMatchFieldType interface {
	isEntryMatchFieldType()
}

type EntryMatchFieldExact struct {
	Value0 uint32
}

type EntryMatchFieldLpm struct {
	Value0 string
	Value1 uint32
}

func (*EntryMatchFieldExact) isEntryMatchFieldType() {}

func (*EntryMatchFieldLpm) isEntryMatchFieldType() {}

func DeFaultAction(defact bool) TableEntryOptions {
	return func(te *v1.TableEntry) {
		te.IsDefaultAction = defact
	}
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
func (pih *P4InfoHelper) BuildTableEntry(
	tablename string,
	iemft *IsEntryMatchField,
	options ...TableEntryOptions,
) (*v1.TableEntry, error) {
	id, err := pih.GetID(TableEntity, tablename)
	if err != nil {
		return nil, err
	}
	table := new(v1.TableEntry)
	table.TableId = id
	pih.GetMatchFieldPb(tablename, iemft)

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
