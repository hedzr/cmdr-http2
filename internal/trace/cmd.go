// Copyright © 2019 Hedzr Yeh.

package trace

import "github.com/hedzr/cmdr"

func WithTraceEnable(enabled bool) cmdr.ExecOption {
	return func(w *cmdr.ExecWorker) {
		// daemonImpl = daemonImplX

		w.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {

			if enabled {
				// attaches `--trace` to root command
				cmdr.NewBool(false).
					Titles("tr", "trace").
					Description("enable trace mode for tcp/mqtt send/recv data dump", "").
					AttachToRoot(root)
			}
			
		})
	}
}
