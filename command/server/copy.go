package server

import (
	"fmt"
	"sshman/define"

	"github.com/spf13/cobra"
)

func newServerCopyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "copy",
		Short: "copy host aaa",
		Args:  cobra.MinimumNArgs(2),
		Run:   copy,
	}

	return cmd
}

func copy(cmd *cobra.Command, args []string) {
	host := args[0]
	index := serverConfig.Find(host)
	if index < 0 {
		fmt.Println("server not found")
		return
	}
	theOne := serverConfig.Get(index)

	for _, server := range args[1:] {
		serverConfig.Add(&define.Server{
			User:     theOne.User,
			Password: theOne.Password,
			Host:     server,
			Port:     theOne.Port,
			Group:    theOne.Group,
		})
	}

	fmt.Println("copy success")
}
