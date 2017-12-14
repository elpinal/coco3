package complete

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func File(buf []rune, pos int) ([]string, error) {
	i := strings.LastIndexAny(string(buf[:pos]), " '") + 1
	prefix := string(buf[i:pos])
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
	sort.Strings(dirnames)
	for _, name := range dirnames {
		if !strings.HasPrefix(name, pend) {
			continue
		}
		n := name[len(pend):]
		stat, err := os.Stat(filepath.Join(p, name))
		if err != nil {
			return nil, err
		}
		if stat.IsDir() {
			n += "/"
		}
		names = append(names, n)
	}
	return names, nil
}

func FromPath(buf []rune, pos int) ([]string, error) {
	i := strings.LastIndexAny(string(buf[:pos]), " '") + 1
	prefix := string(buf[i:pos])
	path := os.Getenv("PATH")
	names := make([]string, 0, 1)
	for _, dir := range filepath.SplitList(path) {
		if dir == "" {
			dir = "."
		}
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if !strings.HasPrefix(file.Name(), prefix) {
				continue
			}
			names = append(names, file.Name()[len(prefix):])
		}
	}
	return names, nil
}
