# Two Realities

Entry for GitHub Game Off 2026 (theme: **SHIFT**, deadline **Dec 1, 2026**). Go + Ebitengine v2, shipped as WebAssembly on itch.io.

## Concept
2D puzzle-platformer. The player exists in two overlapping realities:
- **Yellow** world: warm, free, open (yellows, warm whites).
- **Grey** world: cold, watched, oppressive (greys, muted reds).

Shift swaps the active reality. Only the active world's tiles collide. The inactive world is drawn as a ghost (hatched outlines, clearly "not solid here") so the player can plan.

### Fragments and asymmetric shift
- Yellow fragments (`Y`) exist only in the grey layer and are collected only while grey is active. Each is collected once per attempt.
- **Yellow → Grey is free. Grey → Yellow costs one fragment.**
- `Session.TryShift()` returns an explicit `ShiftOutcome`: `ShiftDone`, `ShiftBlockedSolid` (player would overlap a solid tile in the target world → screen shake + red flash) or `ShiftBlockedNoFragments` (→ fragment counter blinks red and an "R to restart" hint shows until a fragment is collected or the level restarts). If both apply, no-fragments wins.
- Restart (R or death) resets world, fragments, counter and watchers to the level's start values.

### Hazards
- **Spikes** (`^`, either layer): lethal only while their layer is active (hitbox = bottom 5px of the tile). Inactive spikes are drawn as ghost outlines.
- **Watchers** (`W`, grey layer only): 8x10 body, patrol horizontally at `WatcherSpeed`, turning at grey walls or grey ledges. They always patrol (deterministic, even while yellow is active) but only see/harm the player while grey is active; otherwise they render as ghosts.
  - Line of sight: horizontal ray at eye height in the facing direction, `SightRange` tiles max, stopped by grey solid tiles. The player is seen if the ray crosses their box (jumping above it breaks sight).
  - Alert meter fills one per tick while seen (`AlertTicks` = 0.4s, "!" + meter above the watcher); leaving sight or shifting to yellow resets it. Full meter = caught.
  - **Core dodge:** shifting to yellow before the meter fills saves the player.
- **Death** (spikes, caught, falling out): `DeathFreezeTicks` (0.5s) freeze with white flash, shake and blinking player, input ignored, then restart from spawn.

### Progression
Levels load from `internal/level/levels/*.txt` sorted by name (`LoadAll`). Reaching the exit advances to the next level; the last loops to the first.

Original characters and names only — no band names, logos, or lyrics anywhere.

## Architecture rules
- Game logic (`geom`, `world`, `level`, `physics`, `player`, `game`, `input`) must NOT import ebiten. Ebiten lives only in `internal/ebitenapp`.
- Input is abstracted as `input.Intent` + `input.Source`; logic tests feed `Intent` values directly.
- Package-by-feature under `internal/`. Idiomatic Go, SOLID, DRY, small files, clear names, no long comments or doc blocks. All code and identifiers in English.
- Table-driven tests for parsing, collision, jump timing, shift rules.

## Layout
- `cmd/game` — entry point (desktop and wasm).
- `internal/level` — grid, `TileSize` (10px), header + ASCII parser, embedded levels in `levels/*.txt`.
- `internal/physics` — axis-separated AABB vs tile grid. Sides/top of the map are solid, bottom is open (falling out restarts).
- `internal/player` — movement, gravity, coyote time and jump buffer (`tuning.go`).
- `internal/game` — `Session`: active world, shift outcomes and cost, fragments, watchers, spikes, death freeze, feedback timers, restart hint, level progression.
- `internal/ebitenapp` — keyboard, palettes, scene, cached ghost layer (tiles + spikes, per level), fragments (bobbing diamonds), spikes, watchers, HUD with a built-in 3x5 bitmap font (`font.go`). Logical resolution 320x180, 60 TPS.
- Level tests in `internal/game/levels_test.go` script a full solution for each level and prove the intended mechanic is required.

## Level format
```
fragments: 0

[yellow]
#P....#
[grey]
#..Y.E#
```
- First non-empty line is the header `fragments: N` (non-negative integer, starting fragment count). Missing or invalid header is an error.
- `#` solid, `.` empty, `P` spawn and `E` exit exactly once each across both layers.
- `Y` fragment and `W` watcher: grey layer only (in yellow they are parse errors).
- `^` spike: either layer.
- Both layers must be the same size.
- Keep row 0 solid in both layers: the HUD (world, `LV n/m`, fragments, "R TO RESTART") lives inside that top border.
- Physics reference for design: run 1.4 px/tick, jump ≈ 26px high / ≈ 40px long, so 1-tile corridors block jumping and gaps of 5+ tiles can't be jumped.

## Controls
A/D or arrows move · Space/W/Up jump · Left/Right Shift swap worlds · R restart.

## Commands
`make run` · `make test` · `make lint` (go vet + gofmt) · `make build-wasm` · `make serve` (web/ on :8080, override with `PORT=`). Go version pinned in `.tool-versions` (asdf).

## Roadmap
1. Core loop: movement, two worlds, shift, ghost — **done**
2. Yellow resource: fragments in the grey world, asymmetric shift cost, HUD, hatched ghost — **done**
3. Hazards: spikes, grey-world watchers, death/respawn, 2-level progression, compact HUD — **done**
4. Level progression (5–8 levels) and title screen
5. Art and audio pass
6. Stretch: beat-synced shifting
7. itch.io release polish

Only implement the step currently being worked on.
