package run

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"sshman/define"
	"sshman/service"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

// NewRunCommand ...
//
//	@param engine
//	@return *cobra.Command
func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run command",
		Run:   run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) {
	commands := getCommands(args)
	if len(commands) == 0 {
		fmt.Println("no commands found")
		return
	}
	list := define.GetServers()
	if len(list) == 0 {
		fmt.Println("server not found")
		return
	}

	fmt.Println("press ctrl+c to stop")
	ch := make(chan bool)
	go func() {
		for _, server := range list {

			client, err := service.NewSSHClient(server)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for index, content := range commands {
				fmt.Println("")
				fmt.Println(fmt.Sprintf("#%s [%d/%d]", server.Host, index+1, len(commands)), content)
				err := client.RunCommand(content)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

		}
		ch <- true
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)

	select {
	case <-ch:
	case <-signalChan:
		fmt.Println("user cancel")
	}
}

func getCommands(args []string) []string {
	if len(args) > 0 {
		return args
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	return strings.Split(strings.TrimSpace(string(data)), "\n")
}
