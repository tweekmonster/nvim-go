// Copyright 2016 The nvim-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autocmd

import (
	"sync"

	"nvim-go/commands"
	"nvim-go/context"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

// Autocmd represents a autocmd context.
type Autocmd struct {
	Nvim *nvim.Nvim
	ctxt *context.Context
	cmds *commands.Commands

	bufWritePostChan chan error
	bufWritePreChan  chan interface{}
	mu               sync.Mutex
	wg               sync.WaitGroup

	errors []error
}

// Register register autocmd to nvim.
func Register(p *plugin.Plugin, ctxt *context.Context, cmds *commands.Commands) {
	autocmd := &Autocmd{
		Nvim:             p.Nvim,
		ctxt:             ctxt,
		cmds:             cmds,
		bufWritePreChan:  make(chan interface{}),
		bufWritePostChan: make(chan error),
	}

	p.HandleAutocmd(&plugin.AutocmdOptions{Event: "VimEnter", Pattern: "*.go", Group: "nvim-go", Eval: "*"}, autocmd.VimEnter)
	p.HandleAutocmd(&plugin.AutocmdOptions{Event: "BufWritePre", Pattern: "*.go", Group: "nvim-go", Eval: "[getcwd(), expand('%:p')]"}, autocmd.BufWritePre)
	p.HandleAutocmd(&plugin.AutocmdOptions{Event: "BufWritePost", Pattern: "*.go", Group: "nvim-go", Eval: "[getcwd(), expand('%:p')]"}, autocmd.BufWritePost)
}
