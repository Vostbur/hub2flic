package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"

	git_http "github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Project struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Alias       *string `json:"alias"`
	Language    *string `json:"language"`
	Private     *bool   `json:"private"`
	PushURL     *string
}

func String(v string) *string { return &v }

func Bool(v bool) *bool { return &v }

func NewProject(cfg *Config, r *github.Repository) *Project {
	return &Project{
		Title:       r.Name,
		Description: r.Description,
		Alias:       r.Name,
		Language:    r.Language,
		Private:     Bool(false),
		PushURL:     String(fmt.Sprintf("%s/%s/%s.git", GITFLIC_URL, cfg.GitFlicName, *r.Name)),
	}
}

// check GitFlick project exists
func (p *Project) Exists(cfg *Config) bool {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/%s/%s", GITFLIC_API_URL, cfg.GitFlicName, *p.Alias), nil)
	if err != nil {
		log.Fatalf("\033[31;1m%s\033[0m\n", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+cfg.GitFlicToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("\033[31;1m%s\033[0m\n", err)
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Printf("\033[34;1m%s\033[0m\n",
			fmt.Sprintf("GitFlic project '%s' exists", *p.Title))
		return true
	}

	return false
}

func (p *Project) setupNilOptions() {
	if p.Description == nil {
		p.Description = String("Description")
	}

	if p.Language == nil {
		p.Language = String("Markdown")
	}
}

// create GitFlic project
func (p *Project) Create(cfg *Config) error {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(p)
	if err != nil {
		return err
	}

	p.setupNilOptions()

	req, err := http.NewRequest("POST", GITFLIC_API_URL, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+cfg.GitFlicToken)

	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("GitFlic project '%s' is not created. Status code: %d",
			*p.Title, resp.StatusCode)
	}

	log.Printf("\033[34;1m%s\033[0m\n",
		fmt.Sprintf("GitFlic project '%s' is created", *p.Title))

	return nil
}

// push to GitFlic
func (p *Project) Push(cfg *Config, clonePath string) error {
	r, err := git.PlainOpen(clonePath)
	if err != nil {
		return err
	}

	auth := &git_http.BasicAuth{
		Username: cfg.GitFlicName,
		Password: cfg.GitFlicPass,
	}

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		RemoteURL:  *p.PushURL,
		Auth:       auth,
	})
	if err != nil {
		return err
	}

	log.Printf("\033[34;1m%s\033[0m\n",
		fmt.Sprintf("repository '%s' is pushed", *p.PushURL))

	return nil
}
