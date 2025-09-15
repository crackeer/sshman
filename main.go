package main

import (
	"sshman/command/load"
	"sshman/command/run"
	"sshman/command/server"
	"sshman/command/shell"
	"sshman/command/upload"
	"sshman/define"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "sshman",
	Run: runHelp,
}

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&define.GServers, "server", "s", []string{}, "--server=127.0.0.1 -s 127.0.0.1")
	rootCmd.PersistentFlags().StringVarP(&define.GGroup, "group", "g", "", "--group=xxx -g xxx")

}

func main() {
	rootCmd.AddCommand(server.NewServerCommand("server"))
	rootCmd.AddCommand(upload.NewUploadCommand())
	rootCmd.AddCommand(load.NewImportCommand("k3s"))
	rootCmd.AddCommand(load.NewImportCommand("docker"))
	rootCmd.AddCommand(run.NewRunCommand())
	rootCmd.AddCommand(shell.NewShellCommand("bash"))
	rootCmd.AddCommand(shell.NewShellCommand("sh"))
	rootCmd.Execute()

}

func runHelp(cmd *cobra.Command, args []string) {
	//cmd.Help()
	if len(args) == 0 {
		cmd.Help()
		return
	}
}
