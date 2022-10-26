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
	parser.AddBool(MinOption{Name: "sortbylines", Help: "Sort by lines"})
	parser.AddStrs(MinOption{Name: "language", Help: "language names TODO",
		ValueCount: ZeroOrMore})
	// DEBUG
	for _, opt := range parser.Options {
		fmt.Println(opt)
	}
	fmt.Printf("%#v\n", parser)
	// END DEBUG
	sortbylines, err := parser.GetBool("sortbylines")
	if err != nil {
		t.Errorf("sortbylines expected bool got %s", err)
	}
	if sortbylines {
		t.Errorf("sortbylines expected false got %t", sortbylines)
	}
	// Safe to ignore retval since default is QuitOnError
	parser.ParseLine("-s -l cpp d pas file1.cpp file2.d")
	// TODO test parser fields
}
