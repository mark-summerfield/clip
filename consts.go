// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

const mainSubCommand = ""
const noShortName = 0

type Validator func(string) error

type optionState uint8

const (
	NotGiven optionState = iota
	Given
	HadValue
)

func (me optionState) String() string {
	switch me {
	case NotGiven:
		return "NotGiven"
	case Given:
		return "Given"
	case HadValue:
		return "HadValue"
	default:
		panic("invalid optionState")
	}
}

type ValueCount uint8

const (
	OneOrMoreValues ValueCount = iota
	TwoValues
	ThreeValues
	FourValues
)

func (me ValueCount) String() string {
	switch me {
	case OneOrMoreValues:
		return "one or more"
	case TwoValues:
		return "two"
	case ThreeValues:
		return "three"
	case FourValues:
		return "four"
	default:
		panic("#310: invalid ValueCount")
	}
}

type PositionalCount uint8

const (
	ZeroPositionals PositionalCount = iota
	ZeroOrOnePositionals
	ZeroOrMorePositionals
	OnePositional
	TwoPositionals
	ThreePositionals
	FourPositionals
	OneOrMorePositionals
)

func (me PositionalCount) String() string {
	switch me {
	case ZeroPositionals:
		return "no"
	case ZeroOrOnePositionals:
		return "zero or one"
	case ZeroOrMorePositionals:
		return "zero or more"
	case OnePositional:
		return "one"
	case TwoPositionals:
		return "two"
	case ThreePositionals:
		return "three"
	case FourPositionals:
		return "four"
	case OneOrMorePositionals:
		return "one or more"
	default:
		panic("#320: invalid PositionalCount")
	}
}
