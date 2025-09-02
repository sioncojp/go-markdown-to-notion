package converter

import (
  "fmt"

  "github.com/jomei/notionapi"
  "github.com/sioncojp/go-markdown-to-notion/chunk"
  "github.com/yuin/goldmark/ast"
)

func isHeading(node ast.Node) bool {
  _, ok := node.(*ast.Heading)
  return ok
}

func convertHeading(node *ast.Heading, source []byte, h1Color, h2Color, h3Color string) notionapi.Block {
  // Handle nil node
  if node == nil {
    return nil
  }

  var block notionapi.Block

  // Extract heading text from the node
  headingText := string(node.Lines().Value(source))

  // Handle empty heading text
  if headingText == "" {
    return nil
  }

  if node.Level == 1 {
    block = notionapi.Heading1Block{
      BasicBlock: notionapi.BasicBlock{
        Object: notionapi.ObjectTypeBlock,
        Type:   notionapi.BlockTypeHeading1,
      },
      Heading1: notionapi.Heading{
        RichText:     chunk.RichText(headingText, nil),
        Children:     nil,
        Color:        fmt.Sprintf("%s_background", h1Color),
        IsToggleable: false,
      },
    }
  }

  if node.Level == 2 {
    block = notionapi.Heading2Block{
      BasicBlock: notionapi.BasicBlock{
        Object: notionapi.ObjectTypeBlock,
        Type:   notionapi.BlockTypeHeading2,
      },
      Heading2: notionapi.Heading{
        RichText:     chunk.RichText(headingText, nil),
        Children:     nil,
        Color:        fmt.Sprintf("%s_background", h2Color),
        IsToggleable: false,
      },
    }
  }

  if node.Level == 3 {
    block = notionapi.Heading3Block{
      BasicBlock: notionapi.BasicBlock{
        Object: notionapi.ObjectTypeBlock,
        Type:   notionapi.BlockTypeHeading3,
      },
      Heading3: notionapi.Heading{
        RichText:     chunk.RichText(headingText, nil),
        Children:     nil,
        Color:        fmt.Sprintf("%s_background", h3Color),
        IsToggleable: false,
      },
    }
  }

  return block
}
