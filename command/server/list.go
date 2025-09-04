package server

import (
	"fmt"
	"os"
	"sshman/define"
	"sshman/service"

	"github.com/spf13/cobra"
	"github.com/tomlazar/table"
)

var (
	serverConfig *service.ServerConfig
)

func init() {
	var err error
	serverConfig, err = service.NewServerConfig()
	if err != nil {
		fmt.Println("init server config error:", err)
		os.Exit(1)
	}
}

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
	return cmd
}

func List(cmd *cobra.Command, args []string) {
	list := serverConfig.List()
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
