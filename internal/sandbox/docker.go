package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	dockerimage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

const defaultImage = "ubuntu:22.04"

// DockerSandbox manages a Docker container used as an isolated sandbox.
type DockerSandbox struct {
	cli         *client.Client
	containerID string
}

// NewDockerSandbox creates and starts a new Docker container.
// If image is empty, defaultImage is used.
func NewDockerSandbox(ctx context.Context, image string) (*DockerSandbox, error) {
	if image == "" {
		image = defaultImage
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("docker client: %w", err)
	}

	// Auto-pull image if not present locally
	_, _, inspectErr := cli.ImageInspectWithRaw(ctx, image)
	if inspectErr != nil {
		reader, pullErr := cli.ImagePull(ctx, image, dockerimage.PullOptions{})
		if pullErr != nil {
			cli.Close()
			return nil, fmt.Errorf("image pull %s: %w", image, pullErr)
		}
		io.Copy(io.Discard, reader)
		reader.Close()
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: image,
			Cmd:   []string{"sleep", "3600"},
			Tty:   false,
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:   512 * 1024 * 1024,
				NanoCPUs: 1_000_000_000,
			},
		},
		nil, nil, "",
	)
	if err != nil {
		cli.Close()
		return nil, fmt.Errorf("container create: %w", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		cli.Close()
		return nil, fmt.Errorf("container start: %w", err)
	}

	return &DockerSandbox{cli: cli, containerID: resp.ID}, nil
}

// ContainerID returns the ID of the underlying Docker container.
func (s *DockerSandbox) ContainerID() string { return s.containerID }

// Exec runs a command inside the container and returns stdout, exit code, and any error.
func (s *DockerSandbox) Exec(ctx context.Context, command string) (string, int, error) {
	execCfg := container.ExecOptions{
		Cmd:          []string{"sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
	}

	execResp, err := s.cli.ContainerExecCreate(ctx, s.containerID, execCfg)
	if err != nil {
		return "", -1, fmt.Errorf("exec create: %w", err)
	}

	attachResp, err := s.cli.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", -1, fmt.Errorf("exec attach: %w", err)
	}
	defer attachResp.Close()

	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, attachResp.Reader)
	if err != nil {
		return "", -1, fmt.Errorf("read output: %w", err)
	}

	inspectResp, err := s.cli.ContainerExecInspect(ctx, execResp.ID)
	if err != nil {
		return stdout.String(), -1, fmt.Errorf("exec inspect: %w", err)
	}

	return stdout.String(), inspectResp.ExitCode, nil
}

// Destroy stops and removes the container.
func (s *DockerSandbox) Destroy(ctx context.Context) error {
	err := s.cli.ContainerRemove(ctx, s.containerID, container.RemoveOptions{Force: true})
	s.cli.Close()
	return err
}

// InteractiveShellArgs returns the command arguments to open an interactive shell in the container.
func (s *DockerSandbox) InteractiveShellArgs() []string {
	return []string{"docker", "exec", "-it", s.containerID, "/bin/bash"}
}
