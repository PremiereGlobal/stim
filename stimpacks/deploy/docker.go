package deploy

import (
	"bufio"
	"context"
	"fmt"

	"github.com/PremiereGlobal/stim/pkg/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
)

func (d *Deploy) startDeployContainer(instance *Instance) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		d.log.Fatal("Error creating docker client. {}", err)
	}

	// Pull the deploy image
	image := fmt.Sprintf("%s:%s", d.config.Deployment.Container.Repo, d.config.Deployment.Container.Tag)
	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
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
	cacheDir := fmt.Sprintf("%s/.kube-vault-deploy/bin-cache", home)
	err = utils.CreateDirIfNotExist(cacheDir, utils.UserGroupMode)
	if err != nil {
		d.log.Fatal("Could not create cache directory {}", cacheDir)
	}

	// Create the container spec
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        image,
		Cmd:          []string{fmt.Sprintf("./%s", d.config.Deployment.Script)},
		Tty:          true,
		Env:          envs,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			mount.Mount{
				Type:     mount.TypeBind,
				Source:   d.config.Deployment.fullDirectoryPath,
				Target:   "/scripts",
				ReadOnly: true,
			},
			mount.Mount{
				Type:     mount.TypeBind,
				Source:   cacheDir,
				Target:   "/bin-cache",
				ReadOnly: false,
			},
		},
	}, nil, "")
	if err != nil {
		d.log.Fatal("Error creating deploy container. {}", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		d.log.Fatal("Error starting deploy container. {}", err)
	}

	// Start capturing the logs
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{Follow: true, ShowStdout: true, ShowStderr: true})
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
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
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
