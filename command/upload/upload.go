package upload

import (
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"sshman/define"
	"sshman/service"
	"strings"
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
				remoteDir := parentDir(file.To)
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
		mapping, _ := GetFileMapping(file)
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

func GetFileMapping(path string) (map[string]string, error) {
	retData := map[string]string{}
	err := filepath.Walk(path, func(tmpPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Clean(path) == filepath.Clean(tmpPath) {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		name, _ := filepath.Rel(path, tmpPath)
		retData[name] = tmpPath
		return nil
	})

	return retData, err
}

func parentDir(path string) string {
	parts := strings.Split(path, "/")
	return strings.Join(parts[0:len(parts)-1], "/")
}
