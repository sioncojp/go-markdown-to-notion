package converter

import (
  "github.com/jomei/notionapi"
  "github.com/sioncojp/go-markdown-to-notion/chunk"
  "github.com/yuin/goldmark/ast"
)

func isBlockquote(node ast.Node) bool {
  _, ok := node.(*ast.Blockquote)
  return ok
}

func convertBlockquote(node *ast.Blockquote, source []byte) *notionapi.QuoteBlock {
  if node == nil {
    return nil
  }

  var richTextBlocks []notionapi.RichText
  for child := node.FirstChild(); child != nil; child = child.NextSibling() {
    // Skip nested blockquotes as they will be handled separately
    if _, isBQ := child.(*ast.Blockquote); isBQ {
      continue
    }

    childBlocks := convertChildNodesToRichText(child, source)
    if childBlocks != nil {
      richTextBlocks = append(richTextBlocks, childBlocks...)
    }
  }

  if len(richTextBlocks) == 0 {
    // Extract text directly from the blockquote if no rich text blocks were created
    content := string(source)
    if content != "" {
      richTextBlocks = chunk.RichText(content, nil)
    } else {
      return nil
    }
  }

  return &notionapi.QuoteBlock{
    BasicBlock: notionapi.BasicBlock{
      Object: notionapi.ObjectTypeBlock,
      Type:   notionapi.BlockTypeQuote,
    },
    Quote: notionapi.Quote{
      RichText: richTextBlocks,
      Color:    "default",
    },
  }
}
