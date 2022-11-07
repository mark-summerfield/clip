// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"fmt"
	"github.com/mark-summerfield/clip"
	"os"
	"path"
	"strings"
	"unicode/utf8"
)

func main() {
	config := getConfig(strings.TrimSuffix(path.Base(os.Args[0]), ".exe"),
		"0.2.0")
	fmt.Println(config)
}

func getConfig(appName, version string) config {
	descs := getDescs()
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" ||
		args[0] == "help" {
		if len(args) == 2 { // Adds support for: -h subcommand
			args[0], args[1] = args[1], args[0]
		} else {
			showHelp(descs) // doesn't return
		}
	}
	subcmd := args[0]
	args = args[1:]
	switch subcmd {
	case "-v", "--version":
		fmt.Printf("%s v%s\n", appName, version)
		os.Exit(0)
	case "c", "compare":
		return parseCompare(appName, args, descs[0])
	case "f", "format":
		return parseFormat(appName, args, descs[1])
	case "l", "lint":
		return parseLint(appName, args, descs[2])
	default:
		fmt.Printf("error: invalid subcommand %q: use -h or --help\n",
			subcmd)
		os.Exit(2)
	}
	panic("BUG getConfig")
}

func parseCompare(appName string, args []string, desc string) config {
	parser := clip.NewParser()
	parser.SetAppName(fmt.Sprintf("%s compare", appName))
	parser.PositionalCount = clip.TwoPositionals
	parser.LongDesc = desc
	parser.PositionalHelp = "The two files to compare. Each may " +
		"have any suffix and may be gzip-compressed if it ends with .gz"
	equivOpt := parser.Flag("equivalent",
		"Compare for equivalance rather than for equality")
	if err := parser.ParseArgs(args); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	config := config{
		subcommand: "compare",
		equivalent: equivOpt.Value(),
		files:      parser.Positionals,
	}
	return config
}

func parseFormat(appName string, args []string, desc string) config {
	parser := clip.NewParser()
	parser.SetAppName(fmt.Sprintf("%s format", appName))
	parser.PositionalCount = clip.TwoPositionals
	parser.LongDesc = desc
	parser.PositionalHelp = "The required infile and the required " +
		"outfile; use - to write to stdout or = to overwrite infile"
	lintOpt := parser.Flag("lint",
		"Print lints to stderr. If only lints are wanted use the l or "+
			"lint subcommand")
	standaloneOpt := parser.Flag("standalone",
		"Same as -d|--dropunused and -r|--replaceimports together")
	dropUnusedOpt := parser.Flag("dropunused",
		"Drop unused imports and ttype definitions (best to use "+
			"-s|--standalone)")
	replaceImportsOpt := parser.Flag("replaceimports",
		"Replace imports with ttype definitions for ttypes that are "+
			"actually used to make the outfile standalone (best to use "+
			"-s|--standalone)")
	indentOpt := parser.IntInRange("indent",
		"Indent (0-8 spaces or 9 to use a tab; ignored if -c|--compact "+
			"used) [default: 2]", 0, 9, 2)
	wrapWidthOpt := parser.IntInRange("wrapwidth",
		"Wrapwidth (40-240; ignored if -c|--compact used) [default: 96]",
		40, 240, 96)
	decimalsOpt := parser.IntInRange("decimals",
		"Decimal digits (0-15; 0 means use at least one (even if .0) "+
			"and as many as needed; 1-15 means used that fixed number of "+
			"digits) [default: 0]", 0, 15, 0)
	compactOpt := parser.Flag("compact",
		"Use compact output format (not human friendly; ignores indent "+
			"and wrapwidth)")
	if err := parser.ParseArgs(args); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	dropUnused := dropUnusedOpt.Value()
	replaceImports := replaceImportsOpt.Value()
	if standaloneOpt.Value() {
		dropUnused = true
		replaceImports = true
	}
	i := indentOpt.Value()
	indent := "\t"
	if i < 9 {
		indent = strings.Repeat(" ", i)
	}
	config := config{
		subcommand:     "format",
		lint:           lintOpt.Value(),
		dropUnused:     dropUnused,
		replaceImports: replaceImports,
		indent:         indent,
		wrapWidth:      wrapWidthOpt.Value(),
		decimals:       decimalsOpt.Value(),
		compact:        compactOpt.Value(),
		files:          parser.Positionals,
	}
	return config
}

