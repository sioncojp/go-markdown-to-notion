package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func TestIsLink(t *testing.T) {
	t.Run("is a link node", func(t *testing.T) {
		node := &ast.Link{}
		assert.True(t, isLink(node))
	})

	t.Run("not a link node", func(t *testing.T) {
		node := &ast.Paragraph{}
		assert.False(t, isLink(node))
	})
}

func TestConvertLink(t *testing.T) {
	t.Run("converts link with content", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("[Link text](https://example.com)")

		// Create a link node
		linkNode := ast.NewLink()
		linkNode.Destination = []byte("https://example.com")

		// Create a text node as a child of the link
		textNode := ast.NewTextSegment(text.NewSegment(1, 10)) // "Link text"の部分

		// Set up the AST structure
		linkNode.AppendChild(linkNode, textNode)

		// Execute the function
		richText := convertLink(linkNode, source)

		// Verify the result
		assert.NotNil(t, richText, "Expected rich text to not be nil")
		assert.Greater(t, len(richText), 0, "Expected rich text to not be empty")
		assert.Equal(t, "Link text", richText[0].PlainText, "Expected plain text to be 'Link text'")
		assert.NotNil(t, richText[0].Text.Link, "Expected link to be not nil")
		assert.Equal(t, "https://example.com", richText[0].Text.Link.Url, "Expected link URL to be 'https://example.com'")
	})

	t.Run("handles nil link", func(t *testing.T) {
		richText := convertLink(nil, nil)
		assert.Nil(t, richText, "Expected rich text to be nil for nil link")
	})

	t.Run("handles link with empty destination", func(t *testing.T) {
		// Create a link node with empty destination
		linkNode := ast.NewLink()
		linkNode.Destination = []byte("")

		// Create a text node as a child of the link
		textNode := ast.NewTextSegment(text.NewSegment(1, 10)) // "Link text"の部分

		// Set up the AST structure
		linkNode.AppendChild(linkNode, textNode)

		// Execute the function
		richText := convertLink(linkNode, []byte("[Link text]()"))

		// Since the destination is empty, the function should return nil
		assert.Nil(t, richText, "Expected rich text to be nil for link with empty destination")
	})

	t.Run("handles link with no text content", func(t *testing.T) {
		// Create a link node with no text content
		linkNode := ast.NewLink()
		linkNode.Destination = []byte("https://example.com")

		// Execute the function
		richText := convertLink(linkNode, []byte("[](https://example.com)"))

		// Verify the result
		assert.NotNil(t, richText, "Expected rich text to not be nil")
		assert.Greater(t, len(richText), 0, "Expected rich text to not be empty")
		assert.Equal(t, "https://example.com", richText[0].PlainText, "Expected plain text to be the URL")
		assert.NotNil(t, richText[0].Text.Link, "Expected link to be not nil")
		assert.Equal(t, "https://example.com", richText[0].Text.Link.Url, "Expected link URL to be 'https://example.com'")
	})

	t.Run("handles link with title", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("[Link text](https://example.com \"Link title\")")

		// Create a link node
		linkNode := ast.NewLink()
		linkNode.Destination = []byte("https://example.com")
		linkNode.Title = []byte("Link title")

		// Create a text node as a child of the link
		textNode := ast.NewTextSegment(text.NewSegment(1, 10)) // "Link text"の部分

		// Set up the AST structure
		linkNode.AppendChild(linkNode, textNode)

		// Execute the function
		richText := convertLink(linkNode, source)

		// Verify the result
		assert.NotNil(t, richText, "Expected rich text to not be nil")
		assert.Greater(t, len(richText), 0, "Expected rich text to not be empty")
		assert.Equal(t, "Link text", richText[0].PlainText, "Expected plain text to be 'Link text'")
		assert.NotNil(t, richText[0].Text.Link, "Expected link to be not nil")
		assert.Equal(t, "https://example.com", richText[0].Text.Link.Url, "Expected link URL to be 'https://example.com'")
	})
}
