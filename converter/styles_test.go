package converter

import (
	"testing"

	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func TestIsEmphasis(t *testing.T) {
	t.Run("is an emphasis node with level 1", func(t *testing.T) {
		node := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      1,
		}
		assert.True(t, isEmphasis(node))
	})

	t.Run("is an emphasis node with level 2", func(t *testing.T) {
		node := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      2,
		}
		assert.False(t, isEmphasis(node))
	})

	t.Run("not an emphasis node", func(t *testing.T) {
		node := &ast.Paragraph{}
		assert.False(t, isEmphasis(node))
	})
}

func TestIsStrong(t *testing.T) {
	t.Run("is an emphasis node with level 2", func(t *testing.T) {
		node := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      2,
		}
		assert.True(t, isStrong(node))
	})

	t.Run("is an emphasis node with level 1", func(t *testing.T) {
		node := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      1,
		}
		assert.False(t, isStrong(node))
	})

	t.Run("not an emphasis node", func(t *testing.T) {
		node := &ast.Paragraph{}
		assert.False(t, isStrong(node))
	})
}

func TestIsCodeSpan(t *testing.T) {
	t.Run("is a code span node", func(t *testing.T) {
		node := &ast.CodeSpan{}
		assert.True(t, isCodeSpan(node))
	})

	t.Run("not a code span node", func(t *testing.T) {
		node := &ast.Paragraph{}
		assert.False(t, isCodeSpan(node))
	})
}

func TestConvertEmphasis(t *testing.T) {
	t.Run("converts emphasis node with content", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("*italic text*")

		// Create an emphasis node
		emphasisNode := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      1,
		}

		// Create a text node as a child of the emphasis
		textNode := ast.NewTextSegment(text.NewSegment(1, 12)) // "italic text"の部分

		// Set up the AST structure
		emphasisNode.AppendChild(emphasisNode, textNode)

		// Execute the function
		richText := convertEmphasis(emphasisNode, source)

		// Verify the result
		assert.NotNil(t, richText, "Expected rich text to not be nil")
		assert.Greater(t, len(richText), 0, "Expected rich text to not be empty")
		assert.Equal(t, "italic text", richText[0].PlainText, "Expected plain text to be 'italic text'")
		assert.NotNil(t, richText[0].Annotations, "Expected annotations to be not nil")
		assert.True(t, richText[0].Annotations.Italic, "Expected italic to be true")

		// Explicitly use notionapi package to avoid unused import error
		var _ notionapi.Annotations
	})

	t.Run("handles nil emphasis node", func(t *testing.T) {
		richText := convertEmphasis(nil, nil)
		assert.Nil(t, richText, "Expected rich text to be nil for nil emphasis node")
	})

	t.Run("handles emphasis node with wrong level", func(t *testing.T) {
		emphasisNode := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      2,
		}
		richText := convertEmphasis(emphasisNode, nil)
		assert.Nil(t, richText, "Expected rich text to be nil for emphasis node with wrong level")
	})

	t.Run("handles emphasis node with no content", func(t *testing.T) {
		emphasisNode := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      1,
		}
		richText := convertEmphasis(emphasisNode, nil)
		assert.Nil(t, richText, "Expected rich text to be nil for emphasis node with no content")
	})
}

func TestConvertStrong(t *testing.T) {
	t.Run("converts strong node with content", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("**bold text**")

		// Create a strong node (emphasis with level 2)
		strongNode := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      2,
		}

		// Create a text node as a child of the strong
		textNode := ast.NewTextSegment(text.NewSegment(2, 11)) // "bold text"の部分

		// Set up the AST structure
		strongNode.AppendChild(strongNode, textNode)

		// Execute the function
		richText := convertStrong(strongNode, source)

		// Verify the result
		assert.NotNil(t, richText, "Expected rich text to not be nil")
		assert.Greater(t, len(richText), 0, "Expected rich text to not be empty")
		assert.Equal(t, "bold text", richText[0].PlainText, "Expected plain text to be 'bold text'")
		assert.NotNil(t, richText[0].Annotations, "Expected annotations to be not nil")
		assert.True(t, richText[0].Annotations.Bold, "Expected bold to be true")
	})

	t.Run("handles nil strong node", func(t *testing.T) {
		richText := convertStrong(nil, nil)
		assert.Nil(t, richText, "Expected rich text to be nil for nil strong node")
	})

	t.Run("handles strong node with wrong level", func(t *testing.T) {
		strongNode := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      1,
		}
		richText := convertStrong(strongNode, nil)
		assert.Nil(t, richText, "Expected rich text to be nil for strong node with wrong level")
	})

	t.Run("handles strong node with no content", func(t *testing.T) {
		strongNode := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      2,
		}
		richText := convertStrong(strongNode, nil)
		assert.Nil(t, richText, "Expected rich text to be nil for strong node with no content")
	})
}

