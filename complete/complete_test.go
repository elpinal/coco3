package complete

import (
	"reflect"
	"testing"
)

func TestFile(t *testing.T) {
	tests := []struct {
		input string
		want  [][]rune
	}{
		{
			input: "testdata/file",
			want:  [][]rune{[]rune("1.txt"), []rune("2.txt")},
		},
		{
			input: "testdata/",
			want:  [][]rune{[]rune("dir1/"), []rune("file1.txt"), []rune("file2.txt"), []rune("ó.txt")},
		},
		{
			input: "testdata/ó",
			want:  [][]rune{[]rune(".txt")},
		},
	}
	for _, test := range tests {
		input := []rune(test.input)
		list, err := File([]rune(input), len(input))
		if err != nil {
			t.Errorf("File: %v", err)
		}
		if !reflect.DeepEqual(list, test.want) {
			t.Errorf("File: want %v, got %v", test.want, list)
		}
	}
}

func testReadDir(dir string) ([][]rune, error) {
	m := map[string][][]rune{
		"a": {[]rune("x1"), []rune("x2"), []rune("y1")},
		"b": {[]rune("x3"), []rune("x4"), []rune("y2")},
	}
	return m[dir], nil
}

func TestFromPath(t *testing.T) {
	tests := []struct {
		input string
		want  [][]rune
	}{
		{
			input: "x",
			want:  [][]rune{[]rune("1"), []rune("2"), []rune("3"), []rune("4")},
		},
		{
			input: "y",
			want:  [][]rune{[]rune("1"), []rune("2")},
		},
		{
			input: "z",
			want:  [][]rune{},
		},
	}
	for _, test := range tests {
		list, err := fromPath([]rune(test.input), len(test.input), "a:b", testReadDir)
		if err != nil {
			t.Errorf("fromPath: %v", err)
		}
		if !reflect.DeepEqual(list, test.want) {
			t.Errorf("fromPath: want %v, got %v", test.want, list)
		}
	}
}
