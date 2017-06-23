package complete

import (
	"os"
	"reflect"
	"testing"
)

func TestFile(t *testing.T) {
	os.Chdir("testdata")
	list, err := File(nil, 0)
	if err != nil {
		t.Errorf("File: %v", err)
	}
	if want := []string{"file1.txt", "file2.txt"}; !reflect.DeepEqual(list, want) {
		t.Errorf("File: want %v, got %v", want, list)
	}
}
