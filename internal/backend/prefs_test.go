package backend

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/4regab/Hover/internal/core"
)

func TestDisabledToolsSaveKeepOneOnAndFallBack(t *testing.T) {
	file := filepath.Join(t.TempDir(), "agents.json")
	d := LoadDisabledTools(file)
	if d.Off(core.Kiro) || d.FirstOn(core.Kiro) != core.Kiro {
		t.Fatal("none is off with no file")
	}
	d.Set(core.Kiro, true)
	d.Set(core.OpenCode, true)
	if again := LoadDisabledTools(file); !again.Off(core.Kiro) || !again.Off(core.OpenCode) || again.Off(core.Codex) {
		t.Fatal("not read back")
	}
	if d.FirstOn(core.Kiro) != core.Codex {
		t.Fatalf("fallback %v", d.FirstOn(core.Kiro))
	}
	// Every other one off: the last stays on.
	for _, x := range core.AllTools {
		if x != core.Claude {
			d.Set(x, true)
		}
	}
	if d.Set(core.Claude, false) != true || d.Set(core.Claude, true) || d.Off(core.Claude) {
		t.Fatal("the last one on was switched off")
	}
	d.Set(core.Kiro, false)
	if LoadDisabledTools(file).Off(core.Kiro) {
		t.Fatal("switching on was not saved")
	}
	// A hand-written file with every agent off reads as all on.
	var all []core.JSON
	for _, x := range core.AllTools {
		all = append(all, core.JStr(x.ID()))
	}
	os.WriteFile(file, []byte(core.JObj(core.P("Disabled", core.JArr(all...))).Compact()), 0o644)
	if LoadDisabledTools(file).Off(core.Codex) {
		t.Fatal("all off was kept")
	}
}

func TestStaleCheckpointsSkipOpenAndRecentChats(t *testing.T) {
	now := time.Now()
	cut := now.AddDate(0, 0, -7)
	updated := map[string]time.Time{"old": now.AddDate(0, 0, -30), "recent": now.AddDate(0, 0, -1), "busy": now.AddDate(0, 0, -30)}
	got := staleCheckpoints([]string{"old.git", "recent.git", "busy.git", "orphan.git", "notes"}, map[string]bool{"busy": true}, updated, cut)
	if !slices.Equal(got, []string{"old", "orphan"}) {
		t.Fatal(got)
	}
}

func TestKeepDaysTakesOnlyTheChoicesAndIsSaved(t *testing.T) {
	defer core.SetKeepDays(14)
	file := filepath.Join(t.TempDir(), "storage.json")
	SaveKeepDays(file, 30)
	core.SetKeepDays(7)
	LoadKeepDays(file)
	if core.KeepDays() != 30 {
		t.Fatal(core.KeepDays())
	}
	core.SetKeepDays(3)
	if core.KeepDays() != 30 {
		t.Fatal("3 days was taken")
	}
}
