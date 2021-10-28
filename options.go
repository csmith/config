package config

import (
	"os"
	"path/filepath"
)

// An Option alters the default configuration of the config provider.
type Option func(*options)

// DirectoryName changes the name of the directory that the config file should be stored in.
//
// This is usually the application's name, and if not specified will default to os.Args[0].
func DirectoryName(n string) Option {
	return func(o *options) {
		o.directory = n
	}
}

// FileName changes the name of the file that the config is stored in.
//
// If not specified, defaults to "config.yml".
func FileName(n string) Option {
	return func(o *options) {
		o.filename = n
	}
}

// Permissions sets the filesystem permissions that will be set on newly created directories and files.
// Existing permissions will not be modified.
//
// If not specified, defaults to 0700 for directories and 0600 for files.
func Permissions(directoryMode os.FileMode, fileMode os.FileMode) Option {
	return func(o *options) {
		o.directoryMode = directoryMode
		o.fileMode = fileMode
	}
}

type options struct {
	directory     string
	directoryMode os.FileMode

	filename string
	fileMode os.FileMode
}

func (o *options) path() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, o.directory, o.filename), nil
}

func newOptions(o []Option) *options {
	c := &options{
		filename:      "config.yml",
		fileMode:      os.FileMode(0600),
		directory:     filepath.Base(os.Args[0]),
		directoryMode: os.FileMode(0700),
	}

	for i := range o {
		o[i](c)
	}

	return c
}
