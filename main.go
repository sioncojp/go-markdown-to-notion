package main

import (
  "context"
  "fmt"
  "log"
  "os"
  "os/signal"

  "github.com/sioncojp/go-markdown-to-notion/converter"
  "github.com/urfave/cli/v3"
  "golang.org/x/sys/unix"
)

var (
  NotionAPIToken   string
  NotionBlockID    string
  NotionPageID     string
  SourceMdFilePath string
  H1Color          string
  H2Color          string
  H3Color          string
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
            Name:     "notion-block-id",
            Usage:    "output below this notion block id",
            Required: true,
          },
          &cli.StringFlag{
            Name:     "source-md-filepath",
            Usage:    "source markdown file path",
            Required: true,
          },
          &cli.StringFlag{
            Name:  "h1-color",
            Usage: "h1 color",
            Value: "blue",
          },
          &cli.StringFlag{
            Name:  "h2-color",
            Usage: "h2 color",
            Value: "orange",
          },
          &cli.StringFlag{
            Name:  "h3-color",
            Usage: "h3 color",
            Value: "yellow",
          },
        },
        Action: func(ctx context.Context, cmd *cli.Command) error {
          NotionBlockID = cmd.String("notion-block-id")
          SourceMdFilePath = cmd.String("source-md-filepath")
          H1Color = cmd.String("h1-color")
          H2Color = cmd.String("h2-color")
          H3Color = cmd.String("h3-color")

          c := &converter.Converter{
            MarkdownFilePath: SourceMdFilePath,
            H1Color:          H1Color,
            H2Color:          H2Color,
            H3Color:          H3Color,
          }

          blocks, err := converter.Convert(c)
          if err != nil {
            return fmt.Errorf("failed to convert markdown to notion: %w", err)
          }
          //          for i, b := range blocks {
          //            fmt.Printf("Block %d: %+v\n", i, b)
          //          }

          if err := notion.InsertBlocks(ctx, NotionBlockID, blocks); err != nil {
            return fmt.Errorf("failed to insert blocks: %w", err)
          }

          return nil
        },
      },
      // subcommand: delete-all-blocks
      {
        Name:  "delete-all-blocks",
        Usage: "delete all blocks in a Notion page",
        Flags: []cli.Flag{
          &cli.StringFlag{
            Name:     "notion-page-or-block-id",
            Usage:    "delete all blocks in this notion page id",
            Required: true,
          },
        },
        Action: func(ctx context.Context, cmd *cli.Command) error {
          NotionBlockID = cmd.String("notion-page-or-block-id")
          if err := notion.DeleteAllBlocks(ctx, NotionBlockID); err != nil {
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
