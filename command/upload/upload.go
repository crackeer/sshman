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
	targetDir string
)

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

	cmd.Flags().StringVarP(&targetDir, "remote-dir", "d", "/tmp/sshman/upload", "--remote-dir=xxx -d xxx")
	return cmd
}

func upload(cmd *cobra.Command, args []string) {
	uploadFiles := []UploadFileConfig{}
	for _, file := range args {
		files, err := collectFiles(file, targetDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		uploadFiles = append(uploadFiles, files...)
	}
	if len(uploadFiles) == 0 {
		fmt.Println("no files found")
		return
	}
	list := define.GetServers()
	if len(list) == 0 {
		fmt.Println("servers not found")
		return
	}
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

			for index, file := range uploadFiles {
				fmt.Println("")
				fmt.Println(fmt.Sprintf("[%d/%d]", index+1, len(uploadFiles)), server.Host)
				fmt.Println("from:", file.From)
				fmt.Println("  to:", file.To)
				remoteDir := service.ParentDir(file.To)
				err = client.Mkdir(remoteDir)
				if err != nil {
					fmt.Println("mkdir error:", err, server.Host, targetDir)
					os.Exit(1)
				}

				err = client.UploadTo(file.From, file.To)
				if err != nil {
					fmt.Println("upload error:", err, server.Host, file.From, file.To)
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
		fmt.Println("")
		fmt.Println("all files upload success.")
	case <-signalChan:
		fmt.Println("user cancel")
	}

}

type UploadFileConfig struct {
	From string
	To   string
}

func collectFiles(file string, targetDir string) ([]UploadFileConfig, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	retData := []UploadFileConfig{}
	if stat.IsDir() {
		mapping, _ := service.GetFileMapping(file)
		for key, value := range mapping {
			retData = append(retData, UploadFileConfig{
				From: value,
				To:   targetDir + "/" + key,
			})
		}
		return retData, nil
	}
	retData = append(retData, UploadFileConfig{
		From: file,
		To:   targetDir + "/" + filepath.Base(file),
	})
	return retData, nil
}
