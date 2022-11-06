// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package clip

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
		return fmt.Sprintf("[%s1 [%s2 ...]]", varName, varName)
	case OnePositional:
		return fmt.Sprintf("<%s>", varName)
	case OneOrMorePositionals:
		return fmt.Sprintf("<%s1> [%s2 [%s3 ...]]", varName, varName,
			varName)
	case TwoPositionals:
		return fmt.Sprintf("<%s1> <%s2>", varName, varName)
	case ThreePositionals:
		return fmt.Sprintf("<%s1> <%s2> <%s3>", varName, varName, varName)
	case FourPositionals:
		return fmt.Sprintf("<%s1> <%s2> <%s3> <%s4>", varName, varName,
			varName, varName)
	}
	panic("BUG: missing PositionalCount case")
}

func valueCountText(count ValueCount, varName string) string {
	switch count {
	case OneOrMoreValues:
		return fmt.Sprintf("<%s1> [%s2 ...]", varName, varName)
	case TwoValues:
		return fmt.Sprintf("<%s1> <%s2>", varName, varName)
	case ThreeValues:
		return fmt.Sprintf("<%s1> <%s2> <%s3>", varName, varName, varName)
	case FourValues:
		return fmt.Sprintf("<%s1> <%s2> <%s3> <%s4>", varName, varName,
			varName, varName)
	}
	panic("BUG: missing ValueCount case")
}

func ArgHelp(argWidth, width int, desc string) string {
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

func GetWidth() int {
	size, err := tsize.GetSize()
	if err == nil && size.Width >= 38 {
		return size.Width
	}
	return 80
}

func initialArgText(option optioner) (int, string) {
	short := 0
	arg := "--" + option.LongName()
	if option.ShortName() != NoShortName {
		arg = fmt.Sprintf("%s-%c, %s", columnGap, option.ShortName(),
			arg)
		short = 1
	} else {
		arg = columnGap + arg
	}
	return short, arg
}

func optArgText(option optioner) string {
	switch opt := option.(type) {
	case *IntOption:
		if opt.AllowImplicit {
			return " [" + opt.VarName() + "]"
		} else {
			return " " + opt.VarName()
		}
	case *RealOption:
		if opt.AllowImplicit {
			return " [" + opt.VarName() + "]"
		} else {
			return " " + opt.VarName()
		}
	case *StrOption:
		if opt.AllowImplicit {
			return " [" + opt.VarName() + "]"
		} else {
			return " " + opt.VarName()
		}
	case *IntsOption:
		return " " + valueCountText(opt.ValueCount, opt.VarName())
	case *RealsOption:
		return " " + valueCountText(opt.ValueCount, opt.VarName())
	case *StrsOption:
		return " " + valueCountText(opt.ValueCount, opt.VarName())
	}
	return ""
}

func prepareOptionsData(maxLeft, gapWidth, width, shorts int,
	data []datum) bool {
	padOnlyLong := false
	if shorts != 0 && shorts != len(data) { // some longs without shorts
		padOnlyLong = true
	}
	allFit := true
	for i := 0; i < len(data); i++ {
		datum := &data[i]
		if maxLeft+gapWidth+utf8.RuneCountInString(datum.help) > width {
			allFit = false
		}
		if padOnlyLong && strings.HasPrefix(datum.arg, columnGap+"--") {
			datum.arg = columnGap + strings.TrimSpace(datum.arg)
		}
	}
	return allFit
}

func optionsDataText(allFit bool, maxLeft, gapWidth, width int,
	data []datum) string {
	text := ""
	for _, datum := range data {
		text += datum.arg
		if datum.help != "" {
			if allFit {
				text += strings.Repeat(" ", maxLeft-datum.lenArg) +
					columnGap + datum.help + "\n"
			} else {
				if datum.lenArg+gapWidth+utf8.RuneCountInString(
					datum.help) > width && datum.lenArg < maxLeft {
					text += strings.Repeat(" ", maxLeft-datum.lenArg)
				}
				text += columnGap + ArgHelp(maxLeft, width, datum.help)
			}
		}
	}
	return text
}
