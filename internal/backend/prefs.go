package backend

import (
	"os"
	"sync"
	"sync/atomic"

	"github.com/4regab/Hover/internal/core"
)

// Settings.MaxRunning: how many tasks run at once (1 to 6, 3 unless changed), which the
// Mac's Settings offers. settings.json has no such key, so it is kept here, in a file of the
// backend's own beside settings.json, and handed to the sessions at start and on every
// change.
const (
	defaultRunning = 3
	minRunning     = 1
	maxRunning     = 6
)

type MaxRunning struct {
	file string
	n    atomic.Int64
}

// LoadMaxRunning reads the file (3 when there is none, or it can't be read).
func LoadMaxRunning(file string) *MaxRunning {
	m := &MaxRunning{file: file}
	n := int64(defaultRunning)
	if b, err := os.ReadFile(file); err == nil {
		if v, err := core.ParseJSON(core.TextOf(b)); err == nil {
			if x, ok := v.Get("MaxRunning"); ok && x.Kind() == core.NumKind {
				if k, err := x.I64(); err == nil {
					n = max(minRunning, min(k, maxRunning))
				}
			}
		}
	}
	m.n.Store(n)
	return m
}

func (m *MaxRunning) Get() int { return int(m.n.Load()) }

// Set is Math.Clamp(value, 1, 6), kept.
func (m *MaxRunning) Set(n int32) {
	k := int64(max(minRunning, min(int(n), maxRunning)))
	if m.n.Swap(k) == k {
		if _, err := os.Stat(m.file); err == nil {
			return
		}
	}
	text := core.JObj(core.P("MaxRunning", core.JInt(k))).Compact()
	tmp := m.file[:len(m.file)-len(".json")] + ".json.tmp"
	if err := os.WriteFile(tmp, []byte(text), 0o644); err == nil {
		err = os.Rename(tmp, m.file)
		if err == nil {
			return
		}
	}
	core.Logf("backend settings save failed")
}

// DisabledTools are the agents switched off in the Mac's Settings: left out of the office's
// and voice's lists, never checked or started. Its own file beside backend.json (agents.json),
// as settings.json keeps its legacy keys. One stays on, so the office always has an agent.
type DisabledTools struct {
	file string
	mu   sync.Mutex
	off  map[core.AgentTool]bool
}

// LoadDisabledTools reads the file (none off when there is none, or it can't be read).
func LoadDisabledTools(file string) *DisabledTools {
	d := &DisabledTools{file: file, off: map[core.AgentTool]bool{}}
	if b, err := os.ReadFile(file); err == nil {
		if v, err := core.ParseJSON(core.TextOf(b)); err == nil {
			if x, ok := v.Get("Disabled"); ok && x.Kind() == core.ArrKind {
				items, _ := x.Items()
				for _, it := range items {
					id, _ := it.AsStr()
					if t, ok := core.ParseTool(&id); ok {
						d.off[t] = true
					}
				}
			}
		}
	}
	// A file that switched every one off (written by hand): all on again.
	if len(d.off) >= len(core.AllTools) {
		d.off = map[core.AgentTool]bool{}
	}
	return d
}

func (d *DisabledTools) Off(t core.AgentTool) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.off[t]
}

// FirstOn is t when it is on, else the first tool that is.
func (d *DisabledTools) FirstOn(t core.AgentTool) core.AgentTool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.off[t] {
		return t
	}
	for _, x := range core.AllTools {
		if !d.off[x] {
			return x
		}
	}
	return t
}

// Set switches t on or off and saves; false when it would switch off the last one on.
func (d *DisabledTools) Set(t core.AgentTool, off bool) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.off[t] == off {
		return true
	}
	if off && len(d.off)+1 >= len(core.AllTools) {
		return false
	}
	if off {
		d.off[t] = true
	} else {
		delete(d.off, t)
	}
	var ids []core.JSON
	for _, x := range core.AllTools {
		if d.off[x] {
			ids = append(ids, core.JStr(x.ID()))
		}
	}
	text := core.JObj(core.P("Disabled", core.JArr(ids...))).Compact()
	tmp := d.file + ".tmp"
	if err := os.WriteFile(tmp, []byte(text), 0o644); err == nil {
		if err = os.Rename(tmp, d.file); err == nil {
			return true
		}
	}
	core.Logf("agents.json save failed")
	return true
}

// Settings → Storage: how many days pasted images and attached files are kept (7, 14 or
// 30). Its own file beside backend.json (storage.json), handed to core at start and on change.
func LoadKeepDays(file string) {
	if b, err := os.ReadFile(file); err == nil {
		if v, err := core.ParseJSON(core.TextOf(b)); err == nil {
			if x, ok := v.Get("KeepDays"); ok && x.Kind() == core.NumKind {
				if k, err := x.I64(); err == nil {
					core.SetKeepDays(int(k))
				}
			}
		}
	}
}

// SaveKeepDays takes one of core.KeepChoices and writes it.
func SaveKeepDays(file string, days int) {
	core.SetKeepDays(days)
	text := core.JObj(core.P("KeepDays", core.JInt(int64(core.KeepDays())))).Compact()
	if err := os.WriteFile(file+".tmp", []byte(text), 0o644); err == nil && os.Rename(file+".tmp", file) == nil {
		return
	}
	core.Logf("storage.json save failed")
}
