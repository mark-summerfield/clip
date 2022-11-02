// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

func realEqual(x, y float64) bool {
	return math.Abs(x-y) < 0.0001
}

type multi interface {
	int | float64 | string
}

func expectEqualSlice[T multi](expected, actuals []T, what string) string {
	if !reflect.DeepEqual(actuals, expected) {
		return fmt.Sprintf("expected %s=%v, got %v", what, expected,
			actuals)
	}
	return ""
}

func expectEmptySlice[T multi](slice []T, what string) string {
	if slice != nil {
		return fmt.Sprintf("expected %s=nil, got %v", what, slice)
	}
	return ""
}

func testingExitFunc(exitCode int, msg string) {
	panic(fmt.Errorf("exit=%d msg=%q", exitCode, msg))
}

func expectPanic(code int, t *testing.T) {
	exitFunc = defaultExitFunc // restore original
	perr := recover()
	if perr == nil {
		t.Errorf("expected panic with code #%d", code)
	} else {
		if err, ok := perr.(error); ok {
			rx := regexp.MustCompile(`#(\d+):`)
			matches := rx.FindStringSubmatch(err.Error())
			if len(matches) < 2 || matches[1] != strconv.Itoa(code) {
				t.Errorf("expected error #%d, got %s", code, err)
			}
		} else {
			t.Errorf("internal error got %v", perr)
		}
	}
}

func createTestParser1(t *testing.T) (Parser, *FlagOption, *IntOption,
	*IntOption) {
	parser := NewParser()
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	verboseOpt := parser.Int("verbose", "verbosity -v or -vN", 1)
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
	return parser, summaryOpt, verboseOpt, maxWidthOpt
}

func createTestParser2(t *testing.T) (Parser, *FlagOption, *IntOption,
	*IntOption, *StrsOption, *StrsOption, *StrsOption, *StrsOption,
	*FlagOption) {
	parser, summaryOpt, verboseOpt, maxWidthOpt := createTestParser1(t)
	languageOpt := parser.Strs("language", "lang help")
	skipLanguageOpt := parser.Strs("skiplanguage", "skip lang help")
	skipLanguageOpt.SetShortName('L')
	excludeOpt := parser.Strs("exclude", "exclude help")
	includeOpt := parser.Strs("include", "include help")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	return parser, summaryOpt, verboseOpt, maxWidthOpt, languageOpt,
		skipLanguageOpt, excludeOpt, includeOpt, sortByLinesOpt
}

func Test001(t *testing.T) {
	parser := NewParserUser("garg.test", "1.0.0")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	line := ""
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if summaryOpt.Value() {
		t.Error("expected false, got true")
	}
	if parser.AppName() != "garg.test" {
		t.Errorf("expected appname=garg.test, got %s", parser.AppName())
	}
}

func Test002(t *testing.T) {
	parser := NewParserUser("myapp", "1.0.0")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	line := "-S"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
}

func Test003(t *testing.T) {
	parser := NewParserUser("myapp", "1.0.0")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	line := "--summary"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
}

func Test004(t *testing.T) {
	parser := NewParserUser("", "")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	verboseOpt := parser.Int("verbose", "verbosity -v or -vN", 1)
	verboseOpt.SetShortName(0)
	line := "-S --verbose 3 filename.txt"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	v := verboseOpt.Value() // Single valued options always have a default
	if v != 3 {
		t.Errorf("expected verbose=3, got %d", v)
	}
	if e := expectEqualSlice([]string{"filename.txt"}, parser.Positionals,
		"positionals"); e != "" {
		t.Error(e)
	}
}

func Test005(t *testing.T) {
	parser := NewParserUser("", "0.1.0")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	line := "-sS"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.Value() {
		t.Error("expected sortbylines=true, got false")
	}
	if parser.AppName() != "garg.test" {
		t.Errorf("expected appname=myapp, got %s", parser.AppName())
	}
	if parser.Version() != "0.1.0" {
		t.Errorf("expected version=0.1.0, got %s", parser.Version())
	}
}

func Test006(t *testing.T) {
	parser, summaryOpt, _, maxWidthOpt := createTestParser1(t)
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	line := "-sSm60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.Value() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.Value()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test007(t *testing.T) {
	parser, summaryOpt, _, maxWidthOpt := createTestParser1(t)
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	line := "-sSm=60 file1.txt file2.dat file3.uxf"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.Value() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.Value()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
	if e := expectEqualSlice([]string{"file1.txt", "file2.dat",
		"file3.uxf"}, parser.Positionals, "positionals"); e != "" {
		t.Error(e)
	}
}

