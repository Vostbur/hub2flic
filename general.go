package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

func validateRepo(r *github.Repository) {
	if r.Description == nil {
		r.Description = String("Description")
	}

	if r.Language == nil {
		r.Language = String("Markdown")
	}
}

// TODO: move create new project to func
// TODO: refactor gistSingle
func transferRepo(cfg *Config, gh *GitHub, name, cloneUrl, descr, lang string, private bool) bool {
	clonePath := cfg.ClonePath + name

	if err := gh.Clone(name, cloneUrl); err != nil {
		log.Printf(INFO, err) // continue if repo exists localy
	} else {
		log.Printf(INFO, fmt.Sprintf("'%s' is cloned", name))
	}

	// TODO: call validation function
	// GitFlic limits for project name and alias
	if len(name) <= 3 {
		name = fmt.Sprintf("github-%s", name)
	}
	name = strings.ReplaceAll(name, ".", "")
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ToLower(name)

	gf := NewProject(cfg, name, descr, lang, private)

	ok, err := gf.Exists(cfg)
	if err != nil {
		log.Fatalf(ERROR, err)
	}

	if !ok {
		if err := gf.Create(cfg); err != nil {
			log.Printf(ERROR, err)
			cleanUp(clonePath)
			return false
		}
		log.Printf(INFO, fmt.Sprintf("GitFlic project '%s' is created", name))
	}

	if err := gf.Push(cfg, clonePath); err != nil {
		log.Printf(ERROR, err)
		cleanUp(clonePath)
		return false
	}

	log.Printf(INFO, fmt.Sprintf("repository '%s' is pushed", name))

	cleanUp(clonePath)

	return true
}

func gistSingle(cfg *Config, gh *GitHub) (count uint) {
	name := "gists"

	// "single" repo is always Private on GitFlic
	// to make Public change 'true' to 'false'
	gf := NewProject(cfg, name, "all GitHub gists in one repository", "Markdown", true)

	ok, err := gf.Exists(cfg)
	if err != nil {
		log.Fatalf(ERROR, err)
	}

	if !ok {
		if err := gf.Create(cfg); err != nil {
			log.Printf(ERROR, err)
		}
		log.Printf(INFO, fmt.Sprintf("gists GitFlic project '%s' is created", name))
	}

	clonePath := cfg.ClonePath + name

	gists, err := gh.GistsList()
	if err != nil {
		log.Fatalf(ERROR, err)
	}

	log.Printf(MESSAGE, fmt.Sprintf("received %d gists from GitHub", len(gists)))

	for _, gist := range gists {
		gst, _, err := gh.Client.Gists.Get(gh.Ctx, *gist.ID)
		if err != nil {
			log.Printf(ERROR, err)
		}

		for fname, details := range gst.Files {
			if err := os.MkdirAll(
				fmt.Sprintf("%s/%s", clonePath, *gst.ID),
				0755); err != nil {
				log.Fatalf(ERROR, err)
			}

			if err := os.WriteFile(
				fmt.Sprintf("%s/%s/%s", clonePath, *gst.ID, fname),
				[]byte(*details.Content),
				0644); err != nil {
				log.Printf(ERROR, err)
			}
		}

		count++
	}

	if err := gf.InitCommit(name); err != nil {
		cleanUp(clonePath)
		log.Fatalf(ERROR, err)
	}

	if err := gf.Push(cfg, clonePath); err != nil {
		cleanUp(clonePath)
		log.Fatalf(ERROR, err)
	}

	log.Printf(INFO, fmt.Sprintf("repository '%s' is pushed", name))

	cleanUp(clonePath)

	return
}
