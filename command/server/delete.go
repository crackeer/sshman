package server

import (
	"fmt"
	"sshman/define"

	"github.com/spf13/cobra"
)

func newServerDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete server",
		Args:  cobra.MinimumNArgs(1),
		Run:   del,
	}

	return cmd
}

func del(cmd *cobra.Command, args []string) {
	for _, server := range args {
		fmt.Println("--> delete server:", server)
		define.GServerConfig.DeleteByHost(server)
	}

	fmt.Println("delete success")
}
