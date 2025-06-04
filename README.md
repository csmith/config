# Simple configuration for Go desktop apps

This package provides a simple configuration wrapper that marshals and unmarshals YAML to the user's config directory
(e.g. `~/.config/appname` on Linux or `%AppData\appname` on Windows).

Features:

* Only one dependency - `github.com/goccy/go-yaml` for yaml marshalling
* Stores files in the correct config directory by default, across OSes
* Sensible defaults let you get on with more important things

## Example

```go
package main

import "github.com/csmith/config"

type MyAppConfig struct {
	Name         string
	LikesMarmite bool
}

func main() {
	conf := &MyAppConfig{}
	
	// config.Load() can take options to customise the directory, filename, etc.
	// If you don't want to immediately load the config, you can use config.New(),
	// then call .Load on the config struct at a later time.
	c, err := config.Load(conf)
	if err != nil {
		panic(err)
	}

	if conf.Name == "" {
		// Prompt the user for missing settings, or use default values.
		// More advanced apps might want to populate a version field and
		// apply migrations/data entry in a more structured way.
		conf.Name = "Bob"
		conf.LikesMarmite = false
	}

	defer func() {
		// config.Save will write out the config to disk, creating directories if necessary.
		if err := c.Save(conf); err != nil {
			panic(err)
		}
	}()
}
```

## Options

A small set of options can be passed to `config.New()`:

### `config.DirectoryName(string)`

Sets the name of the directory to use. This is always rooted below the user's configuration directory.
For example, setting the directory name to "myapp" would store a config file at `~/.config/myapp/config.yml`
on Linux.

If not specified, defaults to `argv[0]`, i.e. the name of the binary being executed.

### `config.FileName(string)`

Sets the name of the file used. If not specified, defaults to `config.yml`. This can be useful if you
need to load multiple configuration files for one application.

### `config.Permissions(os.FileMode, os.FileMode)`

Changes the permissions that will be used for newly created directories and files, respectively.
Existing directories and files will not have their permissions changed.

If not specified, defaults to 0700 for directories and 0600 for files (i.e., read/write for the
owner only).

## Contributes

Feedback, questions, pull requests and bug reports are welcome! Please raise an issue on GitHub. 
