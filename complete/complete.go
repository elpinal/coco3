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

func ioReadDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(files))
	for i := range files {
		names[i] = files[i].Name()
	}
	return names, nil
}

func FromPath(buf []rune, pos int) ([]string, error) {
	return fromPath(buf, pos, os.Getenv("PATH"), ioReadDir)
}

type dirReader func(string) ([]string, error)

func fromPath(buf []rune, pos int, path string, readDir dirReader) ([]string, error) {
	// TODO: Parsing buf will make it more convenient to complete.
	// Currently implemented as simple and fast, but not solid.
	i := strings.LastIndexAny(string(buf[:pos]), " '!") + 1
	prefix := string(buf[i:pos])
	names := make([]string, 0, 1)
	for _, dir := range filepath.SplitList(path) {
		if dir == "" {
			dir = "."
		}
		files, err := readDir(dir)
		if err != nil {
			return nil, err
		}
		for _, name := range files {
			if !strings.HasPrefix(name, prefix) {
				continue
			}
			names = append(names, name[len(prefix):])
		}
	}
	return names, nil
}
