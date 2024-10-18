// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package clip

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/mark-summerfield/uterm"
)

//go:embed Version.dat
var Version string // This module's version

func appName() string {
	if len(os.Args) > 0 {
		name := strings.TrimSuffix(os.Args[0], ".exe")
		i := strings.LastIndexAny(name, "\\/:")
		if i > -1 {
			name = name[i+1:]
		}
		return name
	}
	return "<app>"
}

func defaultExitFunc(exitCode int, msg string) {
	if exitCode == 0 {
		fmt.Printf(uterm.Red(msg))
	} else {
		fmt.Fprintln(os.Stderr, uterm.Red(msg))
		fmt.Fprintln(os.Stderr, uterm.Red(fmt.Sprintf(
			"for help run: %s --help", appName())))
	}
	os.Exit(exitCode)
}

var exitFunc = defaultExitFunc
