package converter

import (
	"regexp"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/sioncojp/go-markdown-to-notion/chunk"
	gast "github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

// isTable checks if a node is a table.
func isTable(node gast.Node) bool {
	_, ok := node.(*east.Table)
	return ok
}

// isTableHeader checks if a node is a table header.
func isTableHeader(node gast.Node) bool {
	_, ok := node.(*east.TableHeader)
	return ok
}

// convertTable converts a table node to a Notion table block.
func convertTable(node *east.Table, source []byte) *notionapi.TableBlock {
	if node == nil {
		return nil
	}

	// Create table block
	tableBlock := &notionapi.TableBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockType("table"),
		},
		Table: notionapi.Table{
			HasColumnHeader: true,
			HasRowHeader:    false,
			Children:        []notionapi.Block{},
		},
	}

	// Process rows
	var headerProcessed bool
	var columnCount int

	// Process all rows, checking for header row
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		// Check if this is a header row
		if child.Kind() == east.KindTableHeader {
			// Process header row
			header := child.(*east.TableHeader)

			// Convert header cells to rich text
			var headerCells [][]notionapi.RichText
			for headerCell := header.FirstChild(); headerCell != nil; headerCell = headerCell.NextSibling() {
				tableCell, ok := headerCell.(*east.TableCell)
				if !ok {
					continue
				}

				// Convert cell content to rich text
				richText := convertChildNodesToRichText(tableCell, source)
				if richText == nil {
					// Add empty cell
					headerCells = append(headerCells, []notionapi.RichText{})
				} else {
					headerCells = append(headerCells, richText)
				}
			}

			if len(headerCells) == 0 {
				continue
			}

			// Set column count based on the header row
			columnCount = len(headerCells)
			tableBlock.Table.TableWidth = columnCount
			headerProcessed = true

			// Create header row block
			headerRowBlock := &notionapi.TableRowBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: notionapi.ObjectTypeBlock,
					Type:   notionapi.BlockType("table_row"),
				},
				TableRow: notionapi.TableRow{
					Cells: headerCells,
				},
			}

			tableBlock.Table.Children = append(tableBlock.Table.Children, headerRowBlock)
		} else if child.Kind() == east.KindTableRow {
			// Process regular row
			tableRow := child.(*east.TableRow)

			// Convert row to table row block
			rowCells := convertTableRow(tableRow, source)
			if len(rowCells) == 0 {
				continue
			}

			// Set column count based on the first row if header wasn't processed
			if !headerProcessed {
				columnCount = len(rowCells)
				tableBlock.Table.TableWidth = columnCount
				headerProcessed = true
			}

			// Create table row block
			rowBlock := &notionapi.TableRowBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: notionapi.ObjectTypeBlock,
					Type:   notionapi.BlockType("table_row"),
				},
				TableRow: notionapi.TableRow{
					Cells: rowCells,
				},
			}

			tableBlock.Table.Children = append(tableBlock.Table.Children, rowBlock)
		}
	}

	// If no rows were added, return nil
	if len(tableBlock.Table.Children) == 0 {
		return nil
	}

	return tableBlock
}

// convertTableRow converts a table row node to rich text cells.
func convertTableRow(node *east.TableRow, source []byte) [][]notionapi.RichText {
	if node == nil {
		return nil
	}

	var cells [][]notionapi.RichText

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		tableCell, ok := child.(*east.TableCell)
		if !ok {
			continue
		}

		// Convert cell content to rich text
		richText := convertChildNodesToRichText(tableCell, source)
		if richText == nil {
			// Add empty cell
			cells = append(cells, []notionapi.RichText{})
		} else {
			cells = append(cells, richText)
		}
	}

	return cells
}


// parseTable parses a Markdown table from the source text and converts it to a Notion table block.
func parseTable(source []byte) *notionapi.TableBlock {
	if source == nil || len(source) == 0 {
		return nil
	}

	// Convert source to string for easier processing
	text := string(source)

	// Regular expression to match a Markdown table
	// This regex matches:
	// 1. A line starting with | and ending with | (header row)
	// 2. A line with |---|---| format (separator row)
	// 3. One or more lines starting with | and ending with | (data rows)
	tableRegex := regexp.MustCompile(`(?m)^\|(.+)\|\s*$\n^\|(\s*[-:]+[-|\s:]*)\|\s*$(\n^\|(.+)\|\s*$)*`)

	// Find all tables in the source
	tables := tableRegex.FindAllString(text, -1)
	if len(tables) == 0 {
		return nil
	}

	// Process the first table found
	tableText := tables[0]

	// Split the table into rows
	rows := strings.Split(tableText, "\n")
	if len(rows) < 3 { // Need at least header, separator, and one data row
		return nil
	}

	// Parse header row
	headerRow := parseTableRow(rows[0])
	if len(headerRow) == 0 {
		return nil
	}

	// Skip separator row (rows[1])

	// Parse data rows
	var dataRows [][]string
	for i := 2; i < len(rows); i++ {
		if rows[i] == "" {
			continue
		}
		dataRow := parseTableRow(rows[i])
		if len(dataRow) > 0 {
			dataRows = append(dataRows, dataRow)
		}
	}

	// Create table block
	tableBlock := &notionapi.TableBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockType("table"),
		},
		Table: notionapi.Table{
			TableWidth:      len(headerRow),
			HasColumnHeader: true,
			HasRowHeader:    false,
			Children:        []notionapi.Block{},
		},
	}

	// Add header row
	headerRowBlock := createTableRowBlock(headerRow)
	if headerRowBlock != nil {
		tableBlock.Table.Children = append(tableBlock.Table.Children, headerRowBlock)
	}

	// Add data rows
	for _, dataRow := range dataRows {
		// Ensure the row has the same number of columns as the header
		for len(dataRow) < len(headerRow) {
			dataRow = append(dataRow, "")
		}
		dataRowBlock := createTableRowBlock(dataRow)
		if dataRowBlock != nil {
			tableBlock.Table.Children = append(tableBlock.Table.Children, dataRowBlock)
		}
	}

	// If no rows were added, return nil
	if len(tableBlock.Table.Children) == 0 {
		return nil
	}

	return tableBlock
}

// parseTableRow parses a Markdown table row and returns the cell values.
func parseTableRow(row string) []string {
	// Remove leading and trailing | characters
	row = strings.TrimSpace(row)
	if len(row) < 2 || !strings.HasPrefix(row, "|") || !strings.HasSuffix(row, "|") {
		return nil
	}
	row = row[1 : len(row)-1]

	// Split by | character
	cells := strings.Split(row, "|")

	// Trim whitespace from each cell
	for i, cell := range cells {
		cells[i] = strings.TrimSpace(cell)
	}

	return cells
}

// createTableRowBlock creates a Notion table row block from cell values.
func createTableRowBlock(cells []string) *notionapi.TableRowBlock {
	if len(cells) == 0 {
		return nil
	}

	// Create cells array for the row block
	richTextCells := make([][]notionapi.RichText, len(cells))

	// Convert each cell value to rich text
	for i, cell := range cells {
		richTextCells[i] = chunk.RichText(cell, nil)
	}

	// Create table row block
	return &notionapi.TableRowBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockType("table_row"),
		},
		TableRow: notionapi.TableRow{
			Cells: richTextCells,
		},
	}
}
