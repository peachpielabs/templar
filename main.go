/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package main

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/peachpielabs/gitformer/cmd/gitformer"
	"github.com/peachpielabs/gitformer/pkg/playbook"
)

var version = "No release version provided"
var dsn = ""
var environment = ""
var debug = false

func main() {
	if dsn == "" {
		log.Fatal("DSN was not set during the build process.")
	}
	if environment == "" {
		log.Fatal("environment was not set during the build process.")
	}
	if environment == "development" {
		debug = true
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		AttachStacktrace: true,
		Debug:            debug,
		SampleRate:       1.0,
		Environment:      environment,
		Release:          version,
	})
	if err != nil {
		playbook.CaptureError(err)
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	// sentry.CaptureMessage("It works!")
	defer sentry.Recover()

	gitformer.Execute()
}
