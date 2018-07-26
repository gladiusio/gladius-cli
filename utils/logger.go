package utils

import (
	"os"

	"github.com/gladiusio/gladius-utils/config"
	log "github.com/sirupsen/logrus"
)

// LogLevel - What kind of logs to show (1 = Debug and above, 2 = Info and above, 3 = Warnings and above, 4 = Fatal)
var LogLevel int

// LogFile - Where the logs are stored
var LogFile *os.File

// SetLogLevel - Sets the appropriate logging level.
// 1 = Debug < , 2 = Info <, 3 = Warning <, 4 = Fatal.
func SetLogLevel(level int) {
	switch level {
	case 1:
		log.SetLevel(log.DebugLevel)
	case 2:
		log.SetLevel(log.InfoLevel)
	case 3:
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.FatalLevel)
	}
}

// SetupLogger - Clears the previous file, and creates log file ready for writing
func SetupLogger() error {
	logPath := config.GetString("DirLogs")

	// clear previous log file
	os.Remove(logPath + "/log")

	os.MkdirAll(logPath, os.ModePerm)

	err := os.Chdir(logPath)
	if err != nil {
		PrintError(err)
	}

	LogFile, err := os.OpenFile("log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Warning("Failed to log to file, using default stderr")
		return err
	}

	log.SetOutput(LogFile)

	return nil
}
