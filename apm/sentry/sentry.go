package sentry

import (
	goSentry "github.com/getsentry/sentry-go"
)

func New() {
	goSentry.Init(goSentry.ClientOptions{
		Dsn: "your-public-dsn",
	})
}
