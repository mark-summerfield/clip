// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package clip

import "fmt"

type tokenState struct {
	optionForLongName  map[string]optioner
	optionForShortName map[string]optioner
}

type tokenKind uint8

const (
	nameTokenKind tokenKind = iota
	valueTokenKind
	positionalsFollowTokenKind
	helpTokenKind
)

type token struct {
	text   string
	option optioner
	kind   tokenKind
}

func (me token) String() string {
	if me.kind == positionalsFollowTokenKind {
		return "--"
	}
	if me.kind == valueTokenKind {
		return fmt.Sprintf("%q", me.text)
	}
	if len(me.text) == 1 {
		return fmt.Sprintf("-%s", me.text)
	}
	return fmt.Sprintf("--%s", me.text)
}

func newNameToken(text string, option optioner) token {
	option.setGiven()
	return token{text: text, option: option, kind: nameTokenKind}
}

func newValueToken(text string) token {
	return token{text: text, kind: valueTokenKind}
}

func newPositionalsFollowToken() token {
	return token{kind: positionalsFollowTokenKind}
}

func newHelpToken() token {
	return token{kind: helpTokenKind}
}
