/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmd

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr-http2/cli/server"
	"github.com/sirupsen/logrus"
)

func buildRootCmd() (rootCmd *cmdr.RootCommand) {

	// var cmd *Command

	// cmdr.Root("aa", "1.0.1").
	// 	Header("sds").
	// 	NewSubCommand().
	// 	Titles("ms", "microservice").
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

	// xy-print

	root.NewSubCommand().
		Titles("xy", "xy-print").
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

	// mx-test

	mx := root.NewSubCommand().
		Titles("mx", "mx-test").
		Description("test new features", "test new features,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			fmt.Printf("*** Got pp: %s\n", cmdr.GetString("app.mx-test.password"))
			fmt.Printf("*** Got msg: %s\n", cmdr.GetString("app.mx-test.message"))
			return
		})
	mx.NewFlagV("").
		Titles("pp", "password").
		Description("the password requesting.", "").
		Group("").
		Placeholder("PASSWORD").
		ExternalTool(cmdr.ExternalToolPasswordInput)
	mx.NewFlagV("").
		Titles("m", "message", "msg").
		Description("the message requesting.", "").
		Group("").
		Placeholder("MESG").
		ExternalTool(cmdr.ExternalToolEditor)

	// http 2 client

	root.NewSubCommand().
		Titles("h2", "h2-test").
		Description("test http 2 client", "test http 2 client,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			server.RunClient()
			return
		})

	// kv

	kvCmd := root.NewSubCommand().
		Titles("kv", "kvstore").
		Description("consul kv store operations...", ``)

	attachConsulConnectFlags(kvCmd)

	kvBackupCmd := kvCmd.NewSubCommand().
		Titles("b", "backup", "bk", "bf", "bkp").
		Description("Dump Consul's KV database to a JSON/YAML file", ``).
		Action(kvBackup)
	kvBackupCmd.NewFlagV("consul-backup.json").
		Titles("o", "output").
		Description("Write output to a file (*.json / *.yml)", ``).
		Placeholder("FILE")

	kvRestoreCmd := kvCmd.NewSubCommand().
		Titles("r", "restore").
		Description("restore to Consul's KV store, from a a JSON/YAML backup file", ``).
		Action(kvRestore)
	kvRestoreCmd.NewFlagV("consul-backup.json").
		Titles("i", "input").
		Description("Read the input file (*.json / *.yml)", ``).
		Placeholder("FILE")

	// ms

	msCmd := root.NewSubCommand().
		Titles("ms", "micro-service", "microservice").
		Description("micro-service operations... [fake]", "").
		Group("")

	msCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("n", "name").
		Description("name of the service", ``).
		DefaultValue("", "NAME").
		Group("service spec")
	msCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("i", "id", "ID").
		Description("unique id of the service", ``).
		DefaultValue("", "ID").
		Group("service spec")
	msCmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("a", "all").
		Description("all services", ``).
		DefaultValue(false, "").
		Group("service spec")

	msCmd.NewFlag(cmdr.OptFlagTypeUint).
		Titles("t", "retry").
		Description("", "").
		Group("").
		DefaultValue(3, "RETRY").
		Group("Z tests")
	msCmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("mm", "money").
		Description("A placeholder flag.", "").
		DefaultValue(false, "").
		Group("Z tests")

	// ms ls

	msCmd.NewSubCommand().
		Titles("ls", "list", "l", "lst", "dir").
		Description("list tags [fake]", "").
		Group("2333.List").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			logrus.Error("implement me")
			return
		})

	// ms tags

	msTagsCmd := msCmd.NewSubCommand().
		Titles("t", "tags").
		Description("tags operations of a micro-service [fake]", "").
		Group("")

	// cTags.NewFlag(cmdr.OptFlagTypeString).
	// 	Titles("n", "name").
	// 	Description("name of the service", "").
	// 	Group("").
	// 	DefaultValue("", "NAME")
	//
	// cTags.NewFlag(cmdr.OptFlagTypeString).
	// 	Titles("i", "id").
	// 	Description("unique id of the service", "").
	// 	Group("").
	// 	DefaultValue("", "ID")
	//
	// cTags.NewFlag(cmdr.OptFlagTypeString).
	// 	Titles("a", "addr").
	// 	Description("", "").
	// 	Group("").
	// 	DefaultValue("consul.ops.local", "ADDR")

	attachConsulConnectFlags(msTagsCmd)

	// ms tags ls

	msTagsCmd.NewSubCommand().
		Titles("ls", "list", "l", "lst", "dir").
		Description("list tags [fake]", "").
		Group("2333.List").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	// ms tags add

	tagsAdd := msTagsCmd.NewSubCommand().
		Titles("a", "add", "new", "create").
		Description("add tags [fake]", "").
		Deprecated("0.2.1").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	tagsAdd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("ls", "list", "l", "lst", "dir").
		Description("a comma list to be added", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	c1 := tagsAdd.NewSubCommand().
		Titles("c", "check", "chk").
		Description("[sub] check [fake]", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	c2 := c1.NewSubCommand().
		Titles("pt", "check-point", "chk-pt").
		Description("[sub][sub] checkpoint [fake]", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	c2.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("a", "add", "add-list").
		Description("a comma list to be added.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")
	c2.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("r", "remove", "rm-list", "rm", "del", "delete").
		Description("a comma list to be removed.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")

	c3 := c1.NewSubCommand().
		Titles("in", "check-in", "chk-in").
		Description("[sub][sub] check-in [fake]", "").
		Group("")

	c3.NewFlag(cmdr.OptFlagTypeString).
		Titles("n", "name").
		Description("a string to be added.", ``).
		DefaultValue("", "")

	c3.NewSubCommand().
		Titles("d1", "demo-1").
		Description("[sub][sub] check-in sub", "").
		Group("")

	c3.NewSubCommand().
		Titles("d2", "demo-2").
		Description("[sub][sub] check-in sub", "").
		Group("")

	c3.NewSubCommand().
		Titles("d3", "demo-3").
		Description("[sub][sub] check-in sub", "").
		Group("")

	c1.NewSubCommand().
		Titles("out", "check-out", "chk-out").
		Description("[sub][sub] check-out [fake]", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	// ms tags rm

	tagsRm := msTagsCmd.NewSubCommand().
		Titles("r", "rm", "remove", "delete", "del", "erase").
		Description("remove tags [fake]", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	tagsRm.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("ls", "list", "l", "lst", "dir").
		Description("a comma list to be added", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	// ms tags modify

	msTagsModifyCmd := msTagsCmd.NewSubCommand().
		Titles("m", "modify", "mod", "modi", "update", "change").
		Description("modify tags of a service. [fake]", ``).
		Action(msTagsModify)

	attachModifyFlags(msTagsModifyCmd)

	msTagsModifyCmd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("a", "add", "add-list").
		Description("a comma list to be added.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")
	msTagsModifyCmd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("r", "remove", "rm-list", "rm", "del", "delete").
		Description("a comma list to be removed.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")

	// ms tags toggle

	tagsTog := msTagsCmd.NewSubCommand().
		Titles("t", "toggle", "tog", "switch").
		Description("toggle tags [fake]", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	attachModifyFlags(tagsTog)

	tagsTog.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("s", "set").
		Description("a comma list to be set", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	tagsTog.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("u", "unset", "un").
		Description("a comma list to be unset", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	tagsTog.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("a", "address", "addr").
		Description("the address of the service (by id or name)", ``).
		DefaultValue("", "HOST:PORT")

	//
	//

	// server.OnBuildCmd(rootCmd)

	return
}

func attachModifyFlags(cmd cmdr.OptCmd) {
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("d", "delim").
		Description("delimitor char in `non-plain` mode.", ``).
		DefaultValue("=", "")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("c", "clear").
		Description("clear all tags.", ``).
		DefaultValue(false, "").
		Group("Operate")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("g", "string", "string-mode").
		Description("In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("m", "meta", "meta-mode").
		Description("In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false).", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("2", "both", "both-mode").
		Description("In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("p", "plain", "plain-mode").
		Description("In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("t", "tag", "tag-mode").
		Description("In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		DefaultValue(true, "").
		Group("Mode")

}

func attachConsulConnectFlags(cmd cmdr.OptCmd) {
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("a", "addr").
		Description("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')", ``).
		DefaultValue("localhost", "HOST[:PORT]").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeInt).
		Titles("p", "port").
		Description("Consul port", ``).
		DefaultValue(8500, "PORT").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("K", "insecure").
		Description("Skip TLS host verification", ``).
		DefaultValue(true, "").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("px", "prefix").
		Description("Root key prefix", ``).
		DefaultValue("/", "ROOT").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("", "cacert").
		Description("Consul Client CA cert)", ``).
		DefaultValue("", "FILE").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("", "cert").
		Description("Consul Client cert", ``).
		DefaultValue("", "FILE").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("", "scheme").
		Description("Consul connection protocol", ``).
		DefaultValue("http", "SCHEME").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("u", "username", "user", "usr", "uid").
		Description("HTTP Basic auth user", ``).
		DefaultValue("", "USERNAME").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("pw", "password", "passwd", "pass", "pwd").
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
