/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmd

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr-http2/cli/server"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/sirupsen/logrus"
)

// Entry is app main entry
func Entry() {

	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	// logex.EnableWith(logrus.DebugLevel)

	if err := cmdr.Exec(buildRootCmd(),
		// To disable internal commands and flags, uncomment the following codes
		// cmdr.WithBuiltinCommands(false, false, false, false, false),

		daemon.WithDaemon(server.NewDaemon(), modifier, onAppStart, onAppExit),
		
		server.WithHook(),
		
		cmdr.WithLogex(logrus.DebugLevel),
		cmdr.WithLogexPrefix("logger"),

		cmdr.WithWatchMainConfigFileToo(true),
		cmdr.WithOptionMergeModifying(func(keyPath string, value, oldVal interface{}) {
			logrus.Infof("%%-> -> %q: %v -> %v", keyPath, oldVal, value)
		}),
		cmdr.WithOptionModifying(func(keyPath string, value, oldVal interface{}) {
			logrus.Infof("%%-> -> %q: %v -> %v", keyPath, oldVal, value)
		}),

		cmdr.WithHelpTabStop(50),
	); err != nil {
		logrus.Errorf("Error: %v", err)
	}

}
