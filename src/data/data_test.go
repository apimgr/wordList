package data

import "testing"

func TestReadFile(t *testing.T) {
	entries, err := EmbeddedData.ReadDir(".")
	if err != nil {
		t.Fatalf("ReadDir(.) error = %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one embedded data file")
	}

	if _, err := ReadFile(entries[0].Name()); err != nil {
		t.Fatalf("ReadFile(%q) error = %v", entries[0].Name(), err)
	}

	if _, err := ReadFile("does-not-exist.json"); err == nil {
		t.Error("expected error for missing file")
	}
}
