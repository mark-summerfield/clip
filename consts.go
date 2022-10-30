// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

const mainSubCommand = ""
const noShortName = 0

type Validator func(string) error

type ValueCount uint8

const (
	Zero       ValueCount = iota // for flags & for no positionals allowed
	ZeroOrOne                    // for options with default and positionals
	ZeroOrMore                   // default for positionals
	One                          // default for Int, Real, Str
	Two                          // for positionals
	Three                        // for positionals
	OneOrMore                    // default for Strs
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
