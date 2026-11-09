# Hover — core map

Desktop app that drives local coding agents (Kiro, Codex, Cursor, OpenCode, Claude Code, Antigravity; Mac-only in this fork: Oh My Pi ×3, Devin) from a notch/office UI. One Go module `github.com/4regab/Hover` at repo root.

- Windows: C# / WinUI 3 shell (`windows/`) on the shared Go backend; `build.ps1 legacy` keeps the Go Windows app. Linux uses the Go / Gio app (`cmd/hover`).
- macOS: Swift shell (`macos/Sources`) + web office page (`web/office`) + Go backend (`cmd/hover-backend`, `internal/backend`) speaking JSON lines over stdio. Details: `mem:macos/core`.
- Agent protocols, sessions, permissions, sandbox, setup: `internal/agents` (one big package on purpose). Details: `mem:agents/core`.
- Settings, history (encrypted), paths, attachments, JSON: `internal/core`; `core` imports nothing else from the module.

Invariants
- `core.AgentTool` enum order is persisted (settings/history store names/ids): only append at the end. `Custom` is a removed slot kept for old chats.
- `settings.json` must keep legacy bytes: new Mac-only prefs go in their own files beside it (`backend.json` MaxRunning, `agents.json` disabled agents, `storage.json` KeepDays).
- History is encrypted; never replace an unreadable key/index.
- Keep prompts off command lines; agent processes must stop with Hover.
- Feature support differs per OS; explain unsupported features rather than implying they work.
- Upstream reference docs: `AGENTS.md`, `docs/development/architecture.md`, `docs/development/testing.md`, `docs/MACOS.md`.
- Version 5.2.1 integrates upstream `76f81a2` (206 commits). Upstream `0179820` removed automated tests; this fork retains them and their build entry points. GitHub Actions remain disabled.
- Backend replies now run through `tasks.go` (`reply`); keep `withFiles(text, saveFiles(m))` there so Mac attachments survive alongside upstream reply chips. `TestReplyTakesChipsAndChecksThem` checks both.

See `mem:tech_stack`, `mem:suggested_commands`, `mem:conventions`, `mem:task_completion`.
