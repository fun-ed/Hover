# macOS app

- `macos/Sources/Hover.swift` is the entry (compiled as `main.swift`); `OfficeHost.swift` hosts the WKWebView office and runs `hover-guardian` → `hover-backend`; `Settings.swift` is the SwiftUI Settings window (pages enum `SettingsModel.Page`); `Notch.swift` island drawing/animation; `MenuBar.swift`; `Voice*.swift`; `Marks.swift` logos (WPF path strings, `F0` even-odd / `F1` nonzero, 24-unit box).
- `ShellEnvironment.swift` merges login-shell PATH so GUI launches find mise/bun/brew tools.
- Backend messages dispatched in `Hover.swift` `switch m["type"]`; office page messages default to `backend.send(m)`.
- Office page (`web/office/main.js`, `page.html`): render loop capped at 60 fps on Mac (`RETINA` = host present), 10 fps idle, paused when hidden; agent picker shows `tools` filtered by `enabled`; attachments: images as data URLs (≤4), other files `{name,size,data}` (≤10, 20 MB each, 50 MB total).
- Backend host line limit 128 MiB (`internal/backend/wire.go`), sized for attachments.
- APIs newer than macOS 13 are wrapped in `#available(macOS 14.0, *)` (display link, inactiveSchedulingPolicy, sectionHeader, SCScreenshotManager — Screen tab blank on 13).
- Builds are ad-hoc signed, not notarized; bundle id `dev.hover.desktop`.
