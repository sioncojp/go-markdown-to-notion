package converter

import (
  "testing"

  "github.com/jomei/notionapi"
  "github.com/stretchr/testify/assert"
  "github.com/yuin/goldmark/ast"
  "github.com/yuin/goldmark/text"
)

func TestIsBlockquote(t *testing.T) {
  t.Run("is a blockquote node", func(t *testing.T) {
    node := &ast.Blockquote{}
    assert.True(t, isBlockquote(node))
  })

  t.Run("not a blockquote node", func(t *testing.T) {
    node := &ast.Paragraph{}
    assert.False(t, isBlockquote(node))
  })
}

func TestConvertBlockquote(t *testing.T) {
  t.Run("converts blockquote with content", func(t *testing.T) {
    // Create a sample test source code
    source := []byte("> This is a blockquote")

    // Create a blockquote node
    blockquoteNode := &ast.Blockquote{
      BaseBlock: ast.BaseBlock{},
    }

    // Create a paragraph as a child of the blockquote
    paragraph := &ast.Paragraph{
      BaseBlock: ast.BaseBlock{},
    }

    // Set up the lines for the paragraph
    lines := text.NewSegments()
    lines.Append(text.NewSegment(2, 22)) // "This is a blockquote"の部分
    paragraph.SetLines(lines)

    // Set up the AST structure
    blockquoteNode.AppendChild(blockquoteNode, paragraph)

    // Execute the function
    quoteBlock := convertBlockquote(blockquoteNode, source)

    // Verify the result
    assert.NotNil(t, quoteBlock, "Expected quote block to not be nil")
    assert.Equal(t, notionapi.BlockTypeQuote, quoteBlock.Type, "Expected block type to be quote")
    assert.Equal(t, "default", quoteBlock.Quote.Color, "Expected color to be default")
    assert.Greater(t, len(quoteBlock.Quote.RichText), 0, "Expected rich text to not be empty")
  })

  t.Run("handles nil blockquote", func(t *testing.T) {
    quoteBlock := convertBlockquote(nil, nil)
    assert.Nil(t, quoteBlock, "Expected quote block to be nil for nil blockquote")
  })

  t.Run("handles empty blockquote", func(t *testing.T) {
    // Create an empty blockquote node
    blockquoteNode := &ast.Blockquote{
      BaseBlock: ast.BaseBlock{},
    }

    // Execute the function
    quoteBlock := convertBlockquote(blockquoteNode, []byte{})

    // Since there's no content, the function should return nil
    assert.Nil(t, quoteBlock, "Expected quote block to be nil for empty blockquote")
  })

  t.Run("handles nested blockquotes", func(t *testing.T) {
    // Create a sample test source code
    source := []byte("> Outer blockquote\n>> Nested blockquote")

    // Create an outer blockquote node
    outerBlockquote := &ast.Blockquote{
      BaseBlock: ast.BaseBlock{},
    }

    // Create a paragraph for the outer blockquote
    outerParagraph := &ast.Paragraph{
      BaseBlock: ast.BaseBlock{},
    }

    // Set up the lines for the outer paragraph
    outerLines := text.NewSegments()
    outerLines.Append(text.NewSegment(2, 18)) // "Outer blockquote"の部分
    outerParagraph.SetLines(outerLines)

    // Create a nested blockquote
    nestedBlockquote := &ast.Blockquote{
      BaseBlock: ast.BaseBlock{},
    }

    // Create a paragraph for the nested blockquote
    nestedParagraph := &ast.Paragraph{
      BaseBlock: ast.BaseBlock{},
    }

    // Set up the lines for the nested paragraph
    nestedLines := text.NewSegments()
    nestedLines.Append(text.NewSegment(21, 38)) // "Nested blockquote"の部分
    nestedParagraph.SetLines(nestedLines)

    // Set up the AST structure
    outerBlockquote.AppendChild(outerBlockquote, outerParagraph)
    nestedBlockquote.AppendChild(nestedBlockquote, nestedParagraph)
    outerBlockquote.AppendChild(outerBlockquote, nestedBlockquote)

    // Execute the function
    quoteBlock := convertBlockquote(outerBlockquote, source)

    // Verify the result
    assert.NotNil(t, quoteBlock, "Expected quote block to not be nil")
    assert.Equal(t, notionapi.BlockTypeQuote, quoteBlock.Type, "Expected block type to be quote")

  		// Check if the content is correctly extracted
    richTextContent := ""
    for _, rt := range quoteBlock.Quote.RichText {
      richTextContent += rt.PlainText
    }
  		assert.Contains(t, richTextContent, "Outer blockquote", "Expected rich text to contain outer blockquote text")
  		// Note: In the current implementation, nested blockquotes are included in the rich text
  		assert.Contains(t, richTextContent, "Nested blockquote", "Expected rich text to contain nested blockquote text")
  })

  t.Run("handles blockquote with direct text content", func(t *testing.T) {
    // Create a sample test source code with direct text content
    source := []byte("> Direct text content")

    // Create a blockquote node
    blockquoteNode := &ast.Blockquote{
      BaseBlock: ast.BaseBlock{},
    }

    // Set up the lines for the blockquote
    lines := text.NewSegments()
    lines.Append(text.NewSegment(2, 21)) // "Direct text content"の部分
    blockquoteNode.SetLines(lines)

    // Execute the function
    quoteBlock := convertBlockquote(blockquoteNode, source)

    // Verify the result
    assert.NotNil(t, quoteBlock, "Expected quote block to not be nil")
    assert.Equal(t, notionapi.BlockTypeQuote, quoteBlock.Type, "Expected block type to be quote")
    assert.Greater(t, len(quoteBlock.Quote.RichText), 0, "Expected rich text to not be empty")

    // Check if the content is correctly extracted
    richTextContent := ""
    for _, rt := range quoteBlock.Quote.RichText {
      richTextContent += rt.PlainText
    }
    assert.Contains(t, richTextContent, "Direct text content", "Expected rich text to contain the direct text content")
  })
}
