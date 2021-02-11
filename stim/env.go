package stim

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PremiereGlobal/stim/pkg/downloader"
	"github.com/PremiereGlobal/stim/pkg/env"
	"github.com/PremiereGlobal/stim/pkg/kubernetes"
	"github.com/PremiereGlobal/vault-to-envs/pkg/vaulttoenvs"
)

// EnvConfig represets a environment configuration
type EnvConfig struct {

	// EnvVars slice of environment variables to be set
	// Values look like "FOO=BAR"
	EnvVars []string

	// Kubernetes env config
	Kubernetes *EnvConfigKubernetes

	// Kubernetes env config
	Vault *EnvConfigVault

	// WorkDir sets the working directory where any shell commands will be executed
	WorkDir string

	// Tools should contains a list of supported binary tools to install and link
	Tools map[string]EnvTool
}

// EnvConfig represets a environment's Kubernetes configuration
type EnvConfigKubernetes struct {

	// Cluster is the name of the Kubernetes cluster to set up
	Cluster string

	// ServiceAccount to use when connecting to Kubernetes
	ServiceAccount string

	// DefaultNamespace to use when setting up Kubernetes
	DefaultNamespace string
}

// EnvConfig represets a environment's Vault configuration
type EnvConfigVault struct {

	// SecretItems to load into the environment
	SecretItems []*vaulttoenvs.SecretItem
}

// EnvTool contains the configuration for a CLI tool
type EnvTool struct {
	Version string `yaml:"version"`
	Unset   bool   `yaml:"unset"`
}

// Env sets up an environment based on the given config
// Shell commands can be executed against the environment
func (stim *Stim) Env(config *EnvConfig) *env.Env {

	e, err := env.New(env.Config{})
	if err != nil {
		stim.log.Fatal("Stim: Error creating new environment.", err)
	}

	e.SetWorkDir(config.WorkDir)
	e.AddEnvVars(config.EnvVars...)

	// If requiring Kubernetes, set things up
	var kc *kubernetes.Config
	if config.Kubernetes != nil {

		// This is the path where the kubeconfig will be written
		kubeConfigFilePath := filepath.Join(e.GetPath(), "kubeconfig")

		vault := stim.Vault()

		// Get the Kubernetes creds from Vault
		secretValues, err := vault.GetSecretKeys("secret/kubernetes/" + config.Kubernetes.Cluster + "/" + config.Kubernetes.ServiceAccount + "/kube-config")
		if err != nil {
			stim.log.Fatal("Stim: Error getting kubeconfig secrets for environment. {}", err)
		}

		// If namespace not set use the default from Vault
		defaultNamespace := config.Kubernetes.DefaultNamespace
		if defaultNamespace == "" {
			defaultNamespace = secretValues["default-namespace"]
		}

		// Build the Kube config options
		kubeConfigOptions := &kubernetes.ConfigOptions{
			ClusterName:             config.Kubernetes.Cluster,
			ClusterServer:           secretValues["cluster-server"],
			ClusterCA:               secretValues["cluster-ca"],
			AuthName:                config.Kubernetes.Cluster + "-" + config.Kubernetes.ServiceAccount,
			AuthToken:               secretValues["user-token"],
			ContextName:             config.Kubernetes.Cluster,
			ContextSetCurrent:       true,
			ContextDefaultNamespace: defaultNamespace,
		}

		kc = kubernetes.NewConfigFromPath(kubeConfigFilePath)
		err = kc.Modify(kubeConfigOptions)
		if err != nil {
			stim.log.Fatal("Stim: Error writing kubeconfig for environment. {}", err)
		}

		// Tell the environment to use the kubeconfig in the environment PATH
		e.AddEnvVars([]string{fmt.Sprintf("%s=%s", "KUBECONFIG", kubeConfigFilePath)}...)
	}

	// If requiring secrets, set those up
	if config.Vault != nil && len(config.Vault.SecretItems) > 0 {

		vault := stim.Vault()

		vaultAddress, err := vault.GetAddress()
		if err != nil {
			stim.log.Fatal("Stim: Unable to get Vault address for environment. {}", err)
		}

		vaultToken, err := vault.GetToken()
		if err != nil {
			stim.log.Fatal("Stim: Unable to get Vault token for environment. {}", err)
		}

		v2e := vaulttoenvs.NewVaultToEnvs(&vaulttoenvs.Config{
			VaultAddr: vaultAddress,
		})
		v2e.SetVaultToken(vaultToken)
		v2e.AddSecretItems(config.Vault.SecretItems...)

		sleepTime := time.Duration(time.Second)

		var secretEnvs []string

		for i := 0; i < 3; i++ {
			secretEnvs, err = v2e.GetEnvs()
			if err != nil {
				if stim.ConfigGetBool("vault.retryOnThrottle") && strings.Contains(err.Error(), "Throttling: Rate exceeded") {
					stim.log.Info("Stim: Got Throttling error waiting {} then trying again, try number:{}", sleepTime, i+1)
					time.Sleep(sleepTime)
					sleepTime += time.Duration(time.Second)
					continue
				}
				stim.log.Fatal("Stim: Unable to get Vault secrets for environment. {}", err)
			}
			break
		}

		e.AddEnvVars(secretEnvs...)
	}

	// if requiring any CLI tools, download and link them here
	for toolName, toolParams := range config.Tools {

		version := toolParams.Version
		if version == "" {
			stim.log.Debug("Detecting tool version for: {}", toolName)
		} else {
			stim.log.Debug("Setting tool version {}/{} based on configuration", toolName, version)
		}

		var dl downloader.Downloader
		cacheDir := stim.ConfigGetCacheDir(filepath.Join("bin", runtime.GOOS))
		switch toolName {
		case "vault":
			if version == "" {
				version, err = stim.Vault().Version()
				if err != nil {
					stim.log.Fatal("Unable to determine version for {}: {}", toolName, err)
				}
			}
			dl = downloader.NewVaultDownloader(version, cacheDir)
		case "kubectl":
			if version == "" {
				if kc == nil {
					stim.log.Fatal("Kubernetes server not specified, cannot determine version")
				}
				k, err := kubernetes.New(kc)
				if err != nil {
					stim.log.Fatal("Unable to load Kube config, cannot determine version")
				}
				version, err = k.Version()
				if err != nil {
					stim.log.Fatal("Unable to determine version for {}: {}", toolName, err)
				}
			}
			dl = downloader.NewKubeDownloader(version, cacheDir)
		case "helm":
			if version == "" {
				stim.log.Fatal("Version detection not supported for helm, please specify a version in the config")
			}
			dl = downloader.NewHelmDownloader(version, cacheDir)
		default:
			stim.log.Fatal("Unknown deploy tool: {}", toolName)
		}

		result, err := dl.Download()
		if err != nil {
			stim.log.Fatal("Download failed: {} {}", result, err)
		}
		if !result.FileExists {
			stim.log.Debug("Downloaded {} in {}", result.RenderedURL, result.DownloadDuration)
		}
		stim.log.Debug("Linking binary from {} to PATH location {}/{}", dl.GetBinPath(), e.GetPath(), toolName)
		e.Link(dl.GetBinPath(), toolName)
	}

	return e
}
