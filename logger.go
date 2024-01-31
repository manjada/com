package mjd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strings"
)

var log = logrus.New()

type LogFields map[string]interface{}

func init() {
	log.Info("Start Load Logger")
	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	log.Out = os.Stdout

	// You could set this to any `io.Writer` such as a file
	var logPath string
	//get operating system
	if runtime.GOOS == "windows" {
		logPath = GetConfig().LogFile.PathWindows
	} else {
		logPath = GetConfig().LogFile.PathUnix
	}

	file, err := os.OpenFile(fmt.Sprintf(logPath+"logger.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Error(err)
	}

	mw := io.MultiWriter(os.Stdout, log.Out)
	log.SetOutput(mw)
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})

	// Only log the Info severity or above.
	logLevel, err := logrus.ParseLevel(GetConfig().LogFile.Level)
	if err != nil {
		log.Error(err)
	}
	log.SetLevel(logLevel)
}

// Info logs a message at level Info on the standard logger.
func Info(message string, fields ...interface{}) {
	entry := log.WithFields(logrus.Fields{"detail": fields})
	entry.Data["file"] = fileInfo(2)
	entry.Info(message)
}

func Error(err error, fields ...interface{}) {
	entry := log.WithFields(logrus.Fields{"detail": fields})
	entry.Data["file"] = fileInfo(2)
	entry.Error(err)
}

func Fatal(value string, fields ...interface{}) {
	entry := log.WithFields(logrus.Fields{"detail": fields})
	entry.Data["file"] = fileInfo(2)
	entry.Fatal(value)
}

func Panic(err error, fields ...interface{}) {
	entry := log.WithFields(logrus.Fields{"detail": fields})
	entry.Data["file"] = fileInfo(2)
	entry.Panic(err)
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
