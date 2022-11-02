// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"github.com/mark-summerfield/garg"
)

func main() {
	parser := garg.NewParserVersion("0.1.0")
	parser.Parse()
}
