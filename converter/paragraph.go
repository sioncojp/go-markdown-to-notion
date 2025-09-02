package converter

import (
	"github.com/jomei/notionapi"
	"github.com/sioncojp/go-markdown-to-notion/chunk"
	"github.com/yuin/goldmark/ast"
)

// isParagraph checks if a node is a paragraph.
func isParagraph(node ast.Node) bool {
	_, ok := node.(*ast.Paragraph)
	return ok
}

// convertParagraph converts a paragraph node to a Notion paragraph block.
func convertParagraph(node *ast.Paragraph, source []byte) *notionapi.ParagraphBlock {
	if node == nil {
		return nil
	}

	// Extract text content from the paragraph node
	var richTextBlocks []notionapi.RichText

	// Process child nodes to handle inline elements like links, emphasis, etc.
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		// Handle different types of inline elements
		if isLink(child) {
			// Convert link node
			linkRichText := convertLink(child.(*ast.Link), source)
			if linkRichText != nil {
				richTextBlocks = append(richTextBlocks, linkRichText...)
			}
		} else if isEmphasis(child) || isStrong(child) || isCodeSpan(child) {
			// Convert style nodes (emphasis, strong, code span)
			styleRichText := convertStyle(child, source)
			if styleRichText != nil {
				richTextBlocks = append(richTextBlocks, styleRichText...)
			}
		} else if text, ok := child.(*ast.Text); ok {
			// Convert plain text
			content := string(text.Segment.Value(source))
			if content != "" {
				richTextBlocks = append(richTextBlocks, chunk.RichText(content, nil)...)
			}
		}
	}

	// If no rich text blocks were created, try to extract text directly from the paragraph
	if len(richTextBlocks) == 0 {
		content := string(node.Text(source))
		if content != "" {
			richTextBlocks = chunk.RichText(content, nil)
		} else {
			// Return nil for empty paragraphs
			return nil
		}
	}

	// Create and return the paragraph block
	return &notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeParagraph,
		},
		Paragraph: notionapi.Paragraph{
			RichText: richTextBlocks,
			Color:    "default",
		},
	}
}
