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
  blocks, err := n.Client.Block.GetChildren(ctx, notionapi.BlockID(pageOrBlockID), &notionapi.Pagination{})
  if err != nil {
    return err
  }

  for _, b := range blocks.Results {
    if _, err := n.Client.Block.Delete(ctx, b.GetID()); err != nil {
      return err
    }
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
