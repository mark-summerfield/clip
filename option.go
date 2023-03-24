// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package clip

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
	MustSetVarName(string)
	Help() string
	Hide()
	isHidden() bool
	addValue(string) string
	wantsValue() bool
	setGiven()
	check() string
}

type commonOption struct {
	longName  string
	shortName rune
	help      string
	varName   string // e.g., -o|--outfile FILE
	hidden    bool
	state     optionState
}

// LongName returns the option's long name.
func (me *commonOption) LongName() string {
	return me.longName
}

// ShortName returns the option's short name which could be 0
// ([NoShortName]).
func (me *commonOption) ShortName() rune {
	return me.shortName
}

// SetShortName sets the option's short name—or clears it if [NoShortName]
// is passed.
func (me *commonOption) SetShortName(c rune) {
	me.shortName = c
}

// Help returns the option's help text.
func (me *commonOption) Help() string {
	return me.help
}

// Hide sets the option to be hidden: the user can use it normally, but it
// won't show up when -h or --help is given.
func (me *commonOption) Hide() {
	me.hidden = true
}

func (me *commonOption) isHidden() bool {
	return me.hidden
}

// VarName returns the name used for the option's variables: by default the
// option's long name uppercased. (This is never used by FlagOptions.)
func (me *commonOption) VarName() string {
	if me.varName == "" {
		return strings.ToUpper(me.longName)
	}
	return me.varName
}

// SetVarName is used to set the option's variable name. [See also
// MustSetVarName].
func (me *commonOption) SetVarName(name string) error {
	if err := checkName(name, "option var"); err != nil {
		return err
	}
	me.varName = name
	return nil
}

// MustSetVarName is used to set the option's variable name. Panics on
// error. See also [SetVarName]
func (me *commonOption) MustSetVarName(name string) {
	if err := me.SetVarName(name); err != nil {
		panic(err)
	}
}

// Given returns true if (after the parse) the option was given; otherwise
// returns false.
func (me *commonOption) Given() bool {
	return me.state != notGiven
}

func (me *commonOption) setGiven() {
	if me.state == notGiven {
		me.state = given
	}
}

// FlagOption is an option for a flag (i.e., an option that is either
// present or absent).
type FlagOption struct {
	*commonOption
	value bool
}

// Always returns a *FlagOption; _and_ either nil or error.
func newFlagOption(name, help string) (*FlagOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &FlagOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: notGiven}}, err
}

// Value returns true if the flag was given; otherwise false.
func (me FlagOption) Value() bool {
	return me.value
}

func (me FlagOption) wantsValue() bool {
	return false
}

func (me FlagOption) check() string {
	if me.state == hadValue {
		return fmt.Sprintf("#%d:BUG: a flag with a value", eBug)
	}
	return ""
}

func (me *FlagOption) addValue(value string) string {
	return "flag " + me.LongName() + " can't accept a value"
}

// IntOption is an option for accepting a single int.
type IntOption struct {
	*commonOption
	TheDefault    int          // The options default value.
	AllowImplicit bool         // If true, giving the option with no value means use the default.
	Validator     IntValidator // A validation function.
	value         int
}

// Always returns a *IntOption; _and_ either nil or error.
func newIntOption(name, help string, theDefault int) (*IntOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &IntOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: notGiven},
		TheDefault: theDefault, Validator: makeDefaultIntValidator()}, err
}

// Value returns the given value or if the option wasn't given, the default
// value.
func (me IntOption) Value() int {
	if me.state == hadValue {
		return me.value
	}
	return me.TheDefault
}

func (me IntOption) wantsValue() bool {
	return me.state == given
}

func (me IntOption) check() string {
	if me.state == given {
		if me.AllowImplicit {
			return ""
		} else {
			return "expected exactly one value for " + me.LongName() +
				", got none"
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
	me.state = hadValue
	return ""
}

// RealOption is an option for accepting a single real.
type RealOption struct {
	*commonOption
	TheDefault    float64       // The options default value.
	AllowImplicit bool          // If true, giving the option with no value means use the default.
	Validator     RealValidator // A validation function.
	value         float64
}

// Always returns a *RealOption; _and_ either nil or error.
func newRealOption(name, help string, theDefault float64) (*RealOption,
	error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &RealOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: notGiven},
		TheDefault: theDefault, Validator: makeDefaultRealValidator()}, err
}

