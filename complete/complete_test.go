package complete

import (
	"reflect"
	"testing"
)

func TestFile(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			input: "testdata/file",
			want:  []string{"1.txt", "2.txt"},
		},
		{
			input: "testdata/",
			want:  []string{"dir1/", "file1.txt", "file2.txt"},
		},
	}
	for _, test := range tests {
		list, err := File([]rune(test.input), len(test.input))
		if err != nil {
			t.Errorf("File: %v", err)
		}
		if !reflect.DeepEqual(list, test.want) {
			t.Errorf("File: want %v, got %v", test.want, list)
		}
	}
}

func testReadDir(dir string) ([]string, error) {
	m := map[string][]string{
		"a": []string{"x1", "x2", "y1"},
		"b": []string{"x3", "x4", "y2"},
	}
	return m[dir], nil
}

func TestFromPath(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			input: "x",
			want:  []string{"1", "2", "3", "4"},
		},
		{
			input: "y",
			want:  []string{"1", "2"},
		},
		{
			input: "z",
			want:  []string{},
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
