package env

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/PremiereGlobal/stim/pkg/shell"
	"github.com/PremiereGlobal/stim/pkg/utils"
)

// Env represents the Env type
type Env struct {
	config Config
}

// Config represents the an environment configuration
type Config struct {
	Path    *Path
	EnvVars []string
	WorkDir string
}

// Path represents the environment's PATH configuration
type Path struct {
	Directory     string
	Prefix        string
	RemoveOnClose bool
}

// New configures a new environment
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

// AddEnvVars adds a list of environment variables
func (e *Env) AddEnvVars(envvars ...string) {
	e.config.EnvVars = append(e.config.EnvVars, envvars...)
}

// GetEnvVar returns the value of the given environment variable in the current environment
func (e *Env) GetEnvVar(name string) string {
	envs := e.GetEnvVars()
	for _, v := range envs {
		parts := strings.Split(v, "=")
		if parts[0] == name {
			return parts[1]
		}
	}

	return ""
}

// GetPath returns the current PATH directory
func (e *Env) GetPath() string {
	return e.config.Path.Directory
}

// Link creates a symlink in the environment path to the source file
// if linkName is empty, will use the same filename
func (e *Env) Link(source string, linkFileName string) error {

	if linkFileName == "" {
		linkFileName = filepath.Base(source)
	}

	err := utils.EnsureLink(source, filepath.Join(e.GetPath(), linkFileName))

	return err
}

// GetEnvVars provides all the environment variables (including the generated PATH)
func (e *Env) GetEnvVars() []string {
	path := fmt.Sprintf("PATH=%s:%s", e.GetPath(), os.Getenv("PATH"))
	envs := append([]string{path}, e.config.EnvVars...)
	return envs
}

// Run runs a shell command in the environment
func (e *Env) Run(cmdString string) (string, error) {

	fullCmd := fmt.Sprintf("cd %s && %s", e.config.WorkDir, cmdString)
	s, err := shell.Run(shell.ShellCommand{
		Command: []string{fullCmd},
		Envs:    e.GetEnvVars(),
	})

	return s, err
}

// SetWorkDir sets the current working directory
func (e *Env) SetWorkDir(workDir string) {
	e.config.WorkDir = workDir
}

// Close cleans up resources created by the env
func (e *Env) Close() {
	if e.config.Path.RemoveOnClose {
		os.Remove(e.config.Path.Directory)
	}
}
