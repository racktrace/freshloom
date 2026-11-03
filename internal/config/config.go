package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Packages Packages `yaml:"packages"`
	Git      Git      `yaml:"git"`
}

type Packages struct {
	Apt []string `yaml:"apt"`
}

type Git struct {
	Name          string `yaml:"name"`
	Email         string `yaml:"email"`
	DefaultBranch string `yaml:"default_branch"`
}

// Load reads the config file, parses the YAML, and validates its contents.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true) // reject unknown fields (e.g. typos)

	if err := dec.Decode(&cfg); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("config %s is empty", path)
		}
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config %s: %w", path, err)
	}
	return &cfg, nil
}

// Validate checks rules that the YAML parser cannot enforce on its own.
func (c *Config) Validate() error {
	seen := make(map[string]bool)
	for i, pkg := range c.Packages.Apt {
		pkg = strings.TrimSpace(pkg)
		if pkg == "" {
			return fmt.Errorf("packages.apt[%d] is empty", i)
		}
		if strings.ContainsAny(pkg, " \t") {
			return fmt.Errorf("packages.apt[%d] %q contains whitespace", i, pkg)
		}
		if seen[pkg] {
			return fmt.Errorf("package %q is listed more than once", pkg)
		}
		seen[pkg] = true
	}

	if c.Git.Email != "" && !strings.Contains(c.Git.Email, "@") {
		return fmt.Errorf("git.email %q does not look like an email address", c.Git.Email)
	}
	return nil
}