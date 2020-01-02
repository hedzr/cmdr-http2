// Copyright © 2019 Hedzr Yeh.

package sample

import "github.com/hedzr/cmdr"

func WithSampleCmdrOption() cmdr.ExecOption {
	return func(w *cmdr.ExecWorker) {
		// daemonImpl = daemonImplX

		w.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {
			// the following statements show you how to attach an option to a sub-command
			serverCmd := cmdr.FindSubCommandRecursive("server", nil)
			serverStartCmd := cmdr.FindSubCommand("start", serverCmd)
			cmdr.NewInt(5100).
				Titles("vnc", "vnc-server").
				Description("start as a vnc server (just a demo)", "").
				Placeholder("PORT").
				AttachTo(cmdr.NewCmdFrom(serverStartCmd))
		})
	}
}
