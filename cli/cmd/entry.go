// Copyright © 2020 Hedzr Yeh.

package cmd

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr-http2/cli/server"
	"github.com/hedzr/cmdr-http2/internal/shell"
	"github.com/hedzr/cmdr-http2/internal/trace"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

// Entry is app main entry
func Entry() {

	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	// logex.EnableWith(logrus.DebugLevel)

	if err := cmdr.Exec(buildRootCmd(),
		// To disable internal commands and flags, uncomment the following codes
		// cmdr.WithBuiltinCommands(false, false, false, false, false),

		// daemon.WithDaemon(server.NewDaemon(), modifier, onAppStart, onAppExit),
		server.WithCmdrDaemonSupport(),
		server.WithCmdrHook(),

		cmdr.WithLogex(cmdr.DebugLevel),
		cmdr.WithLogexPrefix("logger"),

		cmdr.WithWatchMainConfigFileToo(true),
		cmdr.WithNoWatchConfigFiles(false),
		cmdr.WithOptionMergeModifying(func(keyPath string, value, oldVal interface{}) {
			logrus.Debugf("%%-> -> %q: %v -> %v", keyPath, oldVal, value)
			if strings.HasSuffix(keyPath, ".mqtt.server.stats.enabled") {
				// mqttlib.FindServer().EnableSysStats(!vxconf.ToBool(value))
			}
			if strings.HasSuffix(keyPath, ".mqtt.server.stats.log.enabled") {
				// mqttlib.FindServer().EnableSysStatsLog(!vxconf.ToBool(value))
			}
		}),
		cmdr.WithOptionModifying(func(keyPath string, value, oldVal interface{}) {
			logrus.Infof("%%-> -> %q: %v -> %v", keyPath, oldVal, value)
		}),

		cmdr.WithHelpTabStop(43),

		// sample.WithSampleCmdrOption(),
		trace.WithTraceEnable(true),

		cmdr.WithUnknownOptionHandler(onUnknownOptionHandler),
		cmdr.WithUnhandledErrorHandler(onUnhandledErrorHandler),

		shell.WithShellModule(),
	); err != nil {
		logrus.Errorf("Error: %v", err)
	}

}

func onUnknownOptionHandler(isFlag bool, title string, cmd *cmdr.Command, args []string) (fallbackToDefaultDetector bool) {
	return true
}

// to make the h2 server recoverable
func onUnhandledErrorHandler(err interface{}) {
	// debug.PrintStack()
	// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	if e, ok := err.(error); ok {
		logrus.Errorf("%+v", e)
	} else {
		logrus.Errorf("%+v", err)
		dumpStacks()
	}
}

func dumpStacks() {
	buf := make([]byte, 16384)
	buf = buf[:runtime.Stack(buf, false)]
	fmt.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===\n", buf)
}
