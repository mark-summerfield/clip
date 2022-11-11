// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"bufio"
	"fmt"
	"github.com/mark-summerfield/clip"
	"github.com/mark-summerfield/gset"
	"golang.org/x/exp/maps"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const lineCountWidth = 11

func main() {
	config := getConfig("0.3.0")
	fmt.Println(config)
}

func getConfig(version string) config {
	excludes := gset.New("__pycache__", "build", "build.rs", "CVS", "dist",
		"setup.py", "target")
	dataForLang := make(dataForLangMap)
	initializeDataForLang(dataForLang)
	readConfigFiles(dataForLang)
	allLangs := maps.Keys(dataForLang)
	sort.Strings(allLangs)
	parser := clip.NewParserVersion(version)
	parser.ShortDesc = "Counts the lines in the code files for the " +
		"languages processed"
	parser.LongDesc = fmt.Sprintf("The counting ignores  . folders.\n\n"+
		"Supported language names: %s.", strings.Join(allLangs, " "))
	parser.EndDesc = "Some end notes just to show that EndDesc works."
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
			"and %s]", strings.Join(excludes.ToSortedSlice(), " ")))
	_ = excludeOpt.SetVarName("EXCL")
	includeOpt := parser.Strs("include",
		"The files to include (e.g., those without suffixes)")
	_ = includeOpt.SetVarName("INCL")
	width := 80
	maxWidthOpt := parser.IntInRange("maxwidth",
		"Max line width to use (e.g., for redirected output)", 0, 10000,
		width)
	debugOpt := parser.Flag("debug", "Special debug flag")
	debugOpt.SetShortName('D')
	debugOpt.Hide()
	sortByLinesOpt := parser.Flag("sortbylines",
		"Sort by lines. Default: sort by names")
	summaryOpt := parser.Flag("summary",
		"Output per-language totals and total time if > 0.1 sec. "+
			"Default: output per-language and per-file totals")
	summaryOpt.SetShortName('S')
	if err := parser.Parse(); err != nil {
		parser.OnError(err)
	}
	langs := gset.New[string]()
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
	includes := gset.New[string]()
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
		Debug:       debugOpt.Value(),
	}
	return config
}

func getPaths(positionals []string) gset.Set[string] {
	files := gset.New[string]()
	if len(positionals) == 0 {
		addPath(".", files)
	} else {
		for _, name := range positionals {
			addPath(name, files)
		}
	}
	return files
}

func addPath(filename string, files gset.Set[string]) {
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
			for i := 0; i < len(exts); i++ {
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
	Language    gset.Set[string]
	Exclude     gset.Set[string]
	Include     gset.Set[string]
	MaxWidth    int
	SortByLines bool
	Summary     bool
	File        gset.Set[string]
	DataForLang dataForLangMap
	Debug       bool
}

func (me config) String() string {
	return fmt.Sprintf("Language=[%s]\nExclude=[%s]\nInclude=[%s]\n"+
		"MaxWidth=%d\nSortByLines=%t\nSummary=%t\nFile=[%s]\nDebug=%t",
		strings.Join(me.Language.ToSortedSlice(), " "),
		strings.Join(me.Exclude.ToSortedSlice(), " "),
		strings.Join(me.Include.ToSortedSlice(), " "),
		me.MaxWidth, me.SortByLines, me.Summary,
		strings.Join(me.File.ToSortedSlice(), " "), me.Debug)
}

type dataForLangMap map[string]langData

type langData struct {
	Name string
	Exts gset.Set[string]
}

func newLangData(name string, exts ...string) langData {
	return langData{Name: name, Exts: gset.New(exts...)}
}
