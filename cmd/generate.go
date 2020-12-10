package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
	"github.com/iancoleman/strcase"
	"github.com/mtsmfm/sqlcodegen/utils"
)

func RunGenerate() error {
	configPath, err := utils.FindConfigPath()
	if err != nil {
		return err
	}

	config, err := utils.LoadConfig(configPath)
	if err != nil {
		return err
	}

	rootDir := filepath.Dir(configPath)

	sqlStringLiterals, err := extractAllSQLStringLiterals(rootDir)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	buf.WriteString("// Code generated by sqlcodegen, DO NOT EDIT.\n")
	buf.WriteString("\n")
	buf.WriteString("package " + config.Out.Package + "\n")

	for _, p := range config.Imports {
		buf.WriteString(fmt.Sprintf("import \"%s\"\n", p))
	}

	structNameMap := make(map[string]bool)

	for _, sqlStr := range sqlStringLiterals {
		columns, err := utils.AnalyzeSQL(sqlStr.Value, os.Getenv("DATABASE_URL"))
		if err != nil {
			return err
		}

		structName, err := buildStructName(config, columns, sqlStr)
		if err != nil {
			return err
		}

		if structNameMap[structName] {
			continue
		}

		structNameMap[structName] = true

		for i, c1 := range columns {
			for j, c2 := range columns {
				if c1.Name == c2.Name && i != j {
					return fmt.Errorf("SQL `%s` returns ambigious column `%s`\n-- SQL --\n%s\n---------\n", structName, c1.Name, sqlStr.Value)
				}
			}
		}

		buf.WriteString(utils.Comment(sqlStr.Value))
		buf.WriteString("type " + structName + " struct {" + "\n")

		for _, col := range columns {
			t, err := goTypeFor(config, col.Type)
			if err != nil {
				return err
			}

			tags := make([]string, len(config.Tags))

			for i, k := range config.Tags {
				tags[i] = fmt.Sprintf("%s:\"%s\"", k, col.Name)
			}

			tagString := ""

			if len(tags) > 0 {
				tagString = "`" + strings.Join(tags, " ") + "`"
			}

			buf.WriteString(strcase.ToCamel(col.Name) + " " + t + tagString + "\n")
		}

		buf.WriteString("}" + "\n")
	}

	targetFilePath := filepath.Join(rootDir, config.Out.File)
	err = os.MkdirAll(filepath.Dir(targetFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(targetFilePath)
	file.Write(utils.Format(targetFilePath, buf))
	if err != nil {
		return err
	}

	return nil
}

func extractAllSQLStringLiterals(dir string) ([]utils.StringLiteral, error) {
	var sqlLiterals []utils.StringLiteral

	files, err := doublestar.Glob(filepath.Join(dir, "**", "*.go"))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		xs, err := utils.ExtractCommentedStringLiterals(f, regexp.MustCompile(`// sqlcodegen (\w+)`))
		if err != nil {
			return nil, err
		}

		sqlLiterals = append(sqlLiterals, xs...)
	}

	return sqlLiterals, nil
}

func buildStructName(cfg *utils.Config, columns []utils.SqlResultColumn, sql utils.StringLiteral) (string, error) {
	result := regexp.MustCompile(`// sqlcodegen (\w+)`).FindStringSubmatch(sql.Comment)
	if len(result) == 2 {
		return result[1], nil
	}

	return "", fmt.Errorf("Struct name is wrong %+v", result)
}

func goTypeFor(cfg *utils.Config, sqlType string) (string, error) {
	result, ok := cfg.Typemap[sqlType]

	if !ok {
		return "", fmt.Errorf("go type mapping for %s is not defined", sqlType)
	}

	return result, nil
}
