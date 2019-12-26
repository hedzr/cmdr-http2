/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmd

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr-http2/cli/server"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/sirupsen/logrus"
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

		daemon.WithDaemon(server.NewDaemon(), modifier, onAppStart, onAppExit),

		server.WithHook(),

		cmdr.WithLogex(logrus.DebugLevel),
		cmdr.WithLogexPrefix("logger"),

		cmdr.WithWatchMainConfigFileToo(true),
		cmdr.WithNoWatchConfigFiles(false),
		cmdr.WithOptionMergeModifying(func(keyPath string, value, oldVal interface{}) {
			logrus.Infof("%%-> -> %q: %v -> %v", keyPath, oldVal, value)
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

		cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
			// the following statements show you how to attach an option to a sub-command
			serverCmd := cmdr.FindSubCommandRecursive("server", nil)
			serverStartCmd := cmdr.FindSubCommand("start", serverCmd)
			cmdr.NewInt(5100).
				Titles("vnc", "vnc-server").
				Description("start as a vnc server (just a demo)", "").
				Placeholder("PORT").
				AttachTo(cmdr.NewCmdFrom(serverStartCmd))
		}, nil),

		cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
			// attaches `--trace` to root command
			cmdr.NewBool(false).
				Titles("tr", "trace").
				Description("enable trace mode for tcp/mqtt send/recv data dump", "").
				AttachToRoot(root)
		}, nil),

	); err != nil {
		logrus.Errorf("Error: %v", err)
	}

}
