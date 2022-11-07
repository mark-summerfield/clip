// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package clip

const NoShortName = 0 // Use this for options that don't have short names
const columnGap = "  "

// These take an option's name and the given string value and return a
// valid value and "" or the type's zero value and an error message.
type IntValidator func(string, string) (int, string)
type RealValidator func(string, string) (float64, string)
type StrValidator func(string, string) (string, string)

type optionState uint8

const (
	notGiven optionState = iota
	given
	hadValue
)

func (me optionState) String() string {
	switch me {
	case notGiven:
		return "not given"
	case given:
		return "given"
	case hadValue:
		return "had value"
	default:
		return "BUG: invalid optionState"
	}
}

// This specifies how many value *must* be present—if the option is given at
// all. So even if the ValueCount is TwoValues, if the option isn't given
// the option's Value will be empty. But if it _is_ given, then either it
// will have exactly two values, or there will be a Parser error.
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
		return "BUG: invalid ValueCount"
	}
}

// This specifies how many positionals *must* be present.
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
		return "BUG: invalid PositionalCount"
	}
}

type datum struct {
	arg    string
	lenArg int
	help   string
}

const (
	eUser                   = iota + 100
	eMissing                // 101
	eInvalidValue           // 102
	eInvalidHelpOption      // 103
	eInvalidVersionOption   // 104
	eEmptyVarName           // 105
	eUnrecognizedOption     // 106
	eUnexpectedValue        // 107
	eWrongPositionalCount   // 108
	eInvalidName            // 109
	eEmptyPositionalVarName // 110
	eBug                    = 999
)
