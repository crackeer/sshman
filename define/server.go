package define

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Server struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Group    string `json:"group"`
}

var (
	GServerConfig *ServerConfig
	GServers      []string
	GGroup        string
)

func init() {
	var err error
	sc, err := NewServerConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	GServerConfig = sc
}

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
	list       []*Server
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
	serverList := []*Server{}
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

func (s *ServerConfig) Add(server *Server) error {
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

func (s *ServerConfig) FindByGroup(group string) []*Server {
	retData := []*Server{}
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

func (s *ServerConfig) UpdateByIndex(index int, server *Server) error {
	s.list[index] = server
	return s.save()
}

func (s *ServerConfig) List() []*Server {
	return s.list
}

func (s *ServerConfig) Get(index int) *Server {
	if index < 0 || index >= len(s.list) {
		return nil
	}
	return s.list[index]
}

// Group
//
//	@param group
//	@param hosts
//	@return error
func (s *ServerConfig) Group(group string, hosts []string) error {
	for i, server := range s.list {
		for _, host := range hosts {
			if server.Host == host {
				s.list[i].Group = group
			}
		}
	}
	return s.save()
}

func GetServers() []*Server {
	list := []*Server{}
	fmt.Println("get servers...", GGroup)
	if len(GGroup) > 0 {
		list = GServerConfig.FindByGroup(GGroup)
	} else {
		for _, server := range GServers {
			if index := GServerConfig.Find(server); index >= 0 {
				list = append(list, GServerConfig.Get(index))
				continue
			}
		}
	}
	return list
}

func ReadFileTo(path string, value interface{}) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}
