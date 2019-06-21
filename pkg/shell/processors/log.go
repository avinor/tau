package processors

import (
	"github.com/apex/log"
)

type Log struct {
}

func (l *Log) WriteStdout(line string) {
	log.Info(line)
}

func (l *Log) WriteStderr(line string) {
	log.Error(line)
}
