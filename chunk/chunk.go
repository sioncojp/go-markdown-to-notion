package chunk

import (
  "github.com/jomei/notionapi"
)

const (
  // BlockLimit ... maximum number of blocks that can be sent in a single request.
  BlockLimit = 100

  // CharacterLimit ... maximum number of characters that can be sent in a single rich text block.
  CharacterLimit = 2000
)

// Blocks ... Split blocks into chunks of BlockLimit size.
func Blocks(blocks []notionapi.Block) [][]notionapi.Block {
  var chunks [][]notionapi.Block
  chunkSize := BlockLimit

  for i := 0; i < len(blocks); i += chunkSize {
    end := i + chunkSize
    if end > len(blocks) {
      end = len(blocks)
    }
    chunks = append(chunks, blocks[i:end])
  }

  return chunks
}

// RichText ... Split rich text into chunks of CharacterLimit size.
func RichText(content string, annotations *notionapi.Annotations) []notionapi.RichText {
  var blocks []notionapi.RichText

  if len(content) <= CharacterLimit {
    richText := notionapi.RichText{
      Type: notionapi.ObjectTypeText,
      Text: &notionapi.Text{
        Content: content,
      },
      PlainText:   content,
      Annotations: annotations,
    }

    blocks = append(blocks, richText)
  } else {
    for i := 0; i < len(content); i += CharacterLimit {
      end := i + CharacterLimit
      if end > len(content) {
        end = len(content)
      }

      chunk := content[i:end]
      richText := notionapi.RichText{
        Type: notionapi.ObjectTypeText,
        Text: &notionapi.Text{
          Content: chunk,
        },
        PlainText:   chunk,
        Annotations: annotations,
      }

      blocks = append(blocks, richText)
    }
  }

  return blocks
}

// RichTextWithLink ... Split rich text with link into chunks of CharacterLimit size.
func RichTextWithLink(content string, link string) []notionapi.RichText {
  var blocks []notionapi.RichText

  if len(content) <= CharacterLimit {
    richText := notionapi.RichText{
      Type: notionapi.ObjectTypeText,
      Text: &notionapi.Text{
        Content: content,
        Link: &notionapi.Link{
          Url: link,
        },
      },
      PlainText: content,
    }

    blocks = append(blocks, richText)
  } else {
    for i := 0; i < len(content); i += CharacterLimit {
      end := i + CharacterLimit
      if end > len(content) {
        end = len(content)
      }

      chunk := content[i:end]
      richText := notionapi.RichText{
        Type: notionapi.ObjectTypeText,
        Text: &notionapi.Text{
          Content: chunk,
          Link: &notionapi.Link{
            Url: link,
          },
        },
        PlainText: chunk,
      }

      blocks = append(blocks, richText)
    }
  }

  return blocks
}
