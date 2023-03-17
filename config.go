package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	GitHubName   string `yaml:"github_name"`
	GitHubToken  string `yaml:"github_token"`
	PerPage      int    `yaml:"per_page"`
	GitFlicName  string `yaml:"gitflic_name"`
	GitFlicPass  string `yaml:"gitflic_pass"`
	GitFlicToken string `yaml:"gitflic_token"`
	ClonePath    string `yaml:"clone_path"`
}

// set up configuration
func (cfg *Config) Set(fname string) (err error) {
	buf, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buf, cfg)
	if err != nil {
		return err
	}

	if len(cfg.ClonePath) > 0 && cfg.ClonePath[len(cfg.ClonePath)-1] != '/' {
		cfg.ClonePath += "/"
	}

	return
}
