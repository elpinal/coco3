package complete

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

func File(buf []rune, pos int) ([][]rune, error) {
	i := strings.LastIndexAny(string(buf[:pos]), " '") + 1
	prefix := buf[i:pos]
	p, pend := filepath.Split(string(prefix)) // TODO: introduce a []rune version of filepath.Split.
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
	names := make([][]rune, 0, len(pend))
	dirnames, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	sort.Strings(dirnames)
	for _, name := range dirnames {
		if !strings.HasPrefix(name, pend) {
			continue
		}
		n := []rune(name)[len([]rune(pend)):]
		stat, err := os.Stat(filepath.Join(p, name))
		if err != nil {
			return nil, err
		}
		if stat.IsDir() {
			n = append(n, '/')
		}
		names = append(names, n)
	}
	return names, nil
}

func ioReadDir(dir string) ([][]rune, error) {
	_, err := os.Stat(dir)
	if err != nil {
		// Return due to a kind of errors where a directory does not exist.
		return nil, nil
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make([][]rune, len(files))
	for i := range files {
		names[i] = []rune(files[i].Name())
	}
	return names, nil
}

func FromPath(buf []rune, pos int) ([][]rune, error) {
	return fromPath(buf, pos, os.Getenv("PATH"), ioReadDir)
}

type dirReader func(string) ([][]rune, error)

func fromPath(buf []rune, pos int, path string, readDir dirReader) ([][]rune, error) {
	// TODO: Parsing buf will make it more convenient to complete.
	// Currently implemented as simple and fast, but not solid.
	i := strings.LastIndexAny(string(buf[:pos]), " '!") + 1
	prefix := string(buf[i:pos])
	names := make([][]rune, 0, 1)
	for _, dir := range filepath.SplitList(path) {
		if dir == "" {
			dir = "."
		}
		files, err := readDir(dir)
		if err != nil {
			return nil, errors.Wrap(err, "PATH completion failed")
		}
		for _, name := range files {
			if !strings.HasPrefix(string(name), prefix) {
				continue
			}
			names = append(names, name[len(prefix):])
		}
	}
	return names, nil
}
