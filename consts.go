// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

const MainSubCommand = ""
const NoShortName = 0

// A validator should return whether the given value is acceptable
type Validator func(any) bool

// TODO provide default function makers for use as validators

type Number interface {
	int | float64
}

type ValueType uint8

const (
	Bool ValueType = iota
	Int
	Real
	Str
	Strs
)

func (me ValueType) String() string {
	switch me {
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
