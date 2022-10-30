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

func expectEqualSlice(expected, actuals []string, what string) string {
	if !reflect.DeepEqual(actuals, expected) {
		return fmt.Sprintf("expected %s=%s, got %s", what, expected,
			actuals)
	}
	return ""
}

func expectEmptySlice(slice []string, what string) string {
	if slice != nil {
		return fmt.Sprintf("expected %s=nil, got %s", what, slice)
	}
	return ""
}

func expectError(code int, err error) string {
	rx := regexp.MustCompile(`#(\d+):`)
	matches := rx.FindStringSubmatch(err.Error())
	if len(matches) < 2 || matches[1] != strconv.Itoa(code) {
		return fmt.Sprintf("expected error #%d, got %s", code, err)
	}
	return ""
}

func Test001(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
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
	parser := NewParser()
	parser.DontExit = true // for testing
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
	parser := NewParser()
	parser.DontExit = true // for testing
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
	parser := NewParser()
	parser.SetAppName("myapp")
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	line := "-S4"
	e := 16
	if err := parser.ParseLine(line); err != nil {
		if e := expectError(e, err); e != "" {
			t.Error(e)
		}
	} else {
		t.Errorf("expected error #%d, got nil", e)
	}
	if parser.AppName() != "myapp" {
		t.Errorf("expected appname=myapp, got %s", parser.AppName())
	}
}

func Test005(t *testing.T) {
	parser := NewParser()
	parser.SetAppName("myapp")
	parser.SetVersion("1.0.0")
	parser.DontExit = true // for testing
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
	if parser.AppName() != "myapp" {
		t.Errorf("expected appname=myapp, got %s", parser.AppName())
	}
	if parser.Version() != "1.0.0" {
		t.Errorf("expected version=1.0.0, got %s", parser.Version())
	}
}

func Test006(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
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
	parser := NewParser()
	parser.DontExit = true // for testing
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
	parser := NewParser()
	parser.DontExit = true // for testing
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
	parser := NewParser()
	parser.DontExit = true // for testing
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
	parser := NewParser()
	parser.DontExit = true // for testing
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
	parser := NewParser()
	parser.DontExit = true // for testing
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
	parser.DontExit = true // for testing
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
	parser.DontExit = true // for testing
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
	parser.DontExit = true // for testing
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
	parser := NewParser()
	parser.DontExit = true // for testing
	languageOpt := parser.Strs("language", "lang help")
	skipLanguageOpt := parser.Strs("skiplanguage", "skip lang help")
	skipLanguageOpt.SetShortName('L')
	excludeOpt := parser.Strs("exclude", "exclude help")
	includeOpt := parser.Strs("include", "include help")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 56)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	parser.Strs("language", "lang help")
	skipLanguageOpt := parser.Strs("skiplanguage", "skip lang help")
	skipLanguageOpt.SetShortName('L')
	excludeOpt := parser.Strs("exclude", "exclude help")
	includeOpt := parser.Strs("include", "include help")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	parser.IntInRange("maxwidth", "max width help", 20, 10000, 80)
	specialSubCommand := parser.SubCommand("special", "Special help")
	extraOpt := specialSubCommand.Flag("extra", "extra help")
	maxWidthOpt := specialSubCommand.IntInRange("maxwidth",
		"max width help", 20, 10000, 80)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	languageOpt := parser.Strs("language", "lang help")
	skipLanguageOpt := parser.Strs("skiplanguage", "skip lang help")
	skipLanguageOpt.SetShortName('L')
	excludeOpt := parser.Strs("exclude", "exclude help")
	includeOpt := parser.Strs("include", "include help")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	languageOpt := parser.Strs("language", "lang help")
	skipLanguageOpt := parser.Strs("skiplanguage", "skip lang help")
	skipLanguageOpt.SetShortName('L')
	excludeOpt := parser.Strs("exclude", "exclude help")
	includeOpt := parser.Strs("include", "include help")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	languageOpt := parser.Strs("language", "lang help")
	skipLanguageOpt := parser.Strs("skiplanguage", "skip lang help")
	skipLanguageOpt.SetShortName('L')
	excludeOpt := parser.Strs("exclude", "exclude help")
	includeOpt := parser.Strs("include", "include help")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 80)
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
	parser := NewParser()
	parser.DontExit = true // for testing
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
	parser := NewParser()
	parser.DontExit = true // for testing
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 45)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 45)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	parser.IntInRange("maxwidth", "max width help", 20, 10000, 45)
	line := "--maxwidth -s"
	e := 16
	if err := parser.ParseLine(line); err != nil {
		if e := expectError(e, err); e != "" {
			t.Error(e)
		}
	} else {
		t.Errorf("expected error #%d, got nil", e)
	}
}

