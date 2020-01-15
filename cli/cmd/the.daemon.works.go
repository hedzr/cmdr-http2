// Copyright © 2020 Hedzr Yeh.

package cmd

import (
	"fmt"
	"os"
	"strings"
)

// ESC CSI escape code
const ESC = 27

var clear = fmt.Sprintf("%c[%dA%c[2K", ESC, 1, ESC)

func clearLines(lineCount int) {
	_, _ = fmt.Fprint(os.Stdout, strings.Repeat(clear, lineCount))
}

// earlierInitLogger is deprecated by cmdr.WithLogex()
// func earlierInitLogger() {
// 	l := "OFF"
// 	if !vxconf.IsProd() {
// 		l = "DEBUG"
// 	}
// 	l = vxconf.GetStringR("server.logger.level", l)
// 	logrus.SetLevel(stringToLevel(l))
// 	if l == "OFF" {
// 		logrus.SetOutput(ioutil.Discard)
// 	}
// }

// func stringToLevel(s string) logrus.Level {
// 	s = strings.ToUpper(s)
// 	switch s {
// 	case "TRACE":
// 		return logrus.TraceLevel
// 	case "DEBUG", "devel", "dev":
// 		return logrus.DebugLevel
// 	case "INFO":
// 		return logrus.InfoLevel
// 	case "WARN":
// 		return logrus.WarnLevel
// 	case "ERROR":
// 		return logrus.ErrorLevel
// 	case "FATAL":
// 		return logrus.FatalLevel
// 	case "PANIC":
// 		return logrus.PanicLevel
// 	default:
// 		return logrus.FatalLevel
// 	}
// }
