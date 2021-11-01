package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config provides an easy way to read and write configuration files for desktop applications.
//
// Config files are stored in the appropriate location (per os.UserConfigDir()), and automatically marshalled to and
// from YAML.
type Config struct {
	*options
	path string
}

// New creates a new Config that can be used to load and save configuration information.
// Options can be passed to customise behaviour, but the defaults are designed to be usable out of the box.
func New(o ...Option) (*Config, error) {
	opts := newOptions(o)

	p, err := opts.path()
	if err != nil {
		return nil, err
	}

	return &Config{
		options: opts,
		path:    p,
	}, nil
}

// Load reads the config file from disk, if it exists, and unmarshals it into the target struct.
// If the file does not exist on disk, no changes will be made to the target but no error will be returned.
func (c *Config) Load(target interface{}) error {
	b, err := os.ReadFile(c.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	return yaml.Unmarshal(b, target)
}

// Save marshals and writes the config file to disk. If the config directory does not exist it will be created.
func (c *Config) Save(i interface{}) error {
	b, err := yaml.Marshal(i)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(c.path), os.FileMode(0700)); err != nil {
		return err
	}

	return os.WriteFile(c.path, b, os.FileMode(0600))
}

// Directory returns the directory where the config file is stored.
func (c *Config) Directory() string {
    return filepath.Dir(c.path)
}

// Load is a convenience method that creates a new Config with the given options, and then immediately loads
// the config into the target struct.
func Load(target interface{}, o... Option) (*Config, error) {
	c, err := New(o...)
	if err != nil {
		return nil, err
	}

	if err := c.Load(target); err != nil {
		return nil, err
	}

	return c, nil
}
