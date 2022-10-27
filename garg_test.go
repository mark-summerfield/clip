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
	summaryOpt := parser.Flag("summary", "summary help TODO")
	summaryOpt.ShortName = 'S'
	sortByLinesOpt := parser.Flag("sortbylines", "Sort by lines")
	specialSubCommand := parser.SubCommand("special", "Special help")
	extraOpt := specialSubCommand.Flag("extra", "extra help")

	parser.Debug()

	parser.ParseLine("-sS -l cpp d pas file1.cpp file2.d")

	summary := summaryOpt.AsBool()
	sortByLines := sortByLinesOpt.AsBool()
	fmt.Println()
	fmt.Printf("summary=%t sortbylines=%t extra=%t\n", summary, sortByLines,
		extraOpt.AsBool())
	fmt.Println()

	// TODO
}
