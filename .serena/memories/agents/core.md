# internal/agents

- `agents.go`: per-tool `Exe`, `Arguments`, `InstallHint`, `SignInHint`, readiness `look()` (status command; omp has none → ready if installed and `~/.omp/profiles/<p>` exists; Devin `devin auth status`).
- `acp.go`: ACP host. `configure` maps Hover access (Full/Ask first/Ask always/Read only) onto each agent's `mode` option per tool; model via category `model`, effort via `thought_level`. `promptBlocks`: `Attached image (read it from this file): ` lines → image blocks only for `sendsImages(t)` agents that advertise `promptCapabilities.image`; `Attached file (read it from this file): ` lines stay and add a `resource` (text ≤256 KB when `embeddedContext`) or `resource_link` (ACP baseline).
- `claude.go`: Claude Code via stream-json; `claudeContent` turns attachments into Anthropic image/document blocks; Default model/effort = user's `~/.claude/settings.json` (Hover passes no flags).
- `sandbox.go`: srt sandbox domains (`ToolDomains`) and writable state (`ToolState`) per tool; reads of `kiro-images`/`attachments` allowed.
- `setup.go`: one-click install plans and sign-in commands per tool.
- Backend side (`internal/backend`): `DisabledTools` (agents.json) hides agents from office/voice, skips checks/discovery, refuses `new`; `storage` message drives Settings → Storage (KeepDays 7/14/30, sweeps every 6 h, clean attachments/checkpoints/log).
