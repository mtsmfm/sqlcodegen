package utils

import (
	"io"
	"os/exec"
	"strings"
)

type SqlResultColumn struct {
	Name string
	Type string
}

func AnalyzeSQL(sql string, url string) ([]SqlResultColumn, error) {
	cmd := exec.Command("psql", url, "-A", "-F,", "-t")
	// To suppress perl: warning: Setting locale failed.
	cmd.Env = []string{
		"LANG=C.UTF-8",
		"LANGUAGE=C.UTF-8",
		"LC_ALL=C.UTF-8",
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
	var result []SqlResultColumn

	for _, line := range lines {
		xs := strings.Split(line, ",")

		if len(xs) == 2 {
			result = append(result, SqlResultColumn{Name: xs[0], Type: xs[1]})
		}
	}

	return result, nil
}
