package sentry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MrAndreID/gohelpers"
	"github.com/getsentry/sentry-go"
	goSentry "github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	logDump "github.com/sirupsen/logrus"
)

var (
	userId = ""
)

func New(config goSentry.ClientOptions) {
	goSentry.Init(goSentry.ClientOptions{
		Dsn:        config.Dsn,
		HTTPProxy:  config.HTTPProxy,
		HTTPSProxy: config.HTTPSProxy,
		TracesSampler: goSentry.TracesSamplerFunc(func(ctx goSentry.SamplingContext) goSentry.Sampled {
			return goSentry.SampledTrue
		}),
		BeforeSend: func(event *goSentry.Event, hint *goSentry.EventHint) *goSentry.Event {
			if hint.Context != nil {
				if req, ok := hint.Context.Value(goSentry.RequestContextKey).(*http.Request); ok {
					// You have access to the original Request
					logDump.Info(req)
				}
			}
			logDump.Info(event)
			return event
		},
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
				userId = fmt.Sprintf("%v", fmt.Sprintf("%v", ctx.Get("UserId")))
			}

			hub.Scope().SetUser(sentry.User{
				ID: userId,
				//	IPAddress: m.GetLocalIP(),
			})
			hub.Scope().SetRequest(ctx.Request())

			// hub.ConfigureScope(func(scope *sentry.Scope) {
			// 	scope.SetLevel(sentry.LevelError)
			// 	scope.SetUser(sentry.User{
			// 		ID: userId,
			// 		//IPAddress: h.GetLocalIP(),
			// 	})
			// 	//	scope.AddBreadcrumb(&sentry.Breadcrumb{Message: "auth", Type: "info", Level: "error", Data: map[string]interface{}{"response": data}}, 1)
			// })

			sentry.Logger.SetFlags(time.Now().Minute())
			sentry.Logger.SetPrefix("[sentry SDK]")

		}

		return next(ctx)
	}

}

func SentryLog(c echo.Context, data logDump.Fields, message interface{}) {
	if c != nil {
		if c.Get("UserId") != nil {
			userId = fmt.Sprintf("%v", fmt.Sprintf("%v", c.Get("UserId")))
		} else {
			userId = fmt.Sprintf("%v", c.Get("RequestID"))
		}

		if hub := sentryecho.GetHubFromContext(c); hub != nil {
			hub := sentryecho.GetHubFromContext(c)

			// dataBreadcumb := breadcumb

			// dataBreadcumb.Data = data
			// dataBreadcumb.Message = fmt.Sprintf("%v", message)

			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetLevel(sentry.LevelError)
				scope.SetUser(sentry.User{
					ID: userId,
					//IPAddress: h.GetLocalIP(),
				})
				//	scope.AddBreadcrumb(&sentry.Breadcrumb{Message: "auth", Type: "info", Level: "error", Data: map[string]interface{}{"response": data}}, 1)
			})

			hub.CaptureMessage(string(gohelpers.JSONEncode(data)))
			//sentry.AddBreadcrumb(&dataBreadcumb)
		}
	}

}
