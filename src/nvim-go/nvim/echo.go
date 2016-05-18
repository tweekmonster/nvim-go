// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nvim

import (
	"fmt"
	"os"
	"strings"

	"github.com/garyburd/neovim-go/vim"
	"github.com/pkg/errors"
)

var (
	errorColor = "Identifier"
	progress   = "Identifier"
	success    = "Function"
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
	return v.Command("echomsg '" + strings.TrimSpace(fmt.Sprintln(a...)) + "'")
}

// Echoerr provide the vim 'echoerr' command.
func Echoerr(v *vim.Vim, format string, a ...interface{}) error {
	return v.Command("echoerr '" + fmt.Sprintf(format, a...) + "'")
}

// EchoerrWrap splits the errors.Wrap's cause and error messages,
// and provide the vim 'echo' message with 'echohl' highlighting to cause text.
func EchoerrWrap(v *vim.Vim, err error) error {
	v.Command("redraw")
	er := strings.SplitAfterN(fmt.Sprintf("%s", err), ": ", 2)
	if os.Getenv("NVIM_GO_DEBUG") != "" {
		errors.Fprint(os.Stderr, err)
	}
	return v.Command("echo \"" + er[0] + "\" | echohl " + errorColor + " | echon \"" + er[1] + "\" | echohl None")
}

// EchohlErr provide the vim 'echo' command with the 'echohl' highlighting prefix text.
func EchohlErr(v *vim.Vim, prefix string, a ...interface{}) error {
	v.Command("redraw")
	if prefix != "" {
		prefix += ": "
	}
	er := fmt.Sprintf("%s", a...)
	return v.Command("echo '" + prefix + "' | echohl " + errorColor + " | echon \"" + er + "\" | echohl None")
}

// EchohlBefore provide the vim 'echo' command with the 'echohl' highlighting prefix text.
func EchohlBefore(v *vim.Vim, prefix string, highlight string, format string, a ...interface{}) error {
	v.Command("redraw")
	suffix := "' | echohl None | echon '"
	if prefix != "" {
		suffix += ": "
	}
	return v.Command("echohl " + highlight + " | echo '" + prefix + suffix + fmt.Sprintf(format, a...) + "'")
}

// EchohlAfter provide the vim 'echo' command with the 'echohl' highlighting message text.
func EchohlAfter(v *vim.Vim, prefix string, highlight string, format string, a ...interface{}) error {
	v.Command("redraw")
	if prefix != "" {
		prefix += ": "
	}
	return v.Command("echo '" + prefix + "' | echohl " + highlight + " | echon '" + fmt.Sprintf(format, a...) + "' | echohl None")
}

// EchoProgress displays a command progress message to echo area.
func EchoProgress(v *vim.Vim, prefix, format string, a ...interface{}) error {
	v.Command("redraw")
	msg := fmt.Sprintf(format, a...)
	return v.Command(fmt.Sprintf("echon '%s: ' | echohl %s | echon '%s ...' | echohl None", prefix, progress, msg))
}

// EchoSuccess displays the success of the command to echo area.
func EchoSuccess(v *vim.Vim, prefix string, msg string) error {
	v.Command("redraw")
	if msg != "" {
		msg = " " + msg
	}
	return v.Command(fmt.Sprintf("echon '%s: ' | echohl %s | echon 'SUCCESS' | echohl None | echon '%s'", prefix, success, msg))
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
