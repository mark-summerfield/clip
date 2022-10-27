// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPLv3

package garg

import "fmt"

type tokenState struct {
	subcommand         *SubCommand
	subCommandForName  map[string]*SubCommand
	optionForLongName  map[string]*Option
	optionForShortName map[string]*Option
	hasSubCommands     bool
	hadSubCommand      bool
}

type token struct {
	text              string
	isValue           bool
	positionalsFollow bool
}

func (me token) String() string {
	if me.positionalsFollow {
		return fmt.Sprintf("--")
	}
	if me.isValue {
		return fmt.Sprintf("%#v", me.text)
	}
	if len(me.text) == 1 {
		return fmt.Sprintf("-%s", me.text)
	}
	return fmt.Sprintf("--%s", me.text)
}

func newNameToken(text string) token {
	return token{text: text}
}

func newValueToken(text string) token {
	return token{text: text, isValue: true}
}

func newPositionalsFollowToken() token {
	return token{positionalsFollow: true}
}
