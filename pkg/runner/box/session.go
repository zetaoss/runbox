package box

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/distribution/reference"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"k8s.io/utils/ptr"
)

const staleAgeLimitSeconds int = 300

type Session struct {
	opts      *Opts
	cli       *client.Client
	ctx       context.Context
	id        string
	startTime time.Time
	startCPU  int
	result    Result
}

func NewSession(cli *client.Client, opts *Opts) *Session {
	if opts.CollectStats == nil {
		opts.CollectStats = ptr.To(true)
	}
	if opts.CollectImagesCount == 0 {
		opts.CollectImagesCount = 2
	}
	if opts.PullImageIfNotPresent == nil {
		opts.PullImageIfNotPresent = ptr.To(true)
	}
	if opts.Shell == "" {
		opts.Shell = "sh"
	}
	if opts.Timeout == 0 {
		opts.Timeout = 60000 // 60s
	}
	return &Session{
		cli:  cli,
		opts: opts,
		ctx:  context.Background(),
	}
}

func (s *Session) pruneStaleContainers() {
	log.Println("Starting to prune stale containers...")

	ctx := context.Background()
	containers, err := s.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		log.Printf("Failed to list containers: %v", err)
		return
	}

	for _, ct := range containers {
		containerID := ct.ID[:10]
		containerAge := time.Since(time.Unix(ct.Created, 0))

		if ct.State == "removing" || containerAge > time.Duration(staleAgeLimitSeconds)*time.Second {
			containerName := "<unknown>"
			if len(ct.Names) > 0 {
				containerName = ct.Names[0]
			}
			log.Printf("Removing container ID: %s, Name: %s, Age: %v...", containerID, containerName, containerAge)

			// Force remove container
			err := s.cli.ContainerRemove(ctx, ct.ID, container.RemoveOptions{Force: true})
			if err != nil {
				log.Printf("Failed to remove container ID: %s, Name: %s, Error: %v", containerID, containerName, err)
			} else {
				log.Printf("Successfully removed container ID: %s, Name: %s", containerID, containerName)
			}
		}
	}
	log.Println("Container prune operation completed.")
}

func (s *Session) run() error {
	if err := s.checkImage(); err != nil {
		return fmt.Errorf("checkImage err: %w", err)
	}
	if err := s.createContainer(); err != nil {
		return fmt.Errorf("createContainer err: %w", err)
	}
	if err := s.copyFiles(); err != nil {
		return fmt.Errorf("copyFiles err: %w", err)
	}
	if err := s.cli.ContainerStart(s.ctx, s.id, container.StartOptions{}); err != nil {
		return fmt.Errorf("ContainerStart err: %w", err)
	}
	defer func() {
		_ = s.cli.ContainerRemove(s.ctx, s.id, container.RemoveOptions{Force: true})
	}()
	if err := s.execute(); err != nil {
		return fmt.Errorf("execute err: %w", err)
	}
	if err := s.collectImages(); err != nil {
		return fmt.Errorf("getImages err: %w", err)
	}
	return nil
}

func (s *Session) checkImage() error {
	ok, err := s.findImage()
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	if *s.opts.PullImageIfNotPresent {
		return s.pullImage()
	}
	return fmt.Errorf("no image: '%s'", s.opts.Image)
}

func (s *Session) findImage() (bool, error) {
	name := s.opts.Image
	_, err := reference.ParseNormalizedNamed(name)
	if err != nil {
		return false, err // invalid reference format
	}
	if !strings.Contains(name, ":") {
		name = name + ":latest"
	}
	images, err := s.cli.ImageList(s.ctx, image.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == name {
				return true, nil
			}
		}
	}
	return false, nil
}

