package ui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/bgentry/speakeasy"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// CliHandler is the default handler and will write to stdout / stderr
type CliHandler struct {
	Reader      io.Reader
	Writer      io.Writer
	ErrorWriter io.Writer

	previousLine string
}

// Ask user to input
func (hnd *CliHandler) Ask(query string) (string, error) {
	return hnd.ask(query, false)
}

// AskSecret as for secret input
func (hnd *CliHandler) AskSecret(query string) (string, error) {
	return hnd.ask(query, true)
}

// Debug prints a debug message
func (hnd *CliHandler) Debug(msg string, args ...interface{}) {
	hnd.printLine(hnd.Writer, msg, args...)
}

// Info prints message to standard out
func (hnd *CliHandler) Info(msg string, args ...interface{}) {
	hnd.printLine(hnd.Writer, msg, args...)
}

// Warn prints message
func (hnd *CliHandler) Warn(msg string, args ...interface{}) {
	hnd.printLine(hnd.Writer, color.YellowString(msg, args))
}

// Error prints message
func (hnd *CliHandler) Error(msg string, args ...interface{}) {
	hnd.printLine(hnd.ErrorWriter, msg, args...)
}

// Fatal prints message and makes sure to fail
func (hnd *CliHandler) Fatal(msg string, args ...interface{}) {
	hnd.printLine(hnd.ErrorWriter, color.RedString(msg, args))
	os.Exit(1)
}

// Header does nothing
func (hnd *CliHandler) Header(msg string) {
	hnd.NewLine()
	hnd.printLine(hnd.Writer, color.New(color.Bold).Sprint(msg))
}

// Separator does nothing
func (hnd *CliHandler) Separator() {
	hnd.printLine(hnd.Writer, "")
	hnd.printLine(hnd.Writer, "------------------------------------------------------------------------")
	hnd.printLine(hnd.Writer, "")
}

// NewLine creates a new line, except if previous line was a new line.
// Do not want 2 new lines after another
func (hnd *CliHandler) NewLine() {
	hnd.printLine(hnd.Writer, "")
}

// printLine writes a line to writer and saves it in temporary variable. It will
// never print 2 empty lines after another.
func (hnd *CliHandler) printLine(writer io.Writer, msg string, args ...interface{}) {
	if hnd.previousLine == "" && msg == "" {
		return
	}

	line := fmt.Sprintf(msg, args...)
	hnd.previousLine = line

	fmt.Fprint(writer, line)
	fmt.Fprintln(writer)
}

// ask for user input
// implementation based on https://github.com/mitchellh/cli/blob/master/ui.go
func (hnd *CliHandler) ask(query string, secret bool) (string, error) {
	if _, err := fmt.Fprint(os.Stdout, query+" "); err != nil {
		return "", err
	}

	// Register for interrupts so that we can catch it and immediately
	// return...
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	// Ask for input in a go-routine so that we can ignore it.
	errCh := make(chan error, 1)
	lineCh := make(chan string, 1)
	go func() {
		var line string
		var err error
		if secret && isatty.IsTerminal(os.Stdin.Fd()) {
			line, err = speakeasy.Ask("")
		} else {
			r := bufio.NewReader(hnd.Reader)
			line, err = r.ReadString('\n')
		}
		if err != nil {
			errCh <- err
			return
		}

		lineCh <- strings.TrimRight(line, "\r\n")
	}()

	select {
	case err := <-errCh:
		return "", err
	case line := <-lineCh:
		return line, nil
	case <-sigCh:
		// Print a newline so that any further output starts properly
		// on a new line.
		fmt.Fprintln(hnd.Writer)

		return "", errors.New("interrupted")
	}
}
