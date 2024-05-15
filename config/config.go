package config

import (
	"github.com/BurntSushi/toml"
	"path/filepath"
)

// Section is a section of the configuration file.
// T is always going to be the type of this section.
type Section[T any] interface {
	// Defaults completes the section with default values, set values are not replaced.
	Defaults() T
}

// Config is a struct representation of the TOML configuration file.
type Config struct {
	// HTTP is the "http" configuration section.
	HTTP *HTTP `toml:"http"`
	// Repos is the collection of repository configuration, keyed by their ID.
	Repos map[string]*Repo `toml:"repos"`
}

// Defaults completes the configuration with default values.
func (c *Config) Defaults() *Config {
	c.HTTP = c.HTTP.Defaults()
	for k, v := range c.Repos {
		c.Repos[k] = v.Defaults()
	}

	return c
}

// HTTP is an HTTP configuration section of the configuration file.
type HTTP struct {
	// Nero is the nero API configuration section.
	Nero *HTTPNero `toml:"nero"`
	// Nekos is the nekos API configuration section.
	Nekos *HTTPNekos `toml:"nekos"`
}

// Defaults completes the section with default values.
func (h *HTTP) Defaults() *HTTP {
	h.Nero = h.Nero.Defaults()
	h.Nekos = h.Nekos.Defaults()

	return h
}

// HTTPNero is a nero HTTP API configuration section of the configuration file.
type HTTPNero struct {
	// Host is the host string, used for http.ListenAndServe.
	Host string `toml:"host"`
}

// Defaults completes the section with default values.
func (hn *HTTPNero) Defaults() *HTTPNero {
	return hn
}

// Enabled returns whether a host was specified.
func (hn *HTTPNero) Enabled() bool {
	return hn.Host != ""
}

// HTTPNekos is a nekos HTTP API configuration section of the configuration file.
type HTTPNekos struct {
	// Host is the host string, used for http.ListenAndServe.
	Host string `toml:"host"`
}

// Defaults completes the section with default values.
func (hn *HTTPNekos) Defaults() *HTTPNekos {
	return hn
}

// Enabled returns whether a host was specified.
func (hn *HTTPNekos) Enabled() bool {
	return hn.Host != ""
}

// Repo is a base repository configuration.
type Repo struct {
	// Path is the relative or absolute path of the repository's directory.
	Path string `toml:"path"`
	// LockPath is the relative or absolute path of the repository's lock file.
	LockPath string `toml:"lock_path"`
	// Meta is the repository metadata.
	Meta map[string]string `toml:"meta"`
}

// Defaults completes the configuration with default values.
func (r *Repo) Defaults() *Repo {
	if r.LockPath == "" {
		r.LockPath = filepath.Join(r.Path, "nero.lock")
	}

	return r
}

// Parse parses the configuration from a file.
func Parse(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(filepath.Clean(path), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ParseWithDefaults parses the configuration from a file and completes it with default values (Section.Defaults).
func ParseWithDefaults(path string) (*Config, error) {
	cfg, err := Parse(path)
	if err != nil {
		return nil, err
	}

	return cfg.Defaults(), nil
}
