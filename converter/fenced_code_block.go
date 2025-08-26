package converter

import (
  "github.com/jomei/notionapi"
  "github.com/sioncojp/go-markdown-to-notion/chunk"
  "github.com/yuin/goldmark/ast"
)

// isCodeBlock checks if a node is a code block.
func isCodeBlock(node ast.Node) bool {
  _, ok := node.(*ast.FencedCodeBlock)
  return ok
}

// extractLanguage extracts the language from a code block node.
func extractLanguage(node *ast.FencedCodeBlock, source []byte) string {
  if node == nil {
    return ""
  }

  language := string(node.Language(source))
  valid := validateLanguage(language)
  if !valid {
    return ""
  }

  return language
}

func convertFencedCodeBlock(node *ast.FencedCodeBlock, source []byte) *notionapi.CodeBlock {
  if node == nil {
    return nil
  }

  content := string(node.Lines().Value(source))
  if content == "" {
    return nil
  }

  result := &notionapi.CodeBlock{
    BasicBlock: notionapi.BasicBlock{
      Object: notionapi.ObjectTypeBlock,
      Type:   notionapi.BlockTypeCode,
    },
    Code: notionapi.Code{
      RichText: chunk.RichText(content, nil),
    },
  }

  language := extractLanguage(node, source)
  if language != "" {
    result.Code.Language = language
  }

  return result
}

// validateLanguage checks if the code language is a valid option
// for Notion's code block.
//
// https://developers.notion.com/reference/block#code
func validateLanguage(language string) bool {
  validLanguages := map[string]bool{
    "abap":          true,
    "arduino":       true,
    "bash":          true,
    "basic":         true,
    "c":             true,
    "clojure":       true,
    "coffeescript":  true,
    "c++":           true,
    "c#":            true,
    "css":           true,
    "dart":          true,
    "diff":          true,
    "docker":        true,
    "elixir":        true,
    "elm":           true,
    "erlang":        true,
    "flow":          true,
    "fortran":       true,
    "f#":            true,
    "gherkin":       true,
    "glsl":          true,
    "go":            true,
    "graphql":       true,
    "groovy":        true,
    "haskell":       true,
    "html":          true,
    "java":          true,
    "javascript":    true,
    "json":          true,
    "julia":         true,
    "kotlin":        true,
    "latex":         true,
    "less":          true,
    "lisp":          true,
    "livescript":    true,
    "lua":           true,
    "makefile":      true,
    "markdown":      true,
    "markup":        true,
    "matlab":        true,
    "mermaid":       true,
    "nix":           true,
    "objective-c":   true,
    "ocaml":         true,
    "pascal":        true,
    "perl":          true,
    "php":           true,
    "plain text":    true,
    "powershell":    true,
    "prolog":        true,
    "protobuf":      true,
    "python":        true,
    "r":             true,
    "reason":        true,
    "ruby":          true,
    "rust":          true,
    "sass":          true,
    "scala":         true,
    "scheme":        true,
    "scss":          true,
    "shell":         true,
    "sql":           true,
    "swift":         true,
    "typescript":    true,
    "vb.net":        true,
    "verilog":       true,
    "vhdl":          true,
    "visual basic":  true,
    "webassembly":   true,
    "xml":           true,
    "yaml":          true,
    "java/c/c++/c#": true,
  }

  _, ok := validLanguages[language]
  return ok
}
