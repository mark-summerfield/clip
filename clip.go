// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package clip

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"
)

//go:embed Version.dat
var Version string // This module's version

func appName() string {
	if len(os.Args) > 0 {
		return strings.TrimSuffix(path.Base(os.Args[0]), ".exe")
	}
	return "<app>"
}

func defaultExitFunc(exitCode int, msg string) {
	if exitCode == 0 {
		fmt.Println(msg)
	} else {
		fmt.Fprintln(os.Stderr, msg)
	}
	os.Exit(exitCode)
}

var exitFunc = defaultExitFunc
