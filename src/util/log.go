package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

import (
	"code.google.com/p/log4go"
)

var Logger log4go.Logger
var initialized bool = false

func logDirCreate(logDir string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func filenameGen(progName, logDir string, isErrLog bool) string {
	strings.TrimSuffix(logDir, "/")

	var fileName string
	if isErrLog {
		/* for log file of warning, error, critical  */
		fileName = filepath.Join(logDir, progName+".wf.log")
	} else {
		/* for log file of all log  */
		fileName = filepath.Join(logDir, progName+".log")
	}

	return fileName
}

func stringToLevel(str string) log4go.LevelType {
	var level log4go.LevelType

	str = strings.ToUpper(str)

	switch str {
	case "DEBUG":
		level = log4go.DEBUG
	case "TRACE":
		level = log4go.TRACE
	case "INFO":
		level = log4go.INFO
	case "WARNING":
		level = log4go.WARNING
	case "ERROR":
		level = log4go.ERROR
	case "CRITICAL":
		level = log4go.CRITICAL
	default:
		level = log4go.INFO
	}
	return level
}

/*
* Init - initialize log lib
*
* PARAMS:
*   - progName: program name. Name of log file will be progName.log
*   - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
*   - logDir: directory for log. It will be created if noexist
*   - hasStdOut: whether to have stdout output
*   - when:
*       "M", minute
*       "H", hour
*       "D", day
*       "MIDNIGHT", roll over at midnight
*   - backupCount: If backupCount is > 0, when rollover is done, no more than
*       backupCount files are kept - the oldest ones are deleted.
*
* RETURNS:
*   nil, if succeed
*   error, if fail
 */
func LogInit(progName string, levelStr string, logDir string,
	hasStdOut bool, when string, backupCount int) error {
	if initialized {
		return errors.New("Initialized Already")
	}

	if !log4go.WhenIsValid(when) {
		return fmt.Errorf("invalid value of when: %s", when)
	}

	if err := logDirCreate(logDir); err != nil {
		log4go.Error("Init(), in logDirCreate(%s)", logDir)
		return err
	}

	level := stringToLevel(levelStr)

	Logger = make(log4go.Logger)
	if hasStdOut {
		Logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}
	fileName := filenameGen(progName, logDir, false)
	logWriter := log4go.NewTimeFileLogWriter(fileName, when, backupCount)
	if logWriter == nil {
		return fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileName)
	}
	logWriter.SetFormat(log4go.LogFormat)
	Logger.AddFilter("log", level, logWriter)

	/* create file writer for warning and fatal log */
	fileNameWf := filenameGen(progName, logDir, true)
	logWriter = log4go.NewTimeFileLogWriter(fileNameWf, when, backupCount)
	if logWriter == nil {
		return fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileNameWf)
	}
	logWriter.SetFormat(log4go.LogFormat)
	Logger.AddFilter("log_wf", log4go.WARNING, logWriter)

	initialized = true
	return nil
}
