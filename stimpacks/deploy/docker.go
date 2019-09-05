package deploy

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	// "github.com/docker/docker/pkg/stdcopy"
	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"os"
)

func (d *Deploy) startDeployContainer(cluster *Cluster) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	image := fmt.Sprintf("%s:%s", d.config.Container.Repo, d.config.Container.Tag)

	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	color.Set(color.FgYellow)
	io.Copy(os.Stdout, reader)

	// Get the raw Vault token
	vault := d.stim.Vault()
	token, err := vault.GetToken()
	if err != nil {
		panic(err)
	}

	// Build the secret config
	secretConfig, err := d.makeSecretConfig(cluster)
	if err != nil {
		panic(err)
	}

	envs := make([]string, len(cluster.EnvSpec.EnvironmentVars)+2)
	envs[0] = fmt.Sprintf("VAULT_TOKEN=%s", token)
	envs[1] = fmt.Sprintf("SECRET_CONFIG=%s", secretConfig)
	for i, e := range cluster.EnvSpec.EnvironmentVars {
		envs[i+2] = fmt.Sprintf("%s=%s", e.Name, e.Value)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{fmt.Sprintf("./%s", d.config.Deployment.Script)},
		Tty:   true,
		Env:   envs,
	}, &container.HostConfig{
		// AutoRemove: true,
		Mounts: []mount.Mount{
			mount.Mount{
				Type:     mount.TypeBind,
				Source:   d.config.Deployment.fullDirectoryPath,
				Target:   "/scripts",
				ReadOnly: true,
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
		spew.Dump("Container done")
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	color.Set(color.FgCyan)
	_, err = io.Copy(os.Stdout, out)
	if err != nil && err != io.EOF {
		panic(err)
	}

	// stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	// spew.Dump(out)
	// containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	//
	// for _, container := range containers {
	// 	fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	// }
}
