package cmd

import (
	"os"
	"strings"

	"github.com/mtsmfm/sqlcodegen/utils"
)

func RunInit() error {
	file, err := os.Create(utils.ConfigFileName)

	if err != nil {
		return err
	}

	file.WriteString(strings.TrimSpace(`
out:
  package: sqlstructs
  file: sqlstructs/sqlstructs.go
tags:
  - db
  - json
  #- yaml
  #- toml
imports:
  - github.com/lib/pq
  - database/sql
typemap:
  bigint: sql.NullInt64
  integer: sql.NullInt64
  uuid: sql.NullString
  text: sql.NullString
  text[]: pq.StringArray
`))

	file.WriteString("\n")

	file.Close()

	return nil
}
