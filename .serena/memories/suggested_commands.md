# Commands (Darwin)

- Go tests: `mise exec -- go test ./cmd/... ./internal/...` (never `./...`: web build leaves non-Go dirs). Linux needs `-tags nowayland,nox11,novulkan`.
- Mac app: `mise exec -- scripts/build-macos.sh` → `dist/macos-osx-arm64/Hover.app`; `HOVER_ARCH=x64` for Intel; `HOVER_SETTINGS_TESTS=1 HOVER_APP_OUTPUT=/abs/Hover.app` builds the Settings smoke variant. Script runs `npm ci` for the office.
- Swift-only check (fast): copy `macos/Sources/Hover.swift` to `main.swift` in a temp dir, then `xcrun swiftc -swift-version 5 -typecheck -target arm64-apple-macos13.0 <main.swift + other Sources/*.swift>` (compiling Hover.swift under its own name gives bogus top-level errors).
- DMG like CI: stage `Hover.app` + `Applications` symlink, `hdiutil create -volname Hover -srcfolder <stage> -format UDZO`, then `hdiutil verify`.
- Check binary floor: `vtool -show-build <bin> | grep minos`.
- Probe an ACP agent without a model call: pipe `initialize` then `session/new` JSON-RPC lines into `<agent> acp` and read `promptCapabilities`/`configOptions`.
- Running Hover data: `~/Library/Application Support/Hover` (`hover.log`). No `timeout` binary on this Mac.
