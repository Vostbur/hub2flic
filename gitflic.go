package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"

	git_http "github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Project struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Alias       *string `json:"alias"`
	Language    *string `json:"language"`
	Private     *bool   `json:"private"`
	PushURL     *string
	ClonePath   *string
}

func String(v string) *string { return &v }

func Bool(v bool) *bool { return &v }

func NewProject(cfg *Config, name string, description string, language string, private bool) *Project {
	return &Project{
		Title:       String(name),
		Description: String(description),
		Alias:       String(name),
		Language:    String(language),
		Private:     Bool(private),
		PushURL:     String(fmt.Sprintf("%s/%s/%s.git", GITFLIC_URL, cfg.GitFlicName, name)),
		ClonePath:   String(cfg.ClonePath),
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

func (p *Project) InitCommit(pth string) error {
	path := *p.ClonePath + pth

	r, err := git.PlainInit(path, false)
	if err != nil {
		return err
	}

	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{*p.PushURL},
	})
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	_, err = w.Commit("Added my new file", &git.CommitOptions{})
	if err != nil {
		return err
	}

	return nil
}

// push to GitFlic
func (p *Project) Push(cfg *Config, path string) error {
	r, err := git.PlainOpen(path)
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
