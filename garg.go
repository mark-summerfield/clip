// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// A validator should return whether the given value is acceptable
type Validator func(any) bool

type Parser struct {
	AppName         string
	AppVersion      string
	QuitOnError     bool
	Options         map[string]Option
	PositionalCount ValueCount
	Positionals     []string
}

func NewParser(appname, version string) Parser {
	return Parser{AppName: appname, AppVersion: version,
		QuitOnError: true, Options: make(map[string]Option)}
}

func (me *Parser) Parse() error {
	return me.ParseArgs(os.Args[1:])
}

func (me *Parser) ParseLine(line string) error {
	return me.ParseArgs(strings.Fields(line))
}

func (me *Parser) ParseArgs(args []string) error {
	var err error
	// TODO
	if err != nil && me.QuitOnError {
		fmt.Println(err)
		os.Exit(2)
	}
	return err
}

func (me *Parser) AddBool(option MinOption) {
	opt := option.ToOption()
	opt.ValueCount = ZeroOrOne
	opt.ValueType = Bool
	me.Options[opt.LongName] = opt
}

func (me *Parser) AddBoolOpt(option Option) {
	option.ValueType = Bool
	me.Options[option.LongName] = option
}

// TODO func (me *Parser) AddInt(option MinOption)

func (me *Parser) AddIntOpt(option Option) {
	option.ValueType = Int
	me.Options[option.LongName] = option
}

// TODO func (me *Parser) AddReal(option MinOption)

func (me *Parser) AddRealOpt(option Option) {
	option.ValueType = Real
	me.Options[option.LongName] = option
}

// TODO func (me *Parser) AddChoice(name, help string, choices []string)
// should create a validator to check that the given value is one of the
// choices

// TODO func (me *Parser) AddStr(option MinOption)

func (me *Parser) AddStrOpt(option Option) {
	option.ValueType = Str
	me.Options[option.LongName] = option
}

func (me *Parser) AddStrs(option MinOption) {
	opt := option.ToOption()
	opt.ValueCount = ZeroOrMore
	opt.ValueType = Strs
	me.Options[opt.LongName] = opt
}

func (me *Parser) AddStrsOpt(option Option) {
	option.ValueType = Strs
	me.Options[option.LongName] = option
}

func (me *Parser) GetBool(name string) (bool, error) {
	opt := me.get(name)
	v, ok := opt.Value.(bool)
	if !ok {
		me.handleError(fmt.Sprintf("GetBool() %s has an invalid bool %t",
			name, v))
	}
	return v, nil
}

func (me *Parser) get(name string) Option {
	opt, ok := me.Options[name]
	if !ok {
		me.handleError(fmt.Sprintf("Get*() has no option called %s", name))
	}
	return opt
}

func (me *Parser) GetInt(name string) (int, error) {
	opt := me.get(name)
	v, ok := opt.Value.(int)
	if !ok {
		me.handleError(fmt.Sprintf("GetInt() %s has an invalid int %d",
			name, v))
	}
	return v, nil
}

func (me *Parser) GetReal(name string) (float64, error) {
	opt := me.get(name)
	v, ok := opt.Value.(float64)
	if !ok {
		me.handleError(fmt.Sprintf("GetReal() %s has an invalid bool %f",
			name, v))
	}
	return v, nil
}

func (me *Parser) GetStr(name string) (string, error) {
	opt := me.get(name)
	v, ok := opt.Value.(string)
	if !ok {
		me.handleError(fmt.Sprintf("GetStr() %s has an invalid str %s",
			name, v))
	}
	return v, nil
}

func (me *Parser) GetStrs(name string) ([]string, error) {
	opt := me.get(name)
	v, ok := opt.Value.([]string)
	if !ok {
		me.handleError(fmt.Sprintf("GetStrs() %s has an invalid bool %v",
			name, v))
	}
	return v, nil
}

func (me *Parser) handleError(msg string) error {
	msg = fmt.Sprintf("error: %s", msg)
	if me.QuitOnError {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(2)
	}
	return errors.New(msg)
}

type MinOption struct {
	Name         string
	Help         string
	ValueCount   ValueCount
	DefaultValue any
	Value        any
}

func (me MinOption) ToOption() Option {
	shortName, longName := namesForName(me.Name)
	hasDefaultValue := false
	if me.DefaultValue != nil {
		hasDefaultValue = true
	}
	return Option{
		LongName:        longName,
		ShortName:       shortName,
		Help:            me.Help,
		ValueCount:      me.ValueCount,
		HasDefaultValue: hasDefaultValue,
		DefaultValue:    me.DefaultValue,
	}
}

type Option struct {
	LongName        string
	ShortName       string
	Help            string
	Required        bool
	ValueCount      ValueCount
	VarName         string // e.g., -o|--outfile FILENAME
	HasDefaultValue bool
	DefaultValue    any
	Value           any
	ValueType       ValueType
	Validator       Validator
}

func (opt Option) String() string {
	return fmt.Sprintf("-%s|--%s req=%t vc=%v hasd=%t def=%v val=%v vt=%s",
		opt.ShortName, opt.LongName, opt.Required, opt.ValueCount,
		opt.HasDefaultValue, opt.DefaultValue, opt.Value, opt.ValueType)
}

func namesForName(name string) (string, string) {
	var shortName string
	for _, c := range name {
		shortName = string(c)
		break
	}
	return shortName, name
}

// TODO provide default function makers for use as validators

type Number interface {
	int | float64
}

// TODO see if this will work
func MakeRangeValidator[V Number](minimum, maximum V) func(V) bool {
	return func(x V) bool {
		return minimum <= x && x <= maximum
	}
}

type ValueType uint8

const (
	Bool ValueType = iota
	Int
	Real
	Str
	Strs
)

func (vt ValueType) String() string {
	switch vt {
	case Bool:
		return "bool"
	case Int:
		return "int"
	case Real:
		return "float64"
	case Str:
		return "string"
	case Strs:
		return "[]string"
	default:
		panic("invalid ValueType")
	}
}

type ValueCount uint8

const (
	Zero      ValueCount = iota // for flags; for no positionals allowed
	ZeroOrOne                   // i.e., optional
	ZeroOrMore
	One
	OneOrMore
)
