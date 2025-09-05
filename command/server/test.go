package server

import (
	"fmt"
	"sshman/define"
	"sshman/service"

	"github.com/spf13/cobra"
)

func newServerTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "test server",
		Args:  cobra.MinimumNArgs(1),
		Run:   test,
	}

	return cmd
}

func test(cmd *cobra.Command, args []string) {
	for _, server := range args {
		fmt.Println("--> test server:", server)
		index := define.GServerConfig.Find(server)
		if index < 0 {
			fmt.Println("server not found")
			return
		}
		config := define.GServerConfig.Get(index)
		client, err := service.NewSSHClient(config)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("connect success.")
		client.Close()
	}
}
