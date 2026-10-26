# Tech stack

- Go 1.27 (`go.mod` min 1.27.0; `mise.toml` pins go 1.27.0, node 22). Run Go/node via `mise exec --`.
- Gio is vendored-replaced: `gioui.org => ./third_party/gioui.org`.
- macOS: Swift compiled directly with `xcrun swiftc` (no Xcode project), Swift 5 mode; deployment target macOS 13.0 in this fork (was 14). `guardian.c` via clang. Office page: three.js bundled by `web/office/build.mjs` (esbuild) into `kiro-office.html`.
- Win: no cgo; office draws via wgpu-native DLL. Linux: cgo + EGL headers, tags `nowayland,nox11,novulkan`.
- Go 1.27 binaries are minos 13.0 on darwin, so macOS 13 is the floor without downgrading Go.
- Agent transports: ACP over stdio (Kiro, Codex via `codex-acp`, Cursor, Antigravity, `omp [--profile=X] acp`, `devin acp`), OpenCode HTTP server (v1 API only; OpenCode v2 is incompatible), Claude Code stream-json + control protocol.
