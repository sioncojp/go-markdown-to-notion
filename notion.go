package main

import (
  "context"

  "github.com/jomei/notionapi"
)

// Notion ... Store Notion client
type Notion struct {
  Client *notionapi.Client
}

// NewNotionClient ... Create a new Notion client
func NewNotionClient() *Notion {
  client := &Notion{
    Client: notionapi.NewClient(notionapi.Token(NotionAPIToken)),
  }
  return client
}

// DeleteAllBlocks ... Delete all blocks in a Notion page
func (n *Notion) DeleteAllBlocks(ctx context.Context, pageOrBlockID string) error {
  var hasError error
  hasMore := true
  startCursor := notionapi.Cursor("")

  for hasMore {
    blocks, err := n.Client.Block.GetChildren(ctx, notionapi.BlockID(pageOrBlockID), &notionapi.Pagination{
      StartCursor: startCursor,
      PageSize:    100,
    })
    if err != nil {
      hasError = err
      break
    }

    for _, b := range blocks.Results {
      if _, err := n.Client.Block.Delete(ctx, b.GetID()); err != nil {
        hasError = err
        break
      }
    }
    if blocks.HasMore {
      startCursor = notionapi.Cursor(blocks.NextCursor)
    } else {
      hasMore = false
    }
  }

  if hasError != nil {
    return hasError
  }
  return nil
}

// InsertBlocks ... Insert blocks into a Notion page
func (n *Notion) InsertBlocks(ctx context.Context, blockID string, blocks []notionapi.Block) error {
  // Notion API has a limit of 100 blocks per request
  const blockChildrenSize = 100

  for i := 0; i < len(blocks); i += blockChildrenSize {
    end := i + blockChildrenSize
    if end > len(blocks) {
      end = len(blocks)
    }

    if _, err := n.Client.Block.AppendChildren(ctx, notionapi.BlockID(blockID), &notionapi.AppendBlockChildrenRequest{
      Children: blocks[i:end],
    }); err != nil {
      return err
    }

  }

  return nil
}
