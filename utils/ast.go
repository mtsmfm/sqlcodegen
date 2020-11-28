package utils

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type StringLiteral struct {
	Value   string
	Comment string
}

func ExtractStringLiterals(filePath string, re *regexp.Regexp) ([]StringLiteral, error) {
	fset := token.NewFileSet()
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	file, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)

	if err != nil {
		return nil, err
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

				if re.MatchString(str) {
					line := fset.Position(v.Pos()).Line
					var comment string

					for _, cg := range file.Comments {
						for _, c := range cg.List {
							if fset.Position(c.Pos()).Line == line-1 {
								comment = c.Text
							}
						}
					}

					results = append(results, StringLiteral{Value: str, Comment: comment})
				}
			}
		}
		return true
	}, nil)

	return results, nil
}

func Gofmt(buf bytes.Buffer) []byte {
	src, err := format.Source(buf.Bytes())

	if err != nil {
		log.Printf("gofmt failed: %+v", err)
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
