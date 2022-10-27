// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

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

type parserState struct {
	subcommand         *SubCommand
	subCommandForName  map[string]*SubCommand
	optionForLongName  map[string]*Option
	optionForShortName map[string]*Option
	hasSubCommands     bool
	hadSubCommand      bool
	args               []string
	index              int
}

// Returns the current arg and increments the index to point at the next
func (me *parserState) next() string {
	if me.index < len(me.args) {
		arg := me.args[me.index]
		me.index++
		return arg
	}
	return ""
}
