// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nvim

import (
	"fmt"

	"github.com/garyburd/neovim-go/vim"
)

var (
	progress = "Identifier"
	success  = "Function"
)

// Echo provide the vim 'echo' command.
func Echo(v *vim.Vim, format string, a ...interface{}) error {
	return v.Command("echo '" + fmt.Sprintf(format, a...) + "'")
}

// EchoRaw provide the raw output vim 'echo' command.
func EchoRaw(v *vim.Vim, a string) error {
	return v.Command("echo \"" + a + "\"")
}

// Echomsg provide the vim 'echomsg' command.
func Echomsg(v *vim.Vim, a ...interface{}) error {
	return v.Command("echomsg '" + fmt.Sprintln(a...) + "'")
}

// Echoerr provide the vim 'echoerr' command.
func Echoerr(v *vim.Vim, format string, a ...interface{}) error {
	return v.Command("echoerr '" + fmt.Sprintf(format, a...) + "'")
}

// EchohlBefore provide the vim 'echohl' command with message prefix and highlighting suffix text.
func EchohlBefore(v *vim.Vim, prefix string, highlight string, format string, a ...interface{}) error {
	v.Command("redraw")
	suffix := "' | echohl None | echon '"
	if prefix != "" {
		suffix += ": "
	}
	return v.Command("echohl " + highlight + " | echo '" + prefix + suffix + fmt.Sprintf(format, a...) + "'")
}

// EchohlAfter provide the vim 'echohl' command with message prefix and highlighting prefix text.
func EchohlAfter(v *vim.Vim, prefix string, highlight string, format string, a ...interface{}) error {
	v.Command("redraw")
	if prefix != "" {
		prefix += ": "
	}
	return v.Command("echo '" + prefix + "' | echohl " + highlight + " | echon '" + fmt.Sprintf(format, a...) + "' | echohl None")
}

// EchoProgress displays a command progress message to echo area.
func EchoProgress(v *vim.Vim, prefix, before, from, to string) error {
	v.Command("redraw")
	if prefix != "" {
		prefix += ": "
	}
	// TODO(zchee): Refactoring because line too long.
	return v.Command(fmt.Sprintf("echon '%s%s ' | echohl %s | echon '%s' | echohl None | echon ' to ' | echohl %s | echon '%s' | echohl None | echon ' ...'", prefix, before, progress, from, progress, to))
}

// EchoSuccess displays the success of the command to echo area.
func EchoSuccess(v *vim.Vim, prefix string) error {
	v.Command("redraw")
	return v.Command(fmt.Sprintf("echo '%s: ' | echohl %s | echon 'SUCCESS' | echohl None", prefix, success))
}

// ReportError output of the accumulated errors report.
// TODO(zchee): research vim.ReportError behavior
// Why it does not immediately display error?
func ReportError(v *vim.Vim, format string, a ...interface{}) error {
	return v.ReportError(fmt.Sprintf(format, a...))
}

// ClearMsg cleanups the echo area.
func ClearMsg(v *vim.Vim) error {
	return v.Command("echon")
}
