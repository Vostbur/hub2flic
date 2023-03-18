package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func gistMulti(cfg *Config, gh *GitHub) (count int) {
	for i, gist := range gh.GistsList() {
		// TODO: FOR TEST ONLY
		if i > 1 {
			break
		}

		clonePath := cfg.ClonePath + *gist.ID

		gh.Clone(gist.ID, gist.GitPullURL)

		gfg := NewProject(cfg, *gist.ID, *gist.Description, "Markdown", false)

		if !gfg.Exists(cfg) {
			if err := gfg.Create(cfg); err != nil {
				log.Printf("\033[31;1m%s\033[0m\n", err)
				cleanUp(clonePath)
				continue
			}
		}

		if err := gfg.Push(cfg, clonePath); err != nil {
			log.Printf("\033[31;1m%s\033[0m\n", err)
			cleanUp(clonePath)
			continue
		}

		cleanUp(clonePath)
		count++
	}

	return
}

func gistSingle(cfg *Config, gh *GitHub) (count int) {
	name := "gists"

	gfg := NewProject(cfg, name, "all GitHub gists in one repository", "Markdown", false)

	if !gfg.Exists(cfg) {
		if err := gfg.Create(cfg); err != nil {
			log.Printf("\033[31;1m%s\033[0m\n", err)
		}
	}

	clonePath := cfg.ClonePath + name

	for i, gist := range gh.GistsList() {
		// TODO: FOR TEST ONLY
		if i > 1 {
			break
		}

		gst, _, err := gh.Client.Gists.Get(gh.Ctx, *gist.ID)
		if err != nil {
			log.Printf("\033[31;1m%s\033[0m\n", err)
		}

		for fname, details := range gst.Files {
			if err := os.MkdirAll(
				fmt.Sprintf("%s/%s", clonePath, *gst.ID),
				0755); err != nil {
				log.Fatalf("\033[31;1m%s\033[0m\n", err)
			}

			if err := os.WriteFile(
				fmt.Sprintf("%s/%s/%s", clonePath, *gst.ID, fname),
				[]byte(*details.Content),
				0644); err != nil {
				log.Printf("\033[31;1m%s\033[0m\n", err)
			}
		}

		count++
	}

	if err := gfg.InitCommit(name); err != nil {
		cleanUp(clonePath)
		log.Fatalf("\033[31;1m%s\033[0m\n", err)
	}

	if err := gfg.Push(cfg, clonePath); err != nil {
		cleanUp(clonePath)
		log.Fatalf("\033[31;1m%s\033[0m\n", err)
	}

	cleanUp(clonePath)

	return
}

func reposGH(cfg *Config, gh *GitHub) (count int) {
	for i, repo := range gh.ReposList() {
		// TODO: FOR TEST ONLY
		if i > 1 {
			break
		}

		clonePath := cfg.ClonePath + *repo.Name

		gh.Clone(repo.Name, repo.CloneURL)

		// GitFlic limits for project name and alias
		if len(*repo.Name) < 3 {
			repo.Name = String(fmt.Sprintf("github_%s", *repo.Name))
		}
		repo.Name = String(strings.ReplaceAll(*repo.Name, ".", ""))

		gf := NewProject(cfg, *repo.Name, *repo.Description, *repo.Language, false)

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

	return
}
