// Copyright (c) Avinor AS. All rights reserved.
// Licensed under the MIT license.

package main

import (
	stdlog "log"
	"os"

	"github.com/apex/log"
	"github.com/avinor/tau/internal/handler/simpletext"
	"github.com/avinor/tau/cmd"
	"github.com/avinor/tau/pkg/getter"
)

func main() {
	log.SetHandler(simpletext.Default)
	stdlog.SetOutput(new(getter.LogParser))

	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
