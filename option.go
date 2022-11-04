// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package clop

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
	Help() string
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

func (me *commonOption) Help() string {
	return me.help
}

func (me *commonOption) VarName() string {
	if me.varName == "" {
		return strings.ToUpper(me.longName)
	}
	return me.varName
}

func (me *commonOption) SetVarName(name string) error {
	if err := checkName(name, "option var"); err != nil {
		return err
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

// Always returns a *FlagOption; _and_ either nil or error
func newFlagOption(name, help string) (*FlagOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &FlagOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: NotGiven}}, err
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
	TheDefault    int
	AllowImplicit bool
	Validator     IntValidator
	value         int
}

// Always returns a *IntOption; _and_ either nil or error
func newIntOption(name, help string, theDefault int) (*IntOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &IntOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: NotGiven},
		TheDefault: theDefault, Validator: makeDefaultIntValidator()}, err
}

func (me IntOption) Value() int {
	if me.state == HadValue {
		return me.value
	}
	return me.TheDefault
}

func (me IntOption) wantsValue() bool {
	return me.state == Given
}

func (me IntOption) check() string {
	if me.state == Given {
		if me.AllowImplicit {
			return ""
		} else {
			return fmt.Sprintf(
				"expected exactly one value for %s, got none",
				me.LongName())
		}
	}
	return ""
}

func (me *IntOption) addValue(value string) string {
	i, msg := me.Validator(me.longName, value)
	if msg != "" {
		return msg
	}
	me.value = i
	me.state = HadValue
	return ""
}

type RealOption struct {
	*commonOption
	TheDefault    float64
	AllowImplicit bool
	Validator     RealValidator
	value         float64
}

// Always returns a *RealOption; _and_ either nil or error
func newRealOption(name, help string, theDefault float64) (*RealOption,
	error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &RealOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: NotGiven},
		TheDefault: theDefault, Validator: makeDefaultRealValidator()}, err
}

func (me RealOption) Value() float64 {
	if me.state == HadValue {
		return me.value
	}
	return me.TheDefault
}

func (me RealOption) wantsValue() bool {
	return me.state == Given
}

func (me RealOption) check() string {
	if me.state == Given {
		if me.AllowImplicit {
			return ""
		} else {
			return fmt.Sprintf(
				"expected exactly one value for %s, got none",
				me.LongName())
		}
	}
	return ""
}

func (me *RealOption) addValue(value string) string {
	r, msg := me.Validator(me.longName, value)
	if msg != "" {
		return msg
	}
	me.value = r
	me.state = HadValue
	return ""
}

type StrOption struct {
	*commonOption
	TheDefault    string
	AllowImplicit bool
	Validator     StrValidator
	value         string
}

// Always returns a *StrOption; _and_ either nil or error
func newStrOption(name, help, theDefault string) (*StrOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &StrOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: NotGiven},
		TheDefault: theDefault, Validator: makeDefaultStrValidator()}, err
}

func (me StrOption) Value() string {
	if me.state == HadValue {
		return me.value
	}
	return me.TheDefault
}

func (me StrOption) wantsValue() bool {
	return me.state == Given
}

func (me StrOption) check() string {
	if me.state == Given {
		if me.AllowImplicit {
			return ""
		} else {
			return fmt.Sprintf(
				"expected exactly one value for %s, got none",
				me.LongName())
		}
	}
	return ""
}

func (me *StrOption) addValue(value string) string {
	s, msg := me.Validator(me.longName, value)
	if msg != "" {
		return msg
	}
	me.value = s
	me.state = HadValue
	return ""
}

type StrsOption struct {
	*commonOption
	ValueCount ValueCount
	Validator  StrValidator
	value      []string
}

// Always returns a *StrsOption; _and_ either nil or error
func newStrsOption(name, help string) (*StrsOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &StrsOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: NotGiven},
		ValueCount: OneOrMoreValues,
		Validator:  makeDefaultStrValidator()}, err
}

func (me StrsOption) Value() []string {
	return me.value
}

func (me StrsOption) wantsValue() bool {
	return me.state != NotGiven
}

func (me StrsOption) check() string {
	return checkMulti(me.LongName(), me.state, me.ValueCount, len(me.value))
}

func (me *StrsOption) addValue(value string) string {
	s, msg := me.Validator(me.longName, value)
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

type IntsOption struct {
	*commonOption
	ValueCount ValueCount
	Validator  IntValidator
	value      []int
}

// Always returns a *IntsOption; _and_ either nil or error
func newIntsOption(name, help string) (*IntsOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &IntsOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: NotGiven},
		ValueCount: OneOrMoreValues,
		Validator:  makeDefaultIntValidator()}, err
}

func (me IntsOption) Value() []int {
	return me.value
}

func (me IntsOption) wantsValue() bool {
	return me.state != NotGiven
}

func (me IntsOption) check() string {
	return checkMulti(me.LongName(), me.state, me.ValueCount, len(me.value))
}

func (me *IntsOption) addValue(value string) string {
	s, msg := me.Validator(me.longName, value)
	if msg != "" {
		return msg
	}
	if me.value == nil {
		me.value = make([]int, 0, 1)
	}
	me.value = append(me.value, s)
	me.state = HadValue
	return ""
}

type RealsOption struct {
	*commonOption
	ValueCount ValueCount
	Validator  RealValidator
	value      []float64
}

// Always returns a *RealsOption; _and_ either nil or error
func newRealsOption(name, help string) (*RealsOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &RealsOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: NotGiven},
		ValueCount: OneOrMoreValues,
		Validator:  makeDefaultRealValidator()}, err
}

func (me RealsOption) Value() []float64 {
	return me.value
}

func (me RealsOption) wantsValue() bool {
	return me.state != NotGiven
}

func (me RealsOption) check() string {
	return checkMulti(me.LongName(), me.state, me.ValueCount, len(me.value))
}

func (me *RealsOption) addValue(value string) string {
	s, msg := me.Validator(me.longName, value)
	if msg != "" {
		return msg
	}
	if me.value == nil {
		me.value = make([]float64, 0, 1)
	}
	me.value = append(me.value, s)
	me.state = HadValue
	return ""
}

func checkName(name, what string) error {
	rx := regexp.MustCompile(`^\pL[\pL\pNd_]*$`)
	if rx.MatchString(name) {
		return nil
	}
	return fmt.Errorf("#%d: expected identifier name for %s, got %s",
		eInvalidName, what, name)
}

func checkMulti(name string, state optionState, valueCount ValueCount,
	count int) string {
	if state == Given {
		return fmt.Sprintf(
			"expected %s values for %s, got none", valueCount, name)
	} else if state == HadValue {
		ok := true
		switch valueCount {
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
				"expected %s values for %s, got %d", valueCount, name,
				count)
		}
	}
	return ""
}
