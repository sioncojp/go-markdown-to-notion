package main

import (
  "context"
  "fmt"
  "log"
  "os"

  "github.com/urfave/cli/v3"
)

var (
  NotionAPIToken string
  NotionPageID   string
  SourceMdFile   string
)

func main() {
  if err := run(); err != nil {
    fmt.Fprintf(os.Stderr, "%v", err)
    os.Exit(1)
  }
}

func run() error {
  NotionAPIToken = os.Getenv("Notion_API_TOKEN")
  cmd := &cli.Command{
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:     "notion-page-id",
        Usage:    "output notion page id",
        Required: true,
      },
      &cli.StringFlag{
        Name:     "source-md-file",
        Usage:    "source markdown file",
        Required: true,
      },
    },
    Action: func(ctx context.Context, cmd *cli.Command) error {
      NotionPageID = cmd.String("notion-page-id")
      SourceMdFile = cmd.String("source-md-file")
      return nil
    },
  }

  if err := cmd.Run(context.Background(), os.Args); err != nil {
    log.Fatal(err)
  }

  return nil
}
