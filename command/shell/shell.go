package shell

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sshman/define"
	"sshman/service"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	defaultShellDir string = "/tmp/sshman/shell-files"
)

// NewBashCommand ...
//
//	@param engine
//	@return *cobra.Command
func NewShellCommand(engine string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   engine,
		Short: "execute shell files",
		Run:   importTar,
		Args:  cobra.MinimumNArgs(1),
		Annotations: map[string]string{
			"engine": engine,
		},
	}
	return cmd
}

func importTar(cmd *cobra.Command, args []string) {
	list := define.GetServers()
	if len(list) == 0 {
		fmt.Println("servers not found")
		return
	}

	bashFile := args[0]

	_, name := filepath.Split(bashFile)
	targetFile := defaultShellDir + "/" + time.Now().Format("20060102150405") + name

	fmt.Println("press ctrl+c to stop")
	ch := make(chan bool)
	go func() {
		for _, server := range list {
			fmt.Println("")
			fmt.Println("--> upload to server:", server.Host)
			client, err := service.NewSSHClient(server)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = client.Mkdir(defaultShellDir)
			if err != nil {
				fmt.Println("mkdir error:", err, server.Host, defaultShellDir)
				os.Exit(1)
			}

			err = client.UploadTo(bashFile, targetFile)
			if err != nil {
				fmt.Println("upload error:", err, server.Host, bashFile, targetFile)
				os.Exit(1)
			}

			engine := cmd.Annotations["engine"]
			fmt.Println("run shell file:", targetFile)

			err = client.Run(fmt.Sprintf("%s %s", engine, targetFile))
			if err != nil {
				fmt.Println("run shell error:", err, server.Host, targetFile)
				os.Exit(1)
			}

		}
		ch <- true
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)

	select {
	case <-ch:
		fmt.Println("")
	case <-signalChan:
		fmt.Println("user cancel")
	}
}
