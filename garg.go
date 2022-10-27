// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Parser struct {
	AppName           string
	AppVersion        string
	QuitOnError       bool
	SubCommands       map[string]*SubCommand
	PositionalCount   ValueCount
	PositionalVarName string
	Positionals       []string
}

func NewParser(appname, version string) Parser {
	subcommands := make(map[string]*SubCommand)
	subcommands[MainSubCommand] = newMainSubCommand()
	return Parser{AppName: appname, AppVersion: version,
		QuitOnError: true, SubCommands: subcommands}
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

func (me *Parser) SubCommand(name, help string) *SubCommand {
	subcommand := newSubCommand(name, help)
	me.SubCommands[name] = subcommand
	return subcommand
}

func (me *Parser) Bool(name, help string) *Option {
	option := newOption(name, help)
	option.ValueCount = ZeroOrOne
	option.ValueType = Bool
	me.SubCommands[MainSubCommand].Options = append(
		me.SubCommands[MainSubCommand].Options, option)
	return option
}

// TODO Int() Real() Str() Strs() Choice()

/*
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
*/

/*
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
*/

func (me *Parser) handleError(msg string) error {
	msg = fmt.Sprintf("error: %s", msg)
	if me.QuitOnError {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(2)
	}
	return errors.New(msg)
}

func (me *Parser) Debug() {
	fmt.Println("Parser")
	fmt.Printf("    %v v%v\n", me.AppName, me.AppVersion)
	fmt.Printf("    QuitOnError=%t\n", me.QuitOnError)
	for name, subcommand := range me.SubCommands {
		if name == MainSubCommand {
			name = "<MAIN>"
		}
		fmt.Printf("    SubCommand=%v\n", name)
		subcommand.Debug(8)
	}
	fmt.Printf("    PositionalCount=%s\n", me.PositionalCount)
	fmt.Printf("    PositionalVarName=%s\n", me.PositionalVarName)
	fmt.Printf("    Positionals=%v\n", me.Positionals)
}
