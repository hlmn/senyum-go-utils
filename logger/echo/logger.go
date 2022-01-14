package echo

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
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

func Info(c echo.Context, breadcumb sentry.Breadcrumb, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Info(message)

	SentryLog(c, breadcumb, data, fmt.Sprintf("%v", message))

}

func Error(c echo.Context, breadcumb sentry.Breadcrumb, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Error(message)
	SentryLog(c, breadcumb, data, fmt.Sprintf("%v", message))

}

func Fatal(c echo.Context, breadcumb sentry.Breadcrumb, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Fatal(message)
	SentryLog(c, breadcumb, data, fmt.Sprintf("%v", message))

}

func Debug(c echo.Context, breadcumb sentry.Breadcrumb, data logDump.Fields, message interface{}) {
	logDump.WithFields(data).Debug(message)
	SentryLog(c, breadcumb, data, fmt.Sprintf("%v", message))

}

func Panic(c echo.Context, breadcumb sentry.Breadcrumb, data logDump.Fields, message interface{}) {

	logDump.WithFields(data).Panic(message)
	SentryLog(c, breadcumb, data, fmt.Sprintf("%v", message))

}

func SentryLog(c echo.Context, breadcumb sentry.Breadcrumb, data logDump.Fields, message string) {
	if c != nil {
		if c.Get("UserId") != nil {
			userId = fmt.Sprintf("%v", fmt.Sprintf("%v", c.Get("UserId")))
		} else {
			userId = fmt.Sprintf("%v", c.Get("RequestID"))
		}

		if hub := sentryecho.GetHubFromContext(c); hub != nil {
			hub := sentryecho.GetHubFromContext(c)

			dataBreadcumb := breadcumb

			dataBreadcumb.Data = data
			dataBreadcumb.Message = message

			hub.CaptureMessage(message)
			hub.AddBreadcrumb(&dataBreadcumb, &sentry.BreadcrumbHint{})
		}
	}

}
