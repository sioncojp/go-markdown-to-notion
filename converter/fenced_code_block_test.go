package converter

import (
  "testing"

  "github.com/jomei/notionapi"
  "github.com/sioncojp/go-markdown-to-notion/chunk"
  "github.com/stretchr/testify/assert"
  "github.com/yuin/goldmark/ast"
  "github.com/yuin/goldmark/text"
)

func TestConvertFencedCodeBlock(t *testing.T) {
  t.Run("converts code block with content and language", func(t *testing.T) {
    // create a sample test source code
    source := []byte("```go\nfmt.Println(\"Hello, World!\")\n```")

    codeBlockNode := &ast.FencedCodeBlock{
      BaseBlock: ast.BaseBlock{},
      Info:      ast.NewText(),
    }

    // setting Info (language info)
    codeBlockNode.Info.Segment = text.NewSegment(3, 5) // "go"の部分

    // setting code block content
    lines := text.NewSegments()
    lines.Append(text.NewSegment(6, 33)) // "fmt.Println(\"Hello, World!\")"の部分
    codeBlockNode.SetLines(lines)

    language := extractLanguage(codeBlockNode, source)
    assert.Equal(t, "go", language, "Expected language to be 'go'")

    expected := notionapi.CodeBlock{
      BasicBlock: notionapi.BasicBlock{
        Type: notionapi.BlockTypeCode,
      },
      Code: notionapi.Code{
        RichText: chunk.RichText("fmt.Println(\"Hello, World!\")", nil),
        Language: language,
      },
    }

    // execute the function
    codeBlock := convertFencedCodeBlock(codeBlockNode, source)
    assert.NotNil(t, codeBlock, "Expected code block to not be nil")
    assert.Equal(t, expected.Code.Language, codeBlock.Code.Language, "Expected language to match")
    assert.Equal(t, len(expected.Code.RichText), len(codeBlock.Code.RichText), "Expected rich text length to match")
  })

  t.Run("handles nil code block", func(t *testing.T) {
    codeBlock := convertFencedCodeBlock(nil, nil)
    assert.Nil(t, codeBlock, "Expected code block to be nil for nil code block")
  })

  t.Run("handles code block with no language", func(t *testing.T) {
    source := []byte("```\nfmt.Println(\"Hello, World!\")\n```")

    codeBlockNode := &ast.FencedCodeBlock{
      BaseBlock: ast.BaseBlock{},
    }

    lines := text.NewSegments()
    lines.Append(text.NewSegment(4, 31)) // "fmt.Println(\"Hello, World!\")"の部分
    codeBlockNode.SetLines(lines)

    codeBlock := convertFencedCodeBlock(codeBlockNode, source)
    assert.NotNil(t, codeBlock, "Expected code block to not be nil")
    assert.Equal(t, "", codeBlock.Code.Language, "Expected language to be empty")
  })

  t.Run("handles code block with empty content", func(t *testing.T) {
    source := []byte("```go\n```")

    codeBlockNode := &ast.FencedCodeBlock{
      BaseBlock: ast.BaseBlock{},
      Info:      ast.NewText(),
    }

    codeBlockNode.Info.Segment = text.NewSegment(3, 5) // "go"の部分

    lines := text.NewSegments()
    codeBlockNode.SetLines(lines)

    codeBlock := convertFencedCodeBlock(codeBlockNode, source)
    assert.Nil(t, codeBlock, "Expected code block to be nil for empty content")
  })

  t.Run("handles code block with invalid language", func(t *testing.T) {
    source := []byte("```invalid-language\nfmt.Println(\"Hello, World!\")\n```")

    codeBlockNode := &ast.FencedCodeBlock{
      BaseBlock: ast.BaseBlock{},
      Info:      ast.NewText(),
    }

    codeBlockNode.Info.Segment = text.NewSegment(3, 18) // "invalid-language"の部分

    lines := text.NewSegments()
    lines.Append(text.NewSegment(19, 46)) // "fmt.Println(\"Hello, World!\")"の部分
    codeBlockNode.SetLines(lines)

    language := extractLanguage(codeBlockNode, source)
    assert.Equal(t, "", language, "Expected language to be empty for invalid language")

    codeBlock := convertFencedCodeBlock(codeBlockNode, source)
    assert.NotNil(t, codeBlock, "Expected code block to not be nil")
    assert.Equal(t, "", codeBlock.Code.Language, "Expected language to be empty for invalid language")
  })
}
