package confutil

import (
	"os"
	"path"
	"path/filepath"

	"github.com/docker/cli/cli/command"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ConfigDir will look for correct configuration store path;
// if `$BUILDX_CONFIG` is set - use it, otherwise use parent directory
// of Docker config file (i.e. `${DOCKER_CONFIG}/buildx`)
func ConfigDir(dockerCli command.Cli) string {
	if buildxConfig := os.Getenv("BUILDX_CONFIG"); buildxConfig != "" {
		logrus.Debugf("using config store %q based in \"$BUILDX_CONFIG\" environment variable", buildxConfig)
		return buildxConfig
	}

	buildxConfig := filepath.Join(filepath.Dir(dockerCli.ConfigFile().Filename), "buildx")
	logrus.Debugf("using default config store %q", buildxConfig)
	return buildxConfig
}

// DefaultConfigFile returns the default BuildKit configuration file path
func DefaultConfigFile(dockerCli command.Cli) (string, bool) {
	f := path.Join(ConfigDir(dockerCli), "buildkitd.default.toml")
	if _, err := os.Stat(f); err == nil {
		return f, true
	}
	return "", false
}

// loadConfigTree loads BuildKit config toml tree
func loadConfigTree(fp string) (*toml.Tree, error) {
	f, err := os.Open(fp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to load config from %s", fp)
	}
	defer f.Close()
	t, err := toml.LoadReader(f)
	if err != nil {
		return t, errors.Wrap(err, "failed to parse config")
	}
	return t, nil
}
