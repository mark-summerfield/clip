// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

func namesForName(name string) (rune, string) {
	var shortName rune
	for _, c := range name {
		shortName = c
		break
	}
	return shortName, name
}

// TODO provide default function makers for use as validators

// TODO see if this will work
func MakeRangeValidator[V Number](minimum, maximum V) func(V) bool {
	return func(x V) bool {
		return minimum <= x && x <= maximum
	}
}
