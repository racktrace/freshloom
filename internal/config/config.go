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

// Load membaca file config, mem-parsing YAML, lalu memvalidasi isinya.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("membaca config %s: %w", path, err)
	}

	var cfg Config
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true) // tolak field yang tidak dikenal (misal salah ketik)

	if err := dec.Decode(&cfg); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("config %s kosong", path)
		}
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config %s tidak valid: %w", path, err)
	}
	return &cfg, nil
}

// Validate memeriksa isi config yang tidak bisa dicek oleh parser YAML.
func (c *Config) Validate() error {
	seen := make(map[string]bool)
	for i, pkg := range c.Packages.Apt {
		pkg = strings.TrimSpace(pkg)
		if pkg == "" {
			return fmt.Errorf("packages.apt[%d] kosong", i)
		}
		if strings.ContainsAny(pkg, " \t") {
			return fmt.Errorf("packages.apt[%d] %q mengandung spasi", i, pkg)
		}
		if seen[pkg] {
			return fmt.Errorf("paket %q ditulis lebih dari sekali", pkg)
		}
		seen[pkg] = true
	}

	if c.Git.Email != "" && !strings.Contains(c.Git.Email, "@") {
		return fmt.Errorf("git.email %q tidak terlihat seperti email", c.Git.Email)
	}
	return nil
}
