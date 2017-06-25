package complete

import (
	"os"
	"path/filepath"
	"strings"
)

func File(buf []rune, pos int) ([]string, error) {
	words := strings.Split(string(buf), " ")
	prefix := words[len(words)-1]
	p, pend := filepath.Split(prefix)
	if p == "" {
		p = "."
	}
	if strings.HasPrefix(p, "~") {
		p = os.Getenv("HOME") + p[1:]
	}
	dir, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	names := make([]string, 0, len(pend))
	dirnames, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	for _, name := range dirnames {
		if !strings.HasPrefix(name, pend) {
			continue
		}
		names = append(names, name[len(pend):])
	}
	return names, nil
}
