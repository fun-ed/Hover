# When a task is done

1. `gofmt -l internal cmd` → empty (gofmt may realign struct literals; check the diff).
2. `mise exec -- go vet` on touched packages; `mise exec -- go test ./cmd/... ./internal/...`.
   - Known local-only failure: `TestNotInstalledIsNoServerEvenWhenOn` fails when Cua Driver is installed at `/Applications/CuaDriver.app`.
   - `TestPathsWithSpacesUnicodeAndShellCharactersArriveAsTheyAre` can flake under full parallel load; rerun alone.
   - Several agents tests leave `hover-handoff-*`, `hover-kagents-*`, `hover-queue-*`, `hover-test-data-*`, `hover-usable-*` dirs in `$TMPDIR`.
3. JS changes: `mise exec -- node --check web/office/main.js`; beware TDZ — top-level `const` used by render code must not be read before its declaration runs.
4. Swift changes: macOS 13 typecheck (see `mem:suggested_commands`), then `scripts/build-macos.sh` for arm64 and x64.
5. CRLF files: no bare LF. Update README and CHANGELOG `[Unreleased]`/version section for user-visible changes.
6. Real-agent and UI behaviour needs a manual run; say what was not run.
