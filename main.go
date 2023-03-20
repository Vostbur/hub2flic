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
	ERROR           = "\033[31;1m%s\033[0m\n"
	MESSAGE         = "\033[34;43;1m%s\033[0m\n"
	INFO            = "\033[34;1m%s\033[0m\n"
)

// vars for CLI keys
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
func newConfig() (*Config, error) {
	cfg := new(Config)

	if isFlagSet("config") {
		if err := cfg.Set(configFile); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("config filename is empty. Use with flag -config or -help")
	}

	if !isFlagSet("gist") {
		gistOpt = "no"
	}

	return cfg, nil
}

func main() {
	flag.Parse()

	cfg, err := newConfig()
	if err != nil {
		log.Fatalf(ERROR, err)
	}

	log.Printf(MESSAGE, "configuration set up successfully")

	gh := newGitHub(cfg)
	var countRepo, countGist uint

	if isFlagSet("repo") {

		repo, err := gh.GetRepoByName(cfg.GitHubName, repoOpt)
		if err != nil {
			log.Fatalf(ERROR, err)
		}

		log.Printf(MESSAGE, fmt.Sprintf("received '%s' repository from GitHub", *repo.Name))

		validateRepo(repo)

		if isSuccess := transferRepo(
			cfg,
			gh,
			*repo.Name,
			*repo.CloneURL,
			*repo.Description,
			*repo.Language,
			*repo.Private); isSuccess {

			countRepo++
		}

	} else {

		repos, err := gh.ReposList()
		if err != nil {
			log.Fatalf(ERROR, err)
		}

		log.Printf(MESSAGE, fmt.Sprintf("received %d repositories from GitHub", len(repos)))

		for _, repo := range repos {

			validateRepo(repo)

			if isSuccess := transferRepo(
				cfg,
				gh,
				*repo.Name,
				*repo.CloneURL,
				*repo.Description,
				*repo.Language,
				*repo.Private); isSuccess {

				countRepo++
			}
		}

	}

	log.Printf(MESSAGE, fmt.Sprintf("moved %d repositories from GitHub to GitFlic", countRepo))

	switch gistOpt {
	case "yes":
		gists, err := gh.GistsList()
		if err != nil {
			log.Fatalf(ERROR, err)
		}

		log.Printf(MESSAGE, fmt.Sprintf("received %d gists from GitHub", len(gists)))
		
		for _, gist := range gists {
			if isSuccess := transferRepo(
				cfg,
				gh,
				*gist.ID,
				*gist.GitPullURL,
				*gist.Description,
				"Markdown",
				!*gist.Public); isSuccess {

				countGist++
			}
		}

	case "single":
		countGist = gistSingle(cfg, gh)

	case "no":
		fallthrough

	default:
		countGist = 0
	}

	log.Printf(MESSAGE, fmt.Sprintf("moved %d gists from GitHub to GitFlic", countGist))
}
