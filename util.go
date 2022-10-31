// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
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

func makeDefaultIntValidator() func(string, string) (int, string) {
	return func(name, value string) (int, string) {
		i, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Sprintf("option %s expected an int value, got %s",
				name, value)
		}
		return i, ""
	}
}

func makeIntRangeValidator(minimum, maximum int) func(string, string) (int,
	string) {
	return func(name, value string) (int, string) {
		i, err := strconv.Atoi(value)
		if err != nil {
			return 0, err.Error()
		}
		if minimum <= i && i <= maximum {
			return i, ""
		}
		if i < minimum {
			return 0, fmt.Sprintf("option %s's minimum is %d, got %d",
				name, minimum, i)
		}
		return 0, fmt.Sprintf("option %s's maximum is %d, got %d",
			name, maximum, i)
	}
}

func makeDefaultRealValidator() func(string, string) (float64, string) {
	return func(name, value string) (float64, string) {
		r, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, fmt.Sprintf("option %s expected a real value, got %s",
				name, value)
		}
		return r, ""
	}
}

func makeRealRangeValidator(minimum, maximum float64) func(string,
	string) (float64, string) {
	return func(name, value string) (float64, string) {
		r, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, err.Error()
		}
		if minimum <= r && r <= maximum {
			return r, ""
		}
		if r < minimum {
			return 0, fmt.Sprintf("option %s's minimum is %g, got %g",
				name, minimum, r)
		}
		return 0, fmt.Sprintf("option %s's maximum is %g, got %g",
			name, maximum, r)
	}
}

func makeDefaultStrValidator() func(string, string) (string, string) {
	return func(name, value string) (string, string) {
		if value == "" {
			return "", fmt.Sprintf("option %s expected a nonempty string",
				name)
		}
		return value, ""
	}
}

func makeChoiceValidator(choices []string) func(string, string) (string,
	string) {
	return func(name, value string) (string, string) {
		for _, choice := range choices {
			if value == choice {
				return value, ""
			}
		}
		return "", fmt.Sprintf("option %s's value of %q is not one of "+
			"the valid choices", name, value)
	}
}
