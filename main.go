package main

import (
	"sshman/command/server"
	"sshman/command/upload"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "sshman",
		Run: runHelp,
	}

	rootCmd.AddCommand(server.NewServerCommand("server"))
	rootCmd.AddCommand(upload.NewUploadCommand())
	rootCmd.Execute()

}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
