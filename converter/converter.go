package converter

import (
  "fmt"
  "os"

  "github.com/jomei/notionapi"
  "github.com/sioncojp/go-markdown-to-notion/chunk"
  "github.com/yuin/goldmark"
  "github.com/yuin/goldmark/ast"
  "github.com/yuin/goldmark/extension"
  east "github.com/yuin/goldmark/extension/ast"
  "github.com/yuin/goldmark/text"
)

type Converter struct {
  MarkdownFilePath string
  H1Color          string
  H2Color          string
  H3Color          string
}

func Convert(c *Converter) ([]notionapi.Block, error) {
  // Read the markdown file
  source, err := os.ReadFile(c.MarkdownFilePath)
  if err != nil {
    return nil, fmt.Errorf("failed to read markdown file: %w", err)
  }

  // Create a new goldmark instance with table extension
  md := goldmark.New(
    goldmark.WithExtensions(extension.Table),
  )
  document := md.Parser().Parse(text.NewReader(source))

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

    if isHeading(node) {
      block := convertHeading(node.(*ast.Heading), source, c.H1Color, c.H2Color, c.H3Color)
      if block != nil {
        blocks = append(blocks, block)
      }
      return ast.WalkSkipChildren, nil
    }

    if isList(node) {
      listBlocks := convertList(node.(*ast.List), source)
      if len(listBlocks) > 0 {
        blocks = append(blocks, listBlocks...)
      }
      return ast.WalkSkipChildren, nil
    }

    if isBlockquote(node) {
      quoteBlock := convertBlockquote(node.(*ast.Blockquote), source)
      if quoteBlock != nil {
        blocks = append(blocks, quoteBlock)
      }
      return ast.WalkSkipChildren, nil
    }

    if isParagraph(node) {
      paragraphBlock := convertParagraph(node.(*ast.Paragraph), source)
      if paragraphBlock != nil {
        blocks = append(blocks, paragraphBlock)
      }
      return ast.WalkSkipChildren, nil
    }

    if isTable(node) {
      tableBlock := convertTable(node.(*east.Table), source)
      if tableBlock != nil {
        blocks = append(blocks, tableBlock)
      }
      return ast.WalkSkipChildren, nil
    }

    // Handle other node types here
    // Note: Image nodes are not implemented yet

    return ast.WalkContinue, nil
  })

  return blocks, nil
}

// convertChildNodesToRichText converts the child nodes of a given AST node to Notion rich text blocks.
func convertChildNodesToRichText(node ast.Node, source []byte) []notionapi.RichText {
  if node == nil {
    return nil
  }

  var blocks []notionapi.RichText
  for child := node.FirstChild(); child != nil; child = child.NextSibling() {
    // Handle different types of inline elements
    if isLink(child) {
      // Convert link node
      linkRichText := convertLink(child.(*ast.Link), source)
      if linkRichText != nil {
        blocks = append(blocks, linkRichText...)
      }
    } else if isEmphasis(child) || isStrong(child) || isCodeSpan(child) {
      // Convert style nodes (emphasis, strong, code span)
      styleRichText := convertStyle(child, source)
      if styleRichText != nil {
        blocks = append(blocks, styleRichText...)
      }
    } else if text, ok := child.(*ast.Text); ok {
      // Convert plain text
      content := string(text.Segment.Value(source))
      if content != "" {
        blocks = append(blocks, chunk.RichText(content, nil)...)
      }
    } else {
      // Recursively process other node types
      childBlocks := convertChildNodesToRichText(child, source)
      if childBlocks != nil {
        blocks = append(blocks, childBlocks...)
      }
    }
  }

  return blocks
}
