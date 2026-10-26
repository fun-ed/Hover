package agents

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/4regab/Hover/internal/core"
)

func TestVersionsCompareAsDotnetDoes(t *testing.T) {
	v := func(s string) [4]int64 {
		p, ok := ParseVersion(s)
		if !ok {
			t.Fatalf("%q doesn't parse", s)
		}
		return p
	}
	if !versionLess(v(OpenCodeMinVersion), v("1.18.31")) || !versionLess(v("1.14.2"), v("1.14.19")) || !versionLess(v("1.14"), v("1.14.0")) {
		t.Error("order")
	}
	for _, s := range []string{"1", "v1.2.3", "1.2.3.4.5"} {
		if _, ok := ParseVersion(s); ok {
			t.Errorf("%q parses", s)
		}
	}
}

func TestOhMyPiStartsOnItsProfileAndNeedsThatProfile(t *testing.T) {
	if got := Arguments(core.Omp); !slices.Equal(got, []string{"acp"}) {
		t.Error(got)
	}
	if got := Arguments(core.OmpOpus); !slices.Equal(got, []string{"--profile=opus", "acp"}) {
		t.Error(got)
	}
	// Devin runs its local agent, never the --cloud relay.
	if got := Arguments(core.Devin); !slices.Equal(got, []string{"acp"}) {
		t.Error(got)
	}
	home, bin := t.TempDir(), t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PATH", bin)
	if r, _ := look(core.OmpSol); r.Installed {
		t.Fatal("installed without omp")
	}
	os.WriteFile(filepath.Join(bin, "omp"), nil, 0o777)
	if r, _ := look(core.Omp); !r.OK() {
		t.Error("the default profile needs no folder", r)
	}
	if r, _ := look(core.OmpSol); !r.Installed || r.SignedIn || r.Hint == "" {
		t.Error("a profile never made counts as signed in", r)
	}
	os.MkdirAll(filepath.Join(home, ".omp", "profiles", "sol"), 0o755)
	if r, _ := look(core.OmpSol); !r.OK() {
		t.Error("the sol profile is there", r)
	}
}

func TestOhMyPiAndDevinGetPicturesAsImageBlocks(t *testing.T) {
	pic := filepath.Join(t.TempDir(), "shot.png")
	os.WriteFile(pic, []byte("\x89PNG\r\n\x1a\n"), 0o644)
	prompt := "What is wrong here?\n" + Attached + pic
	for _, tool := range []core.AgentTool{core.Kiro, core.Omp, core.OmpSol, core.Devin} {
		got := promptBlocks(prompt, sendsImages(tool), false).Compact()
		if !strings.Contains(got, `"type":"image"`) || strings.Contains(got, Attached) {
			t.Error(tool, got)
		}
	}
	if got := promptBlocks(prompt, sendsImages(core.Cursor), false).Compact(); strings.Contains(got, `"type":"image"`) {
		t.Error("cursor", got)
	}
	if !sendsImages(core.Codex) {
		t.Error("codex says it takes images")
	}
}

func TestAttachedFilesGoAsResourcesAndKeepTheirLine(t *testing.T) {
	dir := t.TempDir()
	notes := filepath.Join(dir, "20261026T160000-12345678-notes.md")
	os.WriteFile(notes, []byte("# Plan\nship it"), 0o644)
	blob := filepath.Join(dir, "20261026T160000-12345678-report.pdf")
	os.WriteFile(blob, []byte("%PDF\x00\x01"), 0o644)
	prompt := "Read these.\n" + AttachedFile + notes + "\n" + AttachedFile + blob
	got := promptBlocks(prompt, false, true).Compact()
	for _, want := range []string{`"type":"resource"`, `"text":"# Plan\nship it"`, `"type":"resource_link"`, `"name":"report.pdf"`, `"mimeType":"application/pdf"`, AttachedFile} {
		if !strings.Contains(got, want) {
			t.Error(want, got)
		}
	}
	// Without embedded context both are links.
	if got := promptBlocks(prompt, false, false).Compact(); strings.Contains(got, `"type":"resource"`) || strings.Count(got, `"type":"resource_link"`) != 2 {
		t.Error(got)
	}
}

func TestClaudeGetsPicturesAndDocumentsAsItsOwnBlocks(t *testing.T) {
	dir := t.TempDir()
	pic := filepath.Join(dir, "shot.png")
	os.WriteFile(pic, []byte("\x89PNG\r\n\x1a\n"), 0o644)
	notes := filepath.Join(dir, "20261026T160000-12345678-notes.md")
	os.WriteFile(notes, []byte("ship it"), 0o644)
	pdf := filepath.Join(dir, "20261026T160000-12345678-report.pdf")
	os.WriteFile(pdf, []byte("%PDF-1.7\x00"), 0o644)
	zip := filepath.Join(dir, "20261026T160000-12345678-logs.zip")
	os.WriteFile(zip, []byte("PK\x03\x04\x00"), 0o644)
	prompt := "Look.\n" + Attached + pic + "\n" + AttachedFile + notes + "\n" + AttachedFile + pdf + "\n" + AttachedFile + zip
	got := claudeContent(prompt).Compact()
	for _, want := range []string{`"type":"image","source":{"type":"base64","media_type":"image/png"`,
		`"source":{"type":"text","media_type":"text/plain","data":"ship it"},"title":"notes.md"`,
		`"media_type":"application/pdf"`, `"title":"report.pdf"`, AttachedFile + zip} {
		if !strings.Contains(got, want) {
			t.Error(want, got)
		}
	}
	if strings.Contains(got, Attached+pic) || strings.Count(got, `"type":"document"`) != 2 {
		t.Error(got)
	}
	if got := claudeContent("plain").Compact(); got != `[{"type":"text","text":"plain"}]` {
		t.Error(got)
	}
}
