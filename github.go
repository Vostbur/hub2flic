package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitHub struct {
	Ctx       context.Context
	Ts        oauth2.TokenSource
	Tc        *http.Client
	Client    *github.Client
	Lops      github.ListOptions
	ClonePath string
}

func (g *GitHub) Set(cfg *Config) {
	g.Ctx = context.Background()
	g.Ts = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	g.Tc = oauth2.NewClient(g.Ctx, g.Ts)
	g.Client = github.NewClient(g.Tc)
	g.Lops = github.ListOptions{PerPage: cfg.PerPage}
	g.ClonePath = cfg.ClonePath
}

// list all repositories for the authenticated user
func (g *GitHub) List() []*github.Repository {
	repos, resp, err := g.Client.Repositories.List(
		g.Ctx,
		"",
		&github.RepositoryListOptions{ListOptions: g.Lops})
	if err != nil {
		log.Fatalf("\033[31;1m%s\033[0m\n", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("\031[34;1m%s\033[0m\n",
			fmt.Sprintf("no repositories received from GitHub. Status code: %d", resp.StatusCode))
	}

	log.Printf("\033[34;43;1m%s\033[0m\n",
		fmt.Sprintf("received %d repositories from GitHub", len(repos)))

	return repos
}

// clone repository from GitHub
func (g *GitHub) Clone(r *github.Repository) {
	clonePath := g.ClonePath + *r.Name

	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:      *r.CloneURL,
		Progress: nil,
	})
	if err != nil {
		log.Printf("\033[31;1m%s\033[0m\n", err) // normal continuetion if repo exists localy
	} else {
		log.Printf("\033[34;1m%s\033[0m\n", fmt.Sprintf("repository '%s' is cloned", *r.FullName))
	}
}
