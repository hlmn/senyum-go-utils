package echo

import (
	"io"
	"os"
	"time"

	sentryUmi "github.com/fonysaputra/go-utils/apm/sentry"

	"github.com/labstack/echo/v4"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logDump "github.com/sirupsen/logrus"
)

var (
	userId = ""
)

func InitBodyDumpLog() (err error) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	logf, err := rotatelogs.New(
		dir+"/logs/RequestResponseDump.log.%Y%m%d",
		rotatelogs.WithLinkName(dir+"/logs/RequestResponseDump.log"),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(365),
	)

	logDump.SetFormatter(&logDump.JSONFormatter{DisableHTMLEscape: true})
	logDump.SetOutput(io.MultiWriter(os.Stdout, logf))
	logDump.SetLevel(logDump.InfoLevel)
	logDump.SetReportCaller(true)

	return
}

func Info(c echo.Context, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Info(message)

	sentryUmi.SentryLog(c, data, message)

}

func Error(c echo.Context, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Error(message)
	sentryUmi.SentryLog(c, data, message)

}

func Fatal(c echo.Context, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Fatal(message)
	sentryUmi.SentryLog(c, data, message)

}

func Debug(c echo.Context, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Debug(message)
	sentryUmi.SentryLog(c, data, message)

}

func Panic(c echo.Context, data logDump.Fields, message interface{}) {

	logDump.WithFields(data).Panic(message)
	sentryUmi.SentryLog(c, data, message)

}
