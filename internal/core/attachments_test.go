package core

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSaveFilesKeepsNamesAndHoldsTheLimits(t *testing.T) {
	dir := t.TempDir()
	file := func(name string, n int) JSON {
		data := "data:application/octet-stream;base64," + base64.StdEncoding.EncodeToString(make([]byte, n))
		return JObj(P("name", JStr(name)), P("data", JStr(data)))
	}
	saved := SaveFiles([]JSON{file(`..\evil/notes.txt`, 3), JStr("not a file"), file("big.bin", MaxFileBytes+1)}, dir)
	if len(saved) != 1 || filepath.Dir(saved[0]) != dir || !strings.HasSuffix(saved[0], "-notes.txt") {
		t.Fatal(saved)
	}
	if b, _ := os.ReadFile(saved[0]); len(b) != 3 {
		t.Fatal("contents")
	}
	// Three of 20 MiB pass one by one, but only two fit in the 50 MiB of a message.
	if got := SaveFiles([]JSON{file("a", MaxFileBytes), file("b", MaxFileBytes), file("c", MaxFileBytes)}, dir); len(got) != 2 {
		t.Fatal(len(got))
	}
	var many []JSON
	for range MaxFiles + 2 {
		many = append(many, file("x.txt", 1))
	}
	if got := SaveFiles(many, dir); len(got) != MaxFiles {
		t.Fatal(len(got))
	}
}

func TestSweepOlderCountsWhatGoesAndUsageWhatStays(t *testing.T) {
	dir := t.TempDir()
	old, fresh := filepath.Join(dir, "old.png"), filepath.Join(dir, "fresh.txt")
	os.WriteFile(old, make([]byte, 10), 0o644)
	os.WriteFile(fresh, make([]byte, 3), 0o644)
	past := time.Now().AddDate(0, 0, -20)
	os.Chtimes(old, past, past)
	if n, size := SweepOlder(dir, time.Now().AddDate(0, 0, -7)); n != 1 || size != 10 {
		t.Fatal(n, size)
	}
	if size, n := Usage(dir); size != 3 || n != 1 {
		t.Fatal(size, n)
	}
}
