package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
	"golang.org/x/tools/go/ast/astutil"
)

type visitFn func(node ast.Node)

func (fn visitFn) Visit(node ast.Node) ast.Visitor {
	fn(node)
	return fn
}

type sqlResultColumn struct {
	Name string
	Type string
}

func analyzeSQL(sql string) ([]sqlResultColumn, error) {
	cmd := exec.Command("psql", "-A", "-F,", "-t")
	cmd.Env = []string{
		"PGHOST=postgres",
		"PGUSER=postgres",
		"PGPASSWORD=password",
		"PGDATABASE=postgres",
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, sql+"\\gdesc")
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var result []sqlResultColumn

	for _, line := range lines {
		xs := strings.Split(line, ",")

		if len(xs) == 2 {
			c := sqlResultColumn{Name: strcase.ToCamel(xs[0])}

			switch xs[1] {
			case "uuid":
				c.Type = "string"
			case "text":
				c.Type = "string"
			case "bigint":
				c.Type = "int64"
			default:
				return nil, fmt.Errorf("unknown type %+v for key %+v", xs[1], xs[0])
			}

			result = append(result, c)
		}
	}

	return result, nil
}

func run(filename string) error {
	fset := token.NewFileSet()
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	file, err := parser.ParseFile(fset, filename, src, parser.AllErrors)

	if err != nil {
		return err
	}

	var mostRecentDeclStmtPos token.Pos
	type data struct {
		Name string
		SQL  string
	}
	targets := make(map[token.Pos]data)

	astutil.Apply(file, func(cr *astutil.Cursor) bool {
		switch v := cr.Node().(type) {
		case *ast.DeclStmt:
			mostRecentDeclStmtPos = v.Pos()
			break
		case *ast.CallExpr:
			selector, ok := v.Fun.(*ast.SelectorExpr)
			if !ok {
				break
			}
			ident, ok := selector.X.(*ast.Ident)
			if !ok {
				break
			}

			if ident.Name == "db" && selector.Sel.Name == "Select" {
				if len(v.Args) < 2 {
					break
				}

				desc, ok := v.Args[0].(*ast.UnaryExpr)
				if !ok {
					break
				}
				descIdent, ok := desc.X.(*ast.Ident)
				if !ok {
					break
				}

				sql, ok := v.Args[1].(*ast.BasicLit)
				if !ok {
					break
				}

				sqlStr, err := strconv.Unquote(sql.Value)
				if err != nil {
					break
				}

				targets[mostRecentDeclStmtPos] = data{Name: descIdent.Name, SQL: sqlStr}
			}
			break
		}

		return true
	}, nil)

	n := astutil.Apply(file, func(cr *astutil.Cursor) bool {
		switch v := cr.Node().(type) {
		case *ast.DeclStmt:
			if t, ok := targets[v.Pos()]; ok {
				columns, err := analyzeSQL(t.SQL)
				if err != nil {
					log.Fatalf("%+v", err)
				}

				list := make([]*ast.Field, len(columns))

				for i, col := range columns {
					list[i] = &ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: col.Name,
							},
						},
						Type: &ast.Ident{
							Name: col.Type,
						},
					}
				}

				cr.Replace(&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Specs: []ast.Spec{&ast.ValueSpec{
							Type: &ast.ArrayType{
								Elt: &ast.StructType{
									Fields: &ast.FieldList{
										List: list,
									},
								},
							},
							Names: []*ast.Ident{
								&ast.Ident{
									Name: t.Name,
								},
							}}},
						Tok: token.VAR,
					},
				})
			}
			break
		}

		return true
	}, nil)

	if err := format.Node(os.Stdout, token.NewFileSet(), n); err != nil {
		return err
	}

	return nil
}

func main() {
	app := &cli.App{
		Name: "sqlcodegen",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Usage:    "File to check",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			filename := c.String("file")
			err := run(filename)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
