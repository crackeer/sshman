package upload

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sshman/define"
	"sshman/service"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	servers   []string
	group     string
	targetDir string

	serverConfig *service.ServerConfig
)

func init() {
	var err error
	serverConfig, err = service.NewServerConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NewUploadCommand
//
//	@return *cobra.Command
func NewUploadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "upload ",
		Run:   upload,
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.Flags().StringSliceVarP(&servers, "server", "s", []string{}, "--server=127.0.0.1 -s 127.0.0.1")
	cmd.Flags().StringVarP(&group, "group", "g", "", "--group=xxx -g xxx")
	cmd.Flags().StringVarP(&targetDir, "remote-dir", "d", "/tmp/sshman/upload", "--dir=xxx -d xxx")
	cmd.MarkFlagsOneRequired("server", "group")
	return cmd
}

func upload(cmd *cobra.Command, args []string) {
	list := []*define.Server{}
	if len(group) > 0 {
		list = serverConfig.FindByGroup(group)
	} else {
		for _, server := range servers {
			if index := serverConfig.Find(server); index >= 0 {
				list = append(list, serverConfig.Get(index))
				continue
			}
		}
	}
	if len(list) == 0 {
		fmt.Println("servers not found")
		return
	}
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

			err = client.Mkdir(targetDir)
			if err != nil {
				fmt.Println("mkdir error:", err, server.Host, targetDir)
				os.Exit(1)
			}

			for index, file := range args {
				_, name := filepath.Split(file)
				fmt.Println(index+1, ". updload:", file, "->", targetDir+"/"+name)
				err = client.UploadTo(file, targetDir+"/"+name)
				if err != nil {
					fmt.Println("upload error:", err, server.Host, targetDir+"/"+name)
					os.Exit(1)
				}
			}
		}
		ch <- true

	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	fmt.Println("press ctrl+c to stop")
	select {
	case <-ch:
		fmt.Println("")
		fmt.Println("all files upload success.")
	case <-signalChan:
		fmt.Println("user cancel")
	}

}
