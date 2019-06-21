// Copyright (c) Avinor AS. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"os"
	stdlog "log"

	"github.com/avinor/tau/cmd"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/avinor/tau/pkg/getter"
)

func main() {
	handler := cli.Default
	log.SetHandler(handler)
	stdlog.SetOutput(new(getter.LogParser))

	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
