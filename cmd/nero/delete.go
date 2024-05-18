package main

import (
	"fmt"
	"github.com/cephxdev/nero/internal/errors"
	"github.com/cephxdev/nero/server/api"
	v1 "github.com/cephxdev/nero/server/api/v1"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// handleDelete handles the delete sub-command.
func (ac *appContext) handleDelete(cCtx *cli.Context) error {
	c, err := v1.NewClientWithResponses(cCtx.String("url"))
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	uid, err := uuid.Parse(cCtx.String("id"))
	if err != nil {
		return errors.Wrap(err, "could not parse item id")
	}

	res, err := c.DeleteRepoIdWithResponse(
		cCtx.Context,
		cCtx.String("repo"),
		uid,
		&v1.DeleteRepoIdParams{XNeroKey: api.MakeOptString(cCtx.String("key"))},
	)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}

	code := res.StatusCode()
	if code > 399 {
		ac.logger.Error(
			"request completed with errors",
			zap.String("status", res.Status()),
			zap.Int("code", code),
			zap.ByteString("body", res.Body),
		)

		// error out to force an error exit code
		return fmt.Errorf("request completed with error status code %d", code)
	}

	ac.logger.Info("request completed", zap.ByteString("body", res.Body))
	return nil
}