func Test008(t *testing.T) {
	parser, summaryOpt, _, maxWidthOpt := createTestParser1(t)
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	line := "-sSm 60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.Value() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.Value()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test009(t *testing.T) {
	parser, summaryOpt, _, maxWidthOpt := createTestParser1(t)
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	languageOpt := parser.Strs("language", "lang help")
	line := "-sSm=60 -l cpp pas xml -- file1.txt file2.dat file3.uxf"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.Value() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.Value()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
	if e := expectEqualSlice([]string{"cpp", "pas", "xml"},
		languageOpt.Value(), "language"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"file1.txt", "file2.dat",
		"file3.uxf"}, parser.Positionals, "positionals"); e != "" {
		t.Error(e)
	}
}

func Test010(t *testing.T) {
	parser := NewParserVersion("1.2.3")
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
	line := "-m60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	m := maxWidthOpt.Value()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test011(t *testing.T) {
	parser := NewParserVersion("1.2.4")
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
	line := "-m=60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	m := maxWidthOpt.Value()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test012(t *testing.T) {
	parser := NewParserUser("myapp", "1.0.0")
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
	line := "--maxwidth 60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	m := maxWidthOpt.Value()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test013(t *testing.T) {
	parser := NewParserUser("myapp", "1.0.0")
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
	line := "--maxwidth=60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	m := maxWidthOpt.Value()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test014(t *testing.T) {
	parser := NewParserUser("myapp", "1.0.0")
	languageOpt := parser.Strs("language", "lang help")
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
	line := "-lcpp -- file1.txt file2.dat file3.uxf"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	m := maxWidthOpt.Value()
	if m != 80 {
		t.Errorf("expected maxwidth=80, got %d", m)
	}
	if e := expectEqualSlice([]string{"cpp"}, languageOpt.Value(),
		"language"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"file1.txt", "file2.dat",
		"file3.uxf"}, parser.Positionals, "positionals"); e != "" {
		t.Error(e)
	}
}

func Test015(t *testing.T) {
	parser := NewParser()
	languageOpt := parser.Strs("language", "lang help")
	line := "-l=cpp -- file1.txt file2.dat file3.uxf path/to/file4.xml"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEqualSlice([]string{"cpp"}, languageOpt.Value(),
		"language"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"file1.txt", "file2.dat", "file3.uxf",
		"path/to/file4.xml"}, parser.Positionals, "positionals"); e != "" {
		t.Error(e)
	}
}

func Test016(t *testing.T) {
	parser := NewParser()
	languageOpt := parser.Strs("language", "lang help")
	line := "-l cpp -- file1.txt file2.dat"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEqualSlice([]string{"cpp"}, languageOpt.Value(),
		"language"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"file1.txt", "file2.dat"},
		parser.Positionals, "positionals"); e != "" {
		t.Error(e)
	}
}

func Test017(t *testing.T) {
	parser := NewParser()
	languageOpt := parser.Strs("language", "lang help")
	line := "-l cpp pas xml -- file1.txt"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEqualSlice([]string{"cpp", "pas", "xml"},
		languageOpt.Value(), "language"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"file1.txt"}, parser.Positionals,
		"positionals"); e != "" {
		t.Error(e)
	}
}

func Test018(t *testing.T) {
	parser, summaryOpt, _, maxWidthOpt, languageOpt, _, excludeOpt,
		includeOpt, sortByLinesOpt := createTestParser2(t)
	maxWidthOpt.TheDefault = 56
	specialSubCommand := parser.SubCommand("special", "Special help")
	extraOpt := specialSubCommand.Flag("extra", "extra help")
	line := "-sS -l h red -e zOld t -L d -i peek -- file1.cpp file2.d"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEqualSlice([]string{"h", "red"}, languageOpt.Value(),
		"language"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"zOld", "t"}, excludeOpt.Value(),
		"exclude"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"peek"}, includeOpt.Value(),
		"include"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"file1.cpp", "file2.d"},
		parser.Positionals, "positionals"); e != "" {
		t.Error(e)
	}
	if !summaryOpt.Value() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.Value() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.Value()
	if m != 56 {
		t.Errorf("expected maxwidth=56, got %d", m)
	}
	if extraOpt.Value() {
		t.Error("expected extra=false, got true")
	}
}

func Test019(t *testing.T) {
	parser, summaryOpt, _, _, _, skipLanguageOpt, excludeOpt, includeOpt,
		sortByLinesOpt := createTestParser2(t)
	specialSubCommand := parser.SubCommand("special", "Special help")
	extraOpt := specialSubCommand.Flag("extra", "extra help")
	maxWidthOpt := specialSubCommand.IntInRange("maxwidth",
		"special max width help", 20, 10000, 80)
	line := "special -e -m98 file1.cpp file2.d"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEmptySlice(skipLanguageOpt.Value(),
		"skiplanguage"); e != "" {
		t.Error(e)
	}
	if e := expectEmptySlice(excludeOpt.Value(), "exclude"); e != "" {
		t.Error(e)
	}
	if e := expectEmptySlice(includeOpt.Value(), "include"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"file1.cpp", "file2.d"},
		parser.Positionals, "positionals"); e != "" {
		t.Error(e)
	}
	if summaryOpt.Value() {
		t.Error("expected summary=false, got true")
	}
	if sortByLinesOpt.Value() {
		t.Error("expected sortbylines=false, got true")
	}
	if !extraOpt.Value() {
		t.Error("expected extra=true, got false")
	}
	m := maxWidthOpt.Value()
	if m != 98 {
		t.Errorf("expected maxwidth=98, got %d", m)
	}
}

