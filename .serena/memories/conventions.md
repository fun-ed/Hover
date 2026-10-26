# Conventions

- Comments explain why/constraints, in short plain sentences; many reference the old C#/Rust names (e.g. "KiroPage.SaveImages"). Keep that style.
- Line endings per `.gitattributes`: Go, Swift, shell, Python, Makefile, VERSION LF; most other text (README, CHANGELOG, web/office/*.js/.html, docs) CRLF. Edit CRLF files with a newline-preserving tool and check for bare LF afterwards.
- OS-specific Go goes in `_windows.go`/`_linux.go`/build-tagged files; Mac-only agents are added to `core.AllTools` in a darwin-only `init()`.
- Per-agent behaviour is `switch t { case core.X: ... }` spread over `internal/agents/{agents,acp,sandbox,setup,runtime}.go`; adding an agent means touching each, plus Swift `Marks.swift`/`Settings.order`/voice strings and web `TOOLS`/`LOGOS`/CSS `.lg.<id>`. Gio pages (`internal/app/pages.go`) have one page per tool and none for Mac-only ones.
- Backend↔Swift messages are `{type: ...}` JSON lines; add fields, don't change existing ones (office.go note). Swift forwards unknown office messages to the backend unchanged.
- Settings model on Mac: `SettingsModel.request` for saveSettings round-trips; independent pages use their own message type (`computerUse`, `spaces`, `storage`).
- JSON in Go is the hand-written `core.JSON` (`core.JObj`, `core.P`, `JStr`, `JInt`, `JArr`), not encoding/json.
- Tests: plain `testing`, names are sentences (`TestXDoesY`); golden fixture `tests/golden/fixtures/office-state.json` lists exactly the six original tools (tests trim `core.AllTools[:6]`).
- CHANGELOG entries are prose bullets saying what changed and what was not checked.
