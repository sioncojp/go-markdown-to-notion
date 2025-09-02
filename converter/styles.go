package converter

import (
  "github.com/jomei/notionapi"
  "github.com/sioncojp/go-markdown-to-notion/chunk"
  "github.com/yuin/goldmark/ast"
)

// isEmphasis checks if a node is an emphasis node with level 1 (italic).
func isEmphasis(node ast.Node) bool {
  if emphasis, ok := node.(*ast.Emphasis); ok {
    return emphasis.Level == 1
  }
  return false
}

// isStrong checks if a node is an emphasis node with level 2 (bold).
func isStrong(node ast.Node) bool {
  if emphasis, ok := node.(*ast.Emphasis); ok {
    return emphasis.Level == 2
  }
  return false
}

// isCodeSpan checks if a node is a code span.
func isCodeSpan(node ast.Node) bool {
  _, ok := node.(*ast.CodeSpan)
  return ok
}

// convertEmphasis converts an emphasis node (italic) to Notion rich text.
func convertEmphasis(node *ast.Emphasis, source []byte) []notionapi.RichText {
  if node == nil || node.Level != 1 {
    return nil
  }

  // Extract text content from the emphasis node
  var content string
  for child := node.FirstChild(); child != nil; child = child.NextSibling() {
    if text, ok := child.(*ast.Text); ok {
      content += string(text.Segment.Value(source))
    }
  }

  if content == "" {
    return nil
  }

  // Create annotations with italic set to true
  annotations := &notionapi.Annotations{
    Italic: true,
  }

  // Create rich text with italic annotation
  return chunk.RichText(content, annotations)
}

// convertStrong converts an emphasis node (bold) to Notion rich text.
func convertStrong(node *ast.Emphasis, source []byte) []notionapi.RichText {
  if node == nil || node.Level != 2 {
    return nil
  }

  // Extract text content from the strong node
  var content string
  for child := node.FirstChild(); child != nil; child = child.NextSibling() {
    if text, ok := child.(*ast.Text); ok {
      content += string(text.Segment.Value(source))
    }
  }

  if content == "" {
    return nil
  }

  // Create annotations with bold set to true
  annotations := &notionapi.Annotations{
    Bold: true,
  }

  // Create rich text with bold annotation
  return chunk.RichText(content, annotations)
}

// convertCodeSpan converts a code span node to Notion rich text.
func convertCodeSpan(node *ast.CodeSpan, source []byte) []notionapi.RichText {
  if node == nil {
    return nil
  }

  // Extract text content from the code span node
  var content string
  for child := node.FirstChild(); child != nil; child = child.NextSibling() {
    if text, ok := child.(*ast.Text); ok {
      content += string(text.Segment.Value(source))
    }
  }

  if content == "" {
    return nil
  }

  // Create annotations with code set to true
  annotations := &notionapi.Annotations{
    Code: true,
  }

  // Create rich text with code annotation
  return chunk.RichText(content, annotations)
}

// convertStyle converts a style node (emphasis, strong, code span) to Notion rich text.
func convertStyle(node ast.Node, source []byte) []notionapi.RichText {
  if node == nil {
    return nil
  }

  if isEmphasis(node) {
    return convertEmphasis(node.(*ast.Emphasis), source)
  }

  if isStrong(node) {
    return convertStrong(node.(*ast.Emphasis), source)
  }

  if isCodeSpan(node) {
    return convertCodeSpan(node.(*ast.CodeSpan), source)
  }

  return nil
}
