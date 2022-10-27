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
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	specialSubCommand := parser.SubCommand("special", "Special help")
	extraOpt := specialSubCommand.Flag("extra", "extra help")
	maxWidthOpt := parser.IntInRange("maxwidth", "max width help", 20,
		10000)

	parser.ParseLine("-sS -l cpp d pas file1.cpp file2.d")

	parser.Debug()

	summary := summaryOpt.AsBool()
	sortByLines := sortByLinesOpt.AsBool()
	fmt.Println()
	fmt.Printf("summary=%t sortbylines=%t extra=%t maxwidth=%d\n",
		summary, sortByLines, extraOpt.AsBool(), maxWidthOpt.AsInt())
	fmt.Printf("language=%v", languageOpt.AsStrs())
	fmt.Printf("skiplanguage=%v", skipLanguageOpt.AsStrs())
	fmt.Printf("exclude=%v", excludeOpt.AsStrs())
	fmt.Printf("include=%v", includeOpt.AsStrs())
	fmt.Println()

	// TODO
}
