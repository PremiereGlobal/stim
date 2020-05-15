package deploy

import (
	"bufio"
	"context"
	"fmt"

	"github.com/PremiereGlobal/stim/pkg/docker"
	"github.com/PremiereGlobal/stim/pkg/downloader"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
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

	var envs []string
	deprecatedHelmVersionSet := false
	for _, e := range instance.Spec.EnvironmentVars {
		if e.Name == "HELM_VERSION" {
			d.log.Warn("The use of the HELM_VERSION environment variable for specifying Helm versions has been deprecated.  Use the `.spec.tools.helm` configuration for specifying the helm version to use.  See https://github.com/PremiereGlobal/stim/blob/master/docs/DEPLOY.md for more details.")
			deprecatedHelmVersionSet = true
		}
		envs = append(envs, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}

	if _, ok := instance.Spec.Tools["helm"]; ok {
		if !deprecatedHelmVersionSet {
			envs = append(envs, fmt.Sprintf("HELM_VERSION=%s", downloader.GetBaseVersion(instance.Spec.Tools["helm"].Version)))
		} else {
			d.log.Warn("Both `spec.tools.helm` and the deprecated HELM_VERSION environment variable are set.  HELM_VERSION is taking precedence")
		}
	} else if !deprecatedHelmVersionSet {
		d.log.Warn("Auto-detection of Helm v2 versions now deprecated.  Use the `.spec.tools.helm` configuration for specifying the helm version to use. See https://github.com/PremiereGlobal/stim/blob/master/docs/DEPLOY.md for more details.")
		// DEPRECATION: Auto-matching helm version is deprecated and the env variable below should be uncommented
		// once this feature is removed
		// envs = append(envs, "HELM_MATCH_SERVER=false")
	}

	// Since we're using Docker, we need to mount the Linux binaries
	hostCacheDir := d.stim.ConfigGetCacheDir("bin/linux")
	cacheDir := "/bin-cache"
	workDir := "/scripts"
	pathDir := "/stim/path"

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
