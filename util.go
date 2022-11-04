// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package garg

import (
	"fmt"
	tsize "github.com/kopoli/go-terminal-size"
	"github.com/mark-summerfield/gong"
	"strconv"
	"strings"
	"unicode/utf8"
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
			return 0, fmt.Sprintf(
				"option %s's value of %s isn't an int: %s",
				name, value, err)
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
			return 0, fmt.Sprintf(
				"option %s's value of %s isn't a real: %s",
				name, value, err)
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
		colon := ""
		end := fmt.Sprintf("the %d valid choices", len(choices))
		if len(choices) < 9 {
			colon = ":"
			end = strings.Join(choices, " ")
		}
		return "", fmt.Sprintf("option %s's value of %q is not one of%s %s",
			name, value, colon, end)
	}
}

func positionalCountText(count PositionalCount, varName string) string {
	switch count {
	case ZeroPositionals:
		return ""
	case ZeroOrOnePositionals:
		return fmt.Sprintf("[%s]", varName)
	case ZeroOrMorePositionals: // any count is valid
		return fmt.Sprintf("[%s [%s ...]]", varName, varName)
	case OnePositional:
		return fmt.Sprintf("<%s>", varName)
	case OneOrMorePositionals:
		return fmt.Sprintf("<%s> [%s [%s ...]]", varName, varName, varName)
	case TwoPositionals:
		return fmt.Sprintf("<%s> <%s>", varName, varName)
	case ThreePositionals:
		return fmt.Sprintf("<%s> <%s> <%s>", varName, varName, varName)
	case FourPositionals:
		return fmt.Sprintf("<%s> <%s> <%s> <%s>", varName, varName, varName,
			varName)
	}
	panic("BUG: missing PositionalCount case")
}

func valueCountText(count ValueCount, varName string) string {
	switch count {
	case OneOrMoreValues:
		return fmt.Sprintf("<%s> [%s ...]", varName, varName)
	case TwoValues:
		return fmt.Sprintf("<%s> <%s>", varName, varName)
	case ThreeValues:
		return fmt.Sprintf("<%s> <%s> <%s>", varName, varName, varName)
	case FourValues:
		return fmt.Sprintf("<%s> <%s> <%s> <%s>", varName, varName, varName,
			varName)
	}
	panic("BUG: missing ValueCount case")
}

func argHelp(argWidth, width int, desc string) string {
	text := ""
	gapWidth := utf8.RuneCountInString(columnGap)
	argWidth += gapWidth
	descLen := utf8.RuneCountInString(desc)
	if argWidth+gapWidth+descLen <= width { // desc fits
		text += desc
	} else {
		indent := strings.Repeat(columnGap, 4)
		theWidth := width - utf8.RuneCountInString(indent)
		desc := gong.TextWrapIndent(desc, theWidth, indent)
		text += "\n" + strings.Join(desc, "\n")
	}
	if text[len(text)-1] != '\n' {
		text += "\n"
	}
	return text
}

func getWidth() int {
	size, err := tsize.GetSize()
	if err == nil && size.Width >= 38 {
		return size.Width
	}
	return 80
}
