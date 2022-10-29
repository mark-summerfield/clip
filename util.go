// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"errors"
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

func makeIntRangeValidator(minimum, maximum int) func(string) error {
	return func(arg string) error {
		i, err := strconv.Atoi(arg)
		if err != nil {
			return err
		}
		if minimum <= i && i <= maximum {
			return nil
		}
		if i < minimum {
			return fmt.Errorf("%d less than the minimum of %d ", i, minimum)
		}
		return fmt.Errorf("%d greater than the maximum of %d ", i, maximum)
	}
}

func makeRealRangeValidator(minimum, maximum float64) func(string) error {
	return func(arg string) error {
		r, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return err
		}
		if minimum <= r && r <= maximum {
			return nil
		}
		if r < minimum {
			return fmt.Errorf("%g less than the minimum of %g ", r, minimum)
		}
		return fmt.Errorf("%g greater than the maximum of %g ", r, maximum)
	}
}

func makeChoiceValidator(choices []string) func(string) error {
	return func(arg string) error {
		for _, choice := range choices {
			if arg == choice {
				return nil
			}
		}
		return errors.New("not one of the valid choices")
	}
}
