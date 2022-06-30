package main

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"strings"
	"time"
)

type GitClient struct {
	repo   *git.Repository
	origin string
	author string
	email  string
}

func (gitClient *GitClient) CloneRepo(origin string, destination string, username string, accessToken string, author string, email string) error {
	if destination == "" {
		return fmt.Errorf("destination cannot be blank")
	}
	if origin == "" {
		return fmt.Errorf("origin cannot be blank")
	}
	if !strings.HasPrefix(origin, "http") && !strings.HasPrefix(origin, "git") {
		return fmt.Errorf("origin must start with http or git")
	}
	if strings.HasPrefix(origin, "git") {
		fmt.Println("ssh authentication not currently supported, attempting without auth")
	}

	if gitClient.repo == nil {
		err := os.MkdirAll(destination, 0700)
		if err != nil {
			return fmt.Errorf("error making directory: %v", err)
		}
		auth := http.BasicAuth{}
		if username != "" || accessToken != "" {
			auth.Password = accessToken
			auth.Username = username
		}
		r, err := git.PlainCloneContext(context.Background(), destination, false, &git.CloneOptions{
			URL:               origin,
			RecurseSubmodules: git.NoRecurseSubmodules,
			Depth:             1,
			Auth:              &auth,
			Progress:          nil,
		})
		if err != nil {
			return fmt.Errorf("error cloning: %v", err)
		}
		gitClient.repo = r
		gitClient.origin = origin
		gitClient.email = email
		gitClient.author = author
	} else {
		if gitClient.origin != origin {
			return fmt.Errorf("cannot change origin")
		}
	}
	return nil
}

func (gitClient *GitClient) IsClean() (bool, error) {
	if !gitClient.Ready() {
		return false, fmt.Errorf("git client not inititialized")
	}
	worktree, err := gitClient.repo.Worktree()
	if err != nil {
		return false, err
	}
	status, err := worktree.Status()
	if err != nil {
		return false, err
	}
	return status.IsClean(), nil
}

func (gitClient *GitClient) AddAll() error {
	if !gitClient.Ready() {
		return fmt.Errorf("git client not inititialized")
	}
	worktree, err := gitClient.repo.Worktree()
	if err != nil {
		return err
	}
	err = worktree.AddWithOptions(&git.AddOptions{
		All: true,
	})
	if err != nil {
		return fmt.Errorf("error adding: %v", err)
	}
	return nil
}

func (gitClient *GitClient) Ready() bool {
	if gitClient.repo == nil {
		return false
	}

	return true
}

func (gitClient *GitClient) CommitAndPush(accessToken string, username string, force bool) error {
	worktree, err := gitClient.repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree in commit and push: %v", err)
	}
	auth := http.BasicAuth{}
	if username != "" || accessToken != "" {
		auth.Password = accessToken
		auth.Username = username
	}
	commit, err := worktree.Commit(fmt.Sprintf("current time: %v", time.Now()), &git.CommitOptions{
		All: false,
		Author: &object.Signature{
			Name:  gitClient.author,
			Email: gitClient.email,
			When:  time.Now(),
		},
		Committer: &object.Signature{
			Name:  gitClient.author,
			Email: gitClient.email,
			When:  time.Now(),
		},
		Parents: nil,
		SignKey: nil,
	})
	if err != nil {
		return fmt.Errorf("error committing: %v", err)
	}
	config, err := gitClient.repo.Config()
	if err != nil {
		return fmt.Errorf("error getting config: %v", err)
	}
	remotes := config.Remotes
	for k, v := range remotes {
		fmt.Printf("remotes, key: %v, Name: %v, Url: %v\n", k, v.Name, fmt.Sprint(v.URLs))
	}
	err = gitClient.repo.PushContext(context.Background(), &git.PushOptions{
		RemoteName:        "origin",
		RefSpecs:          nil,
		Auth:              &auth,
		Progress:          nil,
		Prune:             false,
		Force:             force,
		InsecureSkipTLS:   false,
		CABundle:          nil,
		RequireRemoteRefs: nil,
	})
	if err != nil {
		return err
	}
	fmt.Printf("commit %v\n", commit.String())
	return nil
}
