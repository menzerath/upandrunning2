package lib

import (
	"github.com/op/go-logging"
	"os"
)

// The application-wide used logging-format.
var format = logging.MustStringFormatter(
	"%{time:15:04:05} %{color}%{level:.4s}%{color:reset} %{message} @ %{shortfunc}",
)

var dockerFormat = logging.MustStringFormatter(
	"%{time:15:04:05} %{level:.4s} %{message} @ %{shortfunc}",
)

// Init the logger.
func SetupLogger() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	if os.Getenv("UAR2_IS_DOCKER") == "true" {
		backendFormatter = logging.NewBackendFormatter(backend, dockerFormat)
	}
	logging.SetBackend(backendFormatter)
}
