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
		panic("#300: invalid ValueType")
	}
}

type ValueCount uint8

const (
	Zero      ValueCount = iota // for flags & for no positionals allowed
	ZeroOrOne                   // for options with default and positionals
	ZeroOrMore
	One
	Two   // for positionals
	Three // for positionals
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
	case Two:
		return "Two"
	case Three:
		return "Three"
	case OneOrMore:
		return "OneOrMore"
	default:
		panic("#310: invalid ValueCount")
	}
}
