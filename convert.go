package main

import (
  "github.com/jomei/notionapi"
  "github.com/sioncojp/go-markdown-to-notion/converter"
)

// Convert ... convert markdown to Notion blocks
func Convert(markdownFilePath string) ([]notionapi.Block, error) {
  return converter.Convert(markdownFilePath)
}
