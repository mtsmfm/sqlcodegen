package utils

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"
)

type StringLiteral struct {
	Value   string
	Comment string
}

func ExtractCommentedStringLiterals(filePath string, re *regexp.Regexp) ([]StringLiteral, error) {
	fset := token.NewFileSet()
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	file, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)

	if err != nil {
		log.Printf("Parse error %+v", err)
		return nil, nil
	}

	var results []StringLiteral

	astutil.Apply(file, func(cr *astutil.Cursor) bool {
		switch v := cr.Node().(type) {
		case *ast.BasicLit:
			if v.Kind == token.STRING {
				str, err := strconv.Unquote(v.Value)
				if err != nil {
					break
				}
				line := fset.Position(v.Pos()).Line
				var comment string

				for _, cg := range file.Comments {
					for _, c := range cg.List {
						if fset.Position(c.Pos()).Line == line-1 {
							comment = c.Text
						}
					}
				}

				if re.MatchString(comment) {
					results = append(results, StringLiteral{Value: str, Comment: comment})
				}
			}
		}
		return true
	}, nil)

	return results, nil
}

func Format(filename string, buf bytes.Buffer) []byte {
	src, err := imports.Process(filename, buf.Bytes(), nil)
	if err != nil {
		log.Printf("format failed: %+v", err)
		return buf.Bytes()
	}

	return src
}

func Comment(str string) string {
	return "/*\n" + indentMultiLine(strip(str), 1) + "\n*/\n"
}

func strip(str string) string {
	lines := strings.Split(str, "\n")
	var result []string

	for _, line := range lines {
		s := strings.TrimSpace(line)
		if s != "" {
			result = append(result, s)
		}
	}

	return strings.Join(result, "\n")
}

func indentMultiLine(str string, indent int) string {
	lines := strings.Split(str, "\n")
	var result []string

	for _, line := range lines {
		for i := 0; i < indent; i++ {
			line = "\t" + line
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
