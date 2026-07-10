# AGENTS.md

## Git workflow: read-only

Agents working in this repo must treat `git` as **read-only** — use it freely
to inspect state (`status`, `diff`, `log`, `show`, etc.), but never to mutate
it. That means no `git add`, `git rm`, `git commit`, `git push`, `git
checkout -b`, staging, or committing on the owner's behalf, even when asked
to "clean up" or "finish" something — do the work as plain filesystem edits
(create/edit/delete files directly) and leave the working tree unstaged. The
owner reviews the diff and merges/commits locally themselves; that's a
deliberate workflow choice (mirrors a no-agentic-git policy from their day
job) so that a human stays accountable for every commit that lands. This
applies across this project, not just to any one task.

Guidance for AI coding agents working in this repository. This documents
**current state as of 2026-07-10**, not the target state — the project is
about to go through an overhaul, so treat everything here as "what exists
today," which may be rewritten or removed.

## What this project is

`cycles` is a small desktop CPU/memory monitor written in Go, using the
[Fyne](https://fyne.io/) GUI toolkit. Single Go module, single `main`
package, no subpackages. It's a personal hobby project ("threw together
while debugging my computer") — expect inconsistent polish and some
aspirational documentation (see "Docs vs. reality" below).

- Module: `cycles` (`go.mod`), Go 1.21.5 toolchain, built with `go1.22.2` locally.
- Linux-only today (reads `/proc/cpuinfo`, `/proc/meminfo` directly).
- Dependencies are vendored in `vendor/` (Go will use vendor mode automatically
  since `vendor/modules.txt` is present and `go.mod` declares `go 1.21.5`).
- Distributed as: plain binary, AppImage, Snap. Flatpak was attempted and
  abandoned (never got it working) — it's explicitly deprioritized now.
  **Snap Store is the priority distribution channel going forward**; Flatpak
  is a nice-to-have at best. An orphaned, never-wired-up Flatpak manifest
  still sits at repo root (`us.tylerc.cycles.yaml`, `us.tylerc.cycles.desktop`,
  `us.tylerc.cycles.appdata.xml`) — parked there deliberately rather than
  deleted, in case Flatpak gets revisited later.

## ⚠️ Current build is broken

`go build ./...` and `go vet ./...` fail right now on `main`:

```
./settingsui.go:150:87: undefined: fyne.Color
./settings.go:98:15: undefined: fyne.VariantLight
./settings.go:100:15: undefined: fyne.VariantDark
./settings.go:102:15: undefined: fyne.VariantDark
```

Root cause: `settings.go` and `settingsui.go` reference symbols that don't
exist in vendored Fyne v2.4.2:
- `fyne.Color` isn't a type in this Fyne version — it should be `color.Color`
  from `image/color` (see how `theme.go` does it correctly).
- `fyne.VariantLight` / `fyne.VariantDark` don't live in the `fyne` package —
  they're `theme.VariantLight` / `theme.VariantDark` in `fyne.io/fyne/v2/theme`
  (confirmed via `grep -rl VariantDark vendor/fyne.io/`).

This was introduced along with the settings/theme system (`git log
settings.go settingsui.go` for history) and appears to have never been
verified with a clean build. **This is Phase 0, item 1 of the overhaul
roadmap below — fix it before anything else.** Any agent asked to "add a
feature" here should not assume `go build` currently succeeds.

## Architecture (current)

Flat file-per-concern layout, no internal packages:

| File | Responsibility |
|---|---|
| `main.go` | Entry point: config/settings load, window+menu setup, tile grids, update goroutines |
| `config.go` | `AppConfig` struct, defaults, CLI flag parsing (`--columns`, `--interval`, `--history`, `--logical`) |
| `settings.go` | `Settings` struct — persisted via Fyne `Preferences` (theme, grid columns, history size, logical cores, update interval). **Currently broken**, see above |
| `settingsui.go` | Settings dialog (Fyne form) + `CustomTheme` wrapper type. **Currently broken**, see above |
| `theme.go` | Graph line colors per utilization band (green/yellow/red) x light/dark, `ApplyTheme()` |
| `tile.go` | `CoreTile` — one CPU core's UI tile (labels + history graph) |
| `memorytile.go` | `MemoryTile` — memory UI tile + byte-size/percent formatters |
| `sysinfo.go` | Reads `/proc/cpuinfo` and `/proc/meminfo` directly (not via gopsutil, deliberately — see DEVELOPER_GUIDE.md "Architecture Decisions"); `UpdateCPUInfo` / `UpdateMemoryInfo` mutate tile slices in place |
| `graphics.go` | `DrawGraph` (line chart into an `image.RGBA`), Bresenham line drawing, label formatters |
| `info.go` | About-dialog text: app version + OS/arch/Go version |
| `*_test.go` | One test file per corresponding source file |

Runtime shape: `main()` builds a config, loads/merges `Settings` from Fyne
preferences (CLI flags override saved settings), builds two Fyne tab pages
(CPU grid, Memory grid), and spins up two independent polling goroutines
(`time.Sleep(config.UpdateInterval)` loops) that mutate tile widgets from a
background goroutine and call `.SetText()` / redraw directly — **no
synchronization/locking around cross-goroutine UI mutation**, relies on
Fyne's internal thread-safety of widget updates.

`gopsutil` (`github.com/shirou/gopsutil`, the v1-style import path, not
`v3`/`v4` module path despite `go.mod` saying `v3.21.11+incompatible`) is
only used for `cpu.Counts()` and `cpu.Percent()`; everything else in
`sysinfo.go` hand-parses `/proc` files.

## Docs vs. reality — read with skepticism

As of 2026-07-10 the four prior-agent planning/retrospective docs that used
to live at repo root (`OVERHAUL_PLAN_V2.md`, `OVERHAUL_SUMMARY.md`,
`SPRINT_SUMMARY.md`, `PACKAGING_PLAN.md`) were deleted — they described a
10-sprint plan targeting v1.2.0 of which only Sprints 1–2 were ever actually
implemented (memory tab, settings system), plus retrospectives and a
packaging roadmap that had gone stale. Don't recreate docs like these
unless the user explicitly asks for a persistent plan file — prefer
conversation-scoped planning (see repo's general "don't create planning
docs unless asked" convention).

Remaining docs, with caveats:
- `README.md` — corrected 2026-07-10 to mention the Memory tab, Settings
  dialog, and theme toggle, and to drop the old "FlatPak coming in 0.5"
  promise. Still contains an intro line about the author being disengaged
  from the project ("not dedicating a lot of time... since I'm not on a
  machine with Snap access anymore") that's now inaccurate given the
  overhaul — left as-is since it's the owner's voice, not a factual claim.
- `DEVELOPER_GUIDE.md` is the most accurate of the remaining docs and
  matches the code layout well, but was written before the settings/theme
  build breakage below — its "Testing Checklist" assumes a clean build.
- `CHANGELOG.md` is real version history, trust it for what shipped when.

When in doubt, **trust the code and `git log`, not these docs.**

## Build / test / run

```bash
make setup      # OS-detecting dependency installer (scripts/setup-dev.sh)
make build      # -> build/cycles
make run        # build + run
make dev        # unoptimized quick build
make test       # go test -v ./...
make check      # fmt + vet + test  (currently fails: vet error above)
```

Native Go equivalents also work (`go build .`, `go test ./...`). System
package dependencies (X11/GL dev headers) are required to build at all — see
`scripts/setup-dev.sh` or the README's "Manual Setup" section for the
per-distro package list.

Test coverage exists for `config.go`, `graphics.go`, `info.go`,
`memorytile.go`, `settings.go` (one `_test.go` each) — but note `go test`
will hit the same compile error as `go build` while `settings.go`/
`settingsui.go` are broken, so **no tests currently run at all**, not even
the ones unrelated to settings.

## Conventions worth preserving

- One file per concern, table-driven tests, exported `PascalCase` /
  unexported `camelCase`, doc comments on exported functions — see
  `DEVELOPER_GUIDE.md` "Code Style Guidelines" for the fuller version the
  project has been trying to follow.
- Version string lives in `config.go` (`DefaultConfig().Version`) and is
  duplicated by hand into `CHANGELOG.md` and CI workflow filenames
  (`appimage.yml` hardcodes `cycles-0.4.1-x86_64.AppImage` — already stale
  against the 0.6.0 version in `config.go`, another drift example).
- `vendor/` is intentionally committed (reproducible builds, offline CI) —
  don't `.gitignore` it or assume `go mod vendor` is optional when adding deps.

## Overhaul roadmap (agreed 2026-07-10)

The stated goal is to get the project working and clean again before
attempting the larger UI overhaul. Roughly in order:

**Phase 0 — get it building again**
1. Fix the `fyne.Color` / `fyne.VariantLight` / `fyne.VariantDark` build break.
2. Get `make check` green; confirm the existing test suite actually passes
   (it's never run successfully — see build break above).
3. Manually smoke-test the running app (tiles render, settings dialog opens
   and persists, theme toggle works) — unverified in a long time.
4. Resolve version-string drift: `config.go` says 0.6.0, `snap/snapcraft.yaml`
   says 0.4.1, `.github/workflows/appimage.yml` hardcodes `0.4.1` into the
   output filename.

**Phase 1 — surface cleanup** (docs done 2026-07-10; rest pending)
5. ~~Purge stale planning docs.~~ Done.
6. ~~Reconcile README against real features.~~ Done.
7. Orphaned Flatpak manifest at repo root — currently just parked, no action
   taken yet.

**Phase 2 — harden Snap (the priority channel over Flatpak)**
8. `snap/snapcraft.yaml` is on `base: core20`; consider moving to a current base.
9. No CI builds/tests the Snap package today — only AppImage + repo-archive
   workflows exist under `.github/workflows/`.
10. Verify whether AppImage is actually broken as older docs claimed, rather
    than assuming — there is a working-looking GH Action for it already.

**Phase 3 — UI overhaul**
11. Current layout is a static Fyne grid that doesn't scale or let users
    rearrange tiles — this is the main motivation for the overhaul. Open
    decision: rebuild the layout within Fyne (responsive/reflowing
    container) vs. reassess whether Fyne is still the right toolkit at all.
    Don't assume either answer — this needs a deliberate call before code.
12. Define what "good" looks like (resizable/reflowing grid? drag-to-rearrange?
    separate panels per metric type?) before implementation.
