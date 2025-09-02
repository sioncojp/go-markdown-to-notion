package converter

import (
	"testing"

	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
)

func TestIsList(t *testing.T) {
	t.Run("is a list node", func(t *testing.T) {
		node := &ast.List{}
		assert.True(t, isList(node))
	})

	t.Run("not a list node", func(t *testing.T) {
		node := &ast.Paragraph{}
		assert.False(t, isList(node))
	})
}

func TestIsListItem(t *testing.T) {
	t.Run("is a list item node", func(t *testing.T) {
		node := &ast.ListItem{}
		assert.True(t, isListItem(node))
	})

	t.Run("not a list item node", func(t *testing.T) {
		node := &ast.Paragraph{}
		assert.False(t, isListItem(node))
	})
}

func TestConvertList(t *testing.T) {
	t.Run("can convert unordered list", func(t *testing.T) {
		// Create a list with one item
		list := &ast.List{
			BaseBlock: ast.BaseBlock{},
			Marker:    '-',
			IsTight:   true,
		}
		
		// Create a list item
		listItem := &ast.ListItem{
			BaseBlock: ast.BaseBlock{},
		}
		
		// Create a paragraph as the content of the list item
		paragraph := &ast.Paragraph{
			BaseBlock: ast.BaseBlock{},
		}
		
		// Set up the AST structure
		list.AppendChild(list, listItem)
		listItem.AppendChild(listItem, paragraph)
		
		// Source text
		source := []byte("Item 1")
		
		// Convert and test
		result := convertList(list, source)
		assert.Equal(t, 1, len(result))
		
		// Check the type and content
		bulletedItem, ok := result[0].(notionapi.BulletedListItemBlock)
		assert.True(t, ok)
		assert.Equal(t, "Item 1", bulletedItem.GetRichTextString())
	})
	
	t.Run("can convert ordered list", func(t *testing.T) {
		// Create a list with one item
		list := &ast.List{
			BaseBlock: ast.BaseBlock{},
			Marker:    '.',
			IsTight:   true,
			Start:     1,
		}
		
		// Create a list item
		listItem := &ast.ListItem{
			BaseBlock: ast.BaseBlock{},
		}
		
		// Create a paragraph as the content of the list item
		paragraph := &ast.Paragraph{
			BaseBlock: ast.BaseBlock{},
		}
		
		// Set up the AST structure
		list.AppendChild(list, listItem)
		listItem.AppendChild(listItem, paragraph)
		
		// Source text
		source := []byte("Item 1")
		
		// Convert and test
		result := convertList(list, source)
		assert.Equal(t, 1, len(result))
		
		// Check the type and content
		numberedItem, ok := result[0].(notionapi.NumberedListItemBlock)
		assert.True(t, ok)
		assert.Equal(t, "Item 1", numberedItem.GetRichTextString())
	})
	
	t.Run("returns nil for nil node", func(t *testing.T) {
		result := convertList(nil, []byte{})
		assert.Nil(t, result)
	})
	
	t.Run("skips empty list items", func(t *testing.T) {
		// Create a list with one item
		list := &ast.List{
			BaseBlock: ast.BaseBlock{},
			Marker:    '-',
			IsTight:   true,
		}
		
		// Create a list item
		listItem := &ast.ListItem{
			BaseBlock: ast.BaseBlock{},
		}
		
		// Create an empty paragraph
		paragraph := &ast.Paragraph{
			BaseBlock: ast.BaseBlock{},
		}
		
		// Set up the AST structure
		list.AppendChild(list, listItem)
		listItem.AppendChild(listItem, paragraph)
		
		// Empty source text
		source := []byte("")
		
		// Convert and test
		result := convertList(list, source)
		assert.Equal(t, 0, len(result))
	})
}
