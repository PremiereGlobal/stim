package env

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/PremiereGlobal/stim/pkg/shell"
)

type Env struct {
	config Config
}

type Config struct {
	Path    *Path
	EnvVars []string
}

type Path struct {
	Directory     string
	Prefix        string
	RemoveOnClose bool
}

func New(config Config) (*Env, error) {

	// If path not set, use a temp dir
	if config.Path == nil {

		dir, err := ioutil.TempDir("", "")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to create env path directory: %v", err))
		}

		config.Path = &Path{
			Directory:     dir,
			RemoveOnClose: true,
		}

	}

	env := Env{config: config}

	return &env, nil
}

func (e *Env) AddEnvVars(envvars ...string) {
	e.config.EnvVars = append(e.config.EnvVars, envvars...)
}

func (e *Env) GetPath() string {
	return e.config.Path.Directory
}

// Link creates a symlink in the environment path to the source file
// if newname is empty, will use the same filename
func (e *Env) Link(source string, newname string) error {

	if newname == "" {
		newname = filepath.Base(source)
	}

	err := os.Symlink(source, filepath.Join(e.GetPath(), newname))

	return err
}

// GetEnvVars combines the user provided envs with the environment PATH
func (e *Env) GetEnvVars() []string {

	// Merge environment PATH with user PATH
	path := fmt.Sprintf("PATH=%s:%s", e.GetPath(), os.Getenv("PATH"))

	envs := append([]string{path}, e.config.EnvVars...)

	return envs
}

// Run a shell command in the environment
func (e *Env) Run(cmdString string) (string, error) {

	s, err := shell.Run(shell.ShellCommand{
		Command: []string{cmdString},
		Envs:    e.GetEnvVars(),
	})

	return s, err
}

// Close cleans up resources created by the env
func (e *Env) Close() {
	if e.config.Path.RemoveOnClose {
		os.Remove(e.config.Path.Directory)
	}
}
