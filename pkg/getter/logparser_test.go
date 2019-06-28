package getter

import (
	"fmt"
	stdlog "log"
	"testing"
	"sync"

	"github.com/apex/log"
	"github.com/stretchr/testify/assert"
)

func TestLogParser(t *testing.T) {
	tests := []struct {
		Message       string
		ExpectLevel   log.Level
		ExpectMessage string
	}{
		{
			"[TRACE] Test message",
			log.DebugLevel,
			"Test message",
		},
		{
			"[DEBUG] Test message",
			log.DebugLevel,
			"Test message",
		},
		{
			"[INFO] Test message",
			log.InfoLevel,
			"Test message",
		},
		{
			"[WARN] Test message",
			log.WarnLevel,
			"Test message",
		},
		{
			"[ERR] Test message",
			log.ErrorLevel,
			"Test message",
		},
		{
			"[ERROR] Test message",
			log.ErrorLevel,
			"Test message",
		},
		{
			"Test message",
			log.InfoLevel,
			"Test message",
		},
	}

	testLogger := &testLogger{}
	logger := &LogParser{
		Logger: testLogger,
	}
	stdlog.SetOutput(logger)
	mux := sync.Mutex{}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			mux.Lock()
			defer mux.Unlock()

			stdlog.Printf(test.Message)

			assert.Equal(t, test.ExpectLevel, testLogger.level)
			assert.Equal(t, test.ExpectMessage, testLogger.message)
		})
	}
}

type testLogger struct {
	level   log.Level
	message string
}

func (l *testLogger) WithFields(fields log.Fielder) *log.Entry           { return nil }
func (l *testLogger) WithField(key string, value interface{}) *log.Entry { return nil }
func (l *testLogger) WithError(err error) *log.Entry                     { return nil }
func (l *testLogger) Debug(msg string)                                   { l.level = log.DebugLevel; l.message = msg }
func (l *testLogger) Info(msg string)                                    { l.level = log.InfoLevel; l.message = msg }
func (l *testLogger) Warn(msg string)                                    { l.level = log.WarnLevel; l.message = msg }
func (l *testLogger) Error(msg string)                                   { l.level = log.ErrorLevel; l.message = msg }
func (l *testLogger) Fatal(msg string)                                   { l.level = log.FatalLevel; l.message = msg }
func (l *testLogger) Debugf(msg string, v ...interface{})                { l.level = log.DebugLevel; l.message = msg }
func (l *testLogger) Infof(msg string, v ...interface{})                 { l.level = log.InfoLevel; l.message = msg }
func (l *testLogger) Warnf(msg string, v ...interface{})                 { l.level = log.WarnLevel; l.message = msg }
func (l *testLogger) Errorf(msg string, v ...interface{})                { l.level = log.ErrorLevel; l.message = msg }
func (l *testLogger) Fatalf(msg string, v ...interface{})                { l.level = log.FatalLevel; l.message = msg }
func (l *testLogger) Trace(msg string) *log.Entry                        { return nil }
