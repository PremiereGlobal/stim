package deploy

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
	// "github.com/fatih/color"
	// "github.com/davecgh/go-spew/spew"
	"os"
)

func (d *Deploy) startDeployContainer(environment *Environment, instance *Instance) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Pull the deploy image
	image := fmt.Sprintf("%s:%s", d.config.Container.Repo, d.config.Container.Tag)
	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	// color.Set(color.FgYellow)
	io.Copy(os.Stdout, reader)

	// Build our env variables
	// comboenvs := instance.EnvSpec.EnvironmentVars
	// comboenvs = append(comboenvs, ...)

	envs := make([]string, len(instance.EnvSpec.EnvironmentVars))
	// envs[0] = fmt.Sprintf("VAULT_TOKEN=%s", token)
	// envs[1] = fmt.Sprintf("SECRET_CONFIG=%s", secretConfig)
	// envs[2] = fmt.Sprintf("DEPLOY_CLUSTER=%s", instance.EnvSpec.Kubernetes.Cluster)
	for i, e := range instance.EnvSpec.EnvironmentVars {
		envs[i] = fmt.Sprintf("%s=%s", e.Name, e.Value)
	}

	// Get the "home" directory for storing cache files
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
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
				Source:   fmt.Sprintf("%s/.kube-vault-deploy/bin-cache", home),
				Target:   "/bin-cache",
				ReadOnly: false,
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// Start capturing the logs
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{Follow: true, ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Stream the logs as they come in
	// color.Set(color.FgCyan)
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case status := <-statusCh:
		if status.Error != nil {
			d.stim.Fatal(errors.New(fmt.Sprintf("Deployment resulted in error. %s. Halting any further deployments...", status.Error.Message)))
		}
		if status.StatusCode != 0 {
			d.stim.Fatal(errors.New(fmt.Sprintf("Deployment to '%s' resulted in non-zero exit code %d. Halting any further deployments...", instance.Name, status.StatusCode)))
		}
	}

}
