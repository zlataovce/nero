package main

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
)

// main is the application entrypoint.
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	appCtx := &appContext{
		logger: logger,
	}
	app := &cli.App{
		Name:  "nero",
		Usage: "CLI interface for the nero server",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "launches the server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "the configuration path, defaults to config.toml",
						Value:   "config.toml",
						EnvVars: []string{"NERO_CONFIG_PATH"},
					},
				},
				Action: appCtx.handleServer,
			},
			{
				Name:  "client",
				Usage: "client commands",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "url",
						Aliases:  []string{"u"},
						Usage:    "the nero v1 API server host",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "repo",
						Aliases:  []string{"r"},
						Usage:    "the target repo",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "key",
						Aliases: []string{"k"},
						Usage:   "the repo authentication key",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "upload",
						Usage: "upload commands",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "path",
								Aliases:  []string{"f"},
								Usage:    "the uploaded file path or remote url",
								Required: true,
							},
						},
						Subcommands: []*cli.Command{
							{
								Name:  "generic",
								Usage: "upload a file with generic metadata",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:  "source",
										Usage: "the source",
									},
									&cli.StringFlag{
										Name:  "artist",
										Usage: "the artist",
									},
									&cli.StringFlag{
										Name:  "artist-link",
										Usage: "a link to the artist",
									},
								},
								Action: appCtx.handleUploadGeneric,
							},
							{
								Name:  "anime",
								Usage: "upload a file with anime metadata",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:  "name",
										Usage: "the anime name",
									},
								},
								Action: appCtx.handleUploadAnime,
							},
						},
					},
					{
						Name:  "delete",
						Usage: "deletes media",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "id",
								Aliases:  []string{"i"},
								Usage:    "the media id to be deleted",
								Required: true,
							},
						},
						Action: appCtx.handleDelete,
					},
				},
			},
			{
				Name:  "config",
				Usage: "generates an example configuration file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "the configuration path, defaults to config.toml",
						Value:   "config.toml",
						EnvVars: []string{"NERO_CONFIG_PATH"},
					},
				},
				Action: appCtx.handleConfig,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal("failed to run cli", zap.Error(err))
	}
}