func TestConvertCodeSpan(t *testing.T) {
	t.Run("converts code span node with content", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("`code text`")

		// Create a code span node
		codeSpanNode := &ast.CodeSpan{
			BaseInline: ast.BaseInline{},
		}

		// Create a text node as a child of the code span
		textNode := ast.NewTextSegment(text.NewSegment(1, 10)) // "code text"の部分

		// Set up the AST structure
		codeSpanNode.AppendChild(codeSpanNode, textNode)

		// Execute the function
		richText := convertCodeSpan(codeSpanNode, source)

		// Verify the result
		assert.NotNil(t, richText, "Expected rich text to not be nil")
		assert.Greater(t, len(richText), 0, "Expected rich text to not be empty")
		assert.Equal(t, "code text", richText[0].PlainText, "Expected plain text to be 'code text'")
		assert.NotNil(t, richText[0].Annotations, "Expected annotations to be not nil")
		assert.True(t, richText[0].Annotations.Code, "Expected code to be true")
	})

	t.Run("handles nil code span node", func(t *testing.T) {
		richText := convertCodeSpan(nil, nil)
		assert.Nil(t, richText, "Expected rich text to be nil for nil code span node")
	})

	t.Run("handles code span node with no content", func(t *testing.T) {
		codeSpanNode := &ast.CodeSpan{
			BaseInline: ast.BaseInline{},
		}
		richText := convertCodeSpan(codeSpanNode, nil)
		assert.Nil(t, richText, "Expected rich text to be nil for code span node with no content")
	})
}

func TestConvertStyle(t *testing.T) {
	t.Run("converts emphasis node", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("*italic text*")

		// Create an emphasis node
		emphasisNode := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      1,
		}

		// Create a text node as a child of the emphasis
		textNode := ast.NewTextSegment(text.NewSegment(1, 12)) // "italic text"の部分

		// Set up the AST structure
		emphasisNode.AppendChild(emphasisNode, textNode)

		// Execute the function
		richText := convertStyle(emphasisNode, source)

		// Verify the result
		assert.NotNil(t, richText, "Expected rich text to not be nil")
		assert.Greater(t, len(richText), 0, "Expected rich text to not be empty")
		assert.Equal(t, "italic text", richText[0].PlainText, "Expected plain text to be 'italic text'")
		assert.NotNil(t, richText[0].Annotations, "Expected annotations to be not nil")
		assert.True(t, richText[0].Annotations.Italic, "Expected italic to be true")
	})

	t.Run("converts strong node", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("**bold text**")

		// Create a strong node (emphasis with level 2)
		strongNode := &ast.Emphasis{
			BaseInline: ast.BaseInline{},
			Level:      2,
		}

		// Create a text node as a child of the strong
		textNode := ast.NewTextSegment(text.NewSegment(2, 11)) // "bold text"の部分

		// Set up the AST structure
		strongNode.AppendChild(strongNode, textNode)

		// Execute the function
		richText := convertStyle(strongNode, source)

		// Verify the result
		assert.NotNil(t, richText, "Expected rich text to not be nil")
		assert.Greater(t, len(richText), 0, "Expected rich text to not be empty")
		assert.Equal(t, "bold text", richText[0].PlainText, "Expected plain text to be 'bold text'")
		assert.NotNil(t, richText[0].Annotations, "Expected annotations to be not nil")
		assert.True(t, richText[0].Annotations.Bold, "Expected bold to be true")
	})

	t.Run("converts code span node", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("`code text`")

		// Create a code span node
		codeSpanNode := &ast.CodeSpan{
			BaseInline: ast.BaseInline{},
		}

		// Create a text node as a child of the code span
		textNode := ast.NewTextSegment(text.NewSegment(1, 10)) // "code text"の部分

		// Set up the AST structure
		codeSpanNode.AppendChild(codeSpanNode, textNode)

		// Execute the function
		richText := convertStyle(codeSpanNode, source)

		// Verify the result
		assert.NotNil(t, richText, "Expected rich text to not be nil")
		assert.Greater(t, len(richText), 0, "Expected rich text to not be empty")
		assert.Equal(t, "code text", richText[0].PlainText, "Expected plain text to be 'code text'")
		assert.NotNil(t, richText[0].Annotations, "Expected annotations to be not nil")
		assert.True(t, richText[0].Annotations.Code, "Expected code to be true")
	})

	t.Run("handles nil node", func(t *testing.T) {
		richText := convertStyle(nil, nil)
		assert.Nil(t, richText, "Expected rich text to be nil for nil node")
	})

	t.Run("handles unsupported node type", func(t *testing.T) {
		node := &ast.Paragraph{}
		richText := convertStyle(node, nil)
		assert.Nil(t, richText, "Expected rich text to be nil for unsupported node type")
	})
}
