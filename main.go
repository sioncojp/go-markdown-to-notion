package main

import (
  "context"
  "fmt"
  "log"
  "os"
  "os/signal"

  "github.com/urfave/cli/v3"
  "golang.org/x/sys/unix"
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
  done := make(chan struct{})
  ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGTERM, unix.SIGINT)
  go func() {
    select {
    case <-ctx.Done():
      done <- struct{}{}
      close(done)
      cancel()
      log.Println("received cancel...")
      os.Exit(1)
    }
  }()

  NotionAPIToken = os.Getenv("NOTION_API_TOKEN")
  if NotionAPIToken == "" {
    return fmt.Errorf("Notion_API_TOKEN environment variable is not set")
  }

  // Create a new Notion client
  notion := NewNotionClient()

  cmd := &cli.Command{
    Commands: []*cli.Command{
      // subcommand: upload
      {
        Name:  "upload",
        Usage: "upload markdown to notion",
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
          fmt.Println("added task: ", cmd.Args().First())
          return nil
        },
      },
      // subcommand: delete-all-blocks
      {
        Name:  "delete-all-blocks",
        Usage: "delete all blocks in a Notion page",
        Flags: []cli.Flag{
          &cli.StringFlag{
            Name:     "notion-page-id",
            Usage:    "delete all blocks in this notion page id",
            Required: true,
          },
        },
        Action: func(ctx context.Context, cmd *cli.Command) error {
          NotionPageID = cmd.String("notion-page-id")
          if err := notion.DeleteAllBlocks(ctx, NotionPageID); err != nil {
            return fmt.Errorf("failed to delete all blocks: %w", err)
          }
          return nil
        },
      },
    },
  }

  // Run the CLI app
  if err := cmd.Run(context.Background(), os.Args); err != nil {
    cancel()
    log.Fatal(err)
  }

  return nil
}
