package load

import (
	"fmt"
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
	targetDir string
)

// NewUploadCommand
//
//	@return *cobra.Command
func NewImportCommand(engine string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   engine + "-import",
		Short: "import files to k3s or docker",
		Run:   importTar,
		Args:  cobra.MinimumNArgs(1),
		Annotations: map[string]string{
			"engine": engine,
		},
	}
	engine = strings.ToLower(engine)
	defaultTargetDir := "/tmp/sshman/image-upload"
	cmd.Flags().StringVarP(&targetDir, "remote-dir", "d", defaultTargetDir, "--dir=xxx -d xxx")
	return cmd
}

func importTar(cmd *cobra.Command, args []string) {
	list := define.GetServers()
	if len(list) == 0 {
		fmt.Println("servers not found")
		return
	}

	uploadFiles := []UploadFileConfig{}
	for _, file := range args {
		files, err := collectTarFiles(file, targetDir)
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
				engine := cmd.Annotations["engine"]
				fmt.Println("loading file:", file.To)
				if engine == "k3s" {
					err = client.K3sImport(file.To)
				}
				if engine == "docker" {
					err = client.DockerLoad(file.To)
				}
				if err != nil {
					fmt.Println("load error:", err, server.Host, file.To)
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

func collectTarFiles(file string, targetDir string) ([]UploadFileConfig, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	retData := []UploadFileConfig{}
	if stat.IsDir() {
		entry, err := os.ReadDir(file)
		if err != nil {
			return nil, err
		}
		for _, value := range entry {
			if value.IsDir() || !strings.HasSuffix(value.Name(), ".tar") {
				continue
			}
			retData = append(retData, UploadFileConfig{
				From: filepath.Join(file, value.Name()),
				To:   targetDir + "/" + value.Name(),
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
