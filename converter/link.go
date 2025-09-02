package converter

import (
	"github.com/jomei/notionapi"
	"github.com/sioncojp/go-markdown-to-notion/chunk"
	"github.com/yuin/goldmark/ast"
)

// isLink checks if a node is a link.
func isLink(node ast.Node) bool {
	_, ok := node.(*ast.Link)
	return ok
}

// convertLink converts a link node to Notion rich text with link.
func convertLink(node *ast.Link, source []byte) []notionapi.RichText {
	if node == nil {
		return nil
	}

	// Get the destination (URL) of the link
	destination := string(node.Destination)
	if destination == "" {
		return nil
	}

	// Get the text content of the link
	var content string
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if text, ok := child.(*ast.Text); ok {
			content += string(text.Segment.Value(source))
		}
	}

	// If no content is found, use the URL as content
	if content == "" {
		content = destination
	}

	// Create rich text with link
	return chunk.RichTextWithLink(content, destination)
}
