package sandbox

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

// drainPullStream consumes the JSON message stream returned by ImagePull and
// surfaces in-stream failures. Docker reports mid-pull errors (network
// interruption, disk full, ...) as {"error": ...} messages inside an HTTP
// 200 stream, so simply discarding the stream would hide them and later
// produce a misleading "No such image" error at container create time.
func drainPullStream(r io.Reader) error {
	dec := json.NewDecoder(r)
	for {
		var msg struct {
			Error string `json:"error"`
		}
		if err := dec.Decode(&msg); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("decode pull progress: %w", err)
		}
		if msg.Error != "" {
			return errors.New(msg.Error)
		}
	}
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
		drainErr := drainPullStream(reader)
		reader.Close()
		if drainErr != nil {
			cli.Close()
			return nil, fmt.Errorf("image pull %s: %w", image, drainErr)
		}
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: image,
			// The container lifetime is owned by Destroy; a finite sleep
			// would kill long learning sessions (and every exec in them)
			// after it elapses.
			Cmd: []string{"sleep", "infinity"},
			Tty: false,
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
		// Remove the created container on an independent context — ctx may
		// already be cancelled, which is exactly when cleanup matters most.
		cleanupCtx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
		cli.ContainerRemove(cleanupCtx, resp.ID, container.RemoveOptions{Force: true})
		cancel()
		cli.Close()
		return nil, fmt.Errorf("container start: %w", err)
	}

	return &DockerSandbox{cli: cli, containerID: resp.ID}, nil
}

// ContainerID returns the ID of the underlying Docker container.
func (s *DockerSandbox) ContainerID() string { return s.containerID }

// Exec runs a command inside the container and returns its combined
// stdout+stderr, exit code, and any error. Output is merged to match the
// CombinedOutput semantics of LocalSandbox and ComposeSandbox.
func (s *DockerSandbox) Exec(ctx context.Context, command string) (string, int, error) {
	execCfg := container.ExecOptions{
		Cmd:          []string{"bash", "-c", command},
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

	// StdCopy blocks on the hijacked connection, which the docker client
	// does not tie to ctx — run it in a goroutine and force-close the
	// connection on cancellation so a blocking command cannot hang us
	// forever and the goroutine cannot leak.
	var output bytes.Buffer
	copyDone := make(chan error, 1)
	go func() {
		// Both writers share one buffer: merge stdout and stderr.
		_, copyErr := stdcopy.StdCopy(&output, &output, attachResp.Reader)
		copyDone <- copyErr
	}()

	select {
	case <-ctx.Done():
		attachResp.Close() // unblocks StdCopy inside the goroutine
		<-copyDone
		return "", -1, ctx.Err()
	case err = <-copyDone:
		if err != nil {
			return "", -1, fmt.Errorf("read output: %w", err)
		}
	}

	inspectResp, err := s.cli.ContainerExecInspect(ctx, execResp.ID)
	if err != nil {
		return output.String(), -1, fmt.Errorf("exec inspect: %w", err)
	}

	return output.String(), inspectResp.ExitCode, nil
}

// Destroy stops and removes the container.
// Cleanup runs on an independent timeout context so it still succeeds when
// the caller's context has already expired or been cancelled — otherwise a
// timed-out challenge would leak its container.
func (s *DockerSandbox) Destroy(_ context.Context) error {
	cleanupCtx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
	defer cancel()
	err := s.cli.ContainerRemove(cleanupCtx, s.containerID, container.RemoveOptions{Force: true})
	s.cli.Close()
	return err
}

// InteractiveShellArgs returns the command arguments to open an interactive shell in the container.
// Uses /tmp/.linuxlab_bashrc if it exists (injected by the TUI before launch).
func (s *DockerSandbox) InteractiveShellArgs() []string {
	return []string{"docker", "exec", "-it", s.containerID, "/bin/bash", "--rcfile", "/tmp/.linuxlab_bashrc"}
}
