package lib

import (
	"github.com/mattn/go-colorable"
	"github.com/op/go-logging"
)

var format = logging.MustStringFormatter(
	"%{time:15:04:05.000} %{color}%{level:.4s}%{color:reset} %{message} @ %{shortfunc}",
)

func SetupLogger() {
	backend := logging.NewLogBackend(colorable.NewColorableStderr(), "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}
