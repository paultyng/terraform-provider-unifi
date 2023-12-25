package docker

import (
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/docker/buildx/driver"
	"github.com/docker/buildx/driver/bkimage"
	"github.com/docker/buildx/util/confutil"
	"github.com/docker/buildx/util/imagetools"
	"github.com/docker/buildx/util/progress"
	"github.com/docker/cli/opts"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	dockerarchive "github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/tracing/detect"
	"github.com/pkg/errors"
)

const (
	volumeStateSuffix   = "_state"
	buildkitdConfigFile = "buildkitd.toml"
)

type Driver struct {
	driver.InitConfig
	factory driver.Factory

	// if you add fields, remember to update docs:
	// https://github.com/docker/docs/blob/main/content/build/drivers/docker-container.md
	netMode      string
	image        string
	memory       opts.MemBytes
	memorySwap   opts.MemSwapBytes
	cpuQuota     int64
	cpuPeriod    int64
	cpuShares    int64
	cpusetCpus   string
	cpusetMems   string
	cgroupParent string
	env          []string
}

func (d *Driver) IsMobyDriver() bool {
	return false
}

func (d *Driver) Config() driver.InitConfig {
	return d.InitConfig
}

func (d *Driver) Bootstrap(ctx context.Context, l progress.Logger) error {
	return progress.Wrap("[internal] booting buildkit", l, func(sub progress.SubLogger) error {
		_, err := d.DockerAPI.ContainerInspect(ctx, d.Name)
		if err != nil {
			if dockerclient.IsErrNotFound(err) {
				return d.create(ctx, sub)
			}
			return err
		}
		return sub.Wrap("starting container "+d.Name, func() error {
			if err := d.start(ctx, sub); err != nil {
				return err
			}
			return d.wait(ctx, sub)
		})
	})
}

func (d *Driver) create(ctx context.Context, l progress.SubLogger) error {
	imageName := bkimage.DefaultImage
	if d.image != "" {
		imageName = d.image
	}

	if err := l.Wrap("pulling image "+imageName, func() error {
		ra, err := imagetools.RegistryAuthForRef(imageName, d.Auth)
		if err != nil {
			return err
		}
		rc, err := d.DockerAPI.ImageCreate(ctx, imageName, dockertypes.ImageCreateOptions{
			RegistryAuth: ra,
		})
		if err != nil {
			return err
		}
		_, err = io.Copy(io.Discard, rc)
		return err
	}); err != nil {
		// image pulling failed, check if it exists in local image store.
		// if not, return pulling error. otherwise log it.
		_, _, errInspect := d.DockerAPI.ImageInspectWithRaw(ctx, imageName)
		if errInspect != nil {
			return err
		}
		l.Wrap("pulling failed, using local image "+imageName, func() error { return nil })
	}

	cfg := &container.Config{
		Image: imageName,
		Env:   d.env,
	}
	cfg.Cmd = getBuildkitFlags(d.InitConfig)

	useInit := true // let it cleanup exited processes created by BuildKit's container API
	return l.Wrap("creating container "+d.Name, func() error {
		hc := &container.HostConfig{
			Privileged: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: d.Name + volumeStateSuffix,
					Target: confutil.DefaultBuildKitStateDir,
				},
			},
			Init: &useInit,
		}
		if d.netMode != "" {
			hc.NetworkMode = container.NetworkMode(d.netMode)
		}
		if d.memory != 0 {
			hc.Resources.Memory = int64(d.memory)
		}
		if d.memorySwap != 0 {
			hc.Resources.MemorySwap = int64(d.memorySwap)
		}
		if d.cpuQuota != 0 {
			hc.Resources.CPUQuota = d.cpuQuota
		}
		if d.cpuPeriod != 0 {
			hc.Resources.CPUPeriod = d.cpuPeriod
		}
		if d.cpuShares != 0 {
			hc.Resources.CPUShares = d.cpuShares
		}
		if d.cpusetCpus != "" {
			hc.Resources.CpusetCpus = d.cpusetCpus
		}
		if d.cpusetMems != "" {
			hc.Resources.CpusetMems = d.cpusetMems
		}
		if info, err := d.DockerAPI.Info(ctx); err == nil {
			if info.CgroupDriver == "cgroupfs" {
				// Place all buildkit containers inside this cgroup by default so limits can be attached
				// to all build activity on the host.
				hc.CgroupParent = "/docker/buildx"
				if d.cgroupParent != "" {
					hc.CgroupParent = d.cgroupParent
				}
			}

			secOpts, err := dockertypes.DecodeSecurityOptions(info.SecurityOptions)
			if err != nil {
				return err
			}
			for _, f := range secOpts {
				if f.Name == "userns" {
					hc.UsernsMode = "host"
					break
				}
			}

		}
		_, err := d.DockerAPI.ContainerCreate(ctx, cfg, hc, &network.NetworkingConfig{}, nil, d.Name)
		if err != nil && !errdefs.IsConflict(err) {
			return err
		}
		if err == nil {
			if err := d.copyToContainer(ctx, d.InitConfig.Files); err != nil {
				return err
			}
			if err := d.start(ctx, l); err != nil {
				return err
			}
		}
		return d.wait(ctx, l)
	})
}

