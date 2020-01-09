// Copyright © 2020 Hedzr Yeh.

package server

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/sirupsen/logrus"
	"net"
)

// AttachToCmdr adds command-line commands and options to cmdr system
func AttachToCmdr(root cmdr.OptCmd) {

	// http 2 client

	root.NewSubCommand().
		Titles("h2", "h2-test").
		Description("test http 2 client", "test http 2 client,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			runClient()
			return
		})

}

// WithCmdrDaemonSupport adds daemon support to cmdr system
func WithCmdrDaemonSupport() cmdr.ExecOption {
	return daemon.WithDaemon(newDaemon(),
		modifier, onAppStart, onAppExit,
		daemon.WithOnGetListener(getOnGetListener),
	)
}

// WithCmdrHook adds hooking to cmdr system so that we can modify daemon `server` command at an appropriate time.
func WithCmdrHook() cmdr.ExecOption {
	return cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {

		// app.server.port
		if cx := root.Command.FindSubCommand("server"); cx != nil {
			cx.Description = "daemon service: HTTP2 server"
			// logrus.Debugf("`server` command found")
			opt := cmdr.NewCmdFrom(cx)
			if flg := cx.FindFlag("port"); flg != nil {
				flg.DefaultValue = defaultPort

			} else {
				opt.NewFlagV(defaultPort).
					Titles("p", "port").
					Description("the port to listen.", "").
					Group("").
					Placeholder("PORT")
			}
		}
	}, func(root *cmdr.RootCommand, args []string) {
		logrus.Debugf("cmd: root=%+v, args: %v", root, args)
	})
}

func modifier(daemonServerCommands *cmdr.Command) *cmdr.Command {
	if startCmd := daemonServerCommands.FindSubCommand("start"); startCmd != nil {
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
	// earlierInitLogger() // deprecated by cmdr.WithLogex()
	logrus.Debug("onServerPreStart")
	return
}

func getOnGetListener() net.Listener {
	// l := h2listener
	// h2listener = nil
	// return l
	return h2listener
}

var h2listener net.Listener
