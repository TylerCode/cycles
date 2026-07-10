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
**current state as of 2026-07-10** (v0.8.0) — the UI overhaul described in
"Overhaul roadmap" below is now complete (all phases 0–3 done); treat this
as the settled architecture unless a new overhaul gets kicked off.

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

## Build/runtime health (fixed 2026-07-10)

As of this date `go build`, `go vet`, `go test`, and a manual smoke-test run
of the built binary are all clean. This section records what was actually
wrong, in case similar bugs get reintroduced during the overhaul:

1. **Compile error**: `settings.go`/`settingsui.go` referenced `fyne.Color`
   (not a real type — vendored Fyne v2.4.2 uses `color.Color` from
   `image/color`, see `theme.go` for the correct pattern) and
   `fyne.VariantLight`/`fyne.VariantDark` (those constants live in
   `fyne.io/fyne/v2/theme`, not the `fyne` package itself).
2. **Preferences silently didn't persist**: `main.go` called `app.New()`
   instead of `app.NewWithID("us.tylerc.cycles")`. Fyne's Preferences API
   (used by `settings.go` for save/load) requires a unique app ID or it logs
   a warning and settings never actually hit disk — the "settings persist
   between runs" feature was never true at runtime despite compiling and
   passing unit tests. Only caught by actually launching the app.
3. **Infinite recursion → stack overflow on startup**: `CustomTheme` in
   `settingsui.go` (installed via `ApplyTheme()` →
   `app.Settings().SetTheme(&CustomTheme{...})`) delegated its
   `Color`/`Font`/`Icon`/`Size` methods to
   `fyne.CurrentApp().Settings().Theme()` — but once installed, that call
   returns the `CustomTheme` itself, so every color lookup recursed into
   itself forever. Fixed by delegating to the stable `theme.DefaultTheme()`
   base theme instead. This is why a build-clean / vet-clean / tests-passing
   state is not sufficient proof the app works — **always smoke-test the
   actual running binary** (`fyne.io/fyne/v2/test` covers widget
   construction in unit tests, e.g. the `TestMain` in `memorytile_test.go`,
   but not full app wiring like `ApplyTheme`).
4. Version-string drift across `config.go` (was 0.6.0 at the time),
   `snap/snapcraft.yaml` (was 0.4.1), and `.github/workflows/appimage.yml`
   (hardcoded 0.4.1 into the artifact filename) — resolved by pinning
   snapcraft to match `config.go` and making the AppImage workflow derive
   `VERSION` from `config.go` at build time instead of hardcoding it, so the
   AppImage side can't drift again silently. `snap/snapcraft.yaml` still
   needs a manual bump alongside `config.go` on every release — see
   "Conventions worth preserving" below. Current version: 0.8.0.

## Architecture (current, post-overhaul)

Flat file-per-concern layout, no internal packages:

| File | Responsibility |
|---|---|
| `main.go` | Entry point: config/settings load, window+menu setup, CPU tab + Memory dashboard wiring, theme-change listener, update goroutines |
| `config.go` | `AppConfig` struct, defaults, CLI flag parsing (`--columns`, `--interval`, `--history`, `--logical`) |
| `settings.go` | `Settings` struct — persisted via Fyne `Preferences` (theme, grid columns, history size, logical cores, update interval, CPU tab view mode) |
| `settingsui.go` | Settings dialog (Fyne form) + `CustomTheme` wrapper type |
| `theme.go` | `GetGraphLineColor`/`GetSeriesColor` (green/yellow/red status colors plus named chart series colors) x light/dark, `UtilizationStatus()` thresholds, `ApplyTheme()` |
| `cputab.go` | `CPUTab` — owns the whole CPU tab: aggregate stats strip (Threads/Avg util/Peak core/Max clock), Tiles/List view toggle, both view containers |
| `tile.go` | `CoreState` (per-core data, shared by both views) + `CoreTile` (Tiles view card, responsive `GridWrap` cell) + `CoreListRow` (List view row with inline `Bar`) |
| `memorytile.go` | `MemoryDashboard` — the whole Memory tab: radial usage gauge, used/cached/buffers/free/swap breakdown rows, full-width history area chart; byte-size/percent formatters |
| `bar.go` | `Bar` — a `fyne.Layout`-driven horizontal progress bar with a fixed fill color, shared by the CPU list view and the Memory breakdown rows |
| `sysinfo.go` | Reads `/proc/cpuinfo` and `/proc/meminfo` directly (not via gopsutil, deliberately — see DEVELOPER_GUIDE.md "Architecture Decisions"), including swap (`SwapTotal`/`SwapFree`); `UpdateCPUInfo` / `UpdateMemoryInfo` push samples into `CoreState`/`MemoryDashboard` |
| `graphics.go` | `DrawSparkline`, `DrawAreaChart` (memory+swap dual-series), `DrawRadialGauge` — all render into an `image.RGBA` at their actual requested size (no more fixed-bitmap graphs), Bresenham line drawing, label formatters |
| `info.go` | About-dialog text: app version + OS/arch/Go version |
| `*_test.go` | One test file per corresponding source file |

