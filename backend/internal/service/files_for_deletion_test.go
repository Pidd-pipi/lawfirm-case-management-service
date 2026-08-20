package service

import (
	"reflect"
	"testing"
)

func TestFilesForDeletionStripsPrefix(t *testing.T) {
	got := filesForDeletion([]string{"/uploads/a.pdf", "b.pdf"})
	want := []string{"a.pdf"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("filesForDeletion = %v, want %v", got, want)
	}
}
