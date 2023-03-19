package main

import (
	"errors"
	"flag"
	"fmt"

	"log"
)

const (
	GITFLIC_API_URL = "https://api.gitflic.ru/project"
	GITFLIC_URL     = "https://gitflic.ru/project"
)

var configFile, gistOpt, repoOpt string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to YAML config file")
	flag.StringVar(&repoOpt, "repo", "", "Repository name")
	flag.StringVar(&gistOpt, "gist", "no", "Clone gists: 'no' (or without key), 'yes' or 'single'")
}

func isFlagSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// Configure program options with CLI args
func setup(cfg *Config) error {
	if isFlagSet("config") {
		if err := cfg.Set(configFile); err != nil {
			return err
		}
	} else {
		return errors.New("config filename is empty. Use with flag -config or -help")
	}
	
	if !isFlagSet("gist") {
		gistOpt = "no"
	}

	return nil
}

func main() {
	flag.Parse()

	cfg := new(Config)

	if err := setup(cfg); err != nil {
		log.Fatalf("\033[31;1m%s\033[0m\n", err)
	}

	log.Printf("\033[34;43;1m%s\033[0m\n", "configuration set up successfully")

	gh := new(GitHub)
	gh.Set(cfg)

	var count uint
	
	if isFlagSet("repo") {
		if err := getRepoByName(cfg, gh, repoOpt); err != nil {
			log.Fatalf("\033[31;1m%s\033[0m\n", err)
		}
		count = 1
	} else {
		count = reposGH(cfg, gh)
	}

	log.Printf("\033[34;43;1m%s\033[0m\n",
		fmt.Sprintf("moved %d repositories from GitHub to GitFlic", count))

	switch gistOpt {
	case "yes":
		count = gistMulti(cfg, gh)
	case "single":
		count = gistSingle(cfg, gh)
	case "no":
		fallthrough
	default:
		count = 0
	}

	log.Printf("\033[34;43;1m%s\033[0m\n",
		fmt.Sprintf("moved %d gists from GitHub to GitFlic", count))
}
