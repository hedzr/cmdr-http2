// Copyright © 2020 Hedzr Yeh.

package cmd

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr-http2/cli/server"
	"github.com/hedzr/cmdr/tool"
	"github.com/sirupsen/logrus"
)

func buildRootCmd() (rootCmd *cmdr.RootCommand) {

	// var cmd *Command

	// cmdr.Root("aa", "1.0.1").
	// 	Header("sds").
	// 	NewSubCommand().
	// 	Titles("microservice", "ms").
	// 	Description("", "").
	// 	Group("").
	// 	Action(func(cmd *cmdr.Command, args []string) (err error) {
	// 		return
	// 	})

	// root

	root := cmdr.Root(appName, "1.1.1").
		// Header("cmdr-http2 - An HTTP2 server - no version - hedzr").
		Copyright("cmdr-http2 - An HTTP2 server", "Hedzr").
		Description(desc, longDesc).
		Examples(examples)
	rootCmd = root.RootCommand()

	server.AttachToCmdr(root)

	soundex(root)
	xy(root)
	mxTest(root)
	kv(root)
	ms(root)

	// // http 2 client
	//
	// root.NewSubCommand().
	// 	Titles("h2-test", "h2").
	// 	Description("test http 2 client", "test http 2 client,\nverbose long descriptions here.").
	// 	Group("Test").
	// 	Action(func(cmd *cmdr.Command, args []string) (err error) {
	// 		server.RunClient()
	// 		return
	// 	})

	//
	//

	// server.OnBuildCmd(rootCmd)

	return
}

func soundex(root cmdr.OptCmd) {
	// soundex

	root.NewSubCommand("soundex", "snd", "sndx", "sound").
		Description("soundex test").
		Group("Test").
		TailPlaceholder("[text1, text2, ...]").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			for ix, s := range args {
				fmt.Printf("%5d. %s => %s\n", ix, s, tool.Soundex(s))
			}
			return
		})
}

func xy(root cmdr.OptCmd) {
	// xy-print

	root.NewSubCommand().
		Titles("xy-print", "xy").
		Description("test terminal control sequences", "test terminal control sequences,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			//
			// https://en.wikipedia.org/wiki/ANSI_escape_code
			// https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97
			// https://en.wikipedia.org/wiki/POSIX_terminal_interface
			//

			fmt.Println("\x1b[2J") // clear screen

			for i, s := range args {
				fmt.Printf("\x1b[s\x1b[%d;%dH%s\x1b[u", 15+i, 30, s)
			}

			return
		})
}

func mxTest(root cmdr.OptCmd) {
	// mx-test

	mx := root.NewSubCommand().
		Titles("mx-test", "mx").
		Description("test new features", "test new features,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			fmt.Printf("*** Got pp: %s\n", cmdr.GetString("app.mx-test.password"))
			fmt.Printf("*** Got msg: %s\n", cmdr.GetString("app.mx-test.message"))
			return
		})
	mx.NewFlagV("").
		Titles("password", "pp").
		Description("the password requesting.", "").
		Group("").
		Placeholder("PASSWORD").
		ExternalTool(cmdr.ExternalToolPasswordInput)
	mx.NewFlagV("").
		Titles("message", "m", "msg").
		Description("the message requesting.", "").
		Group("").
		Placeholder("MESG").
		ExternalTool(cmdr.ExternalToolEditor)
}

func kv(root cmdr.OptCmd) {
	// kv

	kvCmd := root.NewSubCommand().
		Titles("kvstore", "kv").
		Description("consul kv store operations...", ``)

	attachConsulConnectFlags(kvCmd)

	kvBackupCmd := kvCmd.NewSubCommand().
		Titles("backup", "b", "bk", "bf", "bkp").
		Description("Dump Consul's KV database to a JSON/YAML file", ``).
		Action(kvBackup)
	kvBackupCmd.NewFlagV("consul-backup.json").
		Titles("output", "o").
		Description("Write output to a file (*.json / *.yml)", ``).
		Placeholder("FILE")

	kvRestoreCmd := kvCmd.NewSubCommand().
		Titles("restore", "r").
		Description("restore to Consul's KV store, from a a JSON/YAML backup file", ``).
		Action(kvRestore)
	kvRestoreCmd.NewFlagV("consul-backup.json").
		Titles("input", "i").
		Description("Read the input file (*.json / *.yml)", ``).
		Placeholder("FILE")
}

