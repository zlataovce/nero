package main

import (
	"github.com/cephxdev/nero"
	"github.com/cephxdev/nero/internal/errors"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

// handleConfig handles the config sub-command.
func (ac *appContext) handleConfig(cCtx *cli.Context) error {
	path := filepath.Clean(cCtx.String("config"))

	if _, err := os.Stat(path); err == nil {
		return errors.New("path already exists")
	}
	if err := os.WriteFile(path, nero.ExampleConfig, 0); err != nil {
		return errors.Wrap(err, "failed to save example configuration")
	}

	ac.logger.Info("example configuration saved successfully", zap.String("path", path))
	return nil
}