func Test020(t *testing.T) {
	parser, summaryOpt, _, maxWidthOpt, languageOpt, skipLanguageOpt,
		excludeOpt, includeOpt, sortByLinesOpt := createTestParser2(t)
	line := "-S -l h red -e zOld t -L d -i peek -m 40 file1.cpp file2.d"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEqualSlice([]string{"h", "red"}, languageOpt.Value(),
		"language"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"d"}, skipLanguageOpt.Value(),
		"skiplanguage"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"zOld", "t"}, excludeOpt.Value(),
		"exclude"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"peek"}, includeOpt.Value(),
		"include"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"file1.cpp", "file2.d"},
		parser.Positionals, "positionals"); e != "" {
		t.Error(e)
	}
	if !summaryOpt.Value() {
		t.Error("expected summary=true, got false")
	}
	if sortByLinesOpt.Value() {
		t.Error("expected sortbylines=false, got true")
	}
	m := maxWidthOpt.Value()
	if m != 40 {
		t.Errorf("expected maxwidth=40, got %d", m)
	}
}

func Test021(t *testing.T) {
	parser, summaryOpt, _, maxWidthOpt, languageOpt, skipLanguageOpt,
		excludeOpt, includeOpt, sortByLinesOpt := createTestParser2(t)
	line := "-l h red -e zOld t -L d -i peek -s file1.cpp file2.d"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEqualSlice([]string{"h", "red"}, languageOpt.Value(),
		"language"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"d"}, skipLanguageOpt.Value(),
		"skiplanguage"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"zOld", "t"}, excludeOpt.Value(),
		"exclude"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"peek"}, includeOpt.Value(),
		"include"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"file1.cpp", "file2.d"},
		parser.Positionals, "positionals"); e != "" {
		t.Error(e)
	}
	if summaryOpt.Value() {
		t.Error("expected summary=false, got true")
	}
	if !sortByLinesOpt.Value() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.Value()
	if m != 80 {
		t.Errorf("expected maxwidth=80, got %d", m)
	}
}

func Test022(t *testing.T) {
	parser, summaryOpt, _, maxWidthOpt, languageOpt, skipLanguageOpt,
		excludeOpt, includeOpt, sortByLinesOpt := createTestParser2(t)
	line := "-l go h red -e zOld t test -s -L f77 asm -i peek unz"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEqualSlice([]string{"go", "h", "red"},
		languageOpt.Value(), "language"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"f77", "asm"},
		skipLanguageOpt.Value(), "skiplanguage"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"zOld", "t", "test"},
		excludeOpt.Value(), "exclude"); e != "" {
		t.Error(e)
	}
	if e := expectEqualSlice([]string{"peek", "unz"}, includeOpt.Value(),
		"include"); e != "" {
		t.Error(e)
	}
	if parser.Positionals != nil {
		t.Errorf("expected positionals=nil, got %s", parser.Positionals)
	}
	if summaryOpt.Value() {
		t.Error("expected summary=false, got true")
	}
	if !sortByLinesOpt.Value() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.Value()
	if m != 80 {
		t.Errorf("expected maxwidth=80, got %d", m)
	}
}

func Test023(t *testing.T) {
	parser := NewParserUser("myapp", "1.0.0")
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
	line := "-m60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	m := maxWidthOpt.Value()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test024(t *testing.T) {
	parser := NewParserUser("myapp", "1.0.0")
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 45)
	line := "-m70"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	m := maxWidthOpt.Value()
	if m != 70 {
		t.Errorf("expected maxwidth=70, got %d", m)
	}
}

func Test025(t *testing.T) {
	parser, summaryOpt, _, maxWidthOpt := createTestParser1(t)
	line := "--maxwidth=25"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if summaryOpt.Value() {
		t.Error("expected false, got true")
	}
	m := maxWidthOpt.Value()
	if m != 25 {
		t.Errorf("expected maxwidth=45, got %d", m)
	}
}

func Test026(t *testing.T) {
	parser, summaryOpt, _, maxWidthOpt := createTestParser1(t)
	line := "--maxwidth=99 -S"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	m := maxWidthOpt.Value()
	if m != 99 {
		t.Errorf("expected maxwidth=45, got %d", m)
	}
}