func ms(root cmdr.OptCmd) { // ms

	msCmd := root.NewSubCommand().
		Titles("micro-service", "ms", "microservice").
		Description("micro-service operations... [fake]", "").
		Group("")

	msCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("name", "n").
		Description("name of the service", ``).
		DefaultValue("", "NAME").
		Group("service spec")
	msCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("id", "i", "ID").
		Description("unique id of the service", ``).
		DefaultValue("", "ID").
		Group("service spec")
	msCmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("all", "a").
		Description("all services", ``).
		DefaultValue(false, "").
		Group("service spec")

	msCmd.NewFlag(cmdr.OptFlagTypeUint).
		Titles("retry", "t").
		Description("", "").
		Group("").
		DefaultValue(3, "RETRY").
		Group("Z tests")
	msCmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("money", "mm").
		Description("A placeholder flag.", "").
		DefaultValue(false, "").
		Group("Z tests")

	// ms ls

	msCmd.NewSubCommand().
		Titles("list", "ls", "l", "lst", "dir").
		Description("list tags [fake]", "").
		Group("2333.List").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			logrus.Error("implement me")
			return
		})

	// ms tags

	msTagsCmd := msCmd.NewSubCommand().
		Titles("tags", "t").
		Description("tags operations of a micro-service [fake]", "").
		Group("")

	attachConsulConnectFlags(msTagsCmd)

	// ms tags ls

	msTagsCmd.NewSubCommand().
		Titles("list", "ls", "l", "lst", "dir").
		Description("list tags [fake]", "").
		Group("2333.List").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	// ms tags add

	tagsAdd := msTagsCmd.NewSubCommand().
		Titles("add", "a", "new", "create").
		Description("add tags [fake]", "").
		Deprecated("0.2.1").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	tagsAdd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("list", "ls", "l", "lst", "dir").
		Description("a comma list to be added", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	c1 := tagsAdd.NewSubCommand().
		Titles("check", "c", "chk").
		Description("[sub] check [fake]", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	c2 := c1.NewSubCommand().
		Titles("check-point", "pt", "chk-pt").
		Description("[sub][sub] checkpoint [fake]", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	c2.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("add", "a", "add-list").
		Description("a comma list to be added.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")
	c2.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("remove", "r", "rm-list", "rm", "del", "delete").
		Description("a comma list to be removed.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")

	c3 := c1.NewSubCommand().
		Titles("check-in", "in", "chk-in").
		Description("[sub][sub] check-in [fake]", "").
		Group("")

	c3.NewFlag(cmdr.OptFlagTypeString).
		Titles("name", "n").
		Description("a string to be added.", ``).
		DefaultValue("", "")

	c3.NewSubCommand().
		Titles("demo-1", "d1").
		Description("[sub][sub] check-in sub", "").
		Group("")

	c3.NewSubCommand().
		Titles("demo-2", "d2").
		Description("[sub][sub] check-in sub", "").
		Group("")

	c3.NewSubCommand().
		Titles("demo-3", "d3").
		Description("[sub][sub] check-in sub", "").
		Group("")

	c1.NewSubCommand().
		Titles("check-out", "out", "chk-out").
		Description("[sub][sub] check-out [fake]", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	// ms tags rm

	tagsRm := msTagsCmd.NewSubCommand().
		Titles("rm", "r", "remove", "delete", "del", "erase").
		Description("remove tags [fake]", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	tagsRm.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("list", "ls", "l", "lst", "dir").
		Description("a comma list to be added", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	// ms tags modify

	msTagsModifyCmd := msTagsCmd.NewSubCommand().
		Titles("modify", "m", "mod", "modi", "update", "change").
		Description("modify tags of a service. [fake]", ``).
		Action(msTagsModify)

	attachModifyFlags(msTagsModifyCmd)

	msTagsModifyCmd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("add", "a", "add-list").
		Description("a comma list to be added.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")
	msTagsModifyCmd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("remove", "r", "rm-list", "rm", "del", "delete").
		Description("a comma list to be removed.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")

	// ms tags toggle

	tagsTog := msTagsCmd.NewSubCommand().
		Titles("toggle", "t", "tog", "switch").
		Description("toggle tags [fake]", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	attachModifyFlags(tagsTog)

	tagsTog.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("set", "s").
		Description("a comma list to be set", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	tagsTog.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("unset", "u", "un").
		Description("a comma list to be unset", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	tagsTog.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("address", "a", "addr").
		Description("the address of the service (by id or name)", ``).
		DefaultValue("", "HOST:PORT")
}

func attachModifyFlags(cmd cmdr.OptCmd) {
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("delim", "d").
		Description("delimitor char in `non-plain` mode.", ``).
		DefaultValue("=", "")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("clear", "c").
		Description("clear all tags.", ``).
		DefaultValue(false, "").
		Group("Operate")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("string", "g", "string-mode").
		Description("In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("meta", "m", "meta-mode").
		Description("In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false).", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("both", "2", "both-mode").
		Description("In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("plain", "p", "plain-mode").
		Description("In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("tag", "t", "tag-mode").
		Description("In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		DefaultValue(true, "").
		Group("Mode")

}

func attachConsulConnectFlags(cmd cmdr.OptCmd) {
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("addr", "a").
		Description("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')", ``).
		DefaultValue("localhost", "HOST[:PORT]").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeInt).
		Titles("port", "p").
		Description("Consul port", ``).
		DefaultValue(8500, "PORT").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("insecure", "K").
		Description("Skip TLS host verification", ``).
		DefaultValue(true, "").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("prefix", "px").
		Description("Root key prefix", ``).
		DefaultValue("/", "ROOT").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("cacert", "ca").
		Description("Consul Client CA cert)", ``).
		DefaultValue("", "FILE").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("cert", "").
		Description("Consul Client cert", ``).
		DefaultValue("", "FILE").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("scheme", "").
		Description("Consul connection protocol", ``).
		DefaultValue("http", "SCHEME").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("username", "u", "user", "usr", "uid").
		Description("HTTP Basic auth user", ``).
		DefaultValue("", "USERNAME").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("password", "pw", "passwd", "pass", "pwd").
		Description("HTTP Basic auth password", ``).
		DefaultValue("", "PASSWORD").
		Group("Consul").
		ExternalTool(cmdr.ExternalToolPasswordInput)

}

const (
	appName   = "cmdr-http2"
	copyright = "cmdr-http2 is an demo app"
	desc      = "cmdr-http2 is an demo app. It make an demo application for `cmdr`."
	longDesc  = "cmdr-http2 is an demo app. It make an demo application for `cmdr`."
	examples  = `
$ {{.AppName}} gen shell [--bash|--zsh|--auto]
  generate bash/shell completion scripts
$ {{.AppName}} gen man
  generate linux man page 1
$ {{.AppName}} --help
  show help screen.
`
	overview = ``
)