Runtime shape: `main()` builds a config, loads/merges `Settings` from Fyne
preferences (CLI flags override saved settings), builds the CPU tab
(`NewCPUTab`) and Memory dashboard (`NewMemoryDashboard`) as Fyne tab pages,
and spins up two independent polling goroutines (`time.Sleep
(config.UpdateInterval)` loops) that mutate widgets from a background
goroutine and call `.Refresh()` / redraw directly — **no
synchronization/locking around cross-goroutine UI mutation**, relies on
Fyne's internal thread-safety of widget updates. A separate goroutine
listens for app theme changes and pushes `RefreshTheme()` through both tabs,
since `canvas.Text`/`canvas.Rectangle` primitives only read theme colors
once at construction and don't repaint themselves on a live theme switch the
way built-in widgets do.

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
  matches the code layout well; its "Testing Checklist" now reflects
  reality again since the build/runtime fixes above.
- `CHANGELOG.md` is real version history, trust it for what shipped when.

When in doubt, **trust the code and `git log`, not these docs.**

## Build / test / run

```bash
make setup      # OS-detecting dependency installer (scripts/setup-dev.sh)
make build      # -> build/cycles
make run        # build + run
make dev        # unoptimized quick build
make test       # go test -v ./...
make check      # fmt + vet + test  (all green as of 2026-07-10)
```

Native Go equivalents also work (`go build .`, `go test ./...`). System
package dependencies (X11/GL dev headers) are required to build at all — see
`scripts/setup-dev.sh` or the README's "Manual Setup" section for the
per-distro package list.

Test coverage exists for `config.go`, `graphics.go`, `info.go`,
`memorytile.go`, `settings.go`, `sysinfo.go` (one `_test.go` each), all
passing. `memorytile_test.go` has a package-level `TestMain` that boots a
headless `fyne.io/fyne/v2/test` app — needed because widget constructors
that call `theme.BackgroundColor()` (e.g. `NewMemoryDashboard`,
`NewCoreTile`) panic without a `fyne.CurrentApp()` in scope. Any new test
that constructs a Fyne widget touching theme colors relies on that same
`TestMain`, so keep it if you add or move test files. `tile.go`, `bar.go`,
and `cputab.go` (all added in the Phase 3 UI overhaul) currently have **no**
dedicated test files — their logic is mostly Fyne widget wiring rather than
pure functions, but `CoreState.Update`'s history-trimming and
`ComputeCPUAggregateStats` (tested indirectly via `sysinfo_test.go`) would
be reasonable candidates if gaps here start to matter.

## Conventions worth preserving

- One file per concern, table-driven tests, exported `PascalCase` /
  unexported `camelCase`, doc comments on exported functions — see
  `DEVELOPER_GUIDE.md` "Code Style Guidelines" for the fuller version the
  project has been trying to follow.
- Version string lives in `config.go` (`DefaultConfig().Version`) and is
  still duplicated by hand into `CHANGELOG.md` and `snap/snapcraft.yaml`
  (kept in sync manually as of 2026-07-10 — no automation yet, so bumping
  the version means updating both). The AppImage workflow now derives its
  filename from `config.go` automatically, so that one can't drift again.
- `vendor/` is intentionally committed (reproducible builds, offline CI) —
  don't `.gitignore` it or assume `go mod vendor` is optional when adding deps.

## Overhaul roadmap (agreed 2026-07-10)

The stated goal is to get the project working and clean again before
attempting the larger UI overhaul. Roughly in order:

**Phase 0 — get it building again — done 2026-07-10**
1. ~~Fix the `fyne.Color` / `fyne.VariantLight` / `fyne.VariantDark` build break.~~ Done.
2. ~~Get `make check` green; confirm tests pass.~~ Done — also fixed three
   latent test bugs uncovered once compilation succeeded (stale hardcoded
   version/column expectations, a segfault from missing Fyne test-app
   context, and inverted light/dark constant assumptions). See "Build/runtime
   health" above.
