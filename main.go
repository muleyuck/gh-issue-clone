package main

import (
	"context"
	"fmt"
	"os"

	"github.com/muleyuck/gh-issue-clone/cmd"
	"github.com/muleyuck/gh-issue-clone/versions"
	"github.com/urfave/cli/v3"
)

func main() {
	command := &cli.Command{
		Name:    "issue-clone",
		Usage:   "clone GitHub issues from a given issue.",
		Version: versions.AppVersion,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "template", Aliases: []string{"t"}, Usage: "Issue Template Name"},
		},
		Arguments: []cli.Argument{
			&cli.StringArg{Name: "url"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			return cmd.CloneIssue(ctx, c)
		},
	}
	if err := command.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("✗ Error:%s\n", err)
		os.Exit(1)
	}
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
