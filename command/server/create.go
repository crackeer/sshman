package server

import (
	"fmt"
	"sshman/define"

	prompt "github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

var (
	user     string
	password string
	port     string = "22"
	group    string
)

func newServerCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create server password",
		Args:  cobra.MinimumNArgs(2),
		Run:   create,
	}
	cmd.Flags().StringVarP(&user, "user", "u", "root", "--user=root -u root")
	cmd.Flags().StringVarP(&port, "port", "P", "22", "--port=22 -P 22")
	cmd.Flags().StringVarP(&group, "group", "g", "", "--group=xxx -g xxx")
	return cmd
}

func completer(in prompt.Document) []prompt.Suggest {
	return nil
}

func create(cmd *cobra.Command, args []string) {
	host := args[0]
	password = args[1]
	index := serverConfig.Find(host)
	var (
		err error
	)
	if index >= 0 {
		yesNo := prompt.Input("server already exist, press yes to modify or exit:", completer)
		if yesNo != "yes" {
			return
		}
		err = serverConfig.UpdateByIndex(index, &define.Server{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
			Group:    group,
		})

	} else {
		err = serverConfig.Add(&define.Server{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
			Group:    group,
		})
	}

	if err != nil {
		fmt.Println("save server list error: ", err)
	} else {
		fmt.Println("save server list success")
	}
}