func (d *Driver) wait(ctx context.Context, l progress.SubLogger) error {
	try := 1
	for {
		bufStdout := &bytes.Buffer{}
		bufStderr := &bytes.Buffer{}
		if err := d.run(ctx, []string{"buildctl", "debug", "workers"}, bufStdout, bufStderr); err != nil {
			if try > 15 {
				if err != nil {
					d.copyLogs(context.TODO(), l)
					if bufStdout.Len() != 0 {
						l.Log(1, bufStdout.Bytes())
					}
					if bufStderr.Len() != 0 {
						l.Log(2, bufStderr.Bytes())
					}
				}
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(try*120) * time.Millisecond):
				try++
				continue
			}
		}
		return nil
	}
}

func (d *Driver) copyLogs(ctx context.Context, l progress.SubLogger) error {
	rc, err := d.DockerAPI.ContainerLogs(ctx, d.Name, dockertypes.ContainerLogsOptions{
		ShowStdout: true, ShowStderr: true,
	})
	if err != nil {
		return err
	}
	stdout := &logWriter{logger: l, stream: 1}
	stderr := &logWriter{logger: l, stream: 2}
	if _, err := stdcopy.StdCopy(stdout, stderr, rc); err != nil {
		return err
	}
	return rc.Close()
}

func (d *Driver) copyToContainer(ctx context.Context, files map[string][]byte) error {
	srcPath, err := writeConfigFiles(files)
	if err != nil {
		return err
	}
	if srcPath != "" {
		defer os.RemoveAll(srcPath)
	}
	srcArchive, err := dockerarchive.TarWithOptions(srcPath, &dockerarchive.TarOptions{
		ChownOpts: &idtools.Identity{UID: 0, GID: 0},
	})
	if err != nil {
		return err
	}
	defer srcArchive.Close()

	baseDir := path.Dir(confutil.DefaultBuildKitConfigDir)
	return d.DockerAPI.CopyToContainer(ctx, d.Name, baseDir, srcArchive, dockertypes.CopyToContainerOptions{})
}

func (d *Driver) exec(ctx context.Context, cmd []string) (string, net.Conn, error) {
	execConfig := dockertypes.ExecConfig{
		Cmd:          cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}
	response, err := d.DockerAPI.ContainerExecCreate(ctx, d.Name, execConfig)
	if err != nil {
		return "", nil, err
	}

	execID := response.ID
	if execID == "" {
		return "", nil, errors.New("exec ID empty")
	}

	resp, err := d.DockerAPI.ContainerExecAttach(ctx, execID, dockertypes.ExecStartCheck{})
	if err != nil {
		return "", nil, err
	}
	return execID, resp.Conn, nil
}

func (d *Driver) run(ctx context.Context, cmd []string, stdout, stderr io.Writer) (err error) {
	id, conn, err := d.exec(ctx, cmd)
	if err != nil {
		return err
	}
	if _, err := stdcopy.StdCopy(stdout, stderr, conn); err != nil {
		return err
	}
	conn.Close()
	resp, err := d.DockerAPI.ContainerExecInspect(ctx, id)
	if err != nil {
		return err
	}
	if resp.ExitCode != 0 {
		return errors.Errorf("exit code %d", resp.ExitCode)
	}
	return nil
}

func (d *Driver) start(ctx context.Context, l progress.SubLogger) error {
	return d.DockerAPI.ContainerStart(ctx, d.Name, dockertypes.ContainerStartOptions{})
}

func (d *Driver) Info(ctx context.Context) (*driver.Info, error) {
	ctn, err := d.DockerAPI.ContainerInspect(ctx, d.Name)
	if err != nil {
		if dockerclient.IsErrNotFound(err) {
			return &driver.Info{
				Status: driver.Inactive,
			}, nil
		}
		return nil, err
	}

	if ctn.State.Running {
		return &driver.Info{
			Status: driver.Running,
		}, nil
	}

	return &driver.Info{
		Status: driver.Stopped,
	}, nil
}

