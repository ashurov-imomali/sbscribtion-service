package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Configs struct {
	Srv Server   `json:"srv" yaml:"srv"`
	Pg  Postgres `json:"postgres" yaml:"postgres"`
}

type Server struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}

type Postgres struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	DbName   string `json:"db_name" yaml:"db_name"`
}

func New() (*Configs, error) {
	file, err := os.ReadFile("./config/configs.yaml")
	if err != nil {
		return nil, err
	}
	c := &Configs{}
	return c, yaml.Unmarshal(file, &c)
}
