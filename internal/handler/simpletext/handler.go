package simpletext

import (
	"fmt"
	"os"
	"sync"

	"github.com/apex/log"
	"github.com/fatih/color"
)

// Default handler outputting to stderr.
var Default = &Handler{}

// Colors mapping.
var Colors = [...]*color.Color{
	log.DebugLevel: color.New(color.FgWhite),
	log.InfoLevel:  color.New(color.FgBlue),
	log.WarnLevel:  color.New(color.FgYellow),
	log.ErrorLevel: color.New(color.FgRed),
	log.FatalLevel: color.New(color.FgRed),
}

// Handler implementation.
type Handler struct {
	mu sync.Mutex
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	color := Colors[e.Level]
	names := e.Fields.Names()

	writer := os.Stdout
	if e.Level == log.ErrorLevel || e.Level == log.FatalLevel {
		writer = os.Stderr
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Fprintf(writer, "%-25s", e.Message)

	for _, name := range names {
		if name == "source" {
			continue
		}
		fmt.Fprintf(writer, " %s=%v", color.Sprint(name), e.Fields.Get(name))
	}

	fmt.Fprintln(writer)

	return nil
}