func Test027(t *testing.T) {
	parser, summaryOpt, verboseOpt, _ := createTestParser1(t)
	line := "-S --verbose=3"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	v := verboseOpt.Value() // Single valued options always have a default
	if v != 3 {
		t.Errorf("expected verbose=3, got %d", v)
	}
}

func Test028(t *testing.T) {
	parser, summaryOpt, _, maxWidthOpt := createTestParser1(t)
	maxWidthOpt.TheDefault = 45
	line := ""
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if summaryOpt.Value() {
		t.Error("expected false, got true")
	}
	m := maxWidthOpt.Value()
	if m != 45 {
		t.Errorf("expected maxwidth=45, got %d", m)
	}
}

func Test029(t *testing.T) {
	parser, summaryOpt, _, _ := createTestParser1(t)
	scaleOpt := parser.Real("scale", "max width help", 4.5)
	line := ""
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if summaryOpt.Value() {
		t.Error("expected false, got true")
	}
	s := scaleOpt.Value()
	if !realEqual(4.5, s) {
		t.Errorf("expected scale=4.5, got %f", s)
	}
}

func Test030(t *testing.T) {
	parser, summaryOpt, _, _ := createTestParser1(t)
	scaleOpt := parser.Real("scale", "max width help", 4.5)
	line := "-Ss-3.9"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	s := scaleOpt.Value()
	if !realEqual(-3.9, s) {
		t.Errorf("expected scale=4.5, got %f", s)
	}
}

func Test031(t *testing.T) {
	parser, summaryOpt, _, _ := createTestParser1(t)
	scaleOpt := parser.Real("scale", "max width help", 4.5)
	line := "-Ss3.5"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	s := scaleOpt.Value()
	if !realEqual(3.5, s) {
		t.Errorf("expected scale=3.5, got %f", s)
	}
}

func Test032(t *testing.T) {
	parser, summaryOpt, _, _ := createTestParser1(t)
	scaleOpt := parser.Real("scale", "max width help", 4.5)
	line := "-Ss 3.5"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	s := scaleOpt.Value()
	if !realEqual(3.5, s) {
		t.Errorf("expected scale=3.5, got %f", s)
	}
}

func Test033(t *testing.T) {
	parser, summaryOpt, _, _ := createTestParser1(t)
	scaleOpt := parser.Real("scale", "max width help", 4.5)
	line := "-Ss -2.5"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	s := scaleOpt.Value()
	if !realEqual(-2.5, s) {
		t.Errorf("expected scale=-2.5, got %f", s)
	}
}

func Test034(t *testing.T) {
	parser, summaryOpt, _, _ := createTestParser1(t)
	scaleOpt := parser.Real("scale", "max width help", 4.5)
	line := "-Ss-1.5"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	s := scaleOpt.Value()
	if !realEqual(-1.5, s) {
		t.Errorf("expected scale=-1.5, got %f", s)
	}
}

func Test035(t *testing.T) {
	parser, summaryOpt, _, _ := createTestParser1(t)
	scaleOpt := parser.Real("scale", "max width help", 4.5)
	line := "-Ss88"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	s := scaleOpt.Value()
	if !realEqual(88, s) {
		t.Errorf("expected scale=88, got %f", s)
	}
}

func Test036(t *testing.T) {
	parser, summaryOpt, verboseOpt, _ := createTestParser1(t)
	verboseOpt.AllowImplicit = true
	line := "-Sv"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	v := verboseOpt.Value()
	if v != 1 {
		t.Errorf("expected verbose=1, got %d", v)
	}
}

func Test037(t *testing.T) {
	parser, summaryOpt, verboseOpt, _ := createTestParser1(t)
	line := "-S"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	if verboseOpt.Given() {
		t.Error("expected verbose=!Given")
	}
}

func Test038(t *testing.T) {
	parser, summaryOpt, verboseOpt, _ := createTestParser1(t)
	verboseOpt.SetShortName(0)
	line := "-S"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	if verboseOpt.Given() {
		t.Error("expected verbose=!Given")
	}
	v := verboseOpt.Value() // Single valued options always have a default
	if v != 1 {
		t.Errorf("expected verbose=1, got %d", v)
	}
}

func Test039(t *testing.T) {
	parser, summaryOpt, verboseOpt, _ := createTestParser1(t)
	line := "-Sv2"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.Value() {
		t.Error("expected true, got false")
	}
	if verboseOpt.Given() {
		v := verboseOpt.Value()
		if v != 2 {
			t.Errorf("expected verbose=1, got %d", v)
		}
	} else {
		t.Error("expected verbose=Given")
	}
}

