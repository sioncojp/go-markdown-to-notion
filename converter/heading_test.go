package converter

import (
  "testing"

  "github.com/jomei/notionapi"
  "github.com/stretchr/testify/assert"
  "github.com/yuin/goldmark/ast"
  textm "github.com/yuin/goldmark/text"
)

func TestIsHeading(t *testing.T) {
  t.Run("is a heading node", func(t *testing.T) {
    node := &ast.Heading{}
    assert.True(t, isHeading(node))
  })

  t.Run("not a heading node", func(t *testing.T) {
    node := &ast.Paragraph{}
    assert.False(t, isHeading(node))
  })
}

func TestConvertHeading(t *testing.T) {
  t.Run("handles nil heading", func(t *testing.T) {
    result := convertHeading(nil, []byte(""), "blue", "red", "green")
    assert.Nil(t, result, "Expected nil result for nil heading")
  })

  t.Run("handles empty heading text", func(t *testing.T) {
    node := &ast.Heading{
      BaseBlock: ast.BaseBlock{},
      Level:     1,
    }

    // Set up empty lines for the heading
    segments := textm.NewSegments()
    node.SetLines(segments)

    result := convertHeading(node, []byte(""), "blue", "red", "green")
    assert.Nil(t, result, "Expected nil result for empty heading text")
  })

  t.Run("can convert heading level 1", func(t *testing.T) {
    source := []byte("Heading level 1")
    node := &ast.Heading{
      BaseBlock: ast.BaseBlock{},
      Level:     1,
    }

    // Set up the lines for the heading
    segments := textm.NewSegments()
    segments.Append(textm.NewSegment(0, len(source)))
    node.SetLines(segments)

    expected := notionapi.Heading1Block{
      BasicBlock: notionapi.BasicBlock{},
      Heading1: notionapi.Heading{
        RichText: []notionapi.RichText{
          {
            Type:      notionapi.ObjectTypeText,
            Text:      &notionapi.Text{Content: string(source)},
            PlainText: string(source),
          },
        },
        Color: "blue",
      },
    }

    result := convertHeading(node, source, "blue", "red", "green")
    assertHeadingBlockEqual(t, expected, result)
  })

  t.Run("can convert heading level 2", func(t *testing.T) {
    source := []byte("Heading level 2")
    node := &ast.Heading{
      BaseBlock: ast.BaseBlock{},
      Level:     2,
    }

    // Set up the lines for the heading
    segments := textm.NewSegments()
    segments.Append(textm.NewSegment(0, len(source)))
    node.SetLines(segments)

    expected := notionapi.Heading2Block{
      BasicBlock: notionapi.BasicBlock{},
      Heading2: notionapi.Heading{
        RichText: []notionapi.RichText{
          {
            Type:      notionapi.ObjectTypeText,
            Text:      &notionapi.Text{Content: string(source)},
            PlainText: string(source),
          },
        },
        Color: "red",
      },
    }

    result := convertHeading(node, source, "blue", "red", "green")
    assertHeadingBlockEqual(t, expected, result)
  })

  t.Run("can convert heading level 3", func(t *testing.T) {
    source := []byte("Heading level 3")
    node := &ast.Heading{
      BaseBlock: ast.BaseBlock{},
      Level:     3,
    }

    // Set up the lines for the heading
    segments := textm.NewSegments()
    segments.Append(textm.NewSegment(0, len(source)))
    node.SetLines(segments)

    expected := notionapi.Heading3Block{
      BasicBlock: notionapi.BasicBlock{},
      Heading3: notionapi.Heading{
        RichText: []notionapi.RichText{
          {
            Type:      notionapi.ObjectTypeText,
            Text:      &notionapi.Text{Content: string(source)},
            PlainText: string(source),
          },
        },
        Color: "green",
      },
    }

    result := convertHeading(node, source, "blue", "red", "green")
    assertHeadingBlockEqual(t, expected, result)
  })
}

func assertHeadingBlockEqual(t *testing.T, expected, actual notionapi.Block) {
  switch expected := expected.(type) {
  case notionapi.Heading1Block:
    actual, ok := actual.(notionapi.Heading1Block)
    assert.True(t, ok)
    assert.Equal(t, expected.GetRichTextString(), actual.GetRichTextString())
  case notionapi.Heading2Block:
    actual, ok := actual.(notionapi.Heading2Block)
    assert.True(t, ok)
    assert.Equal(t, expected.GetRichTextString(), actual.GetRichTextString())
  case notionapi.Heading3Block:
    actual, ok := actual.(notionapi.Heading3Block)
    assert.True(t, ok)
    assert.Equal(t, expected.GetRichTextString(), actual.GetRichTextString())
  default:
    t.Fatalf("unexpected block type: %T", expected)
  }
}
