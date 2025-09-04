package server

import (
	"fmt"

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
	index := serverConfig.Find(args[0])

	if index < 0 {
		fmt.Println("server not found")
		return
	}
	theOne := serverConfig.Get(index)
	fmt.Println(theOne)
}
