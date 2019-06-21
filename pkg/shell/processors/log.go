package processors

import (
	"github.com/apex/log"
)

type Log struct {
}

func (l *Log) Stdout(line string) {
	log.Info(line)
}

func (l *Log) Stderr(line string) {
	log.Error(line)
}