func (d *Driver) Version(ctx context.Context) (string, error) {
	bufStdout := &bytes.Buffer{}
	bufStderr := &bytes.Buffer{}
	if err := d.run(ctx, []string{"buildkitd", "--version"}, bufStdout, bufStderr); err != nil {
		if bufStderr.Len() > 0 {
			return "", errors.Wrap(err, bufStderr.String())
		}
		return "", err
	}
	version := strings.Fields(bufStdout.String())
	if len(version) != 4 {
		return "", errors.Errorf("unexpected version format: %s", bufStdout.String())
	}
	return version[2], nil
}

func (d *Driver) Stop(ctx context.Context, force bool) error {
	info, err := d.Info(ctx)
	if err != nil {
		return err
	}
	if info.Status == driver.Running {
		return d.DockerAPI.ContainerStop(ctx, d.Name, container.StopOptions{})
	}
	return nil
}

func (d *Driver) Rm(ctx context.Context, force, rmVolume, rmDaemon bool) error {
	info, err := d.Info(ctx)
	if err != nil {
		return err
	}
	if info.Status != driver.Inactive {
		container, err := d.DockerAPI.ContainerInspect(ctx, d.Name)
		if err != nil {
			return err
		}
		if rmDaemon {
			if err := d.DockerAPI.ContainerRemove(ctx, d.Name, dockertypes.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         force,
			}); err != nil {
				return err
			}
			for _, v := range container.Mounts {
				if v.Name != d.Name+volumeStateSuffix {
					continue
				}
				if rmVolume {
					return d.DockerAPI.VolumeRemove(ctx, d.Name+volumeStateSuffix, false)
				}
			}
		}
	}
	return nil
}

func (d *Driver) Client(ctx context.Context) (*client.Client, error) {
	_, conn, err := d.exec(ctx, []string{"buildctl", "dial-stdio"})
	if err != nil {
		return nil, err
	}

	conn = demuxConn(conn)

	exp, err := detect.Exporter()
	if err != nil {
		return nil, err
	}

	var opts []client.ClientOpt
	var counter int64
	opts = append(opts, client.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		if atomic.AddInt64(&counter, 1) > 1 {
			return nil, net.ErrClosed
		}
		return conn, nil
	}))
	if td, ok := exp.(client.TracerDelegate); ok {
		opts = append(opts, client.WithTracerDelegate(td))
	}
	return client.New(ctx, "", opts...)
}

func (d *Driver) Factory() driver.Factory {
	return d.factory
}

func (d *Driver) Features(ctx context.Context) map[driver.Feature]bool {
	return map[driver.Feature]bool{
		driver.OCIExporter:    true,
		driver.DockerExporter: true,
		driver.CacheExport:    true,
		driver.MultiPlatform:  true,
	}
}

func (d *Driver) HostGatewayIP(ctx context.Context) (net.IP, error) {
	return nil, errors.New("host-gateway is not supported by the docker-container driver")
}

func demuxConn(c net.Conn) net.Conn {
	pr, pw := io.Pipe()
	// TODO: rewrite parser with Reader() to avoid goroutine switch
	go func() {
		_, err := stdcopy.StdCopy(pw, os.Stderr, c)
		pw.CloseWithError(err)
	}()
	return &demux{
		Conn:   c,
		Reader: pr,
	}
}

type demux struct {
	net.Conn
	io.Reader
}

func (d *demux) Read(dt []byte) (int, error) {
	return d.Reader.Read(dt)
}

type logWriter struct {
	logger progress.SubLogger
	stream int
}

func (l *logWriter) Write(dt []byte) (int, error) {
	l.logger.Log(l.stream, dt)
	return len(dt), nil
}

func writeConfigFiles(m map[string][]byte) (_ string, err error) {
	// Temp dir that will be copied to the container
	tmpDir, err := os.MkdirTemp("", "buildkitd-config")
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(tmpDir)
		}
	}()
	configDir := filepath.Base(confutil.DefaultBuildKitConfigDir)
	for f, dt := range m {
		p := filepath.Join(tmpDir, configDir, f)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			return "", err
		}
		if err := os.WriteFile(p, dt, 0644); err != nil {
			return "", err
		}
	}
	return tmpDir, nil
}

func getBuildkitFlags(initConfig driver.InitConfig) []string {
	flags := initConfig.BuildkitFlags
	if _, ok := initConfig.Files[buildkitdConfigFile]; ok {
		// There's no way for us to determine the appropriate default configuration
		// path and the default path can vary depending on if the image is normal
		// or rootless.
		//
		// In order to ensure that --config works, copy to a specific path and
		// specify the location.
		//
		// This should be appended before the user-specified arguments
		// so that this option could be overwritten by the user.
		newFlags := make([]string, 0, len(flags)+2)
		newFlags = append(newFlags, "--config", path.Join("/etc/buildkit", buildkitdConfigFile))
		flags = append(newFlags, flags...)
	}
	return flags
}
