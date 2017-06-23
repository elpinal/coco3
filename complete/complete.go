package complete

import "os"

func File(buf []rune, pos int) ([]string, error) {
	dir, err := os.Open(".")
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	return dir.Readdirnames(0)
}
