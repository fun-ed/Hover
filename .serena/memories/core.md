# Hover — core map

Desktop app that drives local coding agents (Kiro, Codex, Cursor, OpenCode, Claude Code, Antigravity; Mac-only in this fork: Oh My Pi ×3, Devin) from a notch/office UI. One Go module `github.com/4regab/Hover` at repo root.

- Win/Linux: one Go program (`cmd/hover`), Gio UI (`internal/ui`, `internal/app`, `internal/shell`).
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

See `mem:tech_stack`, `mem:suggested_commands`, `mem:conventions`, `mem:task_completion`.
