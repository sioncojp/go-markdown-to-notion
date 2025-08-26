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
func (n *Notion) DeleteAllBlocks(ctx context.Context, pageID string) error {
  blocks, err := n.Client.Block.GetChildren(ctx, notionapi.BlockID(pageID), &notionapi.Pagination{})
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
