// Copyright 2016 The nvim-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"nvim-go/config"
	"nvim-go/nvimutil"

	"github.com/neovim/go-client/nvim"
	"github.com/pkg/errors"
	"golang.org/x/tools/refactor/rename"
)

const pkgRename = "GoRename"

type cmdRenameEval struct {
	Cwd        string `msgpack:",array"`
	File       string
	RenameFrom string
}

func (c *Commands) cmdRename(args []string, bang bool, eval *cmdRenameEval) {
	go c.Rename(args, bang, eval)
}

// Rename rename the current cursor word use golang.org/x/tools/refactor/rename.
func (c *Commands) Rename(args []string, bang bool, eval *cmdRenameEval) error {
	defer nvimutil.Profile(time.Now(), "GoRename")
	dir := filepath.Dir(eval.File)
	defer c.ctx.SetContext(dir)()

	var (
		b nvim.Buffer
		w nvim.Window
	)
	c.Pipeline.CurrentBuffer(&b)
	c.Pipeline.CurrentWindow(&w)
	if err := c.Pipeline.Wait(); err != nil {
		return nvimutil.ErrorWrap(c.Nvim, errors.WithStack(err))
	}

	offset, err := nvimutil.ByteOffset(c.Nvim, b, w)
	if err != nil {
		return nvimutil.ErrorWrap(c.Nvim, errors.WithStack(err))
	}
	pos := fmt.Sprintf("%s:#%d", eval.File, offset)

	var renameTo string
	if len(args) > 0 {
		renameTo = args[0]
	} else {
		askMessage := fmt.Sprintf("%s: Rename '%s' to: ", pkgRename, eval.RenameFrom)
		var toResult interface{}
		if config.RenamePrefill {
			err := c.Nvim.Call("input", &toResult, askMessage, eval.RenameFrom)
			if err != nil {
				return nvimutil.EchohlErr(c.Nvim, pkgRename, "Keyboard interrupt")
			}
		} else {
			err := c.Nvim.Call("input", &toResult, askMessage)
			if err != nil {
				return nvimutil.EchohlErr(c.Nvim, pkgRename, "Keyboard interrupt")
			}
		}
		if toResult.(string) == "" {
			return nvimutil.EchohlErr(c.Nvim, pkgRename, "Not enough arguments for rename destination name")
		}
		renameTo = fmt.Sprintf("%s", toResult)
	}

	c.Nvim.Command(fmt.Sprintf("echo '%s: Renaming ' | echohl Identifier | echon '%s' | echohl None | echon ' to ' | echohl Identifier | echon '%s' | echohl None | echon ' ...'", pkgRename, eval.RenameFrom, renameTo))

	if bang {
		rename.Force = true
	}

	// TODO(zchee): More elegant way
	// save original stdout and stderr
	saveStdout, saveStderr := os.Stdout, os.Stderr
	read, write, _ := os.Pipe()
	// migrate stderr and stdout
	os.Stderr = os.Stdout
	os.Stderr = write
	defer func() {
		os.Stderr = saveStdout
		os.Stderr = saveStderr
	}()

	if err := rename.Main(&build.Default, pos, "", renameTo); err != nil {
		write.Close()
		renameErr, err := ioutil.ReadAll(read)
		if err != nil {
			return nvimutil.ErrorWrap(c.Nvim, errors.WithStack(err))
		}

		log.Printf("er: %+v\n", string(renameErr))
		go func() {
			loclist, _ := nvimutil.ParseError(renameErr, eval.Cwd, &c.ctx.Build)
			nvimutil.SetLoclist(c.Nvim, loclist)
			nvimutil.OpenLoclist(c.Nvim, w, loclist, true)
		}()

		return nvimutil.ErrorWrap(c.Nvim, errors.WithStack(err))
	}

	write.Close()
	out, _ := ioutil.ReadAll(read)
	defer nvimutil.EchoSuccess(c.Nvim, pkgRename, fmt.Sprintf("%s", out))

	// TODO(zchee): 'edit' command is ugly.
	// Should create tempfile and use SetBufferLines.
	return c.Nvim.Command("silent edit")
}
