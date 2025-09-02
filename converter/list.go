package converter

import (
  "github.com/jomei/notionapi"
  "github.com/yuin/goldmark/ast"
)

// isList checks if a node is a list.
func isList(node ast.Node) bool {
  _, ok := node.(*ast.List)
  return ok
}

// isListItem checks if a node is a list item.
func isListItem(node ast.Node) bool {
  _, ok := node.(*ast.ListItem)
  return ok
}

// convertList converts a list node to Notion list blocks.
func convertList(node *ast.List, source []byte) []notionapi.Block {
  if node == nil {
    return nil
  }

  var items []notionapi.Block

  for child := node.FirstChild(); child != nil; child = child.NextSibling() {
    listItem, ok := child.(*ast.ListItem)
    if !ok {
      continue
    }

    var nestedBlocks []notionapi.Block
    for grandChild := listItem.FirstChild(); grandChild != nil; grandChild = grandChild.NextSibling() {
      if nestedList, ok := grandChild.(*ast.List); ok {
        nestedListBlocks := convertList(nestedList, source)
        nestedBlocks = append(nestedBlocks, nestedListBlocks...)
      }
    }

    if node.IsOrdered() {
      item := convertNumberListItem(child.(*ast.ListItem), source, nestedBlocks)
      if item != nil {
        items = append(items, item)
      }
    } else {
      item := convertListItem(child.(*ast.ListItem), source, nestedBlocks)
      if item != nil {
        items = append(items, item)
      }
    }
  }

  return items
}

// convertListItem ... converts a bulleted list item to a Notion block.
func convertListItem(node *ast.ListItem, source []byte, children []notionapi.Block) notionapi.Block {
  if node == nil {
    return nil
  }

  // Create a bulleted list item block by default
  block := notionapi.BulletedListItemBlock{
    BasicBlock: notionapi.BasicBlock{
      Object: notionapi.ObjectTypeBlock,
      Type:   notionapi.BlockTypeBulletedListItem,
    },
    BulletedListItem: notionapi.ListItem{
      RichText: convertListItemContent(node, source),
      Children: children,
    },
  }

  return block
}

// convertNumberListItem ... converts a numbered list item node to a Notion numbered list item block.
func convertNumberListItem(node *ast.ListItem, source []byte, children []notionapi.Block) notionapi.Block {
  if node == nil {
    return nil
  }

  // Create a bulleted list item block by default
  block := notionapi.NumberedListItemBlock{
    BasicBlock: notionapi.BasicBlock{
      Object: notionapi.ObjectTypeBlock,
      Type:   notionapi.BlockTypeNumberedListItem,
    },
    NumberedListItem: notionapi.ListItem{
      RichText: convertListItemContent(node, source),
      Children: children,
    },
  }

  return block
}

// convertListItemContent ... converts the content of a list item to Notion rich text.
func convertListItemContent(node *ast.ListItem, source []byte) []notionapi.RichText {
  if node == nil {
    return nil
  }

  var blocks []notionapi.RichText
  for child := node.FirstChild(); child != nil; child = child.NextSibling() {
    // Skip nested lists
    if _, ok := child.(*ast.List); ok {
      continue
    }

    // Convert other content
    richText := convertChildNodesToRichText(child, source)
    if richText != nil {
      blocks = append(blocks, richText...)
    }
  }

  return blocks
}