func Test040(t *testing.T) {
	version := "0.8.7-alpha"
	parser := NewParserUser("myapp", version)
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	parser.Int("verbose", "verbosity -v or -vN", 1)
	if v := parser.Version(); v != version {
		t.Errorf("Parser.Version is wrong, expected %q, got %q", version, v)
	}
	// TODO if I can capture output
	// _ := parser.ParseLine("-V")

	parser = NewParserUser("myapp", version)
	summaryOpt = parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	parser.Int("verbose", "verbosity -v or -vN", 1)
	if v := parser.Version(); v != version {
		t.Errorf("Parser.Version is wrong, expected %q, got %q", version, v)
	}
	// TODO if I can capture output
	// _ := parser.ParseLine("--version")
}

func Test041(t *testing.T) {
	parser, summaryOpt, verboseOpt, maxWidthOpt := createTestParser1(t)
	maxWidthOpt.TheDefault = 93
	line := "file1.txt file2.dat README.md"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if summaryOpt.Value() {
		t.Error("expected false got true")
	}
	if verboseOpt.Given() {
		t.Error("got unexpected verbose flag")
	}
	v := verboseOpt.Value()
	if v != 1 {
		t.Errorf("expected default verbose=1, got %d", v)
	}
	if e := expectEqualSlice([]string{"file1.txt", "file2.dat",
		"README.md"}, parser.Positionals, "positionals"); e != "" {
		t.Error(e)
	}
}

