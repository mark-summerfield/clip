// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Option struct {
	longName     string
	shortName    rune
	help         string
	given        bool
	valueCount   ValueCount
	varName      string // e.g., -o|--outfile FILENAME
	defaultValue any    // not valid if ValueType == Strs
	value        any
	valueType    ValueType
	validator    Validator
}

// Can't change long name, help, or valueType after creation
func newOption(name, help string, valueType ValueType) *Option {
	if name == "" {
		panic("#100: can't have an unnamed option")
	}
	if strings.HasPrefix(name, "-") {
		panic(fmt.Sprintf(
			"#102: can't have an option name that begins with '-', got %s",
			name))
	}
	if matched, _ := regexp.MatchString(`^\d+`, name); matched {
		panic(fmt.Sprintf("#104: can't have a numeric option name, got %s",
			name))
	}
	shortName, longName := namesForName(name)
	return &Option{longName: longName, shortName: shortName, help: help,
		valueCount: One, valueType: valueType}
}

func (me *Option) SetShortName(c rune) {
	me.shortName = c
}

func (me *Option) SetVarName(name string) {
	if name == "" {
		panic("#106: can't have an empty varname")
	}
	me.varName = name
}

func (me *Option) SetDefault(defaultValue any) {
	if me.valueType == Flag {
		panic("#110: can't set a default value for a flag")
	}
	if me.valueType == Strs {
		panic("#112: can't set a default value for a string list")
	}
	me.valueCount = ZeroOrOne // if option given with no value, use default
	me.defaultValue = defaultValue
}

func (me *Option) SetValidator(validator Validator) {
	if me.valueType == Flag {
		panic("#120: can't set a validator for a flag")
	}
	me.validator = validator
}

func (me *Option) Given() bool {
	return me.given
}

func (me *Option) AsBool() bool {
	if me.valueType != Flag {
		panic(fmt.Sprintf("#130: AsBool() called on type %s", me.valueType))
	}
	if me.value == nil {
		panic(fmt.Sprintf("#132: AsBool() called on type %s with nil value",
			me.valueType))
	}
	return me.value.(bool)
}

func (me *Option) AsInt() int {
	if me.valueType != Int {
		panic(fmt.Sprintf("#140: AsInt() called on type %s", me.valueType))
	}
	if me.value == nil {
		panic(fmt.Sprintf("#142: AsInt() called on type %s with nil value",
			me.valueType))
	}
	return me.value.(int)
}

func (me *Option) AsReal() float64 {
	if me.valueType != Real {
		panic(fmt.Sprintf("#150: AsReal() called on type %s", me.valueType))
	}
	if me.value == nil {
		panic(fmt.Sprintf("#152: AsReal() called on type %s with nil value",
			me.valueType))
	}
	return me.value.(float64)
}

func (me *Option) AsStr() string {
	if me.valueType != Str {
		panic(fmt.Sprintf("#160: AsStr() called on type %s", me.valueType))
	}
	if me.value == nil {
		panic(fmt.Sprintf("#162: AsStr() called on type %s with nil value",
			me.valueType))
	}
	return me.value.(string)
}

func (me *Option) AsStrs() []string {
	if me.valueType != Strs {
		panic(fmt.Sprintf("#170: AsStrs() called on type %s", me.valueType))
	}
	if me.value == nil {
		return nil
	}
	return me.value.([]string)
}

func (me *Option) Count() int {
	if me.value == nil || me.valueType == Flag {
		return 0
	}
	if me.valueType == Strs {
		return len(me.value.([]string))
	}
	return 1
}

func (me *Option) addValue(value string) error {
	if me.validator != nil {
		if err := me.validator(value); err != nil {
			return fmt.Errorf("invalid value of %q for %s: %s", value,
				me.longName, err)
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
		panic("#180: invalid ValueType")
	}
	return nil
}

func (me *Option) setDefaultIfAppropriate() {
	if me.value == nil { // Flag false default is set in ParseArgs()
		if me.valueType == Strs {
			me.value = make([]string, 0)
		} else if me.defaultValue != nil {
			me.value = me.defaultValue
		}
	}
}
