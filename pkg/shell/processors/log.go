package processors

import (
	"github.com/apex/log"
)

type Log struct {
	Debug bool
}

func (l *Log) WriteStdout(line string) {
	if l.Debug {
		log.Debug(line)
	} else {
		log.Info(line)
	}
}

func (l *Log) WriteStderr(line string) {
	log.Error(line)
}
