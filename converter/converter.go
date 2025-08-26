package converter

import (
  "fmt"
  "os"

  "github.com/jomei/notionapi"
  "github.com/yuin/goldmark"
  "github.com/yuin/goldmark/ast"
  "github.com/yuin/goldmark/text"
)

func Convert(markdownFilePath string) ([]notionapi.Block, error) {
  // Read the markdown file
  source, err := os.ReadFile(markdownFilePath)
  if err != nil {
    return nil, fmt.Errorf("failed to read markdown file: %w", err)
  }

  // Create a new goldmark instance
  md := goldmark.New()

  // Parse the markdown document into an AST node
  reader := text.NewReader(source)
  document := md.Parser().Parse(reader)

  // Create a slice to store the Notion blocks
  var blocks []notionapi.Block

  // Walk through the AST and convert each node to a Notion block
  ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
    if !entering {
      return ast.WalkContinue, nil
    }

    // Handle code blocks
    if isCodeBlock(node) {
      codeBlock := convertFencedCodeBlock(node.(*ast.FencedCodeBlock), source)
      if codeBlock != nil {
        blocks = append(blocks, codeBlock)
      }
      return ast.WalkSkipChildren, nil
    }

    // Handle other node types here
    // TODO: Implement conversion for other markdown elements

    return ast.WalkContinue, nil
  })

  return blocks, nil
}
