// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package clip

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	tsize "github.com/kopoli/go-terminal-size"
	"github.com/mark-summerfield/uterm"
)

var (
	tty       bool
	onWindows bool
)

func init() {
	if info, _ := os.Stdout.Stat(); (info.Mode() & os.ModeCharDevice) != 0 {
		tty = true
	}
	onWindows = runtime.GOOS == "windows"
}

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
			return 0, fmt.Sprintf("option %s's value of %q isn't an int",
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
			return 0, fmt.Sprintf("option %s's value of %q isn't an int",
				name, value)
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
			return 0, fmt.Sprintf("option %s's value of %q isn't a real",
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
			return 0, fmt.Sprintf("option %s's value of %q isn't a real",
				name, value)
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
			return "", "option " + name + " expected a nonempty string"
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

func positionalCountText(count PositionalCount, varName1,
	varNameN string) string {
	n := 1
	if varNameN == "" {
		varNameN = varName1
		varName1 += "1"
		n = 2
	}
	switch count {
	case ZeroPositionals:
		return ""
	case ZeroOrOnePositionals:
		return "[" + varName1 + "]"
	case ZeroToTwoPositionals:
		return fmt.Sprintf("[%s [%s%d]]", varName1, varNameN, n)
	case ZeroOrMorePositionals: // any count is valid
		return fmt.Sprintf("[%s [%s%d ...]]", varName1, varNameN, n)
	case OnePositional:
		return "<" + varName1 + ">"
	case OneOrTwoPositionals:
		return fmt.Sprintf("<%s> [%s%d]", varName1, varNameN, n)
	case OneToThreePositionals:
		return fmt.Sprintf("<%s> [%s%d [%s%d]]", varName1, varNameN, n,
			varNameN, n+1)
	case OneOrMorePositionals:
		return fmt.Sprintf("<%s> [%s%d [%s%d ...]]", varName1, varNameN, n,
			varNameN, n+1)
	case TwoPositionals:
		return fmt.Sprintf("<%s> <%s%d>", varName1, varNameN, n)
	case TwoOrThreePositionals:
		return fmt.Sprintf("<%s> <%s%d> [%s%d]", varName1, varNameN, n,
			varNameN, n+1)
	case ThreePositionals:
		return fmt.Sprintf("<%s> <%s%d> <%s%d>", varName1, varNameN, n,
			varNameN, n+1)
	case FourPositionals:
		return fmt.Sprintf("<%s> <%s%d> <%s%d> <%s%d>", varName1, varNameN,
			n, varNameN, n+1, varNameN, n+2)
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

// ArgHelp is used internally by clip, but made public because it can be
// useful for implementing subcommands (see
// `eg/subcommands/subcommands.go`).
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
		text += "\n" + uterm.WrappedIndent(desc, theWidth, indent)
	}
	if text[len(text)-1] != '\n' {
		text += "\n"
	}
	return text
}

// GetWidth returns the terminal width; it is used internally by clip, but
// made public because it can be useful for implementing subcommands (see
// `eg/subcommands/subcommands.go`).
func GetWidth() int {
	size, err := tsize.GetSize()
	if err == nil && size.Width >= 38 {
		return size.Width
	}
	return 80
}

func initialArgText(option optioner) (string, string) {
	arg := "--" + option.LongName()
	displayArg := Strong(arg)
	if option.ShortName() != NoShortName {
		arg = fmt.Sprintf("%s-%c, %s", columnGap, option.ShortName(),
			arg)
		displayArg = columnGap + Strong("-"+string(option.ShortName())) +
			", " + displayArg
	} else {
		arg = fmt.Sprintf("%s    %s", columnGap, arg)
		displayArg = columnGap + "    " + displayArg
	}
	return arg, displayArg
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

func prepareOptionsData(maxLeft, gapWidth, width int, data []datum) bool {
	allFit := true
	for i := 0; i < len(data); i++ {
		datum := &data[i]
		if maxLeft+gapWidth+utf8.RuneCountInString(datum.help) > width {
			allFit = false
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

// Strong returns the given string contained within terminal escape codes to
// make it bold on linux and bold or colored on windows (providing os.Stdout
// is a TTY).
func Strong(s string) string {
	if tty {
		return uterm.Bold(s)
	}
	return s
}

// Bold returns bold text.
//
// Deprecated: Use [Strong] instead.
func Bold(s string) string {
	return Strong(s)
}

// Empth returns the given string contained within terminal escape codes to
// make it italic on linux and underlined on windows (providing os.Stdout is
// a TTY).
func Emph(s string) string {
	if tty {
		if onWindows {
			return uterm.Underline(s)
		}
		return uterm.Italic(s)
	}
	return s
}

// Hint returns the given string contained within terminal escape codes to
// make it underlined on linux and italic on windows (although I've never
// known italics to actually work on windows) (providing os.Stdout is a
// TTY).
func Hint(s string) string {
	if tty {
		if onWindows {
			return uterm.Italic(s)
		}
		return uterm.Underline(s)
	}
	return s
}
