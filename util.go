// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"strconv"
)

func namesForName(name string) (rune, string) {
	var shortName rune
	for _, c := range name {
		shortName = c
		break
	}
	return shortName, name
}

func makeIntRangeValidator(minimum, maximum int) func(string) bool {
	return func(arg string) bool {
		i, err := strconv.Atoi(arg)
		if err != nil {
			return false
		}
		return minimum <= i && i <= maximum
	}
}

func makeRealRangeValidator(minimum, maximum float64) func(string) bool {
	return func(arg string) bool {
		r, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return false
		}
		return minimum <= r && r <= maximum
	}
}

func makeChoiceValidator(choices []string) func(string) bool {
	return func(arg string) bool {
		for _, choice := range choices {
			if arg == choice {
				return true
			}
		}
		return false
	}
}
