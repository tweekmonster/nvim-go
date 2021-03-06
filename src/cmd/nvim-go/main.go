// Copyright 2016 The nvim-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nvim-go is a msgpack remote plugin for Neovim
package main // import "github.com/zchee/nvim-go/src/cmd/nvim-go"

import (
	"log"
	"net/http"
	_ "net/http/pprof" // For pprof debugging.
	"os"
	"runtime"

	"nvim-go/autocmd"
	"nvim-go/commands"
	"nvim-go/commands/delve"
	"nvim-go/context"

	"github.com/google/gops/agent"
	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	register := func(p *plugin.Plugin) error {
		ctxt := context.NewContext()

		c := commands.Register(p, ctxt)
		delve.Register(p, ctxt)

		autocmd.Register(p, ctxt, c)

		return nil
	}
	if os.Getenv("NVIM_GO_DEBUG") != "" {
		// starts the gops agent
		if err := agent.Listen(&agent.Options{NoShutdownCleanup: true}); err != nil {
			log.Fatal(err)
		}

		if os.Getenv("NVIM_GO_PPROF") != "" {
			addr := "localhost:14715" // (n: 14)vim-(g: 7)(o: 15)
			log.Printf("Start the pprof debugging, listen at %s\n", addr)

			// enable the report of goroutine blocking events
			runtime.SetBlockProfileRate(1)
			go func() {
				log.Println(http.ListenAndServe(addr, nil))
			}()
		}
	}

	plugin.Main(register)
}
