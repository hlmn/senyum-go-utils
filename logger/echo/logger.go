package echo

import (
	"io"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryUmi "github.com/hlmn/senyum-go-utils/apm/sentry"

	"github.com/labstack/echo/v4"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logDump "github.com/sirupsen/logrus"
)

func InitBodyDumpLog() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	logf, err := rotatelogs.New(
		dir+"/logs/RequestResponseDump.log.%Y%m%d",
		rotatelogs.WithLinkName(dir+"/logs/RequestResponseDump.log"),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(365),
	)
	if err != nil {
		return err
	}

	logDump.SetFormatter(&logDump.JSONFormatter{DisableHTMLEscape: true})
	logDump.SetOutput(io.MultiWriter(os.Stdout, logf))
	logDump.SetLevel(logDump.InfoLevel)
	logDump.SetReportCaller(true)

	return nil
}

func Info(c echo.Context, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Info(message)

	sentryUmi.SentryLog(c, data, message, sentry.LevelInfo)

}

func Error(c echo.Context, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Error(message)
	sentryUmi.SentryLog(c, data, message, sentry.LevelError)

}

func Fatal(c echo.Context, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Fatal(message)
	sentryUmi.SentryLog(c, data, message, sentry.LevelFatal)

}

func Debug(c echo.Context, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Debug(message)
	sentryUmi.SentryLog(c, data, message, sentry.LevelDebug)

}

func Panic(c echo.Context, data logDump.Fields, message interface{}) {

	logDump.WithFields(data).Panic(message)
	sentryUmi.SentryLog(c, data, message, sentry.LevelError)

}

func Warning(c echo.Context, data logDump.Fields, message interface{}) {

	logDump.WithFields(data).Warning(message)
	sentryUmi.SentryLog(c, data, message, sentry.LevelWarning)

}
