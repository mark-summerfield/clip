// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mark-summerfield/clip"
	"github.com/mark-summerfield/set"
)

const lineCountWidth = 11

func main() {
	config := getConfig("0.1.0")
	fmt.Println(config)
}

func getConfig(version string) config {
	excludes := set.New("__pycache__", "build", "build.rs", "CVS", "dist",
		"setup.py", "target")
	dataForLang := dataForLangMap{}
	initializeDataForLang(dataForLang)
	readConfigFiles(dataForLang)
	allLangs := slices.Sorted(maps.Keys(dataForLang))
	parser := clip.NewParserVersion(version)
	parser.LongDesc = fmt.Sprintf("Counts the lines in the code "+
		"files for the languages processed (ignoring . folders).\n\n"+
		"Supported language names: %s.", strings.Join(allLangs, " "))
	parser.PositionalHelp = "The files to count or the folders " +
		"to recursively search [default: .]"
	languageOpt := parser.Strs("language",
		"The language(s) to count [default: all known]")
	_ = languageOpt.SetVarName("LANG")
	skipLanguageOpt := parser.Strs("skiplanguage",
		"The languages not to count, e.g., '-L d cpp' with no '-l' "+
			"means count all languages except D and C++. Default: none")
	skipLanguageOpt.SetShortName('L')
	_ = skipLanguageOpt.SetVarName("LANG")
	excludeOpt := parser.Strs("exclude",
		fmt.Sprintf("The files and folders to exclude [default: .hidden "+
			"and %s]", strings.Join(slices.Sorted(excludes.All()), " ")))
	_ = excludeOpt.SetVarName("EXCL")
	includeOpt := parser.Strs("include",
		"The files to include (e.g., those without suffixes)")
	_ = includeOpt.SetVarName("INCL")
	width := 80
	maxWidthOpt := parser.IntInRange("maxwidth",
		"Max line width to use (e.g., for redirected output)", 0, 10000,
		width)
	debugOpt := parser.Flag("debug", "Special debug flag")
	debugOpt.SetShortName(clip.NoShortName)
	sortByLinesOpt := parser.Flag("sortbylines",
		"Sort by lines. Default: sort by names")
	summaryOpt := parser.Flag("summary",
		"Output per-language totals and total time if > 0.1 sec. "+
			"Default: output per-language and per-file totals")
	summaryOpt.SetShortName('S')
	if err := parser.Parse(); err != nil {
		parser.OnError(err)
	}
	langs := set.New[string]()
	if languageOpt.Given() {
		langs.Add(languageOpt.Value()...)
	} else {
		langs.Add(allLangs...)
	}
	if skipLanguageOpt.Given() {
		langs.Delete(skipLanguageOpt.Value()...)
	}
	if excludeOpt.Given() {
		excludes.Add(excludeOpt.Value()...)
	}
	includes := set.New[string]()
	if includeOpt.Given() {
		includes.Add(includeOpt.Value()...)
	}
	config := config{
		Language:    langs,
		Exclude:     excludes,
		Include:     includes,
		MaxWidth:    maxWidthOpt.Value() - (lineCountWidth + 2),
		SortByLines: sortByLinesOpt.Value(),
		Summary:     summaryOpt.Value(),
		File:        getPaths(parser.Positionals),
		DataForLang: dataForLang,
	}
	return config
}

func getPaths(positionals []string) set.Set[string] {
	files := set.New[string]()
	if len(positionals) == 0 {
		addPath(".", files)
	} else {
		for _, name := range positionals {
			addPath(name, files)
		}
	}
	return files
}

func addPath(filename string, files set.Set[string]) {
	path, err := filepath.Abs(filename)
	if err == nil {
		files.Add(path)
	} else {
		files.Add(filename)
	}
}

func initializeDataForLang(dataForLang dataForLangMap) {
	dataForLang["c"] = newLangData("C", ".h", ".c")
	dataForLang["cpp"] = newLangData("C++", ".hpp", ".hxx", ".cpp", ".cxx")
	dataForLang["d"] = newLangData("D", ".d")
	dataForLang["go"] = newLangData("Go", ".go")
	dataForLang["java"] = newLangData("Java", ".java")
	dataForLang["jl"] = newLangData("Julia", ".jl")
	dataForLang["nim"] = newLangData("Nim", ".nim")
	dataForLang["pl"] = newLangData("Perl", ".pl", ".pm", ".PL")
	dataForLang["py"] = newLangData("Python", ".pyw", ".py")
	dataForLang["rb"] = newLangData("Ruby", ".rb")
	dataForLang["rs"] = newLangData("Rust", ".rs")
	dataForLang["tcl"] = newLangData("Tcl", ".tcl")
	dataForLang["vala"] = newLangData("Vala", ".vala")
}

func readConfigFiles(dataForLang dataForLangMap) {
	exe, err := os.Executable()
	if err == nil {
		readConfigFile(dataForLang, path.Join(path.Dir(exe), "clc.dat"))
	}
	home, err := os.UserHomeDir()
	if err == nil {
		readConfigFile(dataForLang, path.Join(home, "clc.dat"))
		readConfigFile(dataForLang, path.Join(home, ".config/clc.dat"))
	}
	cwd, err := os.Getwd()
	if err == nil {
		readConfigFile(dataForLang, path.Join(cwd, "clc.dat"))
	}
}

func readConfigFile(dataForLang dataForLangMap, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return // ignore missing or unreadable files
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue // skip blank lines & comments
		}
		parts := strings.SplitN(line, "|", 3)
		if len(parts) == 3 {
			lang := strings.TrimSpace(parts[0])
			name := strings.TrimSpace(parts[1])
			exts := strings.Split(parts[2], " ")
			for i := range exts {
				if exts[i] != "" && exts[i][0] != '.' {
					exts[i] = "." + exts[i]
				}
			}
			dataForLang[lang] = newLangData(name, exts...)
		} else {
			fmt.Fprintf(os.Stderr, "ignoring invalid line from %s: %s",
				filename, line)
		}
	}
}

type config struct {
	Language    set.Set[string]
	Exclude     set.Set[string]
	Include     set.Set[string]
	MaxWidth    int
	SortByLines bool
	Summary     bool
	File        set.Set[string]
	DataForLang dataForLangMap
}

func (me config) String() string {
	return fmt.Sprintf("Language=[%s]\nExclude=[%s]\nInclude=[%s]\n"+
		"MaxWidth=%d\nSortByLines=%t\nSummary=%t\nFile=[%s]",
		strings.Join(slices.Sorted(me.Language.All()), " "),
		strings.Join(slices.Sorted(me.Exclude.All()), " "),
		strings.Join(slices.Sorted(me.Include.All()), " "),
		me.MaxWidth, me.SortByLines, me.Summary,
		strings.Join(slices.Sorted(me.File.All()), " "))
}

type dataForLangMap map[string]langData

type langData struct {
	Name string
	Exts set.Set[string]
}

func newLangData(name string, exts ...string) langData {
	return langData{Name: name, Exts: set.New(exts...)}
}
