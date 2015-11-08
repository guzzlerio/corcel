package main

import (
	"strings"
)

//StateDelegate ...
type StateDelegate func(char rune)

//CommandLineLexer ...
type CommandLineLexer struct {
	args          []string
	state         string
	allowRune     bool
	enclosingRune rune
	delegate      StateDelegate
}

//NewCommandLineLexer ...
func NewCommandLineLexer() *CommandLineLexer {
	return &CommandLineLexer{}
}

//RecordFlag ...
func (instance *CommandLineLexer) RecordFlag(char rune) {
	if char == ' ' || char == ':' {
		instance.args = append(instance.args, instance.state)
		instance.state = ""
		instance.delegate = instance.RecordFlagValue
	} else {
		instance.state += string(char)
	}
}

//RecordFlagValue ...
func (instance *CommandLineLexer) RecordFlagValue(char rune) {
	if char == ' ' {
		if instance.allowRune {
			instance.state += string(char)
		} else {
			if instance.state != "" {
				instance.args = append(instance.args, instance.state)
				instance.state = ""
				instance.delegate = instance.Start
			}
		}
	} else if char == '"' || char == '\'' {
		if !instance.allowRune {
			instance.enclosingRune = char
			instance.allowRune = true
		} else {
			if char == instance.enclosingRune {
				instance.allowRune = false
			} else {
				instance.state += string(char)
			}
		}
	} else {
		instance.state += string(char)
	}
}

//RecordValue ...
func (instance *CommandLineLexer) RecordValue(char rune) {
	if char == ' ' {
		instance.args = append(instance.args, strings.Trim(instance.state, "\""))
		instance.state = ""
		instance.delegate = instance.Start
	} else {
		instance.state += string(char)
	}
}

//Start ...
func (instance *CommandLineLexer) Start(char rune) {
	if char == '-' {
		instance.delegate = instance.RecordFlag
		instance.RecordFlag(char)
	} else if char == ' ' {
		instance.delegate = instance.Start
	} else {
		instance.delegate = instance.RecordValue
		instance.RecordValue(char)
	}
}

//Lex ...
func (instance *CommandLineLexer) Lex(args string) []string {
	instance.state = ""
	instance.args = []string{}
	instance.delegate = instance.Start
	for _, char := range args {
		instance.delegate(char)
	}
	if instance.state != "" {
		instance.args = append(instance.args, instance.state)
	}
	return instance.args
}
