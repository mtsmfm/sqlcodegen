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
typemap:
  bigint: int
  uuid: string
  text: string
`))

	file.WriteString("\n")

	file.Close()

	return nil
}
