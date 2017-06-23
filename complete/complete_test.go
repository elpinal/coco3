package complete

import (
	"reflect"
	"testing"
)

func TestFile(t *testing.T) {
	list, err := File([]rune("testdata/file"), 0)
	if err != nil {
		t.Errorf("File: %v", err)
	}
	if want := []string{"1.txt", "2.txt"}; !reflect.DeepEqual(list, want) {
		t.Errorf("File: want %v, got %v", want, list)
	}
}
