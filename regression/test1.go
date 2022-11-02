// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"fmt"
	"github.com/mark-summerfield/garg"
)

func main() {
	parser := garg.NewParserVersion("0.1.0")
	if err := parser.Parse(); err != nil {
		fmt.Println(err)
	}
	if len(parser.Positionals) > 0 {
		fmt.Println("Got these positional args")
		for i, arg := range parser.Positionals {
			fmt.Printf("#%d: %q\n", i + 1, arg)
		}
	}
}
