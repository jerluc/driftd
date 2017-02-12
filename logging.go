package main

import (
	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("riftd")

// Initializes logging facilities at the
// given level
func InitLogging(levelString string) {
	formatter := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	level, _ := logging.LogLevel(levelString)
	logging.SetFormatter(formatter)
	logging.SetLevel(level, "")
}
