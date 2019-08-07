// Copyright (c) Avinor AS. All rights reserved.
// Licensed under the Apache license.

package main

import (
	"log"
	"os"

	"github.com/avinor/tau/cmd"
	"github.com/avinor/tau/pkg/helpers/ui"
)

func main() {
	log.SetOutput(&ui.Writer{})

	if err := cmd.NewRootCmd().Execute(); err != nil {
		ui.NewLine()
		ui.Fatal("Error: %s", err)
		os.Exit(1)
	}
}