func parseLint(appName string, args []string, desc string) config {
	parser := clip.NewParser()
	parser.SetAppName(fmt.Sprintf("%s lint", appName))
	parser.LongDesc = desc
	parser.PositionalHelp = "The file(s) to lint"
	if err := parser.ParseArgs(args); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return config{subcommand: "lint", files: parser.Positionals}
}

func (me config) String() string {
	files := ""
	if me.files != nil {
		files = strings.Join(me.files, " ")
	}
	switch me.subcommand {
	case "compare":
		return fmt.Sprintf("compare equiv=%t files=[%s]\n", me.equivalent,
			files)
	case "format":
		return fmt.Sprintf("format\n\tlint=%t\n\tdropunused=%t\n\t"+
			"replaceimports=%t\n\tindent=%q\n\twrapwidth=%d\n\t"+
			"decimals=%d\n\tcompact=%t\n\tfiles=[%s]\n", me.lint,
			me.dropUnused, me.replaceImports, me.indent, me.wrapWidth,
			me.decimals, me.compact, files)
	case "lint":
		return fmt.Sprintf("lint files=[%s]\n", files)
	}
	panic(fmt.Sprintf("BUG: config.String subcommand=%s", me.subcommand))
}

func showHelp(descs []string) {
	fmt.Printf("usage: %s [OPTION] [SUBCOMMAND] ...\n\n%s\n",
		clip.Bold(strings.TrimSuffix(path.Base(os.Args[0]), ".exe")),
		clip.Emph("subcommands:"))
	subs := []string{
		clip.Bold("c") + ", " + clip.Bold("compare") +
			" [-e] <FILE1> <FILE2>",
		clip.Bold("f") + ", " + clip.Bold("format") +
			" [OPTIONS] <INFILE> <OUTFILE>",
		clip.Bold("l") + ", " + clip.Bold("lint") + " <FILE> [FILE ...]"}
	leftWidths := make([]int, 0, len(subs))
	descWidths := make([]int, 0, len(descs))
	maxLeft := 0
	maxDesc := 0
	for i, sub := range subs {
		left := utf8.RuneCountInString(sub)
		if left > maxLeft {
			maxLeft = left
		}
		leftWidths = append(leftWidths, left)
		left = utf8.RuneCountInString(descs[i])
		if left > maxDesc {
			maxDesc = left
		}
		descWidths = append(descWidths, left)
	}
	argWidth := maxLeft
	width := clip.GetWidth()
	for i := 0; i < len(subs); i++ {
		if leftWidths[i]+descWidths[i]+4 <= width {
			fmt.Printf("  %s  %s\n", subs[i], descs[i])
		} else {
			fmt.Printf("  %s  ", subs[i])
			fmt.Print(clip.ArgHelp(argWidth, width, descs[i]))
		}
	}
	fmt.Println("\n" + clip.Emph("optional arguments:") + "\n" +
		"  " + clip.Bold("-v") + ", " + clip.Bold("--version") +
		"  Show version and quit\n" +
		"  " + clip.Bold("-h") + ", " + clip.Bold("--help") +
		"  Show help text and quit")
	os.Exit(0)
}

func getDescs() []string {
	descs := []string{}
	for _, desc := range []string{compareDesc, formatDesc, lintDesc} {
		descs = append(descs, strings.Join(strings.Fields(desc), " "))
	}
	return descs
}

const compareDesc = `Compare two UXF files for equality ignoring
insignificant whitespace, or for equivalence (with -e or --equivalent) in
which case the comparison ignores insignificant whitespace, comments, unused
ttypes, and, in effect replaces any imports with the ttypes they define—if
they are used. If a diff is required, format the two UXF files using the
same formatting options (and maybe use the -s --standalone option), then use
a standard diff tool.`

const formatDesc = `Copy the infile to the outfile using the canonical
human-readable format, or with the specified formatting options. This will
alphabetically order any ttype definitions and will order map items by key
(bytes < date < datetime < int < case-insensitive str). However, the order
of imports is preserved (with any duplicates removed) to allow later imports
to override earlier ones. The conversion will also automatically perform
type repairs, e.g., converting strings to dates or ints or reals if that is
the target type, and similar.`

const lintDesc = `Print the repairs that formatting would apply and lint
warnings (if any) to stderr for the given file(s).`

type config struct {
	subcommand     string
	equivalent     bool
	lint           bool
	dropUnused     bool
	replaceImports bool
	indent         string
	wrapWidth      int
	decimals       int
	compact        bool
	files          []string
}
