// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Optioner interface {
	LongName() string
	ShortName() rune
	SetShortName(rune)
	SetVarName(string)
	SetValidator(Validator)
	setGiven()
	Given() bool
	Count() int
	ValueCount() ValueCount
	hasDefault() bool
	defaultValue() any
	beenAdded() bool
}

type commonOption struct {
	longName   string
	shortName  rune
	help       string
	given      bool
	added      bool
	valueCount ValueCount
	varName    string // e.g., -o|--outfile FILENAME
	validator  Validator
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
	return me.given
}

func (me *commonOption) setGiven() {
	me.given = true
}

type FlagOption struct {
	*commonOption
	value bool
}

func newFlagOption(name, help string) *FlagOption {
	validateName(name)
	shortName, longName := namesForName(name)
	return &FlagOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, valueCount: Zero}}
}

func (me FlagOption) Value() bool {
	return me.value
}

func (me FlagOption) defaultValue() any {
	return nil
}

func (me FlagOption) hasDefault() bool {
	return true
}

func (me FlagOption) Count() int {
	return 0
}

func (me FlagOption) beenAdded() bool {
	return false
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
		shortName: shortName, help: help, valueCount: One},
		theDefault: theDefault}
}

func (me IntOption) Value() int {
	if !me.added {
		return me.theDefault
	}
	return me.value
}

func (me *IntOption) AllowImplicit() {
	me.allowImplicit = true
}

func (me IntOption) defaultValue() any {
	return me.theDefault
}

func (me IntOption) hasDefault() bool {
	return true
}

func (me IntOption) Count() int {
	if me.added {
		return 1
	}
	return 0
}

func (me IntOption) beenAdded() bool {
	return me.added
}

func (me *IntOption) addValue(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("option %s expected an int value, got %s",
			me.longName, value)
	}
	me.value = i
	me.added = true
	return nil
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
		shortName: shortName, help: help, valueCount: One},
		theDefault: theDefault}
}

func (me RealOption) Value() float64 {
	if !me.added {
		return me.theDefault
	}
	return me.value
}

func (me *RealOption) AllowImplicit() {
	me.allowImplicit = true
}

func (me RealOption) defaultValue() any {
	return me.theDefault
}

func (me RealOption) hasDefault() bool {
	return true
}

func (me RealOption) Count() int {
	if me.added {
		return 1
	}
	return 0
}

func (me RealOption) beenAdded() bool {
	return me.added
}

func (me *RealOption) addValue(value string) error {
	r, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("option %s expected a real value, got %s",
			me.longName, value)
	}
	me.value = r
	me.added = true
	return nil
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
		shortName: shortName, help: help, valueCount: One},
		theDefault: theDefault}
}

func (me StrOption) Value() string {
	if !me.added {
		return me.theDefault
	}
	return me.value
}

func (me *StrOption) AllowImplicit() {
	me.allowImplicit = true
}

func (me StrOption) defaultValue() any {
	return me.theDefault
}

func (me StrOption) hasDefault() bool {
	return true
}

func (me StrOption) Count() int {
	if me.added {
		return 1
	}
	return 0
}

func (me StrOption) beenAdded() bool {
	return me.added
}

func (me *StrOption) addValue(value string) error {
	me.value = value
	me.added = true
	return nil
}

type StrsOption struct {
	*commonOption
	value []string
}

func newStrsOption(name, help string) *StrsOption {
	validateName(name)
	shortName, longName := namesForName(name)
	return &StrsOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, valueCount: OneOrMore}}
}

func (me StrsOption) Value() []string {
	return me.value
}

func (me StrsOption) defaultValue() any {
	return nil
}

func (me StrsOption) hasDefault() bool {
	return false
}

func (me StrsOption) Count() int {
	return len(me.value)
}

func (me StrsOption) beenAdded() bool {
	return false
}

func (me *StrsOption) addValue(value string) error {
	if me.value == nil {
		me.value = make([]string, 0, 1)
	}
	me.value = append(me.value, value)
	me.added = true
	return nil
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
