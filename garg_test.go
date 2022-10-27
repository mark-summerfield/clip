// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	parser := NewParser("myapp", "1.0.0")
	parser.QuitOnError = false // for testing
	parser.PositionalCount = ZeroOrMore
	parser.PositionalVarName = "FILENAME"
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

	fmt.Println()
	fmt.Printf("language=%v\n", languageOpt.AsStrs())
	fmt.Printf("skiplanguage=%v\n", skipLanguageOpt.AsStrs())
	fmt.Printf("exclude=%v\n", excludeOpt.AsStrs())
	fmt.Printf("include=%v\n", includeOpt.AsStrs())
	fmt.Printf("summary=%v\n", summaryOpt.AsBool())
	fmt.Printf("sortbylines=%v\n", sortByLinesOpt.AsBool())
	fmt.Printf("maxwidth=%v\n", maxWidthOpt.AsInt())
	fmt.Printf("extra=%v\n", extraOpt.AsBool())
}
