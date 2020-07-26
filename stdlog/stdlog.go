// Package stdlog provides simple and fast logging to the standard output
// (stdout) and is optimized for programs launched via a shell or cron. It can
// also be used to log to a file by redirecting the standard output to a file.
// This package is thread-safe.
//
// Basic examples:
//
//     logger := stdlog.GetFromFlags()
//     logger.Info("Connecting to the server...")
//     logger.Errorf("Connection failed: %q", err)
//
// Will output:
//
//     2014-04-02 18:09:15.862 INFO Connecting to the API...
//     2014-04-02 18:10:14.347 ERROR Connection failed (Server is unavailable).
//
// Log*() functions can be used to avoid evaluating arguments when it is
// expensive and unnecessary:
//
//     logger.Debug("Memory usage: %s", getMemoryUsage())
//     if LogDebug() { logger.Debug("Memory usage: %s", getMemoryUsage()) }
//
// If debug logging is off the getMemoryUsage() will be executed on the first
// line while it will not be executed on the second line.
//
// List of command-line arguments:
//
//     -log=info
//         Log events at or above this level are logged.
//     -stderr=false
//         Logs are written to standard error (stderr) instead of standard
//         output.
//     -flushlog=none
//         Until this level is reached nothing is output and logs are stored
//         in the memory. Once a log event is at or above this level, it
//         outputs all logs in memory as well as the future log events. This
//         feature should not be used with long-running processes.
//
// The available levels are the eight ones described in RFC 5424 (debug, info,
// notice, warning, error, critical, alert, emergency) and none.
//
// Some use cases:
//    - By default, all logs except debug ones are output to the stdout. Which
//      is useful to follow the execution of a program launched via a shell.
//    - A program launched by a crontab where the variable `MAILTO` is set
//      with `-debug -flushlog=error` will send all logs generated by the
//      program only if an error happens. When there is no error the email
//      will not be sent.
//    - `my_program > /var/log/my_program/my_program-$(date+%Y-%m-%d-%H%M%S).log`
//      will create a log file in /var/log/my_program each time it is run.
package stdlog

import (
	"flag"
	"io"
	"os"

	"github.com/mehrvarz/log"
	"github.com/mehrvarz/log/buflog"
	"github.com/mehrvarz/log/golog"
)

var (
	logger             log.Logger
	thresholdName      *string
	logToStderr        *bool
	flushThresholdName *string
)

// GetFromFlags returns the logger defined by the command-line flags. This
// function runs flag.Parse() if it has not been run yet.
func GetFromFlags() log.Logger {
	if logger != nil {
		return logger
	}
	if !flag.Parsed() {
		flag.Parse()
	}

	threshold := golog.GetLevelFromName(*thresholdName)
	thresholdName = nil

	out := getStream(*logToStderr)
	logToStderr = nil

	flushThreshold := golog.GetLevelFromName(*flushThresholdName)
	flushThresholdName = nil

	if flushThreshold == log.None {
		logger = golog.New(out, threshold)
	} else {
		logger = buflog.New(out, threshold, flushThreshold)
	}

	return logger
}

func GetFromFlagsDate() log.Logger {
	if logger != nil {
		return logger
	}
	if !flag.Parsed() {
		flag.Parse()
	}

	threshold := golog.GetLevelFromName(*thresholdName)
	thresholdName = nil

	out := getStream(*logToStderr)
	logToStderr = nil

	flushThreshold := golog.GetLevelFromName(*flushThresholdName)
	flushThresholdName = nil

	if flushThreshold == log.None {
		logger = golog.NewDate(out, threshold)
	} else {
		logger = buflog.New(out, threshold, flushThreshold)
	}

	return logger
}

// tmtmtm: added this so I can hand over a writer func
func GetFromFlagsWriter(myWriter func(io.Writer, []byte, log.Level)) log.Logger {
	if logger != nil {
		return logger
	}
	if !flag.Parsed() {
		flag.Parse()
	}

	threshold := golog.GetLevelFromName(*thresholdName)
	thresholdName = nil

	out := getStream(*logToStderr)
	logToStderr = nil

	flushThreshold := golog.GetLevelFromName(*flushThresholdName)
	flushThresholdName = nil

	if flushThreshold == log.None {
		logger = golog.NewWriter(out, threshold, myWriter)
	} else {
		logger = buflog.New(out, threshold, flushThreshold)
	}

	return logger
}

// tmtmtm: added this so I can hand over a writer func
func GetFromFlagsDateWriter(myWriter func(io.Writer, []byte, log.Level)) log.Logger {
	if logger != nil {
		return logger
	}
	if !flag.Parsed() {
		flag.Parse()
	}

	threshold := golog.GetLevelFromName(*thresholdName)
	thresholdName = nil

	out := getStream(*logToStderr)
	logToStderr = nil

	flushThreshold := golog.GetLevelFromName(*flushThresholdName)
	flushThresholdName = nil

	if flushThreshold == log.None {
		logger = golog.NewDateWriter(out, threshold, myWriter)
	} else {
		logger = buflog.New(out, threshold, flushThreshold)
	}

	return logger
}

func init() {
	thresholdName = flag.String("log", "info", "sets the logging threshold")
	logToStderr = flag.Bool("stderr", false, "outputs to standard error (stderr)")
	flushThresholdName = flag.String("flushlog", "none", "sets the flush trigger level")
}

// Stubbed out for testing.
var getStream = func(logToStderr bool) io.Writer {
	if logToStderr {
		return os.Stderr
	}

	return os.Stdout
}
