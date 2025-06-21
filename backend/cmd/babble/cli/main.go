package main

import (
	"babble/backend/internal/auth"
	"babble/backend/internal/config"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	logLevel slog.Level
)

// Init sets up the logging function and customizes verbosity as needed
func init() {
	logLevel = slog.LevelDebug

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true, // Adds source file and line number
	})

	logger := slog.New(textHandler)
	slog.SetDefault(logger)
}

func main() {
	cfg := config.NewConfig()

	app := &cli.App{
		Name: "babble",
		Before: func(ctx *cli.Context) error {
			ctx.App.Metadata["config"] = cfg
			return nil
		},
		Usage: "Handle authentication and management of babble user api",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Options for updating the API",
				Action: func(ctx *cli.Context) error {
					return cli.ShowSubcommandHelp(ctx)
				},
				Subcommands: []*cli.Command{
					{
						Name:  "user",
						Usage: "add a new user",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "username",
							},
							&cli.StringFlag{
								Name: "role",
							},
						},
						Action: func(ctx *cli.Context) error {
							username := ctx.String("username")
							if username == "" {
								fmt.Println("username must be specified, exiting")
								return fmt.Errorf("error in using create user")
							}
							role, err := auth.STRole(ctx.String("role"))
							if err != nil {
								slog.Error("There is no data for this item!")
							}

							apikey, err := auth.CreateUser(cfg.DB, cfg.Cfg.GetString("API_PRIVATE_KEY"), username, role)
							if err != nil {
								return fmt.Errorf("there was an error in creating an account")
							} else {
								fmt.Printf("key: %s", apikey)
							}
							return nil
						},
					},
					{
						Name:  "project",
						Usage: "add a new project name",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "project",
							},
						},
						Action: func(ctx *cli.Context) error {
							project := ctx.String("project")
							if project == "" {
								fmt.Println("project must be specified, exiting")
								os.Exit(1)
							}

							auth.CreateProject(cfg.DB, project)
							return nil
						},
					},
				},
			},
			{
				Name:  "delete",
				Usage: "Options for removal from the API",
				Action: func(ctx *cli.Context) error {
					return cli.ShowSubcommandHelp(ctx)
				},
				Subcommands: []*cli.Command{
					{
						Name:  "user",
						Usage: "delete an existing user",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "username",
							},
						},
						Action: func(ctx *cli.Context) error {
							username := ctx.String("username")
							if username == "" {
								fmt.Println("a username must be specified, exiting")
								os.Exit(1)
							}

							auth.DeleteUser(cfg.DB, username)
							return nil
						},
					},
					{
						Name:  "project",
						Usage: "delete an existing project name",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "project",
							},
						},
						Action: func(ctx *cli.Context) error {
							project := ctx.String("project")
							if project == "" {
								fmt.Println("a project must be specified, exiting")
								os.Exit(1)
							}

							return auth.DeleteProject(cfg.DB, project)
						},
					},
				},
			},
			{
				Name:  "access",
				Usage: "Options for adding or removing access",
				Subcommands: []*cli.Command{
					{
						Name:  "grant",
						Usage: "add project privilege to a user",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "username",
							},
							&cli.StringFlag{
								Name: "project",
							},
						},
						Action: func(ctx *cli.Context) error {
							username := ctx.String("username")
							project := ctx.String("project")

							if username == "" || project == "" {
								fmt.Println("a username and project must be specified in the grant, exiting")
								os.Exit(1)
							}

							return auth.GrantProjectAccess(cfg.DB, username, project)
						},
					},
					{
						Name:  "revoke",
						Usage: "removing a project privilege to a user",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "username",
							},
							&cli.StringFlag{
								Name: "project",
							},
						},
						Action: func(ctx *cli.Context) error {
							username := ctx.String("username")
							project := ctx.String("project")

							if username == "" || project == "" {
								fmt.Println("a username and project must be specified in the grant, exiting")
								os.Exit(1)
							}

							return auth.RevokeProjectAccess(cfg.DB, username, project)
						},
					},
					{
						Name:  "retrieve",
						Usage: "retrieving an api key for a user",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "username",
							},
						},
						Action: func(ctx *cli.Context) error {
							username := ctx.String("username")

							apikey, err := auth.RetrieveAPIKey(cfg.DB, cfg.Cfg.GetString("API_PRIVATE_KEY"), username)
							fmt.Printf("Key: %s", apikey)
							return err
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
