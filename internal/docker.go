package internal

import (
	"bytes"
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func runDockerContainer(ctx context.Context, tag string) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", errors.Wrap(err, "client.NewEnvClient")
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: tag,
		Cmd:   []string{"/bin/app"},
		Tty:   true,
	}, nil, nil, "")
	if err != nil {
		return "", errors.Wrap(err, "cli.ContainerCreate")
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.Wrap(err, "cli.ContainerStart")
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", errors.Wrap(err, "cli.ContainerWait")
		}

	case <-statusCh:
	}

	logs, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", errors.Wrap(err, "cli.ContainerLogs")
	}
	defer logs.Close()

	var out bytes.Buffer
	if _, err := io.Copy(&out, logs); err != nil {
		return "", errors.Wrap(err, "io.Copy")
	}

	return out.String(), nil
}
