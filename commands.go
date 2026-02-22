package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/t-chov/pcloud-uploader/pcloud"
	"github.com/urfave/cli/v2"
)

// contextKey is an unexported type for context keys to avoid collisions.
type contextKey string

// authContextKey is the context key for the pCloud auth token.
const authContextKey contextKey = "auth"

var printFolderFn = color.New(color.FgBlue).FprintlnFunc()

// run creates a default pCloud client and starts the CLI application.
func run() int {
	client := pcloud.NewClient("")
	return runWithClient(client, os.Args)
}

// runWithClient starts the CLI application with the given pCloud API client and args.
// This is the testable entry point that allows dependency injection.
func runWithClient(client pcloud.API, args []string) int {
	app := cli.NewApp()
	app.Name = APP_NAME
	app.Version = VERSION
	app.Before = func(c *cli.Context) error {
		return login(c, client)
	}
	app.Commands = buildCommands(client)
	if err := app.Run(args); err != nil {
		return 1
	}
	return 0
}

// buildCommands returns the CLI command definitions for ls and up.
func buildCommands(client pcloud.API) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "ls",
			Aliases: []string{"listfolder"},
			Usage:   "Receive data for a folder.",
			Action: func(c *cli.Context) error {
				auth, ok := c.Context.Value(authContextKey).(string)
				if !ok {
					return fmt.Errorf("auth not set in context")
				}
				opts := pcloud.ListFolderOptions{
					Recursive:   c.Bool("recursive"),
					ShowDeleted: c.Bool("showdeleted"),
					NoFiles:     c.Bool("nofiles"),
					NoShares:    c.Bool("noshares"),
				}
				result, err := client.ListFolder(auth, c.Args().First(), opts)
				if err != nil {
					return printError(err)
				}
				printListfolder(os.Stdout, "/", result.Metadata)
				return nil
			},
			ArgsUsage: "<PATH>",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "recursive",
					Aliases: []string{"r"},
					Usage:   "full directory tree will be returned",
				},
				&cli.BoolFlag{
					Name:  "showdeleted",
					Usage: "deleted files and folders that can be undeleted will be displayed",
				},
				&cli.BoolFlag{
					Name:  "nofiles",
					Usage: "only the folder (sub)structure will be returned",
				},
				&cli.BoolFlag{
					Name:  "noshares",
					Usage: "only user's own folders and files will be displayed",
				},
			},
		},
		{
			Name:    "up",
			Aliases: []string{"uploadfile"},
			Usage:   "Upload a file.",
			Action: func(c *cli.Context) error {
				auth, ok := c.Context.Value(authContextKey).(string)
				if !ok {
					return fmt.Errorf("auth not set in context")
				}
				src := c.Args().First()
				file, err := os.Open(src)
				if err != nil {
					return printError(fmt.Errorf("cannot open %s", src))
				}
				defer file.Close()

				dest := c.Args().Get(1)
				result, err := client.UploadFile(auth, dest, filepath.Base(src), file)
				if err != nil {
					return printError(err)
				}
				for _, metadata := range result.Metadata {
					fmt.Println(metadata.Hash.String())
				}
				return nil
			},
			ArgsUsage: "<SOURCE_FILE> <DEST_PATH>",
		},
	}
}

// login is a cli.BeforeFunc that authenticates against pCloud and stores
// the session token in the context.
func login(c *cli.Context, client pcloud.API) error {
	username, err := loadFromInput(os.Stdin, ENV_PCLOUD_USERNAME, "username")
	if err != nil {
		return err
	}
	password, err := loadFromInput(os.Stdin, ENV_PCLOUD_PASSWORD, "password")
	if err != nil {
		return err
	}

	auth, err := client.Authenticate(username, password)
	if err != nil {
		return err
	}
	if auth == "" {
		return fmt.Errorf("empty auth")
	}

	c.Context = context.WithValue(c.Context, authContextKey, auth)
	return nil
}

// loadFromInput reads a value from the environment variable envName.
// If the variable is not set, it prompts the user on r with the given title.
func loadFromInput(r io.Reader, envName string, title string) (string, error) {
	value := os.Getenv(envName)
	if value == "" {
		fmt.Printf("%s: ", title)
		if _, err := fmt.Fscan(r, &value); err != nil {
			return "", fmt.Errorf("scan %s: %v", title, err)
		}
	}
	value = strings.TrimSpace(value)
	return value, nil
}

// printListfolder recursively prints a folder tree to w.
// Folders are printed in blue; files are printed in the default color.
func printListfolder(w io.Writer, path string, ls pcloud.FolderMetadata) {
	if ls.Path == "/" {
		printFolderFn(w, "/")
	} else if ls.IsFolder {
		path += ls.Name + "/"
		printFolderFn(w, path)
	} else {
		path += ls.Name
		fmt.Fprintln(w, path)
	}
	for _, content := range ls.Contents {
		printListfolder(w, path, content)
	}
}

// printError prints err to stderr in red and returns it.
func printError(err error) error {
	red := color.New(color.FgRed).FprintfFunc()
	red(os.Stderr, "%v\n", err)
	return err
}
