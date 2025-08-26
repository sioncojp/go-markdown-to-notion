package chunk

import (
  "testing"

  "github.com/jomei/notionapi"
  "github.com/stretchr/testify/assert"
)

func TestBlocks(t *testing.T) {
  t.Run("empty blocks", func(t *testing.T) {
    blocks := []notionapi.Block{}
    result := Blocks(blocks)
    assert.Empty(t, result, "Expected empty result for empty input")
  })

  t.Run("blocks less than limit", func(t *testing.T) {
    // Create a slice with fewer blocks than the limit
    blocks := make([]notionapi.Block, BlockLimit-10)
    for i := 0; i < BlockLimit-10; i++ {
      blocks[i] = &notionapi.ParagraphBlock{
        BasicBlock: notionapi.BasicBlock{
          Type: notionapi.BlockTypeParagraph,
        },
      }
    }

    result := Blocks(blocks)

    assert.Equal(t, 1, len(result), "Expected 1 chunk for blocks less than limit")
    assert.Equal(t, BlockLimit-10, len(result[0]), "Expected all blocks in the first chunk")
  })

  t.Run("blocks equal to limit", func(t *testing.T) {
    // Create a slice with exactly the limit number of blocks
    blocks := make([]notionapi.Block, BlockLimit)
    for i := 0; i < BlockLimit; i++ {
      blocks[i] = &notionapi.ParagraphBlock{
        BasicBlock: notionapi.BasicBlock{
          Type: notionapi.BlockTypeParagraph,
        },
      }
    }

    result := Blocks(blocks)

    assert.Equal(t, 1, len(result), "Expected 1 chunk for blocks equal to limit")
    assert.Equal(t, BlockLimit, len(result[0]), "Expected all blocks in the first chunk")
  })

  t.Run("blocks more than limit", func(t *testing.T) {
    // Create a slice with more blocks than the limit
    totalBlocks := BlockLimit + 50
    blocks := make([]notionapi.Block, totalBlocks)
    for i := 0; i < totalBlocks; i++ {
      blocks[i] = &notionapi.ParagraphBlock{
        BasicBlock: notionapi.BasicBlock{
          Type: notionapi.BlockTypeParagraph,
        },
      }
    }

    result := Blocks(blocks)

    assert.Equal(t, 2, len(result), "Expected 2 chunks for blocks more than limit")
    assert.Equal(t, BlockLimit, len(result[0]), "Expected first chunk to have BlockLimit blocks")
    assert.Equal(t, 50, len(result[1]), "Expected second chunk to have remaining blocks")
  })
}

func TestRichText(t *testing.T) {
  t.Run("can convert rich text under character limit", func(t *testing.T) {
    content := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
    expectedBlocks := []notionapi.RichText{
      {
        Type: notionapi.ObjectTypeText,
        Text: &notionapi.Text{
          Content: content,
        },
        PlainText: content,
      },
    }

    blocks := RichText(content, nil)

    assert.Equal(t, len(expectedBlocks), len(blocks), "Expected %d blocks, but got %d", len(expectedBlocks), len(blocks))

    for i, block := range blocks {
      expectedBlock := expectedBlocks[i]

      assert.Equal(t, expectedBlock.Type, block.Type, "Expected block type %s, but got %s", expectedBlock.Type, block.Type)
      assert.Equal(t, expectedBlock.PlainText, block.PlainText, "Expected plain text %s, but got %s", expectedBlock.PlainText, block.PlainText)
      assert.Equal(t, expectedBlock.Text.Content, block.Text.Content, "Expected text content %s, but got %s", expectedBlock.Text.Content, block.Text.Content)
    }
  })

  t.Run("can convert rich text over character limit", func(t *testing.T) {
    // Create a string over 2000 characters
    var longContent string
    for i := 0; i <= CharacterLimit; i++ {
      longContent += "a"
    }

    expectedBlocks := []notionapi.RichText{
      {
        Type: notionapi.ObjectTypeText,
        Text: &notionapi.Text{
          Content: longContent[:CharacterLimit],
        },
        PlainText: longContent[:CharacterLimit],
      },
      {
        Type: notionapi.ObjectTypeText,
        Text: &notionapi.Text{
          Content: longContent[CharacterLimit:],
        },
        PlainText: longContent[CharacterLimit:],
      },
    }

    result := RichText(longContent, nil)

    assert.Equal(t, len(expectedBlocks), len(result), "Expected %d blocks, but got %d", len(expectedBlocks), len(result))
    for i, block := range result {
      expectedBlock := expectedBlocks[i]

      assert.Equal(t, expectedBlock.Type, block.Type, "Expected block type %s, but got %s", expectedBlock.Type, block.Type)
      assert.Equal(t, expectedBlock.PlainText, block.PlainText, "Expected plain text %s, but got %s", expectedBlock.PlainText, block.PlainText)
      assert.Equal(t, expectedBlock.Text.Content, block.Text.Content, "Expected text content %s, but got %s", expectedBlock.Text.Content, block.Text.Content)
    }
  })
}

