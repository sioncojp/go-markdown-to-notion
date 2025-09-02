package converter

import (
	"testing"

	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func TestIsParagraph(t *testing.T) {
	t.Run("is a paragraph node", func(t *testing.T) {
		node := &ast.Paragraph{}
		assert.True(t, isParagraph(node))
	})

	t.Run("not a paragraph node", func(t *testing.T) {
		node := &ast.Heading{}
		assert.False(t, isParagraph(node))
	})
}

func TestConvertParagraph(t *testing.T) {
	t.Run("converts paragraph with plain text", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("This is a simple paragraph.")

		// Create a paragraph node
		paragraphNode := &ast.Paragraph{
			BaseBlock: ast.BaseBlock{},
		}

		// Create a text node as a child of the paragraph
		textNode := ast.NewTextSegment(text.NewSegment(0, 27)) // "This is a simple paragraph."の部分

		// Set up the AST structure
		paragraphNode.AppendChild(paragraphNode, textNode)

		// Execute the function
		paragraphBlock := convertParagraph(paragraphNode, source)

		// Verify the result
		assert.NotNil(t, paragraphBlock, "Expected paragraph block to not be nil")
		assert.Equal(t, notionapi.BlockTypeParagraph, paragraphBlock.Type, "Expected block type to be paragraph")
		assert.Equal(t, "default", paragraphBlock.Paragraph.Color, "Expected color to be default")
		assert.Greater(t, len(paragraphBlock.Paragraph.RichText), 0, "Expected rich text to not be empty")
		assert.Equal(t, "This is a simple paragraph.", paragraphBlock.Paragraph.RichText[0].PlainText, "Expected plain text to match")
	})

	t.Run("handles nil paragraph", func(t *testing.T) {
		paragraphBlock := convertParagraph(nil, nil)
		assert.Nil(t, paragraphBlock, "Expected paragraph block to be nil for nil paragraph")
	})

	t.Run("handles empty paragraph", func(t *testing.T) {
		// Create an empty paragraph node
		paragraphNode := &ast.Paragraph{
			BaseBlock: ast.BaseBlock{},
		}

		// Execute the function
		paragraphBlock := convertParagraph(paragraphNode, []byte{})

		// Since there's no content, the function should return nil
		assert.Nil(t, paragraphBlock, "Expected paragraph block to be nil for empty paragraph")
	})

	t.Run("handles paragraph with link", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("This is a paragraph with a [link](https://example.com).")

		// Create a paragraph node
		paragraphNode := &ast.Paragraph{
			BaseBlock: ast.BaseBlock{},
		}

		// Create a text node for the first part
		textNode1 := ast.NewTextSegment(text.NewSegment(0, 27)) // "This is a paragraph with a "の部分

		// Create a link node
		linkNode := ast.NewLink()
		linkNode.Destination = []byte("https://example.com")

		// Create a text node as a child of the link
		linkTextNode := ast.NewTextSegment(text.NewSegment(28, 32)) // "link"の部分
		linkNode.AppendChild(linkNode, linkTextNode)

		// Create a text node for the last part
		textNode2 := ast.NewTextSegment(text.NewSegment(54, 55)) // "."の部分

		// Set up the AST structure
		paragraphNode.AppendChild(paragraphNode, textNode1)
		paragraphNode.AppendChild(paragraphNode, linkNode)
		paragraphNode.AppendChild(paragraphNode, textNode2)

		// Execute the function
		paragraphBlock := convertParagraph(paragraphNode, source)

		// Verify the result
		assert.NotNil(t, paragraphBlock, "Expected paragraph block to not be nil")
		assert.Equal(t, notionapi.BlockTypeParagraph, paragraphBlock.Type, "Expected block type to be paragraph")
		assert.Greater(t, len(paragraphBlock.Paragraph.RichText), 0, "Expected rich text to not be empty")

		// Check if the content contains both the text and the link
		var plainText string
		var hasLink bool
		for _, rt := range paragraphBlock.Paragraph.RichText {
			plainText += rt.PlainText
			if rt.Text != nil && rt.Text.Link != nil {
				hasLink = true
				assert.Equal(t, "https://example.com", rt.Text.Link.Url, "Expected link URL to be correct")
			}
		}
		assert.Contains(t, plainText, "This is a paragraph with a ", "Expected plain text to contain the first part")
		assert.Contains(t, plainText, "link", "Expected plain text to contain the link text")
		assert.Contains(t, plainText, ".", "Expected plain text to contain the last part")
		assert.True(t, hasLink, "Expected paragraph to contain a link")
	})

	t.Run("handles paragraph with styled text", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("This is *italic* and **bold** text.")

		// Create a paragraph node
		paragraphNode := &ast.Paragraph{
			BaseBlock: ast.BaseBlock{},
		}

		// Create a text node for the first part
		textNode1 := ast.NewTextSegment(text.NewSegment(0, 8)) // "This is "の部分

		// Create an emphasis node (italic)
		emphasisNode := ast.NewEmphasis(1)
		emphasisTextNode := ast.NewTextSegment(text.NewSegment(9, 15)) // "italic"の部分
		emphasisNode.AppendChild(emphasisNode, emphasisTextNode)

		// Create a text node for the middle part
		textNode2 := ast.NewTextSegment(text.NewSegment(16, 21)) // " and "の部分

		// Create a strong node (bold)
		strongNode := ast.NewEmphasis(2)
		strongTextNode := ast.NewTextSegment(text.NewSegment(23, 27)) // "bold"の部分
		strongNode.AppendChild(strongNode, strongTextNode)

		// Create a text node for the last part
		textNode3 := ast.NewTextSegment(text.NewSegment(29, 34)) // " text."の部分

		// Set up the AST structure
		paragraphNode.AppendChild(paragraphNode, textNode1)
		paragraphNode.AppendChild(paragraphNode, emphasisNode)
		paragraphNode.AppendChild(paragraphNode, textNode2)
		paragraphNode.AppendChild(paragraphNode, strongNode)
		paragraphNode.AppendChild(paragraphNode, textNode3)

		// Execute the function
		paragraphBlock := convertParagraph(paragraphNode, source)

		// Verify the result
		assert.NotNil(t, paragraphBlock, "Expected paragraph block to not be nil")
		assert.Equal(t, notionapi.BlockTypeParagraph, paragraphBlock.Type, "Expected block type to be paragraph")
		assert.Greater(t, len(paragraphBlock.Paragraph.RichText), 0, "Expected rich text to not be empty")

		// Check if the content contains all parts and has proper styling
		var plainText string
		var hasItalic, hasBold bool
		for _, rt := range paragraphBlock.Paragraph.RichText {
			plainText += rt.PlainText
			if rt.Annotations != nil {
				if rt.Annotations.Italic {
					hasItalic = true
					assert.Equal(t, "italic", rt.PlainText, "Expected italic text to be correct")
				}
				if rt.Annotations.Bold {
					hasBold = true
					assert.Equal(t, "bold", rt.PlainText, "Expected bold text to be correct")
				}
			}
		}
		assert.Contains(t, plainText, "This is ", "Expected plain text to contain the first part")
		assert.Contains(t, plainText, "italic", "Expected plain text to contain the italic text")
		assert.Contains(t, plainText, " and ", "Expected plain text to contain the middle part")
		assert.Contains(t, plainText, "bold", "Expected plain text to contain the bold text")
		assert.Contains(t, plainText, "text", "Expected plain text to contain the last part")
		assert.True(t, hasItalic, "Expected paragraph to contain italic text")
		assert.True(t, hasBold, "Expected paragraph to contain bold text")
	})

	t.Run("handles paragraph with code span", func(t *testing.T) {
		// Create a sample test source code
		source := []byte("This is `code` text.")

		// Create a paragraph node
		paragraphNode := &ast.Paragraph{
			BaseBlock: ast.BaseBlock{},
		}

		// Create a text node for the first part
		textNode1 := ast.NewTextSegment(text.NewSegment(0, 8)) // "This is "の部分

		// Create a code span node
		codeSpanNode := ast.NewCodeSpan()
		codeSpanTextNode := ast.NewTextSegment(text.NewSegment(9, 13)) // "code"の部分
		codeSpanNode.AppendChild(codeSpanNode, codeSpanTextNode)

		// Create a text node for the last part
		textNode2 := ast.NewTextSegment(text.NewSegment(14, 20)) // " text."の部分

		// Set up the AST structure
		paragraphNode.AppendChild(paragraphNode, textNode1)
		paragraphNode.AppendChild(paragraphNode, codeSpanNode)
		paragraphNode.AppendChild(paragraphNode, textNode2)

		// Execute the function
		paragraphBlock := convertParagraph(paragraphNode, source)

		// Verify the result
		assert.NotNil(t, paragraphBlock, "Expected paragraph block to not be nil")
		assert.Equal(t, notionapi.BlockTypeParagraph, paragraphBlock.Type, "Expected block type to be paragraph")
		assert.Greater(t, len(paragraphBlock.Paragraph.RichText), 0, "Expected rich text to not be empty")

		// Check if the content contains all parts and has proper styling
		var plainText string
		var hasCode bool
		for _, rt := range paragraphBlock.Paragraph.RichText {
			plainText += rt.PlainText
			if rt.Annotations != nil && rt.Annotations.Code {
				hasCode = true
				assert.Equal(t, "code", rt.PlainText, "Expected code text to be correct")
			}
		}
		assert.Contains(t, plainText, "This is ", "Expected plain text to contain the first part")
		assert.Contains(t, plainText, "code", "Expected plain text to contain the code text")
		assert.Contains(t, plainText, " text.", "Expected plain text to contain the last part")
		assert.True(t, hasCode, "Expected paragraph to contain code text")
	})
}
