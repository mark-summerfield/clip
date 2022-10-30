// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type optioner interface {
	LongName() string
	ShortName() rune
	SetShortName(rune)
	SetVarName(string)
	SetValidator(Validator)
	addValue(string) string
	wantsValue() bool
	setGiven()
	check() string
}

type commonOption struct {
	longName   string
	shortName  rune
	help       string
	valueCount ValueCount
	varName    string // e.g., -o|--outfile FILENAME
	validator  Validator
	state      optionState
}

func (me *commonOption) LongName() string {
	return me.longName
}

func (me *commonOption) ShortName() rune {
	return me.shortName
}

func (me *commonOption) SetShortName(c rune) {
	me.shortName = c
}

func (me *commonOption) SetVarName(name string) {
	if name == "" {
		panic("#100: can't have an empty varname")
	}
	me.varName = name
}

func (me *commonOption) SetValidator(validator Validator) {
	me.validator = validator
}

func (me *commonOption) ValueCount() ValueCount {
	return me.valueCount
}

func (me *commonOption) Given() bool {
	return me.state != NotGiven
}

func (me *commonOption) setGiven() {
	if me.state == NotGiven {
		me.state = Given
	}
}

type FlagOption struct {
	*commonOption
	value bool
}

func newFlagOption(name, help string) *FlagOption {
	validateName(name)
	shortName, longName := namesForName(name)
	return &FlagOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, valueCount: Zero,
		state: NotGiven}}
}

func (me FlagOption) Value() bool {
	return me.value
}

func (me FlagOption) wantsValue() bool {
	return false
}

func (me FlagOption) check() string {
	if me.state == HadValue {
		panic("flag with value logic error")
	}
	return ""
}
func (me *FlagOption) addValue(value string) string {
	return fmt.Sprintf("flag %s can't accept a value", me.LongName())
}

type IntOption struct {
	*commonOption
	theDefault    int
	value         int
	allowImplicit bool
}

func newIntOption(name, help string, theDefault int) *IntOption {
	validateName(name)
	shortName, longName := namesForName(name)
	return &IntOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, valueCount: One, state: NotGiven},
		theDefault: theDefault}
}

func (me IntOption) Value() int {
	if me.state == HadValue {
		return me.value
	}
	return me.theDefault
}

func (me IntOption) wantsValue() bool {
	return me.state == Given
}

func (me IntOption) check() string {
	if me.state == Given {
		if me.allowImplicit {
			return ""
		} else {
			return fmt.Sprintf(
				"expected exactly one value for %s, got none",
				me.LongName())
		}
	}
	return ""
}

func (me *IntOption) AllowImplicit() {
	me.allowImplicit = true
}

func (me *IntOption) addValue(value string) string {
	i, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Sprintf("option %s expected an int value, got %s",
			me.longName, value)
	}
	me.value = i
	me.state = HadValue
	return ""
}

type RealOption struct {
	*commonOption
	theDefault    float64
	value         float64
	allowImplicit bool
}

func newRealOption(name, help string, theDefault float64) *RealOption {
	validateName(name)
	shortName, longName := namesForName(name)
	return &RealOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, valueCount: One, state: NotGiven},
		theDefault: theDefault}
}

func (me RealOption) Value() float64 {
	if me.state == HadValue {
		return me.value
	}
	return me.theDefault
}

func (me RealOption) wantsValue() bool {
	return me.state == Given
}

func (me RealOption) check() string {
	if me.state == Given {
		if me.allowImplicit {
			return ""
		} else {
			return fmt.Sprintf(
				"expected exactly one value for %s, got none",
				me.LongName())
		}
	}
	return ""
}

func (me *RealOption) AllowImplicit() {
	me.allowImplicit = true
}

func (me *RealOption) addValue(value string) string {
	r, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Sprintf("option %s expected a real value, got %s",
			me.longName, value)
	}
	me.value = r
	me.state = HadValue
	return ""
}

type StrOption struct {
	*commonOption
	theDefault    string
	value         string
	allowImplicit bool
}

func newStrOption(name, help, theDefault string) *StrOption {
	validateName(name)
	shortName, longName := namesForName(name)
	return &StrOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, valueCount: One, state: NotGiven},
		theDefault: theDefault}
}

func (me StrOption) Value() string {
	if me.state == HadValue {
		return me.value
	}
	return me.theDefault
}

func (me StrOption) wantsValue() bool {
	return me.state == Given
}

func (me StrOption) check() string {
	if me.state == Given {
		if me.allowImplicit {
			return ""
		} else {
			return fmt.Sprintf(
				"expected exactly one value for %s, got none",
				me.LongName())
		}
	}
	return ""
}

func (me *StrOption) AllowImplicit() {
	me.allowImplicit = true
}

func (me *StrOption) addValue(value string) string {
	me.value = value
	me.state = HadValue
	return ""
}

type StrsOption struct {
	*commonOption
	value []string
}

func newStrsOption(name, help string) *StrsOption {
	validateName(name)
	shortName, longName := namesForName(name)
	return &StrsOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, valueCount: OneOrMore,
		state: NotGiven}}
}

func (me StrsOption) Value() []string {
	return me.value
}

func (me StrsOption) wantsValue() bool {
	return me.state != NotGiven
}

func (me StrsOption) check() string {
	if me.state == Given {
		return fmt.Sprintf(
			"expected exactly at least one value for %s, got none",
			me.LongName())
	}
	return ""
}

func (me *StrsOption) addValue(value string) string {
	if me.value == nil {
		me.value = make([]string, 0, 1)
	}
	me.value = append(me.value, value)
	me.state = HadValue
	return ""
}

func validateName(name string) {
	if name == "" {
		panic("#140: can't have an unnamed option")
	}
	if strings.HasPrefix(name, "-") {
		panic(fmt.Sprintf(
			"#142: can't have an option name that begins with '-', got %s",
			name))
	}
	if matched, _ := regexp.MatchString(`^\d+`, name); matched {
		panic(fmt.Sprintf("#144: can't have a numeric option name, got %s",
			name))
	}
}