func TestRichTextWithLink(t *testing.T) {
  t.Run("can convert rich text with link under character limit", func(t *testing.T) {
    content := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
    link := "https://example.com"
    expectedBlocks := []notionapi.RichText{
      {
        Type: notionapi.ObjectTypeText,
        Text: &notionapi.Text{
          Content: content,
          Link: &notionapi.Link{
            Url: link,
          },
        },
        PlainText: content,
      },
    }

    blocks := RichTextWithLink(content, link)

    assert.Equal(t, len(expectedBlocks), len(blocks), "Expected %d blocks, but got %d", len(expectedBlocks), len(blocks))

    for i, block := range blocks {
      expectedBlock := expectedBlocks[i]

      assert.Equal(t, expectedBlock.Type, block.Type, "Expected block type %s, but got %s", expectedBlock.Type, block.Type)
      assert.Equal(t, expectedBlock.PlainText, block.PlainText, "Expected plain text %s, but got %s", expectedBlock.PlainText, block.PlainText)
      assert.Equal(t, expectedBlock.Text.Content, block.Text.Content, "Expected text content %s, but got %s", expectedBlock.Text.Content, block.Text.Content)
      assert.NotNil(t, block.Text.Link, "Expected link to be not nil")
      assert.Equal(t, expectedBlock.Text.Link.Url, block.Text.Link.Url, "Expected link URL %s, but got %s", expectedBlock.Text.Link.Url, block.Text.Link.Url)
    }
  })

  t.Run("can convert rich text with link over character limit", func(t *testing.T) {
    // Create a string over 2000 characters
    var longContent string
    for i := 0; i <= CharacterLimit; i++ {
      longContent += "a"
    }
    link := "https://example.com"

    expectedBlocks := []notionapi.RichText{
      {
        Type: notionapi.ObjectTypeText,
        Text: &notionapi.Text{
          Content: longContent[:CharacterLimit],
          Link: &notionapi.Link{
            Url: link,
          },
        },
        PlainText: longContent[:CharacterLimit],
      },
      {
        Type: notionapi.ObjectTypeText,
        Text: &notionapi.Text{
          Content: longContent[CharacterLimit:],
          Link: &notionapi.Link{
            Url: link,
          },
        },
        PlainText: longContent[CharacterLimit:],
      },
    }

    result := RichTextWithLink(longContent, link)

    assert.Equal(t, len(expectedBlocks), len(result), "Expected %d blocks, but got %d", len(expectedBlocks), len(result))
    for i, block := range result {
      expectedBlock := expectedBlocks[i]

      assert.Equal(t, expectedBlock.Type, block.Type, "Expected block type %s, but got %s", expectedBlock.Type, block.Type)
      assert.Equal(t, expectedBlock.PlainText, block.PlainText, "Expected plain text %s, but got %s", expectedBlock.PlainText, block.PlainText)
      assert.Equal(t, expectedBlock.Text.Content, block.Text.Content, "Expected text content %s, but got %s", expectedBlock.Text.Content, block.Text.Content)
      assert.NotNil(t, block.Text.Link, "Expected link to be not nil")
      assert.Equal(t, expectedBlock.Text.Link.Url, block.Text.Link.Url, "Expected link URL %s, but got %s", expectedBlock.Text.Link.Url, block.Text.Link.Url)
    }
  })
}
