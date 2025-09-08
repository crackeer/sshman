package service

import (
	"fmt"
	"io"
	"os"
	"sshman/define"

	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
}

// NewSSHClient
//
//	@param server
//	@return *SSHClient
//	@return error
func NewSSHClient(server *define.Server) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User:            server.User,
		Auth:            []ssh.AuthMethod{ssh.Password(server.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", server.Host+":"+server.Port, config)
	if err != nil {
		return nil, fmt.Errorf("dial error: %v", err)
	}

	return &SSHClient{
		client: client,
	}, nil
}

func (s *SSHClient) Close() {
	s.client.Close()
}

func (s *SSHClient) Mkdir(tempDir string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("new ssh session error: %s", err.Error())
	}
	defer session.Close()
	cmd := fmt.Sprintf("mkdir -p %s", tempDir)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("exec ssh command error: %s", err.Error())
	}
	return nil
}

func (s *SSHClient) RemoveDir(tempDir string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("new ssh session error: %s", err.Error())
	}
	defer session.Close()
	cmd := fmt.Sprintf("rm -rf %s", tempDir)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("exec ssh command error: %s", err.Error())
	}
	return nil
}

func (s *SSHClient) Run(cmd string) error {
	session, err := s.client.NewSession()

	if err != nil {
		return fmt.Errorf("new ssh session error: %s", err.Error())
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	defer session.Close()
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("exec ssh command error: %s", err.Error())
	}
	return nil
}

func (s *SSHClient) K3sImport(remoteFile string) error {
	return s.Run(fmt.Sprintf("k3s ctr image import %s", remoteFile))
}

func (s *SSHClient) DockerLoad(remoteFile string) error {
	return s.Run(fmt.Sprintf("docker load -i %s", remoteFile))
}

// UploadTo
//
//	@param localFile
//	@param remoteFile
//	@return error
func (s *SSHClient) UploadTo(localFile, remoteFile string) error {
	sftpClient, err := sftp.NewClient(s.client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("open local file error: %s", err.Error())
	}
	defer srcFile.Close()

	fileStat, err := os.Stat(localFile)
	if err != nil {
		return fmt.Errorf("file stat error: %v", err)
	}

	dstFile, err := sftpClient.Create(remoteFile)
	if err != nil {
		return fmt.Errorf("create remote file error: %s", err.Error())
	}
	defer dstFile.Close()

	// 创建进度条
	bar := pb.Full.Start64(fileStat.Size())
	defer bar.Finish()
	// 创建带进度条的写入器
	writer := bar.NewProxyWriter(dstFile)
	_, err = io.Copy(writer, srcFile)
	return err
}


func (s *SSHClient) DownloadTo(remoteFile, localFile string) error {
	sftpClient, err := sftp.NewClient(s.client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := sftpClient.Open(remoteFile)
	if err != nil {
		return fmt.Errorf("open remote file error: %s", err.Error())
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localFile)
	if err != nil {
		return fmt.Errorf("create local file error: %s", err.Error())
	}
	defer dstFile.Close()

	fileStat, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("file stat error: %v", err)
	}

	// 创建进度条
	bar := pb.Full.Start64(fileStat.Size())
	defer bar.Finish()
	// 创建带进度条的写入器
	writer := bar.NewProxyWriter(dstFile)
	_, err = io.Copy(writer, srcFile)
	return err
}