package log

import (
	"log"
	"os"
)

var Default Logger = &dummyLogger{log.New(os.Stdout, "", log.LstdFlags)}

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
}

func Debug(msg string) { Default.Debug(msg) }
func Info(msg string)  { Default.Info(msg) }
func Warn(msg string)  { Default.Warn(msg) }
func Error(msg string) { Default.Error(msg) }
func Fatal(msg string) { Default.Fatal(msg) }

type dummyLogger struct {
	*log.Logger
}

func (l *dummyLogger) Debug(msg string) { l.Println(msg) }
func (l *dummyLogger) Info(msg string)  { l.Println(msg) }
func (l *dummyLogger) Warn(msg string)  { l.Println(msg) }
func (l *dummyLogger) Error(msg string) { l.Println(msg) }
func (l *dummyLogger) Fatal(msg string) { l.Println(msg) }