func Test042(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, maxWidthOpt, _, _, _, _, _ := createTestParser2(t)
	maxWidthOpt.TheDefault = 93
	if parser.PositionalCount != ZeroOrMorePositionals {
		t.Errorf("expected PositionalCount=%s, got %s",
			ZeroOrMorePositionals, parser.PositionalCount)
	}
	parser.PositionalCount = ZeroPositionals
	line := "-m20"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func Test043(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, maxWidthOpt, _, _, _, _, _ := createTestParser2(t)
	maxWidthOpt.TheDefault = 93
	if parser.PositionalCount != ZeroOrMorePositionals {
		t.Errorf("expected PositionalCount=%s, got %s",
			ZeroOrMorePositionals, parser.PositionalCount)
	}
	parser.PositionalCount = ZeroOrOnePositionals
	line := "-m20"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func Test044(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, maxWidthOpt, _, _, _, _, _ := createTestParser2(t)
	maxWidthOpt.TheDefault = 93
	if parser.PositionalCount != ZeroOrMorePositionals {
		t.Errorf("expected PositionalCount=%s, got %s",
			ZeroOrMorePositionals, parser.PositionalCount)
	}
	parser.PositionalCount = ZeroOrOnePositionals
	line := "-m20 file1.txt"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEqualSlice([]string{"file1.txt"}, parser.Positionals,
		"positionals"); e != "" {
		t.Error(e)
	}
}

func Test045(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, maxWidthOpt, _, _, _, _, _ := createTestParser2(t)
	maxWidthOpt.TheDefault = 93
	if parser.PositionalCount != ZeroOrMorePositionals {
		t.Errorf("expected PositionalCount=%s, got %s",
			ZeroOrMorePositionals, parser.PositionalCount)
	}
	parser.PositionalCount = TwoPositionals
	line := "-m20 a.dat beta.zip"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEqualSlice([]string{"a.dat", "beta.zip"},
		parser.Positionals, "positionals"); e != "" {
		t.Error(e)
	}
}

func Test046(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, _, _, _, _, includeOpt, _ := createTestParser2(t)
	if includeOpt.ValueCount != OneOrMoreValues {
		t.Errorf("expected ValueCount=%s, got %s",
			OneOrMoreValues, includeOpt.ValueCount)
	}
	line := ""
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEmptySlice(includeOpt.Value(), "include"); e != "" {
		t.Error(e)
	}
}

func Test047(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, _, _, _, _, includeOpt, _ := createTestParser2(t)
	if includeOpt.ValueCount != OneOrMoreValues {
		t.Errorf("expected ValueCount=%s, got %s",
			OneOrMoreValues, includeOpt.ValueCount)
	}
	includeOpt.ValueCount = ThreeValues
	line := "-i a bee ceee"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEqualSlice([]string{"a", "bee", "ceee"},
		includeOpt.Value(), "--include"); e != "" {
		t.Error(e)
	}
}

func Test048(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, _, _, _, _, includeOpt, _ := createTestParser2(t)
	if includeOpt.ValueCount != OneOrMoreValues {
		t.Errorf("expected ValueCount=%s, got %s",
			OneOrMoreValues, includeOpt.ValueCount)
	}
	includeOpt.ValueCount = ThreeValues
	line := ""
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEmptySlice(includeOpt.Value(), "include"); e != "" {
		t.Error(e)
	}
}

func Test049(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, verboseOpt, _ := createTestParser1(t)
	line := "-v5"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	v := verboseOpt.Value()
	if v != 5 {
		t.Errorf("expected verbose=5, got %d", v)
	}
}

func Test050(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, verboseOpt, _ := createTestParser1(t)
	verboseOpt.AllowImplicit = true
	line := "-v5"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	v := verboseOpt.Value()
	if v != 5 {
		t.Errorf("expected verbose=5, got %d", v)
	}
}

func Test051(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, verboseOpt, _ := createTestParser1(t)
	verboseOpt.AllowImplicit = true
	verboseOpt.TheDefault = 7
	line := "-v"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	v := verboseOpt.Value()
	if v != 7 {
		t.Errorf("expected default verbose=7, got %d", v)
	}
}

func Test052(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	sizeOpt := parser.RealInRange("size", "size help", -1, 1, 0.5)
	line := ""
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	v := sizeOpt.Value()
	if !realEqual(0.5, v) {
		t.Errorf("expected default size=0.5, got %g", v)
	}
}

func Test053(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	sizeOpt := parser.RealInRange("size", "size help", -1, 1, 0.5)
	line := "--size=-0.19"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	v := sizeOpt.Value()
	if !realEqual(-0.19, v) {
		t.Errorf("expected size=-0.19, got %g", v)
	}
}

func Test054(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	currencyOpt := parser.Choice("currency", "currency help",
		[]string{"USD", "GBP", "EUR"}, "GBP")
	line := ""
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	v := currencyOpt.Value()
	if v != "GBP" {
		t.Errorf("expected default currency=GBP, got %s", v)
	}
}

func Test055(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	currencyOpt := parser.Choice("currency", "currency help",
		[]string{"USD", "GBP", "EUR"}, "GBP")
	line := "-c EUR"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	v := currencyOpt.Value()
	if v != "EUR" {
		t.Errorf("expected currency=EUR, got %s", v)
	}
}

func Test056(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	marginsOpt := parser.Ints("margins", "")
	line := ""
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEmptySlice(marginsOpt.Value(), "margins"); e != "" {
		t.Error(e)
	}
}

func Test057(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	marginsOpt := parser.Ints("margins", "")
	line := "-m 88 0 19 27 42"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if e := expectEqualSlice([]int{88, 0, 19, 27, 42}, marginsOpt.Value(),
		"margins"); e != "" {
		t.Error(e)
	}
}

func Test058(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	marginsOpt := parser.Reals("margins", "")
	line := "-m 0.1 -17.5 63"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	v := marginsOpt.Value()
	if len(v) != 3 {
		t.Errorf("expected 3 float64s got %d: %v", len(v), v)
	}
	expected := []float64{0.1, -17.5, 63}
	for i, r := range v {
		if !realEqual(r, expected[i]) {
			t.Errorf("actual %g != %g expected", r, expected[i])
		}
	}
}

func Test059(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	aOpt := parser.Flag("a", "")
	bOpt := parser.Int("b2", "", 5)
	cOpt := parser.Real("c_3", "", -7.2)
	dOpt := parser.Str("D", "", "x1")
	line := ""
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	a := aOpt.Value()
	if a {
		t.Error("expected a false got true")
	}
	b := bOpt.Value()
	if b != 5 {
		t.Errorf("expected b 5 got %d", b)
	}
	c := cOpt.Value()
	if !realEqual(c, -7.2) {
		t.Errorf("expected c -7.2 got %g", c)
	}
	d := dOpt.Value()
	if d != "x1" {
		t.Errorf("expected d x1 got %s", d)
	}
	if len(parser.Positionals) > 0 {
		t.Errorf("expected no positionals got %d", len(parser.Positionals))
	}
}

func TestPkgDoc001(t *testing.T) {
	parser := NewParserUser("myapp", "1.0.0")
	verboseOpt := parser.Flag("verbose", "whether to show more output")
	if err := parser.ParseLine(""); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	verbose := verboseOpt.Value() // verbose == false
	if verbose {
		t.Error("expected verbose=false, got true")
	}
	verbose = verboseOpt.Given() // verbose == false
	if verbose {
		t.Error("expected verbose=false, got true")
	}

	parser = NewParserUser("myapp", "1.0.0")
	verboseOpt = parser.Flag("verbose", "whether to show more output")
	if err := parser.ParseLine("-v"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	verbose = verboseOpt.Value() // verbose == true
	if !verbose {
		t.Error("expected verbose=true, got false")
	}
	verbose = verboseOpt.Given() // verbose == true
	if !verbose {
		t.Error("expected verbose=true, got false")
	}
}

func createPackageDocParser(t *testing.T) Parser {
	parser := NewParser()
	parser.Flag("verbose", "whether to show more output")
	parser.Flag("xray", "")
	parser.Flag("cat", "")
	return parser
}

func TestPkgDoc002(t *testing.T) {
	parser := createPackageDocParser(t)
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine(""); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	if outfileOpt.Given() {
		t.Error("expected outfile=!Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "" {
		t.Errorf("expected outfile=\"\", got %q", outfile)
	}
}

func TestPkgDoc003(t *testing.T) {
	parser := createPackageDocParser(t)
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine("-v -x -c -o outfile.dat"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "outfile.dat" {
		t.Errorf("expected outfile=\"outfile.dat\", got %q", outfile)
	}
}

func TestPkgDoc004(t *testing.T) {
	parser := createPackageDocParser(t)
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine("-v -x -c -o=outfile.dat"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "outfile.dat" {
		t.Errorf("expected outfile=\"outfile.dat\", got %q", outfile)
	}
}

func TestPkgDoc005(t *testing.T) {
	parser := createPackageDocParser(t)
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine("-vxcooutfile.dat"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "outfile.dat" {
		t.Errorf("expected outfile=\"outfile.dat\", got %q", outfile)
	}
}

func TestPkgDoc006(t *testing.T) {
	parser := createPackageDocParser(t)
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine("-vxco outfile.dat"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "outfile.dat" {
		t.Errorf("expected outfile=\"outfile.dat\", got %q", outfile)
	}
}

func TestPkgDoc007(t *testing.T) {
	parser := createPackageDocParser(t)
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine("-vxco=outfile.dat"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "outfile.dat" {
		t.Errorf("expected outfile=\"outfile.dat\", got %q", outfile)
	}
}

func TestPkgDoc008(t *testing.T) {
	parser := createPackageDocParser(t)
	outfileOpt := parser.Str("outfile", "outfile", "test.dat")
	if err := parser.ParseLine("-vxc"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	if outfileOpt.Given() {
		t.Error("expected outfile=!Given")
	}
	outfile := outfileOpt.Value() // since not given we get the default
	if outfile != "test.dat" {
		t.Errorf("expected outfile=\"test.dat\", got %q", outfile)
	}
}

func TestPkgDoc009(t *testing.T) {
	parser := createPackageDocParser(t)
	outfileOpt := parser.Str("outfile", "outfile", "test.dat")
	outfileOpt.AllowImplicit = true
	if err := parser.ParseLine("-vxco"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value() // since given with no value
	if outfile != "test.dat" {    // we get the default
		t.Errorf("expected outfile=\"test.dat\", got %q", outfile)
	}
}

func TestPkgDoc010(t *testing.T) {
	parser := createPackageDocParser(t)
	outfileOpt := parser.Str("outfile", "outfile", "test.dat")
	if err := parser.ParseLine("-vxcooutfile.txt"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value() // since given with no value
	if outfile != "outfile.txt" { // we get the default
		t.Errorf("expected outfile=\"outfile.txt\", got %q", outfile)
	}
}

func TestPkgDoc011(t *testing.T) {
	parser := createPackageDocParser(t)
	outfileOpt := parser.Str("outfile", "outfile", "test.dat")
	outfileOpt.AllowImplicit = true
	if err := parser.ParseLine("-vxcooutfile.txt"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value() // since given with no value
	if outfile != "outfile.txt" { // we get the default
		t.Errorf("expected outfile=\"outfile.txt\", got %q", outfile)
	}
}

func TestPkgDoc012(t *testing.T) {
	parser := NewParserUser("myapp", "1.0.0")
	verboseOpt := parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine(""); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	verbose := 0
	expected := 1
	if verboseOpt.Given() {
		t.Error("expected verbose=!Given")
	}
	verbose = verboseOpt.Value() // default
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParserUser("myapp", "1.0.0")
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	verboseOpt.AllowImplicit = true
	if err := parser.ParseLine("-v"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	expected = 1
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParserUser("myapp", "1.0.0")
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	verboseOpt.AllowImplicit = true
	if err := parser.ParseLine("--verbose"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	expected = 1
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParserUser("myapp", "1.0.0")
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine("-v1"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	expected = 1
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParserUser("myapp", "1.0.0")
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine("-v=2"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	expected = 2
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParserUser("myapp", "1.0.0")
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine("-v 3"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	expected = 3
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParserUser("myapp", "1.0.0")
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine("--verbose=4"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	expected = 4
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParserUser("myapp", "1.0.0")
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine("--verbose 5"); err != nil {
		t.Errorf("expected successful parse, %s", err)
	}
	expected = 5
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}
}

func TestE001(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParserUser("myapp", "1.0.0")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	line := "-S4"
	e := eUnexpectedValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE002(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParserUser("myapp", "1.0.0")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	parser.IntInRange("maxwidth", "max width help", 20, 10000, 45)
	line := "--maxwidth -s"
	e := eUnexpectedValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE003(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParserUser("myapp", "1.0.0")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	parser.Int("verbose", "verbosity -v or -vN", 1)
	// -v expects either nothing (will use the default of 1) or an int
	line := "-vS"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE004(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParserUser("myapp", "1.0.0")
	parser.Int("maxwidth", "help", 43)
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	line := "-m -S"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE005(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParserUser("myapp", "1.0.0")
	parser.Int("maxwidth", "help", 44)
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	line := "--maxwidth -S"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE006(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParserUser("myapp", "1.0.0")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	parser.IntInRange("maxwidth", "max width help", 20, 10000, 45)
	line := "--maxwidth"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE007(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParserUser("myapp", "1.0.0")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	parser.Real("scale", "max width help", 4.5)
	line := "-Ss"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE008(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParserUser("myapp", "1.0.0")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	parser.IntInRange("maxwidth", "max width help", 20, 10000, 45)
	line := "--maxwidth 11"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE009(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, maxWidthOpt, _, _, _, _, _ := createTestParser2(t)
	maxWidthOpt.TheDefault = 93
	if parser.PositionalCount != ZeroOrMorePositionals {
		t.Errorf("expected PositionalCount=%s, got %s",
			ZeroOrMorePositionals, parser.PositionalCount)
	}
	parser.PositionalCount = ZeroPositionals
	line := "-m20 file1.txt file2.dat README.md"
	e := eWrongPositionalCount
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE010(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, maxWidthOpt, _, _, _, _, _ := createTestParser2(t)
	maxWidthOpt.TheDefault = 93
	if parser.PositionalCount != ZeroOrMorePositionals {
		t.Errorf("expected PositionalCount=%s, got %s",
			ZeroOrMorePositionals, parser.PositionalCount)
	}
	parser.PositionalCount = ZeroOrOnePositionals
	line := "-m20 file1.txt file2.dat README.md"
	e := eWrongPositionalCount
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE011(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, maxWidthOpt, _, _, _, _, _ := createTestParser2(t)
	maxWidthOpt.TheDefault = 93
	if parser.PositionalCount != ZeroOrMorePositionals {
		t.Errorf("expected PositionalCount=%s, got %s",
			ZeroOrMorePositionals, parser.PositionalCount)
	}
	parser.PositionalCount = TwoPositionals
	line := "-m20 file1.txt file2.dat README.md"
	e := eWrongPositionalCount
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE012(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, maxWidthOpt, _, _, _, _, _ := createTestParser2(t)
	maxWidthOpt.TheDefault = 93
	if parser.PositionalCount != ZeroOrMorePositionals {
		t.Errorf("expected PositionalCount=%s, got %s",
			ZeroOrMorePositionals, parser.PositionalCount)
	}
	parser.PositionalCount = TwoPositionals
	line := ""
	e := eWrongPositionalCount
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE013(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, maxWidthOpt, _, _, _, _, _ := createTestParser2(t)
	maxWidthOpt.TheDefault = 93
	if parser.PositionalCount != ZeroOrMorePositionals {
		t.Errorf("expected PositionalCount=%s, got %s",
			ZeroOrMorePositionals, parser.PositionalCount)
	}
	parser.PositionalCount = TwoPositionals
	line := "-m20 README.md"
	e := eWrongPositionalCount
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE014(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, _, _, _, _, includeOpt, _ := createTestParser2(t)
	if includeOpt.ValueCount != OneOrMoreValues {
		t.Errorf("expected ValueCount=%s, got %s",
			OneOrMoreValues, includeOpt.ValueCount)
	}
	includeOpt.ValueCount = ThreeValues
	line := "-i a"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE015(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, _, _, _, _, includeOpt, _ := createTestParser2(t)
	if includeOpt.ValueCount != OneOrMoreValues {
		t.Errorf("expected ValueCount=%s, got %s",
			OneOrMoreValues, includeOpt.ValueCount)
	}
	includeOpt.ValueCount = ThreeValues
	line := "--include x y"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE016(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, _ := createTestParser1(t)
	line := "-v"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE017(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, _ := createTestParser1(t)
	line := "-m9"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE018(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, _ := createTestParser1(t)
	line := "-m10001"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE019(t *testing.T) {
	exitFunc = testingExitFunc
	parser, _, _, _ := createTestParser1(t)
	line := "-m19"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE020(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	parser.RealInRange("size", "size help", -1, 1, 0.5)
	line := "-s-1.1"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE021(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	parser.RealInRange("size", "size help", -1, 1, 0.5)
	line := "-s1.01"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE022(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	parser.Choice("currency", "currency help", []string{"USD", "GBP",
		"EUR"}, "GBP")
	line := "-c OZY"
	e := eInvalidValue
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE023(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	parser.Choice("99", "currency help", []string{"USD", "GBP",
		"EUR"}, "GBP")
	line := ""
	e := eInvalidName
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}

func TestE024(t *testing.T) {
	exitFunc = testingExitFunc
	parser := NewParser()
	parser.SubCommand("", "bad")
	line := ""
	e := eInvalidName
	defer expectPanic(e, t)
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
}
