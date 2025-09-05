package server

import (
	"fmt"
	"sshman/define"

	"github.com/spf13/cobra"
)

func newServerGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "group server",
		Args:  cobra.MinimumNArgs(2),
		Run:   handleGroup,
	}

	return cmd
}

func handleGroup(cmd *cobra.Command, args []string) {
	group := args[0]
	hosts := args[1:]

	err := define.GServerConfig.Group(group, hosts)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ok.")
}
