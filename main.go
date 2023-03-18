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

var configFile, gistOpt string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to YAML config file")
	flag.StringVar(&gistOpt, "gist", "no",
		"Options for clone gists, may be 'no' (default), 'yes' or 'single'")
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

func main() {
	flag.Parse()

	cfg := new(Config)
	if isFlagSet("config") {
		if err := cfg.Set(configFile); err != nil {
			log.Fatalf("\033[31;1m%s\033[0m\n", err)
		}
	} else {
		log.Fatalf("\033[31;1m%s\033[0m\n",
			errors.New("config filename is empty. Use with flag -config or -help"))
	}

	if !isFlagSet("gist") {
		gistOpt = "no"
	}

	log.Printf("\033[34;43;1m%s\033[0m\n", "configuration set up successfully")

	gh := new(GitHub)
	gh.Set(cfg)

	count := reposGH(cfg, gh)

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