func Test028(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 45)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	verboseOpt := parser.Int("verbose", "verbosity -v or -vN", 1)
	verboseOpt.AllowImplicit()
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	verboseOpt := parser.Int("verbose", "verbosity -v or -vN", 1)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	parser.Int("verbose", "verbosity -v or -vN", 1)
	// -v expects either nothing (will use the default of 1) or an int
	line := "-vS"
	e := 8
	if err := parser.ParseLine(line); err != nil {
		if e := expectError(e, err); e != "" {
			t.Error(e)
		}
	} else {
		t.Errorf("expected error #%d, got nil", e)
	}
}

func Test039(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	verboseOpt := parser.Int("verbose", "verbosity -v or -vN", 1)
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
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	verboseOpt := parser.Int("verbose", "verbosity -v or -vN", 1)
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

func Test041(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	verboseOpt := parser.Int("verbose", "verbosity -v or -vN", 1)
	verboseOpt.SetShortName(0)
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

func Test042(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
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

func Test043(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	maxWidthOpt := parser.Int("maxwidth", "help", 43)
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	line := "-m -S"
	e := 34
	if err := parser.ParseLine(line); err != nil {
		if e := expectError(e, err); e != "" {
			t.Error(e)
		}
	} else {
		t.Errorf("expected error #%d, got nil (%d)", e, maxWidthOpt.Value())
	}
}

func Test044(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	maxWidthOpt := parser.Int("maxwidth", "help", 44)
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	line := "--maxwidth -S"
	e := 34
	if err := parser.ParseLine(line); err != nil {
		if e := expectError(e, err); e != "" {
			t.Error(e)
		}
	} else {
		t.Errorf("expected error #%d, got nil (%d)", e, maxWidthOpt.Value())
	}
}

func Test045(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000, 45)
	line := "--maxwidth"
	e := 34
	if err := parser.ParseLine(line); err != nil {
		if e := expectError(e, err); e != "" {
			t.Error(e)
		}
	} else {
		t.Errorf("expected error #%d, got nil (%d)", e, maxWidthOpt.Value())
	}
}

func Test046(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	scaleOpt := parser.Real("scale", "max width help", 4.5)
	line := "-Ss"
	e := 34
	if err := parser.ParseLine(line); err != nil {
		if e := expectError(e, err); e != "" {
			t.Error(e)
		}
	} else {
		t.Errorf("expected error #%d, got nil (%g)", e, scaleOpt.Value())
	}
}

func TestPackageDocFlag1(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	verboseOpt := parser.Flag("verbose", "whether to show more output")
	if err := parser.ParseLine(""); err != nil {
		t.Error("expected successful parse")
	}
	verbose := verboseOpt.Value() // verbose == false
	if verbose {
		t.Error("expected verbose=false, got true")
	}
	verbose = verboseOpt.Given() // verbose == false
	if verbose {
		t.Error("expected verbose=false, got true")
	}
	if err := parser.ParseLine("-v"); err != nil {
		t.Error("expected successful parse")
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

func TestPackageDocFlag2(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	parser.Flag("verbose", "whether to show more output")
	parser.Flag("xray", "")
	parser.Flag("cat", "")
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine(""); err != nil {
		t.Error("expected successful parse")
	}
	if outfileOpt.Given() {
		t.Error("expected outfile=!Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "" {
		t.Errorf("expected outfile=\"\", got %q", outfile)
	}
}

func TestPackageDocFlag3(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	parser.Flag("verbose", "whether to show more output")
	parser.Flag("xray", "")
	parser.Flag("cat", "")
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine("-v -x -c -o outfile.dat"); err != nil {
		t.Error("expected successful parse")
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "outfile.dat" {
		t.Errorf("expected outfile=\"outfile.dat\", got %q", outfile)
	}
}

func TestPackageDocFlag4(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	parser.Flag("verbose", "whether to show more output")
	parser.Flag("xray", "")
	parser.Flag("cat", "")
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine("-v -x -c -o=outfile.dat"); err != nil {
		t.Error("expected successful parse")
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "outfile.dat" {
		t.Errorf("expected outfile=\"outfile.dat\", got %q", outfile)
	}
}

func TestPackageDocFlag5(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	parser.Flag("verbose", "whether to show more output")
	parser.Flag("xray", "")
	parser.Flag("cat", "")
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine("-vxcooutfile.dat"); err != nil {
		t.Error("expected successful parse")
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "outfile.dat" {
		t.Errorf("expected outfile=\"outfile.dat\", got %q", outfile)
	}
}

func TestPackageDocFlag6(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	parser.Flag("verbose", "whether to show more output")
	parser.Flag("xray", "")
	parser.Flag("cat", "")
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine("-vxco outfile.dat"); err != nil {
		t.Error("expected successful parse")
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "outfile.dat" {
		t.Errorf("expected outfile=\"outfile.dat\", got %q", outfile)
	}
}

func TestPackageDocFlag7(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	parser.Flag("verbose", "whether to show more output")
	parser.Flag("xray", "")
	parser.Flag("cat", "")
	outfileOpt := parser.Str("outfile", "outfile", "")
	if err := parser.ParseLine("-vxco=outfile.dat"); err != nil {
		t.Error("expected successful parse")
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value()
	if outfile != "outfile.dat" {
		t.Errorf("expected outfile=\"outfile.dat\", got %q", outfile)
	}
}

func TestPackageDocFlag8(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	parser.Flag("verbose", "whether to show more output")
	parser.Flag("xray", "")
	parser.Flag("cat", "")
	outfileOpt := parser.Str("outfile", "outfile", "test.dat")
	if err := parser.ParseLine("-vxc"); err != nil {
		t.Error("expected successful parse")
	}
	if outfileOpt.Given() {
		t.Error("expected outfile=!Given")
	}
	outfile := outfileOpt.Value() // since not given we get the default
	if outfile != "test.dat" {
		t.Errorf("expected outfile=\"test.dat\", got %q", outfile)
	}
}

func TestPackageDocFlag9(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	parser.Flag("verbose", "whether to show more output")
	parser.Flag("xray", "")
	parser.Flag("cat", "")
	outfileOpt := parser.Str("outfile", "outfile", "test.dat")
	outfileOpt.AllowImplicit()
	if err := parser.ParseLine("-vxco"); err != nil {
		t.Error("expected successful parse")
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value() // since given with no value
	if outfile != "test.dat" {    // we get the default
		t.Errorf("expected outfile=\"test.dat\", got %q", outfile)
	}
}

func TestPackageDocFlag10(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	parser.Flag("verbose", "whether to show more output")
	parser.Flag("xray", "")
	parser.Flag("cat", "")
	outfileOpt := parser.Str("outfile", "outfile", "test.dat")
	if err := parser.ParseLine("-vxcooutfile.txt"); err != nil {
		t.Error("expected successful parse")
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value() // since given with no value
	if outfile != "outfile.txt" {    // we get the default
		t.Errorf("expected outfile=\"outfile.txt\", got %q", outfile)
	}
}

func TestPackageDocFlag11(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	parser.Flag("verbose", "whether to show more output")
	parser.Flag("xray", "")
	parser.Flag("cat", "")
	outfileOpt := parser.Str("outfile", "outfile", "test.dat")
	outfileOpt.AllowImplicit()
	if err := parser.ParseLine("-vxcooutfile.txt"); err != nil {
		t.Error("expected successful parse")
	}
	if !outfileOpt.Given() {
		t.Error("expected outfile=Given")
	}
	outfile := outfileOpt.Value() // since given with no value
	if outfile != "outfile.txt" {    // we get the default
		t.Errorf("expected outfile=\"outfile.txt\", got %q", outfile)
	}
}

func TestPackageDocSingleValue(t *testing.T) {
	parser := NewParser()
	parser.DontExit = true // for testing
	verboseOpt := parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine(""); err != nil {
		t.Error("expected successful parse")
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

	parser = NewParser()
	parser.DontExit = true // for testing
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	verboseOpt.AllowImplicit()
	if err := parser.ParseLine("-v"); err != nil {
		t.Error("expected successful parse")
	}
	expected = 1
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParser()
	parser.DontExit = true // for testing
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	verboseOpt.AllowImplicit()
	if err := parser.ParseLine("--verbose"); err != nil {
		t.Error("expected successful parse")
	}
	expected = 1
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParser()
	parser.DontExit = true // for testing
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine("-v1"); err != nil {
		t.Error("expected successful parse")
	}
	expected = 1
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParser()
	parser.DontExit = true // for testing
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine("-v=2"); err != nil {
		t.Error("expected successful parse")
	}
	expected = 2
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParser()
	parser.DontExit = true // for testing
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine("-v 3"); err != nil {
		t.Error("expected successful parse")
	}
	expected = 3
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParser()
	parser.DontExit = true // for testing
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine("--verbose=4"); err != nil {
		t.Error("expected successful parse")
	}
	expected = 4
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}

	parser = NewParser()
	parser.DontExit = true // for testing
	verboseOpt = parser.Int("verbose", "whether to show more output", 1)
	if err := parser.ParseLine("--verbose 5"); err != nil {
		t.Error("expected successful parse")
	}
	expected = 5
	verbose = verboseOpt.Value()
	if verbose != expected {
		t.Errorf("expected verbose=%d, got %d", expected, verbose)
	}
}
