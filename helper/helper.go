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
	ActionEntity
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
	case ActionEntity:
		for _, v := range pih.P4Info.Actions {
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

func (pih *P4InfoHelper) GetMatchFields(tablename string) []*v1conf.MatchField {
	for _, t := range pih.P4Info.Tables {
		pre := t.Preamble
		if pre.Name == tablename {
			return t.MatchFields
		}
	}
	return nil
}

func (pih *P4InfoHelper) GetMatchFieldWithName(tablename, name string) (*v1conf.MatchField, error) {
	for _, mf := range pih.GetMatchFields(tablename) {
		if mf.Name == name {
			return mf, nil
		}
	}
	return nil, fmt.Errorf("can not Find Match Filed Name is %s in %s Table", name, tablename)
}

func (pih *P4InfoHelper) GetMatchFieldWithId(tablename string, id uint32) (*v1conf.MatchField, error) {
	for _, mf := range pih.GetMatchFields(tablename) {
		if mf.Id == id {
			return mf, nil
		}
	}
	return nil, fmt.Errorf("can not Find Match Filed Name is %d in %s Table", id, tablename)
}

// valueの与え方はもうちょっと考える
// 処理がちょっと不明
func (pih *P4InfoHelper) GetMatchFieldPb(tablename string, iemf *IsEntryMatchField) (*v1.FieldMatch, error) {
	mf, err := pih.GetMatchFieldWithName(tablename, iemf.Name)
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
		values := iemf.EntryMatchField.(*EntryMatchFieldLpm)
		lpm := new(v1.FieldMatch_LPM)
		lpm.Value, err = encode(values.Value0, bitwidth)
		if err != nil {
			return nil, err
		}
		lpm.PrefixLen = values.Value1
		p4rtm.FieldMatchType = &v1.FieldMatch_Lpm{Lpm: lpm}
	case v1conf.MatchField_TERNARY:

	case v1conf.MatchField_RANGE:

	case v1conf.MatchField_OPTIONAL:

	default:
		return nil, fmt.Errorf("MatchFieldType is %s", v1conf.MatchField_UNSPECIFIED.String())
	}
	return p4rtm, nil
}

func (pih *P4InfoHelper) GetMatchFieldValue(mf *v1conf.MatchField, m *v1.FieldMatch) ([]byte, error) {
	switch mf.GetMatchType() {
	case v1conf.MatchField_EXACT:
		v := m.FieldMatchType.(*v1.FieldMatch_Exact_).Exact.Value
		return v, nil
	case v1conf.MatchField_LPM:
		v := m.FieldMatchType.(*v1.FieldMatch_Lpm).Lpm.Value
		return v, nil
	case v1conf.MatchField_TERNARY:
	case v1conf.MatchField_RANGE:
	}
	return nil, fmt.Errorf("FieldMatch Value Error")
}

func (pih *P4InfoHelper) GetActionParams(actionname string) []*v1conf.Action_Param {
	for _, a := range pih.P4Info.Actions {
		pre := a.Preamble
		if pre.Name == actionname {
			return a.Params
		}
	}
	return nil
}

func (pih *P4InfoHelper) GetActionParamWithName(actionname, name string) (*v1conf.Action_Param, error) {
	for _, ap := range pih.GetActionParams(actionname) {
		if ap.Name == name {
			return ap, nil
		}
	}
	return nil, fmt.Errorf("can not Find Action Params Name is %s in %s Table", name, actionname)
}

func (pih *P4InfoHelper) GetActionParamWithId(actionname string, id uint32) (*v1conf.Action_Param, error) {
	for _, ap := range pih.GetActionParams(actionname) {
		if ap.Id == id {
			return ap, nil
		}
	}
	return nil, fmt.Errorf("can not Find Action Params ID is %d in %s Table", id, actionname)
}

func (pih *P4InfoHelper) GetActionParamPb(actionname string, param map[string]interface{}) ([]*v1.Action_Param, error) {
	var aps []*v1.Action_Param
	for k, v := range param {
		param, err := pih.GetActionParamWithName(actionname, k)
		if err != nil {
			return nil, err
		}
		ap := new(v1.Action_Param)
		ap.ParamId = param.Id
		ap.Value, err = encode(v, param.Bitwidth)
		if err != nil {
			return nil, err
		}
		aps = append(aps, ap)
	}
	return aps, nil
}

type TableEntryOptions func(*v1.TableEntry, *P4InfoHelper) error

func Priority(priority int32) TableEntryOptions {
	return func(te *v1.TableEntry, pih *P4InfoHelper) error {
		te.Priority = priority
		return nil
	}
}

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
	Value1 int32
}

func (*EntryMatchFieldExact) isEntryMatchFieldType() {}

func (*EntryMatchFieldLpm) isEntryMatchFieldType() {}

func DeFaultAction(dact bool) TableEntryOptions {
	return func(te *v1.TableEntry, pih *P4InfoHelper) error {
		te.IsDefaultAction = dact
		return nil
	}
}

func Action(actionname string, param map[string]interface{}) TableEntryOptions {
	return func(te *v1.TableEntry, pih *P4InfoHelper) error {
		action := new(v1.Action)
		aid, err := pih.GetID(ActionEntity, actionname)
		if err != nil {
			return err
		}
		action.ActionId = aid
		action.Params, err = pih.GetActionParamPb(actionname, param)
		if err != nil {
			return err
		}
		te.Action.Type = &v1.TableAction_Action{Action: action}
		return nil
	}
}

// 引数の値の与え方は要検討
// とりあえずオプションにしているけど仕様書を見ながら変える
func (pih *P4InfoHelper) BuildTableEntry(
	tablename string,
	iemf *IsEntryMatchField,
	options ...TableEntryOptions,
) (*v1.TableEntry, error) {
	id, err := pih.GetID(TableEntity, tablename)
	if err != nil {
		return nil, err
	}
	te := new(v1.TableEntry)
	te.TableId = id

	if iemf != nil {
		mf, err := pih.GetMatchFieldPb(tablename, iemf)
		if err != nil {
			return nil, err
		}
		te.Match = []*v1.FieldMatch{mf}
	}

	//execute setting Option argument when exist Option Parameter in table entry struct
	for _, o := range options {
		err = o(te, pih)
		if err != nil {
			return nil, err
		}
	}
	return te, nil
}

func (pih *P4InfoHelper) BuidlMulticastGroupEntry() {
}

func (pih *P4InfoHelper) BuildCloneSessionEntry() {

}
