package server

import (
	"fmt"

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
		serverConfig.DeleteByHost(server)
	}

	fmt.Println("delete success")
}
