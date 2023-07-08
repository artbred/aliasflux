package common

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strconv"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			return "", fmt.Sprintf("%s", fileName)
		},
	})
}
