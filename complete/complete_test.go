package complete

import (
	"reflect"
	"testing"
)

func TestFile(t *testing.T) {
	list, err := File([]rune("testdata/"), 0)
	if err != nil {
		t.Errorf("File: %v", err)
	}
	if want := []string{"file1.txt", "file2.txt"}; !reflect.DeepEqual(list, want) {
		t.Errorf("File: want %v, got %v", want, list)
	}
}
