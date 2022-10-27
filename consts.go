// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import "fmt"

const mainSubCommand = ""
const noShortName = 0

type Validator func(string) bool

type ValueType uint8

const (
	Flag ValueType = iota
	Int
	Real
	Str
	Strs
)

func (me ValueType) String() string {
	switch me {
	case Flag:
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

func (me ValueCount) String() string {
	switch me {
	case Zero:
		return "Zero"
	case ZeroOrOne:
		return "ZeroOrOne"
	case ZeroOrMore:
		return "ZeroOrMore"
	case One:
		return "One"
	case OneOrMore:
		return "OneOrMore"
	default:
		panic("invalid ValueCount")
	}
}

type tokenState struct {
	subcommand         *SubCommand
	subCommandForName  map[string]*SubCommand
	optionForLongName  map[string]*Option
	optionForShortName map[string]*Option
	hasSubCommands     bool
	hadSubCommand      bool
}

type Token struct {
	text    string
	isValue bool
}

func (me Token) String() string {
	if me.isValue {
		return fmt.Sprintf("%#v", me.text)
	}
	if len(me.text) == 1 {
		return fmt.Sprintf("-%s", me.text)
	}
	return fmt.Sprintf("--%s", me.text)
}

func NewNameToken(text string) Token {
	return Token{text: text}
}

func NewValueToken(text string) Token {
	return Token{text: text, isValue: true}
}
