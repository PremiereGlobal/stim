package deploy

import (
	"bufio"
	"context"
	"fmt"
	"path/filepath"

	"github.com/PremiereGlobal/stim/pkg/docker"
	"github.com/PremiereGlobal/stim/pkg/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/mitchellh/go-homedir"
	// "github.com/PremiereGlobal/stim/stim"
)

func (d *Deploy) startDeployContainer(instance *Instance) {

	dockerClient, err := docker.NewClient()
	if err != nil {
		d.log.Fatal("Error creating docker client. {}", err)
	}

	ctx := context.Background()

	// Pull the deploy image
	image := fmt.Sprintf("%s:%s", d.config.Deployment.Container.Repo, d.config.Deployment.Container.Tag)
	reader, err := dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		d.log.Fatal("Failed to pull deploy image. {}", err)
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		d.log.Debug(scanner.Text())
	}

	envs := make([]string, len(instance.Spec.EnvironmentVars))
	for i, e := range instance.Spec.EnvironmentVars {
		envs[i] = fmt.Sprintf("%s=%s", e.Name, e.Value)
	}

	// Get the "home" directory for storing cache files
	home, err := homedir.Dir()
	if err != nil {
		d.log.Fatal("Unable to determine home directory. {}", err)
	}
	cacheDir := filepath.Join(home, ".stim/cache/linux")
	err = utils.CreateDirIfNotExist(cacheDir, utils.UserGroupMode)
	if err != nil {
		d.log.Fatal("Could not create cache directory {}", cacheDir)
	}

	hostCacheDir := fmt.Sprintf("%s/.kube-vault-deploy/bin-cache", home)
	hostCacheDir = filepath.Join(home, ".stim/cache/linux")
	cacheDir = "/stim/cache/"
	cacheDir = "/bin-cache"
	workDir := "/stim/deploy"
	workDir = "/scripts"
	pathDir := "/stim/path"

	// Make an "env" on the host to hold the in-container symlinks
	// e := d.stim.Env(&stim.EnvConfig{})
	// defer e.Close()
	// err = e.LinkWithPathOverride(filepath.Join(hostCacheDir, "kubectl-v1.10.8"), cacheDir, "kubectl")
	// if err != nil {
	// 	d.log.Fatal("Unable to link to path. {}", err)
	// }
	// d.log.Warn(e.Run(fmt.Sprintf("ls -alh %s", e.GetPath())))

	// Create the container spec

	cmd := []string{"/bin/sh", "-c", fmt.Sprintf("export PATH=%s:${PATH}; ./%s", pathDir, d.config.Deployment.Script)}
	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image:        image,
		Cmd:          cmd,
		Tty:          true,
		Env:          envs,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   workDir,
	}, &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			mount.Mount{
				Type:     mount.TypeBind,
				Source:   d.config.Deployment.fullDirectoryPath,
				Target:   workDir,
				ReadOnly: false, // This could be set to false when the downloads don't go here
			},
			mount.Mount{
				Type:     mount.TypeBind,
				Source:   hostCacheDir,
				Target:   cacheDir,
				ReadOnly: false,
			},
			// mount.Mount{
			// 	Type:     mount.TypeBind,
			// 	Source:   e.GetPath()+"/",
			// 	Target:   pathDir,
			// 	ReadOnly: true,
			// },
		},
	}, nil, "")
	if err != nil {
		d.log.Fatal("Error creating deploy container. {}", err)
	}

	// Start the container
	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		d.log.Fatal("Error starting deploy container. {}", err)
	}

	// Start capturing the logs
	out, err := dockerClient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{Follow: true, ShowStdout: true, ShowStderr: true})
	if err != nil {
		d.log.Fatal("Error getting container logs. {}", err)
	}
	defer out.Close()

	d.log.Info("--- START Stim deploy - Docker container logs ---")
	scanner = bufio.NewScanner(out)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	d.log.Info("--- END Stim deploy - Docker container logs ---")

	// Wait for the container to finish
	statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			d.log.Fatal("Deploy container error. {}", err)
		}
	case status := <-statusCh:
		if status.Error != nil {
			d.log.Fatal("Deployment resulted in error. {}. Halting any further deployments...", status.Error.Message)
		}
		if status.StatusCode != 0 {
			d.log.Fatal("Deployment to '{}' resulted in non-zero exit code {}. Halting any further deployments...", instance.Name, status.StatusCode)
		}
	}

}
