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
	summaryOpt.SetShortName('S')
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
	summaryOpt.SetShortName('S')
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
	summaryOpt.SetShortName('S')
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
	summaryOpt.SetShortName('S')
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
	summaryOpt.SetShortName('S')
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
	summaryOpt.SetShortName('S')
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
	summaryOpt.SetShortName('S')
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
	summaryOpt.SetShortName('S')
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
	summaryOpt.SetShortName('S')
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

func Test10(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	line := "-m60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	m := maxWidthOpt.AsInt()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test11(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	line := "-m=60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	m := maxWidthOpt.AsInt()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test12(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	line := "--maxwidth 60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	m := maxWidthOpt.AsInt()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test13(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	line := "--maxwidth=60"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	m := maxWidthOpt.AsInt()
	if m != 60 {
		t.Errorf("expected maxwidth=60, got %d", m)
	}
}

func Test14(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	languageOpt := parser.Strs("language", "lang help")
	line := "-lcpp -- file1.txt file2.dat file3.uxf"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	langs := languageOpt.AsStrs()
	if len(langs) != 1 {
		t.Errorf("expected 1 language, got %d", len(langs))
	} else {
		lang := "cpp"
		if langs[0] != lang {
			t.Errorf("expected language %s", lang)
		}
	}
}

func Test15(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	languageOpt := parser.Strs("language", "lang help")
	line := "-l=cpp -- file1.txt file2.dat file3.uxf"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	langs := languageOpt.AsStrs()
	if len(langs) != 1 {
		t.Errorf("expected 1 language, got %d", len(langs))
	} else {
		lang := "cpp"
		if langs[0] != lang {
			t.Errorf("expected language %s", lang)
		}
	}
}

func Test16(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	languageOpt := parser.Strs("language", "lang help")
	line := "-l cpp -- file1.txt file2.dat file3.uxf"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	langs := languageOpt.AsStrs()
	if len(langs) != 1 {
		t.Errorf("expected 1 language, got %d", len(langs))
	} else {
		lang := "cpp"
		if langs[0] != lang {
			t.Errorf("expected language %s", lang)
		}
	}
}

func Test17(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	languageOpt := parser.Strs("language", "lang help")
	line := "-l cpp pas xml -- file1.txt file2.dat file3.uxf"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
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

func Test18(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	languageOpt := parser.Strs("language", "lang help")
	skipLanguageOpt := parser.Strs("skiplanguage", "skip lang help")
	skipLanguageOpt.SetShortName('L')
	excludeOpt := parser.Strs("exclude", "exclude help")
	includeOpt := parser.Strs("include", "include help")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	maxWidthOpt.SetDefaultValue(56)
	specialSubCommand := parser.SubCommand("special", "Special help")
	extraOpt := specialSubCommand.Flag("extra", "extra help")
	line := "-sS -l h red -e zOld t -L d -i peek -- file1.cpp file2.d"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	langs := languageOpt.AsStrs()
	if len(langs) != 2 {
		t.Errorf("expected 2 languages, got %d", len(langs))
	} else {
		lang := "h"
		if langs[0] != lang {
			t.Errorf("expected language %s", lang)
		}
		lang = "red"
		if langs[1] != lang {
			t.Errorf("expected language %s", lang)
		}
	}
	langs = skipLanguageOpt.AsStrs()
	if len(langs) != 1 {
		t.Errorf("expected 1 languages, got %d", len(langs))
	} else {
		lang := "d"
		if langs[0] != lang {
			t.Errorf("expected language %s", lang)
		}
	}
	exclude := excludeOpt.AsStrs()
	if len(exclude) != 2 {
		t.Errorf("expected 2 excludes, got %d", len(langs))
	} else {
		excl := "zOld"
		if exclude[0] != excl {
			t.Errorf("expected exclude %s", excl)
		}
		excl = "t"
		if exclude[1] != excl {
			t.Errorf("expected exclude %s", excl)
		}
	}
	include := includeOpt.AsStrs()
	if len(include) != 1 {
		t.Errorf("expected 1 includes, got %d", len(langs))
	} else {
		incl := "peek"
		if include[0] != incl {
			t.Errorf("expected include %s", incl)
		}
	}
	if !summaryOpt.AsBool() {
		t.Error("expected summary=true, got false")
	}
	if !sortByLinesOpt.AsBool() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.AsInt()
	if m != 56 {
		t.Errorf("expected maxwidth=56, got %d", m)
	}
	if extraOpt.AsBool() {
		t.Error("expected extra=false, got true")
	}
}

func Test19(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	parser.Strs("language", "lang help")
	skipLanguageOpt := parser.Strs("skiplanguage", "skip lang help")
	skipLanguageOpt.SetShortName('L')
	parser.Strs("exclude", "exclude help")
	parser.Strs("include", "include help")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	parser.IntInRange("maxwidth", "max width help", 20, 10000)
	specialSubCommand := parser.SubCommand("special", "Special help")
	extraOpt := specialSubCommand.Flag("extra", "extra help")
	maxWidthOpt := specialSubCommand.IntInRange("maxwidth",
		"max width help", 20, 10000)
	line := "special -e -m98 file1.cpp file2.d"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	if summaryOpt.AsBool() {
		t.Error("expected summary=false, got true")
	}
	if sortByLinesOpt.AsBool() {
		t.Error("expected sortbylines=false, got true")
	}
	if !extraOpt.AsBool() {
		t.Error("expected extra=true, got false")
	}
	m := maxWidthOpt.AsInt()
	if m != 98 {
		t.Errorf("expected maxwidth=98, got %d", m)
	}
}

func Test20(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	languageOpt := parser.Strs("language", "lang help")
	skipLanguageOpt := parser.Strs("skiplanguage", "skip lang help")
	skipLanguageOpt.SetShortName('L')
	excludeOpt := parser.Strs("exclude", "exclude help")
	includeOpt := parser.Strs("include", "include help")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	line := "-S -l h red -e zOld t -L d -i peek -m 40 file1.cpp file2.d"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	langs := languageOpt.AsStrs()
	if len(langs) != 2 {
		t.Errorf("expected 2 languages, got %d", len(langs))
	} else {
		lang := "h"
		if langs[0] != lang {
			t.Errorf("expected language %s", lang)
		}
		lang = "red"
		if langs[1] != lang {
			t.Errorf("expected language %s", lang)
		}
	}
	langs = skipLanguageOpt.AsStrs()
	if len(langs) != 1 {
		t.Errorf("expected 1 languages, got %d", len(langs))
	} else {
		lang := "d"
		if langs[0] != lang {
			t.Errorf("expected language %s", lang)
		}
	}
	exclude := excludeOpt.AsStrs()
	if len(exclude) != 2 {
		t.Errorf("expected 2 excludes, got %d", len(langs))
	} else {
		excl := "zOld"
		if exclude[0] != excl {
			t.Errorf("expected exclude %s", excl)
		}
		excl = "t"
		if exclude[1] != excl {
			t.Errorf("expected exclude %s", excl)
		}
	}
	include := includeOpt.AsStrs()
	if len(include) != 1 {
		t.Errorf("expected 1 includes, got %d", len(langs))
	} else {
		incl := "peek"
		if include[0] != incl {
			t.Errorf("expected include %s", incl)
		}
	}
	if !summaryOpt.AsBool() {
		t.Error("expected summary=true, got false")
	}
	if sortByLinesOpt.AsBool() {
		t.Error("expected sortbylines=false, got true")
	}
	m := maxWidthOpt.AsInt()
	if m != 40 {
		t.Errorf("expected maxwidth=40, got %d", m)
	}
}

func Test21(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	languageOpt := parser.Strs("language", "lang help")
	skipLanguageOpt := parser.Strs("skiplanguage", "skip lang help")
	skipLanguageOpt.SetShortName('L')
	excludeOpt := parser.Strs("exclude", "exclude help")
	includeOpt := parser.Strs("include", "include help")
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.SetShortName('S')
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)
	maxWidthOpt.SetDefaultValue(80)
	line := "-l h red -e zOld t -L d -i peek -s file1.cpp file2.d"
	if err := parser.ParseLine(line); err != nil {
		t.Error(err)
	}
	langs := languageOpt.AsStrs()
	if len(langs) != 2 {
		t.Errorf("expected 2 languages, got %d", len(langs))
	} else {
		lang := "h"
		if langs[0] != lang {
			t.Errorf("expected language %s", lang)
		}
		lang = "red"
		if langs[1] != lang {
			t.Errorf("expected language %s", lang)
		}
	}
	langs = skipLanguageOpt.AsStrs()
	if len(langs) != 1 {
		t.Errorf("expected 1 languages, got %d", len(langs))
	} else {
		lang := "d"
		if langs[0] != lang {
			t.Errorf("expected language %s", lang)
		}
	}
	exclude := excludeOpt.AsStrs()
	if len(exclude) != 2 {
		t.Errorf("expected 2 excludes, got %d", len(langs))
	} else {
		excl := "zOld"
		if exclude[0] != excl {
			t.Errorf("expected exclude %s", excl)
		}
		excl = "t"
		if exclude[1] != excl {
			t.Errorf("expected exclude %s", excl)
		}
	}
	include := includeOpt.AsStrs()
	if len(include) != 1 {
		t.Errorf("expected 1 includes, got %d", len(langs))
	} else {
		incl := "peek"
		if include[0] != incl {
			t.Errorf("expected include %s", incl)
		}
	}
	if summaryOpt.AsBool() {
		t.Error("expected summary=false, got true")
	}
	if !sortByLinesOpt.AsBool() {
		t.Error("expected sortbylines=true, got false")
	}
	m := maxWidthOpt.AsInt()
	if m != 80 {
		t.Errorf("expected maxwidth=80, got %d", m)
	}
}

func expectedError(code int, err error, t *testing.T) {
	rx := regexp.MustCompile(`#(\d+):`)
	matches := rx.FindStringSubmatch(err.Error())
	if len(matches) < 2 || matches[1] != strconv.Itoa(code) {
		t.Errorf("expected error #%d, got %s", code, err)
	}
}