func (s *Session) pullImage() error {
	out, err := s.cli.ImagePull(s.ctx, s.opts.Image, image.PullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.ReadAll(out); err != nil {
		return err
	}
	return nil
}

func (s *Session) createContainer() error {
	resp, err := s.cli.ContainerCreate(s.ctx, &container.Config{
		Image:      s.opts.Image,
		Cmd:        []string{"sleep", strconv.Itoa(staleAgeLimitSeconds)},
		WorkingDir: s.opts.WorkingDir,
		User:       s.opts.User,
	}, &container.HostConfig{
		AutoRemove: true,
		Resources: container.Resources{
			PidsLimit: ptr.To(int64(100)),
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}
	s.id = resp.ID
	return nil
}

func (s *Session) copyFiles() error {
	tarBuffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(tarBuffer)
	defer tarWriter.Close()

	for _, file := range s.opts.Files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0644,
			Size: int64(len(file.Body)),
		}
		if err := tarWriter.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tarWriter.Write([]byte(file.Body)); err != nil {
			return err
		}
	}
	if err := tarWriter.Close(); err != nil {
		return err
	}

	if err := s.cli.CopyToContainer(s.ctx, s.id, "/", tarBuffer, container.CopyToContainerOptions{}); err != nil {
		return err
	}
	return nil
}

func (s *Session) execute() error {
	execOpts := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{s.opts.Shell, "-c", s.opts.Command},
	}
	exec, err := s.cli.ContainerExecCreate(s.ctx, s.id, execOpts)
	if err != nil {
		return err
	}
	if err := s.collectStatsStart(); err != nil {
		return err
	}
	attach, err := s.cli.ContainerExecAttach(s.ctx, exec.ID, container.ExecStartOptions{})
	if err != nil {
		return err
	}
	defer attach.Close()

	stdout := &logWriter{stream: 1, logs: &s.result.Logs}
	stderr := &logWriter{stream: 2, logs: &s.result.Logs}

	ctx, cancel := context.WithTimeout(s.ctx, time.Duration(s.opts.Timeout)*time.Millisecond)
	defer cancel()

	s.startTime = time.Now()
	done := make(chan error, 1)
	go func() {
		_, err := stdcopy.StdCopy(stdout, stderr, attach.Reader)
		done <- err
	}()

	select {
	case <-ctx.Done():
		s.result.Timedout = true
	case err := <-done:
		if err != nil {
			return err
		}
	}

	s.result.Time = int(time.Since(s.startTime).Milliseconds())
	stdout.Close()
	stderr.Close()
	if err := s.collectStatsEnd(); err != nil {
		return err
	}
	resp, err := s.cli.ContainerExecInspect(s.ctx, exec.ID)
	if err != nil {
		return err
	}
	s.result.Code = resp.ExitCode
	return nil
}

func (s *Session) collectStatsStart() error {
	if !*s.opts.CollectStats {
		return nil
	}
	cpu, _, err := s.getStats()
	if err != nil {
		return err
	}
	s.startCPU = cpu
	return nil
}

func (s *Session) collectStatsEnd() error {
	if !*s.opts.CollectStats {
		return nil
	}
	endCPU, mem, err := s.getStats()
	if err != nil {
		return err
	}
	s.result.CPU = endCPU - s.startCPU
	s.result.MEM = mem
	return nil
}

func (s *Session) getStats() (int, int, error) {
	stats, err := s.cli.ContainerStatsOneShot(s.ctx, s.id)
	if err != nil {
		return 0, 0, err
	}
	defer stats.Body.Close()

	var v container.StatsResponse
	if err := json.NewDecoder(stats.Body).Decode(&v); err != nil {
		return 0, 0, err
	}
	cpu := int(v.CPUStats.CPUUsage.TotalUsage / 1000) // core*milliseconds
	mem := int(v.MemoryStats.Usage / 1024)            // kibibytes
	return cpu, mem, nil
}

func (s *Session) collectImages() error {
	if !s.opts.CollectImages {
		return nil
	}
	reader, _, err := s.cli.CopyFromContainer(s.ctx, s.id, s.opts.WorkingDir)
	if err != nil {
		return err
	}
	defer reader.Close()

	tr := tar.NewReader(reader)
	maxFileSize := 100 * 1024 // 100KiB

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if !header.FileInfo().IsDir() && filepath.Ext(header.Name) == ".png" && header.Size <= int64(maxFileSize) {
			var fileContent bytes.Buffer
			if _, err := io.Copy(&fileContent, tr); err != nil {
				return err
			}
			encoded := base64.StdEncoding.EncodeToString(fileContent.Bytes())
			s.result.Images = append(s.result.Images, encoded)
			if len(s.result.Images) >= s.opts.CollectImagesCount {
				break
			}
		}
	}
	return nil
}
