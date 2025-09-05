package server

import (
	"fmt"
	"os"
	"sshman/define"

	"github.com/spf13/cobra"
	"github.com/tomlazar/table"
)

// NewServerCommand
//
//	@param use
//	@return *cobra.Command
func NewServerCommand(use string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: "list server manage",
		Run:   List,
	}
	cmd.AddCommand(newServerCreateCommand())
	cmd.AddCommand(newServerCopyCommand())
	cmd.AddCommand(newServerDeleteCommand())
	cmd.AddCommand(newServerTestCommand())
	cmd.AddCommand(newServerGroupCommand())
	return cmd
}

func List(cmd *cobra.Command, args []string) {
	list := define.GServerConfig.List()
	tableData := []define.Server{}
	for _, server := range list {
		tableData = append(tableData, *server)
	}
	buf, err := table.Marshal(tableData, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(buf))
}
