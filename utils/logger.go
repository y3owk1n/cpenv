package utils

import (
	"os"

	"github.com/charmbracelet/log"
)

var Logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportTimestamp: false,
	Level:           log.DebugLevel,
})