// Value returns the given value or if the option wasn't given, the default
// value.
func (me RealOption) Value() float64 {
	if me.state == hadValue {
		return me.value
	}
	return me.TheDefault
}

func (me RealOption) wantsValue() bool {
	return me.state == given
}

func (me RealOption) check() string {
	if me.state == given {
		if me.AllowImplicit {
			return ""
		} else {
			return "expected exactly one value for " + me.LongName() +
				", got none"
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
	me.state = hadValue
	return ""
}

// StrOption is an option for accepting a single string.
type StrOption struct {
	*commonOption
	TheDefault    string       // The options default value.
	AllowImplicit bool         // If true, giving the option with no value means use the default.
	Validator     StrValidator // A validation function.
	value         string
}

// Always returns a *StrOption; _and_ either nil or error.
func newStrOption(name, help, theDefault string) (*StrOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &StrOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: notGiven},
		TheDefault: theDefault, Validator: makeDefaultStrValidator()}, err
}

// Value returns the given value or if the option wasn't given, the default
// value.
func (me StrOption) Value() string {
	if me.state == hadValue {
		return me.value
	}
	return me.TheDefault
}

func (me StrOption) wantsValue() bool {
	return me.state == given
}

func (me StrOption) check() string {
	if me.state == given {
		if me.AllowImplicit {
			return ""
		} else {
			return "expected exactly one value for " + me.LongName() +
				", got none"
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
	me.state = hadValue
	return ""
}

// StrsOption is an option for accepting a one or more strings.
type StrsOption struct {
	*commonOption
	ValueCount ValueCount   // How many strings are wanted.
	Validator  StrValidator // A validation function.
	value      []string
}

// Always returns a *StrsOption; _and_ either nil or error.
func newStrsOption(name, help string) (*StrsOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &StrsOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: notGiven},
		ValueCount: OneOrMoreValues,
		Validator:  makeDefaultStrValidator()}, err
}

// Value returns the given value(s) or nil.
func (me StrsOption) Value() []string {
	return me.value
}

func (me StrsOption) wantsValue() bool {
	return me.state != notGiven
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
	me.state = hadValue
	return ""
}

// IntsOption is an option for accepting a one or more ints.
type IntsOption struct {
	*commonOption
	ValueCount ValueCount   // How many ints are wanted.
	Validator  IntValidator // A validation function.
	value      []int
}

// Always returns a *IntsOption; _and_ either nil or error.
func newIntsOption(name, help string) (*IntsOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &IntsOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: notGiven},
		ValueCount: OneOrMoreValues,
		Validator:  makeDefaultIntValidator()}, err
}

// Value returns the given value(s) or nil.
func (me IntsOption) Value() []int {
	return me.value
}

func (me IntsOption) wantsValue() bool {
	return me.state != notGiven
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
	me.state = hadValue
	return ""
}

// RealsOption is an option for accepting a one or more reals.
type RealsOption struct {
	*commonOption
	ValueCount ValueCount    // How many strings are wanted.
	Validator  RealValidator // A validation function.
	value      []float64
}

// Always returns a *RealsOption; _and_ either nil or error.
func newRealsOption(name, help string) (*RealsOption, error) {
	err := checkName(name, "option")
	shortName, longName := namesForName(name)
	return &RealsOption{commonOption: &commonOption{longName: longName,
		shortName: shortName, help: help, state: notGiven},
		ValueCount: OneOrMoreValues,
		Validator:  makeDefaultRealValidator()}, err
}

// Value returns the given value(s) or nil.
func (me RealsOption) Value() []float64 {
	return me.value
}

func (me RealsOption) wantsValue() bool {
	return me.state != notGiven
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
	me.state = hadValue
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
	if state == given {
		return fmt.Sprintf(
			"expected %s values for %s, got none", valueCount, name)
	} else if state == hadValue {
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
