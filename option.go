// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
	"strconv"
)

type Option struct {
	longName     string
	shortName    rune
	help         string
	required     bool
	valueCount   ValueCount
	varName      string // e.g., -o|--outfile FILENAME
	defaultValue any    // not valid if ValueType == Strs
	value        any
	valueType    ValueType
	validator    Validator
}

// Can't change long name, help, or valueType after creation
func newOption(name, help string, valueType ValueType) *Option {
	shortName, longName := namesForName(name)
	return &Option{longName: longName, shortName: shortName, help: help,
		valueCount: One, valueType: valueType}
}

func (me *Option) SetShortName(c rune) {
	me.shortName = c
}

func (me *Option) SetRequired() {
	me.required = true
}

func (me *Option) SetVarName(name string) {
	me.varName = name
}

func (me *Option) SetDefaultValue(dv any) {
	if me.valueType == Flag {
		panic("can't set a default value for a flag")
	}
	if me.valueType == Strs {
		panic("can't set a default value for a string list")
	}
	me.defaultValue = dv
}

func (me *Option) SetValidator(vf Validator) {
	if me.valueType == Flag {
		panic("can't set a validator for a flag")
	}
	me.validator = vf
}

func (me *Option) AsBool() bool {
	if me.valueType != Flag {
		panic(fmt.Sprintf("AsBool() called on type %s", me.valueType))
	}
	if me.value == nil {
		panic(fmt.Sprintf("AsBool() called on type %s with nil value",
			me.valueType))
	}
	return me.value.(bool)
}

func (me *Option) AsInt() int {
	if me.valueType != Int {
		panic(fmt.Sprintf("AsInt() called on type %s", me.valueType))
	}
	if me.value == nil {
		panic(fmt.Sprintf("AsInt() called on type %s with nil value",
			me.valueType))
	}
	return me.value.(int)
}

func (me *Option) AsReal() float64 {
	if me.valueType != Real {
		panic(fmt.Sprintf("AsReal() called on type %s", me.valueType))
	}
	if me.value == nil {
		panic(fmt.Sprintf("AsReal() called on type %s with nil value",
			me.valueType))
	}
	return me.value.(float64)
}

func (me *Option) AsStr() string {
	if me.valueType != Str {
		panic(fmt.Sprintf("AsStr() called on type %s", me.valueType))
	}
	if me.value == nil {
		panic(fmt.Sprintf("AsStr() called on type %s with nil value",
			me.valueType))
	}
	return me.value.(string)
}

func (me *Option) AsStrs() []string {
	if me.valueType != Strs {
		panic(fmt.Sprintf("AsStrs() called on type %s", me.valueType))
	}
	if me.value == nil {
		panic(fmt.Sprintf("AsStrs() called on type %s with nil value",
			me.valueType))
	}
	return me.value.([]string)
}

func (me *Option) Size() int {
	if me.value == nil || me.valueType == Flag {
		return 0
	}
	if me.valueType == Strs {
		return len(me.value.([]string))
	}
	return 1
}

func (me *Option) AddValue(value string) error {
	if me.validator != nil {
		if !me.validator(value) {
			return fmt.Errorf("invalid value for %s: %s", me.longName,
				value)
		}
	}
	switch me.valueType {
	case Flag:
		return fmt.Errorf("flag %s got unexpected value %s", me.longName,
			value)
	case Int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("option %s expected an int value, got %s",
				me.longName, value)
		}
		me.value = i
	case Real:
		r, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("option %s expected a real value, got %s",
				me.longName, value)
		}
		me.value = r
	case Str:
		me.value = value
	case Strs:
		if me.value == nil {
			me.value = make([]string, 0, 1)
		}
		me.value = append(me.value.([]string), value)
	default:
		panic("invalid ValueType #2")
	}
	return nil
}

func (me *Option) setDefaultIfAppropriate() {
	if me.value == nil {
		if me.valueType == Flag {
			me.value = false
		} else if me.valueType == Strs {
			me.value = make([]string, 0)
		} else if me.defaultValue != nil {
			me.value = me.defaultValue
		}
	}
}
