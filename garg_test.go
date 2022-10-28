// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"regexp"
	"strconv"
	"testing"
)

func Test1(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	line := ""
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if summaryOpt.AsBool() {
		t.Error("expected false, got true")
	}
}

func Test2(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	line := "-S"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.AsBool() {
		t.Error("expected true, got false")
	}
}

func Test3(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	line := "--summary"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.AsBool() {
		t.Error("expected true, got false")
	}
}

func Test4(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	line := "-S4"
	e := 14
	if err := parser.ParseLine(line); err != nil {
		expectedError(e, err, t)
	} else {
		t.Errorf("expected error #%d, got nil", e)
	}
}

func Test5(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	line := "-sS"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.AsBool() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.AsBool() {
		t.Error("expected sortbylines=true, got false")
	}
}

func Test6(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	line := "-sSm60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.AsBool() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.AsBool() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.AsInt()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test7(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	line := "-sSm=60 file1.txt file2.dat file3.uxf"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.AsBool() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.AsBool() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.AsInt()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test8(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	line := "-sSm 60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.AsBool() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.AsBool() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.AsInt()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test9(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	languageOpt := parser.Strs("language", "lang help")
	line := "-sSm=60 -l cpp pas xml -- file1.txt file2.dat file3.uxf"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if !summaryOpt.AsBool() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.AsBool() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.AsInt()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
	langs := languageOpt.AsStrs()
	if len(langs) != 3 {
		t.Errorf("expected 3 languages, got %d", len(langs))
	} else {
		lang := "cpp"
		if langs[0] != lang {
			t.Errorf("expected language %s", lang)
		}
		lang = "pas"
		if langs[1] != lang {
			t.Errorf("expected language %s", lang)
		}
		lang = "xml"
		if langs[2] != lang {
			t.Errorf("expected language %s", lang)
		}
	}
}

// TODO
// -sS (i.e., -s -S)
// -sSm60 (i.e., -s -S -m60)
// -sSm=60 (i.e., -s -S -m60)
// -m60
// -m=60
// --maxwidth 60
// --maxwidth=60
// -lcpp
// -l=cpp
// -l cpp
// -l cpp pas red
// Then combinations

func expectedError(code int, err error, t *testing.T) {
	rx := regexp.MustCompile(`#(\d+):`)
	matches := rx.FindStringSubmatch(err.Error())
	if len(matches) < 2 || matches[1] != strconv.Itoa(code) {
		t.Errorf("expected error #%d, got %s", code, err)
	}
}

/*
func Test9(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	languageOpt := parser.Strs("language", "lang help")
	skipLanguageOpt := parser.Strs("skiplanguage", "skip lang help")
	skipLanguageOpt.ShortName = 'L'
	excludeOpt := parser.Strs("exclude", "exclude help")
	includeOpt := parser.Strs("include", "include help")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	specialSubCommand := parser.SubCommand("special", "Special help")
	extraOpt := specialSubCommand.Flag("extra", "extra help")
	line := "-sS -l cpp d pas file1.cpp file2.d"
	if err := parser.ParseLine(line); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("language=%v\n", languageOpt.AsStrs())
	fmt.Printf("skiplanguage=%v\n", skipLanguageOpt.AsStrs())
	fmt.Printf("exclude=%v\n", excludeOpt.AsStrs())
	fmt.Printf("include=%v\n", includeOpt.AsStrs())
	fmt.Printf("summary=%v\n", summaryOpt.AsBool())
	fmt.Printf("sortbylines=%v\n", sortByLinesOpt.AsBool())
	fmt.Printf("maxwidth=%v\n", maxWidthOpt.AsInt())
	fmt.Printf("extra=%v\n", extraOpt.AsBool())
	fmt.Println()
}
*/
