# Two Realities

Entry for GitHub Game Off 2026 (theme: **SHIFT**, deadline **Dec 1, 2026**). Go + Ebitengine v2, shipped as WebAssembly on itch.io.

## Concept
2D puzzle-platformer. The player exists in two overlapping realities:
- **Yellow** world: warm, free, open (yellows, warm whites).
- **Grey** world: cold, watched, oppressive (greys, muted reds).

Shift swaps the active reality. Only the active world's tiles collide. The inactive world is drawn as a translucent ghost overlay so the player can plan. A shift is blocked (screen shake + red flash) if the player would end up inside a solid tile of the other world.

Original characters and names only — no band names, logos, or lyrics anywhere.

## Architecture rules
- Game logic (`geom`, `world`, `level`, `physics`, `player`, `game`, `input`) must NOT import ebiten. Ebiten lives only in `internal/ebitenapp`.
- Input is abstracted as `input.Intent` + `input.Source`; logic tests feed `Intent` values directly.
- Package-by-feature under `internal/`. Idiomatic Go, SOLID, DRY, small files, clear names, no long comments or doc blocks. All code and identifiers in English.
- Table-driven tests for parsing, collision, jump timing, shift rules.

## Layout
- `cmd/game` — entry point (desktop and wasm).
- `internal/level` — grid, `TileSize` (10px), ASCII parser, embedded levels in `levels/*.txt`.
- `internal/physics` — axis-separated AABB vs tile grid. Sides/top of the map are solid, bottom is open (falling out restarts).
- `internal/player` — movement, gravity, coyote time and jump buffer (`tuning.go`).
- `internal/game` — `Session`: active world, shift, restart, exit detection.
- `internal/ebitenapp` — keyboard, palettes, drawing, ebiten.Game. Logical resolution 320x180, 60 TPS.

## Level format
```
[yellow]
#P....#
[grey]
#....E#
```
`#` solid, `.` empty, `P` spawn and `E` exit exactly once each across both layers. Both layers must be the same size.

## Controls
A/D or arrows move · Space/W/Up jump · Left/Right Shift swap worlds · R restart.

## Commands
`make run` · `make test` · `make lint` (go vet + gofmt) · `make build-wasm` · `make serve` (web/ on :8080). Go version pinned in `.tool-versions` (asdf).

## Roadmap
1. Core loop: movement, two worlds, shift, ghost — **done**
2. Yellow resource: collect yellow fragments in the grey world; each shift costs one
3. Hazards and "watchers" in the grey world; death/respawn
4. Level progression (5–8 levels) and title screen
5. Art and audio pass
6. Stretch: beat-synced shifting
7. itch.io release polish

Only implement the step currently being worked on.
