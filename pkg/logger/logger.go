package logger

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

var log = logrus.New()
var DEBUG = true

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	log.Level = logrus.TraceLevel
}
func Trace(args ...interface{}) {
	log.Trace(args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	if DEBUG {
		loc := WhereAmI(2)
		log.WithFields(logrus.Fields{
			"trace": loc,
		}).Debug(args...)
	}
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	log.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	log.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) error {
	loc := WhereAmI()
	log.WithFields(logrus.Fields{
		"trace": loc,
	}).Error(args...)
	return errors.New(fmt.Sprint(args...))
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	log.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	loc := WhereAmI()
	log.WithFields(logrus.Fields{
		"trace": loc,
	}).Fatal(args...)
}
func Fatalf(format string, args ...interface{}) {
	loc := WhereAmI()
	log.WithFields(logrus.Fields{
		"trace": loc,
	}).Fatalf(format, args...)
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	if DEBUG {
		loc := WhereAmI()
		log.WithFields(logrus.Fields{
			"trace": loc,
		}).Debugf(format, args...)
	}
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) error {
	loc := WhereAmI()
	log.WithFields(logrus.Fields{
		"trace": loc,
	}).Errorf(format, args...)
	return errors.New(fmt.Sprintf(format, args...))
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func SetDebug(debug bool) {
	DEBUG = debug
}
func WhereAmI(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 2
	} else {
		depth = depthList[0]
	}
	function, file, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("File: %s  Function: %s Line: %d", chopPath(file), runtime.FuncForPC(function).Name(), line)
}
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}
