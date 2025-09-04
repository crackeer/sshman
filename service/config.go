package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sshman/define"
)

// GetConfigDir
//
//	@return string
//	@return error
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "sshman")
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return "", err
	}
	return configDir, nil
}

type ServerConfig struct {
	configPath string
	list       []*define.Server
}

// NewServerConfig
//
//	@return *ServerConfig
//	@return error
func NewServerConfig() (*ServerConfig, error) {
	path, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	serverList := []*define.Server{}
	configPath := filepath.Join(path, "server.json")
	err = ReadFileTo(configPath, &serverList)

	return &ServerConfig{
		configPath: configPath,
		list:       serverList,
	}, nil
}

func (s *ServerConfig) save() error {
	bytes, _ := json.Marshal(s.list)
	return os.WriteFile(s.configPath, bytes, os.ModePerm)
}

func (s *ServerConfig) Add(server *define.Server) error {
	s.list = append(s.list, server)
	return s.save()
}

func (s *ServerConfig) Find(host string) int {
	for index, server := range s.list {
		if server.Host == host {
			return index
		}
	}
	return -1
}

func (s *ServerConfig) FindByGroup(group string) []*define.Server {
	retData := []*define.Server{}
	for _, server := range s.list {
		if server.Group == group {
			retData = append(retData, server)
		}
	}
	return retData
}

func (s *ServerConfig) DeleteByHost(host string) error {
	index := s.Find(host)
	if index < 0 {
		return fmt.Errorf("server not found, host: %s", host)
	}
	s.list = append(s.list[:index], s.list[index+1:]...)
	return s.save()
}

func (s *ServerConfig) UpdateByIndex(index int, server *define.Server) error {
	s.list[index] = server
	return s.save()
}

func (s *ServerConfig) List() []*define.Server {
	return s.list
}

func (s *ServerConfig) Get(index int) *define.Server {
	if index < 0 || index >= len(s.list) {
		return nil
	}
	return s.list[index]
}
