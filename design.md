# Design brief: Cycles UI overhaul (CPU + Memory tabs)

Copy everything below this line into a design-focused Claude conversation,
along with screenshots of the current CPU and Memory tabs (light and dark
if you have them).

---

## What this app is

Cycles is a desktop system monitor (Linux, Go, built on the
[Fyne](https://fyne.io/) v2 GUI toolkit — that choice is final, not up for
debate in this design pass). It currently has two tabs: **CPU** and
**Memory**, each polling system stats every N seconds (configurable, default
2s) and redrawing. There's no push/streaming architecture — assume periodic
redraw, not smooth animation, unless you think animation is worth proposing
as an enhancement.

I'm attaching screenshots of both tabs as they exist today. Use them as
ground truth for the "before" — the description below explains *why* they
look the way they do and what's structurally wrong, which isn't always
obvious from a screenshot alone.

## The core problem

Both tabs currently use the same pattern: a fixed-column `GridWithColumns`
of fixed-minimum-size tiles (CPU tiles are 100x100px, Memory tiles are
150x150px), each containing stacked labels and a **hardcoded 120x50 pixel
bitmap graph** that's drawn with manual Bresenham line-drawing into a
`canvas.Image` — the graph does not resize with its container, the window,
or anything else. Resizing the window just adds dead space or clips
content; it never reflows.

That's a real bug for the CPU tab (see below), but for the **Memory tab
it's the wrong pattern entirely**, not just an unresponsive version of the
right pattern:

- **CPU** is naturally a repeating collection — N identical, structurally
  identical cores (label, utilization %, clock speed, history graph). A
  grid-of-tiles is the *correct* shape for this data; it just needs to be
  responsive (reflow columns based on available width, tiles that scale,
  graphs that redraw at their actual rendered size instead of a fixed
  bitmap).
- **Memory** is a single aggregate metric with several *named,
  non-repeating* sub-components (used, free, cached, and currently-missing
  swap) plus one historical usage graph. There's exactly one tile in a
  2-column grid today, which is why it looks sparse and arbitrary — the
  grid pattern was copy-pasted from the CPU tile rather than designed for
  what memory data actually looks like. This tab should probably be a
  proper **dashboard/breakdown layout** — think: one large usage
  indicator (gauge or big progress bar) + a labeled breakdown of the
  components + a bigger, persistent history graph — not N small
  identical-looking tiles.

## What I want from you

Propose a redesign for **both** tabs:

1. **CPU tab**: a responsive/reflowing layout — tiles (or whatever
   replaces them) that scale with window size and reflow column count
   based on available width, rather than a fixed grid that just clips or
   leaves dead space. Graphs should render at their actual displayed
   size, not a fixed bitmap. Feel free to propose whether "one tile per
   core" is even still the right micro-pattern at high core counts (a
   64-core machine in an 8-column fixed grid is already unusable) — e.g.
   compact rows instead of boxy tiles, sparkline-style graphs, etc.
2. **Memory tab**: rethink the structure from scratch given the data-shape
   mismatch explained above. I'd like to see used/free/cached/buffers and
   swap all represented, a large/prominent overall usage indicator, and a
   history graph that's a first-class element of the layout rather than a
   small afterthought crammed into a tile. Open to genuinely different
   layouts here (e.g. a single wide panel instead of a grid at all).

## Constraints and things to account for

- **Fyne v2 only.** Assume standard Fyne widgets/containers/canvas
  primitives (`widget.ProgressBar`, `container.NewBorder`,
  `container.NewGridWrap`, custom `fyne.Layout` implementations, etc.) —
  no web views, no external UI frameworks.
- **Theme support**: the app already has light/dark theme switching (a
  `CustomTheme` wrapper) and defined green/yellow/red status colors per
  utilization band (currently only green/red are actually wired up — a
  yellow "warning" band exists in code but is unused). Any design should
  work in both themes.
- **Live-updating, not static** — every element that shows a number or a
  graph gets refreshed on each poll tick. Design with that update cadence
  in mind, not as a one-time render.
- **Window can be resized freely**, including much smaller or much larger
  than the current fixed tile sizes assume. The design needs to hold up
  across that whole range, not just look good at one default size.
- Prefer solutions that don't require abandoning Fyne's layout system for
  fully custom pixel-pushed UI, but custom `canvas` drawing (like the
  existing graph rendering) is fair game where it's the right tool.

## What I want back

A design spec, not code — wireframe-level descriptions (ASCII mockups are
fine) for both tabs, called out per breakpoint/size if the layout should
adapt structurally (not just scale) at different window sizes, and a list
of the specific Fyne widgets/layout techniques you'd use for each element
and why. If the CPU and Memory tabs should share an underlying responsive
layout primitive (e.g. a generic reflowing container both tabs build on),
say so explicitly — that reuse question is open and worth an opinion.
