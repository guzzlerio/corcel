package main

import (
	"strings"
)

type StateDelegate func(char rune)

type CommandLineLexer struct {
	args          []string
	state         string
	allowRune     bool
	delegate      StateDelegate
	enclosingRune rune
}

func NewCommandLineLexer() *CommandLineLexer {
	return &CommandLineLexer{}
}

func (instance *CommandLineLexer) recordFlag(char rune) {
	if char == ' ' || char == ':' {
		instance.args = append(instance.args, instance.state)
		instance.state = ""
		instance.delegate = instance.recordFlagValue
	} else {
		instance.state += string(char)
	}
}

func (instance *CommandLineLexer) recordFlagValue(char rune) {
	if char == ' ' {
		if instance.allowRune {
			instance.state += string(char)
		} else {
			if instance.state != ""{
				instance.args = append(instance.args, instance.state)
				instance.state = ""
				instance.delegate = instance.start
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

func (instance *CommandLineLexer) recordValue(char rune) {
	if char == ' ' {
		instance.args = append(instance.args, strings.Trim(instance.state, "\""))
		instance.state = ""
		instance.delegate = instance.start
	} else {
		instance.state += string(char)
	}
}

func (instance *CommandLineLexer) start(char rune) {
	if char == '-' {
		instance.delegate = instance.recordFlag
		instance.recordFlag(char)
	} else if char == ' ' {
		instance.delegate = instance.start
	} else {
		instance.delegate = instance.recordValue
		instance.recordValue(char)
	}
}

func (instance *CommandLineLexer) Lex(args string) []string {
	instance.state = ""
	instance.args = []string{}
	instance.delegate = instance.start
	for _, char := range args {
		instance.delegate(char)
	}
	if instance.state != ""{
		instance.args = append(instance.args, instance.state)
	}
	return instance.args
}