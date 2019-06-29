/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmd

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr-http2/cli/vxconf"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

// ESC CSI escape code
const ESC = 27

var clear = fmt.Sprintf("%c[%dA%c[2K", ESC, 1, ESC)

func clearLines(lineCount int) {
	_, _ = fmt.Fprint(os.Stdout, strings.Repeat(clear, lineCount))
}

func modifier(daemonServerCommands *cmdr.Command) *cmdr.Command {
	if startCmd := cmdr.FindSubCommand("start", daemonServerCommands); startCmd != nil {
		startCmd.PreAction = onServerPreStart
		startCmd.PostAction = onServerPostStop
	}

	return daemonServerCommands
}

func onAppStart(cmd *cmdr.Command, args []string) (err error) {
	logrus.Debug("onAppStart")
	return
}

func onAppExit(cmd *cmdr.Command, args []string) {
	logrus.Debug("onAppExit")
}

func onServerPostStop(cmd *cmdr.Command, args []string) {
	logrus.Debug("onServerPostStop")
}

// onServerPreStart is earlier than onAppStart.
func onServerPreStart(cmd *cmdr.Command, args []string) (err error) {
	earlierInitLogger()
	logrus.Debug("onServerPreStart")
	return
}

func earlierInitLogger() {
	l := "OFF"
	if !vxconf.IsProd() {
		l = "DEBUG"
	}
	l = vxconf.GetStringR("server.logger.level", l)
	logrus.SetLevel(stringToLevel(l))
	if l == "OFF" {
		logrus.SetOutput(ioutil.Discard)
	}
}

func stringToLevel(s string) logrus.Level {
	s = strings.ToUpper(s)
	switch s {
	case "TRACE":
		return logrus.TraceLevel
	case "DEBUG", "devel", "dev":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "FATAL":
		return logrus.FatalLevel
	case "PANIC":
		return logrus.PanicLevel
	default:
		return logrus.FatalLevel
	}
}