3. ~~Manually smoke-test the running app.~~ Done — this caught two runtime
   bugs invisible to build/vet/test: preferences never actually persisting
   (missing app ID) and an infinite-recursion stack overflow in `CustomTheme`
   on startup. Both fixed; see "Build/runtime health" above.
4. ~~Resolve version-string drift.~~ Done — `snap/snapcraft.yaml` bumped to
   0.6.0, `.github/workflows/appimage.yml` now derives its version from
   `config.go` instead of hardcoding it.

**Phase 1 — surface cleanup — done 2026-07-10**
5. ~~Purge stale planning docs.~~ Done.
6. ~~Reconcile README against real features.~~ Done.
7. Orphaned Flatpak manifest at repo root — deliberately deferred by the
   owner ("I'll come back to flatpak eventually"), not a bug, don't touch
   without being asked. (Still true post-overhaul; not part of this phase's
   done/pending accounting.)

**Phase 2 — harden Snap (the priority channel over Flatpak) — done 2026-07-10**
8. ~~`snap/snapcraft.yaml` was on `base: core20`.~~ Done — Launchpad's
   snapcraft (9.0.1) had actually dropped `core20` support entirely
   (`Base 'core20' is not supported by this version of Snapcraft`), so this
   wasn't optional hardening, it was a real outage: Snap builds were failing
   before ever reaching this repo's source. Migrated to `base: core22` with
   the `extensions: [gpu]` app extension replacing the old hand-rolled
   `graphics-core20`/`mesa-core20` content-interface plugs/layout/environment
   blocks, plus an explicit `build-snaps: [go/1.21/stable]` (Go is no longer
   bundled by the plugin under core22). Confirmed working — owner reports a
   successful Launchpad build and working auto-builds.
9. ~~No CI builds/tests the Snap package today.~~ Resolved as a side effect
   of #8 — Launchpad auto-builds are the CI for this channel and are
   confirmed working again.
10. ~~Verify whether AppImage is actually broken.~~ Confirmed fixed per
    GitHub Actions — the version-drift fix from Phase 0 was sufficient,
    it wasn't independently broken.

**Phase 3 — UI overhaul — done 2026-07-10 (shipping as v0.8.0)**
11. ~~Decide Fyne vs. reassessing the toolkit.~~ Decided 2026-07-10: **staying
    on Fyne.** Owner has since built two other substantial Fyne projects and
    has real accumulated expertise there now (though nothing directly
    reusable in this repo) — this isn't a default/inertia choice, it's an
    informed one. Current layout is still a static Fyne grid that doesn't
    scale or let users rearrange tiles; the fix is a better layout within
    Fyne, not a framework change.
12. ~~Define what "good" looks like before implementation.~~ Done — see
    `design.md` (the design brief that was handed to a design-focused
    conversation, with `design/` holding the resulting mockup/reference
    material) for the full spec both tabs below were built against.
13. ~~Implement the Memory tab redesign.~~ Done — `memorytile.go`'s
    `MemoryDashboard` replaced the old single-tile grid with a radial usage
    gauge, a used/cached/buffers/free/**swap** breakdown (swap is now read
    via `parseMemInfo` in `sysinfo.go`), and a full-width dual-series
    (memory + swap) history area chart (`DrawAreaChart` in `graphics.go`).
    The yellow status band in `theme.go` is now actually wired up (swap
    row + CPU utilization bars).
14. ~~Implement the CPU tab redesign.~~ Done — `cputab.go`'s `CPUTab` adds
    an aggregate stats strip (Threads/Avg util/Peak core/Max clock) and a
    Tiles/List view toggle (persisted via `Settings.ViewMode`). Tiles use a
    responsive `container.NewGridWrap` that reflows column count with window
    width instead of a fixed-column grid; List view (`tile.go`'s
    `CoreListRow`) is a compact two-column row layout with an inline `Bar`
    (`bar.go`) for utilization. Graphs (`DrawSparkline`) now render at their
    actual displayed raster size instead of a fixed 120×50px bitmap.
15. Both tabs' custom `canvas.Text`/`canvas.Rectangle` elements needed an
    explicit `RefreshTheme()` pass wired into `main.go`'s theme-change
    listener, since — unlike built-in widgets — they don't repaint
    themselves on a live theme switch. Same underlying issue as the
    `CustomTheme` recursion bug in Phase 0, different symptom.

No further phases are currently planned. If a new UI/feature initiative
starts, prefer adding a new dated section here over resurrecting this one.
