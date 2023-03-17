package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"log"
)

const (
	GITFLIC_API_URL = "https://api.gitflic.ru/project"
	GITFLIC_URL     = "https://gitflic.ru/project"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to YAML config file")
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

	log.Printf("\033[34;43;1m%s\033[0m\n", "configuration set up successfully")

	gh := new(GitHub)
	gh.Set(cfg)

	count := 0
	for i, repo := range gh.List() {
		///// temp
		if i > 3 {
			break
		}
		///// end temp

		clonePath := cfg.ClonePath + *repo.Name

		gh.Clone(repo)

		// GitFlic limits for project name and alias
		if len(*repo.Name) < 3 {
			repo.Name = String(fmt.Sprintf("github_%s", *repo.Name))
		}
		repo.Name = String(strings.ReplaceAll(*repo.Name, ".", ""))

		gf := NewProject(cfg, repo)

		if !gf.Exists(cfg) {
			if err := gf.Create(cfg); err != nil {
				log.Printf("\033[31;1m%s\033[0m\n", err)
				cleanUp(clonePath)
				continue
			}
		}

		if err := gf.Push(cfg, clonePath); err != nil {
			log.Printf("\033[31;1m%s\033[0m\n", err)
			cleanUp(clonePath)
			continue
		}

		cleanUp(clonePath)
		count++
	}

	log.Printf("\033[34;43;1m%s\033[0m\n",
		fmt.Sprintf("moved %d repositories from GitHub to GitFlic", count))

	// list all gists for the authenticated user
	// gists, resp, err := client.Gists.List(
	// 	ctx,
	// 	"",
	// 	&github.GistListOptions{ListOptions: lops})

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Gists:", len(gists), "Status:", resp.StatusCode)

	// for _, gist := range gists {
	// 	_, err := git.PlainClone(cfg.ClonePath+*gist.ID, false, &git.CloneOptions{
	// 		URL:      *gist.GitPullURL,
	// 		Progress: os.Stdout,
	// 	})
	// 	if err != nil {
	// 		log.Println(err)
	// 	}

	// 	// gst, _, err := client.Gists.Get(ctx, *gist.ID)
	// 	// if err != nil {
	// 	// 	log.Println(err)
	// 	// }

	// 	// for fname, details := range gst.Files {
	// 	// 	err := os.WriteFile(
	// 	// 		fmt.Sprintf("%s%s_%s", cfg.ClonePath, *gist.ID, fname),
	// 	// 		[]byte(*details.Content),
	// 	// 		0644)

	// 	// 	if err != nil {
	// 	// 		log.Println(err)
	// 	// 	}
	// 	// }
	// }
}
