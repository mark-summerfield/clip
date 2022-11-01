// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
	"regexp"
	"strings"
)

type optioner interface {
	LongName() string
	ShortName() rune
	SetShortName(rune)
	SetVarName(string) error
	addValue(string) string
	wantsValue() bool
	setGiven()
	check() string
}

type commonOption struct {
	longName  string
	shortName rune
	help      string
	varName   string // e.g., -o|--outfile FILENAME
	state     optionState
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

func (me *commonOption) SetVarName(name string) error {
	if name == "" {
		return fmt.Errorf("#%d: can't have an empty varname", eEmptyVarName)
	}
	me.varName = name
	return nil
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

func newFlagOption(name, help string) (*FlagOption, error) {
	name, err := validatedName(name)
	if err != nil {
		return nil, err
	}
	shortName, longName := namesForName(name)
	return &FlagOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: NotGiven}}, nil
}

func (me FlagOption) Value() bool {
	return me.value
}

func (me FlagOption) wantsValue() bool {
	return false
}

func (me FlagOption) check() string {
	if me.state == HadValue {
		return fmt.Sprintf("#%d:BUG: a flag with a value", eBug)
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
	validator     IntValidator
}

func newIntOption(name, help string, theDefault int) (*IntOption, error) {
	name, err := validatedName(name)
	if err != nil {
		return nil, err
	}
	shortName, longName := namesForName(name)
	return &IntOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: NotGiven},
		theDefault: theDefault, validator: makeDefaultIntValidator()}, nil
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

func (me *IntOption) SetDefault(defaultValue int) {
	me.theDefault = defaultValue
}

func (me *IntOption) SetValidator(validator IntValidator) {
	me.validator = validator
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
	i, msg := me.validator(me.longName, value)
	if msg != "" {
		return msg
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
	validator     RealValidator
}

func newRealOption(name, help string, theDefault float64) (*RealOption,
	error) {
	name, err := validatedName(name)
	if err != nil {
		return nil, err
	}
	shortName, longName := namesForName(name)
	return &RealOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: NotGiven},
		theDefault: theDefault, validator: makeDefaultRealValidator()}, nil
}

func (me RealOption) Value() float64 {
	if me.state == HadValue {
		return me.value
	}
	return me.theDefault
}

func (me *RealOption) SetDefault(defaultValue float64) {
	me.theDefault = defaultValue
}

func (me *RealOption) SetValidator(validator RealValidator) {
	me.validator = validator
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
	r, msg := me.validator(me.longName, value)
	if msg != "" {
		return msg
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
	validator     StrValidator
}

func newStrOption(name, help, theDefault string) (*StrOption, error) {
	name, err := validatedName(name)
	if err != nil {
		return nil, err
	}
	shortName, longName := namesForName(name)
	return &StrOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: NotGiven},
		theDefault: theDefault, validator: makeDefaultStrValidator()}, nil
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

func (me *StrOption) SetDefault(defaultValue string) {
	me.theDefault = defaultValue
}

func (me *StrOption) SetValidator(validator StrValidator) {
	me.validator = validator
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
	s, msg := me.validator(me.longName, value)
	if msg != "" {
		return msg
	}
	me.value = s
	me.state = HadValue
	return ""
}

type StrsOption struct {
	*commonOption
	value      []string
	valueCount ValueCount
	validator  StrValidator
}

func newStrsOption(name, help string) (*StrsOption, error) {
	name, err := validatedName(name)
	if err != nil {
		return nil, err
	}
	shortName, longName := namesForName(name)
	return &StrsOption{commonOption: &commonOption{longName: longName,
			shortName: shortName, help: help, state: NotGiven},
			valueCount: OneOrMoreValues, validator: makeDefaultStrValidator()},
		nil
}

func (me StrsOption) Value() []string {
	return me.value
}

func (me StrsOption) wantsValue() bool {
	return me.state != NotGiven
}

func (me *StrsOption) SetValidator(validator StrValidator) {
	me.validator = validator
}

func (me *StrsOption) SetValueCount(valueCount ValueCount) {
	me.valueCount = valueCount
}

func (me StrsOption) check() string {
	if me.state == Given {
		return fmt.Sprintf(
			"expected %s values for %s, got none", me.valueCount,
			me.LongName())
	} else if me.state == HadValue {
		count := len(me.value)
		ok := true
		switch me.valueCount {
		case OneOrMoreValues:
			if count < 1 {
				ok = false
			}
		case TwoValues:
			if count != 2 {
				ok = false
			}
		case ThreeValues:
			if count != 3 {
				ok = false
			}
		case FourValues:
			if count != 3 {
				ok = false
			}
		default:
			return fmt.Sprintf("#%d:BUG:impossible ValueCount", eBug)
		}
		if !ok {
			return fmt.Sprintf(
				"expected %s values for %s, got %d", me.valueCount,
				me.LongName(), count)
		}
	}
	return ""
}

func (me *StrsOption) addValue(value string) string {
	s, msg := me.validator(me.longName, value)
	if msg != "" {
		return msg
	}
	if me.value == nil {
		me.value = make([]string, 0, 1)
	}
	me.value = append(me.value, s)
	me.state = HadValue
	return ""
}

func validatedName(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("#%d: can't have an empty option name",
			eEmptyOptionName)
	}
	if strings.HasPrefix(name, "-") {
		name = strings.Trim(name, "-")
	}
	if matched, _ := regexp.MatchString(`^\d+`, name); matched {
		return "", fmt.Errorf(
			"#%d: can't have a numeric option name, got %s",
			eNumericOptionName, name)
	}
	return name, nil
}
