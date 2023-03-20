package main

import (
	"context"
	"fmt"
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

// GitHub structure constructor
func newGitHub(cfg *Config) *GitHub {
	g := new(GitHub)
	g.Ctx = context.Background()
	g.Ts = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	g.Tc = oauth2.NewClient(g.Ctx, g.Ts)
	g.Client = github.NewClient(g.Tc)
	g.Lops = github.ListOptions{PerPage: cfg.PerPage}
	g.ClonePath = cfg.ClonePath

	return g
}

// return GitHub repository by name
func (g *GitHub) GetRepoByName(owner, name string) (*github.Repository, error) {
	repo, resp, err := g.Client.Repositories.Get(g.Ctx, owner, name)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("no repositories received from GitHub. Status code: %d",
			resp.StatusCode)
	}

	return repo, nil
}

// list all repositories for the authenticated user
func (g *GitHub) ReposList() ([]*github.Repository, error) {
	repos, resp, err := g.Client.Repositories.List(
		g.Ctx,
		"",
		&github.RepositoryListOptions{ListOptions: g.Lops})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("no repositories received from GitHub. Status code: %d", resp.StatusCode)
	}

	return repos, nil
}

// list all gists for the authenticated user
func (g *GitHub) GistsList() ([]*github.Gist, error) {
	gists, resp, err := g.Client.Gists.List(
		g.Ctx,
		"",
		&github.GistListOptions{ListOptions: g.Lops})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("no gists received from GitHub. Status code: %d", resp.StatusCode)
	}

	return gists, nil
}

// clone repository from GitHub
func (g *GitHub) Clone(name string, url string) error {
	clonePath := g.ClonePath + name

	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:      url,
		Progress: nil,
	})

	return err
}
