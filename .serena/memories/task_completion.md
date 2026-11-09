# When a task is done

1. `gofmt -l internal cmd` → empty (gofmt may realign struct literals; check the diff).
2. Run `mise exec -- go vet ./cmd/... ./internal/...` and `mise exec -- go test ./cmd/... ./internal/...` independently from diff checks, so an earlier failure cannot skip them.
   - Known local-only failure: `TestNotInstalledIsNoServerEvenWhenOn` fails when Cua Driver is installed at `/Applications/CuaDriver.app`.
   - `TestPathsWithSpacesUnicodeAndShellCharactersArriveAsTheyAre` can flake under full parallel load; rerun alone.
   - Several agents tests leave `hover-handoff-*`, `hover-kagents-*`, `hover-queue-*`, `hover-test-data-*`, `hover-usable-*` dirs in `$TMPDIR`.
3. JS changes: `mise exec -- node --check web/office/main.js`; beware TDZ — top-level `const` used by render code must not be read before its declaration runs.
4. Swift changes: macOS 13 typecheck (see `mem:suggested_commands`), then `scripts/build-macos.sh` for arm64 and x64.
5. CRLF files: no bare LF. Update README and CHANGELOG `[Unreleased]`/version section for user-visible changes.
6. Real-agent and UI behaviour needs a manual run; say what was not run.
7. Retain the test suites and build entry points during upstream syncs. Upstream `0179820` removed them; this fork deliberately keeps them.
8. Packaged backend smoke needs `HOVER_SANDBOX_ROOT`, `HOVER_SANDBOXED=1`, an isolated `TMPDIR` and fake-bin. Keep browser enabled, because the test expects `hover-browser`. In the sandbox's `agents.json`, disable every agent except the fake Codex before `ready`, avoiding discovery of installed personal agents.
9. Release 5.2.1: upstream `76f81a2`; full Go suite had only the documented Cua Driver failure. Full Go vet, JS syntax, unresolved-conflict and attachment/chip regression checks passed. Re-run packaged smoke on the final versioned arm64 app; Intel is cross-built, not runtime-tested.
