package main

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Returns Author signature following Git priorities
func GetASignature(repo *git.Repository) (*object.Signature, error) {
	name := os.Getenv("GIT_AUTHOR_NAME")
	email := os.Getenv("GIT_AUTHOR_EMAIL")

	// 1. Local repository config
	if name == "" || email == "" {
		cfg, err := repo.Config()
		if err == nil && cfg.User.Name != "" && cfg.User.Email != "" {
			if name == "" {
				name = cfg.User.Name
			}
			if email == "" {
				email = cfg.User.Email
			}
		}
	}

	// 2. Global config (~/.gitconfig)
	if name == "" || email == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			globalPath := filepath.Join(home, ".gitconfig")
			if data, err := os.ReadFile(globalPath); err == nil {
				globalCfg, err := config.ReadConfig(bytes.NewReader(data))
				if err == nil {
					if name == "" && globalCfg.User.Name != "" {
						name = globalCfg.User.Name
					}
					if email == "" && globalCfg.User.Email != "" {
						email = globalCfg.User.Email
					}
				}
			}
		}
	}

	// 3. Fallback
	if name == "" {
		name = "Unknown User"
	}
	if email == "" {
		email = "unknown@example.com"
	}

	return &object.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}, nil
}
