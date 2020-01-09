// Copyright © 2019 Hedzr Yeh.

package sample

import "github.com/hedzr/cmdr"

// WithSampleCmdrOption show the example howto modify cmdr daemon plugin `server` `start` sub-command at an appropriate time.
func WithSampleCmdrOption() cmdr.ExecOption {
	return func(w *cmdr.ExecWorker) {
		// daemonImpl = daemonImplX

		w.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {
			// the following statements show you how to attach an option to a sub-command
			var serverCmd, serverStartCmd *cmdr.Command
			serverCmd = cmdr.FindSubCommandRecursive("server", nil)
			if serverCmd != nil {
				serverStartCmd = cmdr.FindSubCommand("start", serverCmd)
			}
			if serverStartCmd != nil {
				cmdr.NewInt(5100).
					Titles("vnc", "vnc-server").
					Description("start as a vnc server (just a demo)", "").
					Placeholder("PORT").
					AttachTo(cmdr.NewCmdFrom(serverStartCmd))
			}
		})
	}
}
