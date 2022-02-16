package sentry

import (
	"fmt"
	"time"

	"github.com/MrAndreID/gohelpers"
	"github.com/getsentry/sentry-go"
	goSentry "github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	logDump "github.com/sirupsen/logrus"
)

func New(config goSentry.ClientOptions) {
	goSentry.Init(goSentry.ClientOptions{
		Dsn:        config.Dsn,
		HTTPProxy:  config.HTTPProxy,
		HTTPSProxy: config.HTTPSProxy,
		TracesSampler: goSentry.TracesSamplerFunc(func(ctx goSentry.SamplingContext) goSentry.Sampled {
			return goSentry.SampledTrue
		}),

		Debug:            true,
		AttachStacktrace: true,
	})
	defer goSentry.Flush(2 * time.Second)
}

func MiddlewareSentry(next echo.HandlerFunc) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		span := sentry.StartSpan(ctx.Request().Context(), ctx.Path(), sentry.TransactionName(fmt.Sprintf("%s", ctx.Path())))

		defer span.Finish()

		if hub := sentryecho.GetHubFromContext(ctx); hub != nil {
			var (
				userId = fmt.Sprintf("%v", ctx.Get("RequestID"))
			)

			if ctx.Get("UserId") != nil {
				userId = fmt.Sprintf("%v", ctx.Get("UserId"))
			}

			if ctx.Request().Header.Get("umiUserId") != "" {
				userId = ctx.Request().Header.Get("umiUserId")
			} else if ctx.Request().Header.Get("requestId") != "" {
				userId = ctx.Request().Header.Get("requestId")
			}

			hub.Scope().SetUser(sentry.User{
				ID: userId,
			})
			hub.Scope().SetRequest(ctx.Request())
			sentry.Logger.SetFlags(time.Now().Minute())
			sentry.Logger.SetPrefix("[sentry SDK]")

		}

		return next(ctx)
	}

}

func SentryLog(c echo.Context, data logDump.Fields, message interface{}, level goSentry.Level) {
	userId := ""
	if c != nil {
		if c.Get("UserId") != nil {
			userId = fmt.Sprintf("%v", c.Get("UserId"))
		} else {
			userId = fmt.Sprintf("%v", c.Get("RequestID"))
		}

		if hub := sentryecho.GetHubFromContext(c); hub != nil {
			hub := sentryecho.GetHubFromContext(c)

			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetLevel(level)
				scope.SetUser(sentry.User{
					ID: userId,
				})
			})

			hub.CaptureMessage(string(gohelpers.JSONEncode(data)))
		}
	}

}
