package complete

import (
	"os"
	"path/filepath"
	"strings"
)

func File(buf []rune, pos int) ([]string, error) {
	words := strings.Split(string(buf), " ")
	prefix := words[len(words)-1]
	p := filepath.Dir(prefix)
	dir, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	return dir.Readdirnames(0)
}
