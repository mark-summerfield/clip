// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package parg

import (
	"errors"
	"fmt"
	"github/mark-summerfield/gong"
	"os"
)

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

func (p *Parser) Parse() error {
	return p.ParseArgs(os.Args[1:])
}

func (p *Parser) ParseLine(line string) error {
	return p.ParseArgs(gong.SplitWhitespace(line))
}

func (p *Parser) ParseArgs(args []string) error {
	var err error
	// TODO
	if err != nil && p.QuitOnError {
		fmt.Println(err)
		os.Exit(2)
	}
	return err
}

func (p *Parser) AddBool(option Option) {
	option.ValueType = Bool
	p.Options[option.LongName] = option
}

func (p *Parser) AddInt(option Option) {
	option.ValueType = Int
	p.Options[option.LongName] = option
}

func (p *Parser) AddReal(option Option) {
	option.ValueType = Real
	p.Options[option.LongName] = option
}

func (p *Parser) AddStr(option Option) {
	option.ValueType = Str
	p.Options[option.LongName] = option
}

func (p *Parser) AddStrs(option Option) {
	option.ValueType = Strs
	p.Options[option.LongName] = option
}

func (p *Parser) GetBool(name string) (bool, error) {
	opt, ok := p.Options[name]
	if !ok {
		p.handleError(fmt.Sprintf(
			"GetBool() is for bools, requested for %s", opt.ValueType))
	}
	v, ok := opt.Value.(bool)
	if !ok {
		p.handleError(fmt.Sprintf("GetBool() %s has an invalid bool %t",
			name, v))
	}
	return v, nil
}

func (p *Parser) GetInt(name string) (int, error) {
	opt, ok := p.Options[name]
	if !ok {
		p.handleError(fmt.Sprintf(
			"GetInt() is for int, requested for %s", opt.ValueType))
	}
	v, ok := opt.Value.(int)
	if !ok {
		p.handleError(fmt.Sprintf("GetInt() %s has an invalid int %d",
			name, v))
	}
	return v, nil
}

func (p *Parser) GetReal(name string) (float64, error) {
	opt, ok := p.Options[name]
	if !ok {
		p.handleError(fmt.Sprintf(
			"GetReal() is for reals, requested for %s", opt.ValueType))
	}
	v, ok := opt.Value.(float64)
	if !ok {
		p.handleError(fmt.Sprintf("GetReal() %s has an invalid bool %f",
			name, v))
	}
	return v, nil
}

func (p *Parser) GetStr(name string) (string, error) {
	opt, ok := p.Options[name]
	if !ok {
		p.handleError(fmt.Sprintf(
			"GetStr() is for strings, requested for %s", opt.ValueType))
	}
	v, ok := opt.Value.(string)
	if !ok {
		p.handleError(fmt.Sprintf("GetStr() %s has an invalid str %s",
			name, v))
	}
	return v, nil
}

func (p *Parser) GetStrs(name string) ([]string, error) {
	opt, ok := p.Options[name]
	if !ok {
		p.handleError(fmt.Sprintf(
			"GetStrs() is for slices of strings, requested for %s",
			opt.ValueType))
	}
	v, ok := opt.Value.([]string)
	if !ok {
		p.handleError(fmt.Sprintf("GetStrs() %s has an invalid bool %v",
			name, v))
	}
	return v, nil
}

func (p *Parser) handleError(msg string) error {
	msg = fmt.Sprintf("error: %s", msg)
	if p.QuitOnError {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(2)
	}
	return errors.New(msg)
}

type Option struct {
	LongName        string
	ShortName       string
	Help            string
	Required        bool
	ValueCount      ValueCount
	HasDefaultValue bool
	DefaultValue    any
	Value           any
	ValueType       ValueType
}

func (opt Option) String() string {
	return fmt.Sprintf("-%s|--%s req=%t vc=%v hasd=%t def=%v val=%v vt=%s",
		opt.ShortName, opt.LongName, opt.Required, opt.ValueCount,
		opt.HasDefaultValue, opt.DefaultValue, opt.Value, opt.ValueType)
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
