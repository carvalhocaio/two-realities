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
- **R** restarts the whole level (start values, checkpoint discarded). **Death** respawns at the last checkpoint (or spawn) with that checkpoint's snapshot. Both reset world to yellow and watchers to their start positions.

### Hazards
- **Spikes** (`^`, either layer): lethal only while their layer is active (hitbox = bottom 5px of the tile). Inactive spikes are drawn as ghost outlines.
- **Watchers** (`W`, grey layer only): 8x10 body, patrol horizontally at `WatcherSpeed`, turning at grey walls or grey ledges. They always patrol (deterministic, even while yellow is active) but only see/harm the player while grey is active; otherwise they render as ghosts.
  - Line of sight: horizontal ray at eye height in the facing direction, `SightRange` tiles max, stopped by grey solid tiles. The player is seen if the ray crosses their box (jumping above it breaks sight).
  - Alert meter fills one per tick while seen (`AlertTicks` = 0.4s, "!" + meter above the watcher); leaving sight or shifting to yellow resets it. Full meter = caught.
  - **Core dodge:** shifting to yellow before the meter fills saves the player.
- **Death** (spikes, caught, falling out): `DeathFreezeTicks` (0.5s) freeze with white flash, shake and blinking player, input ignored, then restart from spawn.

### Checkpoints
- `C` marker, either layer, activated by touching it in any world (lit yellow flag). Activation snapshots fragment count and which fragments were collected.
- Death respawns there in yellow with the snapshot restored; spending fragments after the checkpoint is refunded on death. R discards the checkpoint.

### Camera and screen layout
- Screen 320x180: a fixed 10px HUD bar on top (world, `LV n/m` or "R TO RESTART", fragments), the play view is 320x170 below it. 17-row levels fit vertically.
- `internal/camera` (pure, tested): smooth follow (`Smoothing` 0.15) with a 48x32 dead zone, clamped to level bounds; levels smaller than the view are centered. Snaps on respawn/level change (`Session.Epoch()`).
- The adapter renders the world into a level-sized canvas and blits it with the camera offset; the HUD is drawn screen-fixed.

### Flow and progression
- `game.Flow`: title ("press Shift to start") → playing → end screen ("press Shift for title") → title.
- Levels load from `internal/level/levels/*.txt` sorted by name (`LoadAll`). Exit advances; finishing the last level shows the end screen.
- Dev build (`make run-dev`, `-tags dev`): N next level, P previous level.

### Level list (one new idea per level)
| # | File | Teaches | Needs |
|---|------|---------|-------|
| 1 | `01_first_shift` | move, jump a step, first shift over a grey-only bridge | no fragments, never return |
| 2 | `02_fragments` | fragments + asymmetric cost: grab one on the bridge, pay it to pass a grey-only wall | 1 fragment |
| 3 | `03_spikes` | spikes per world: alternate worlds 4 times | 2 fragments, exactly |
| 4 | `04_watcher` | first watcher: hide on the yellow pocket while it passes, return behind it | 1 fragment |
| 5 | `05_planning` | watcher + spikes with exactly 2 fragments: wasting one means R | 2 fragments, exactly |
| 6 | `06_finale` | everything + checkpoints: corridor on top, drop, spikes walking left at the bottom | 3 fragments, 2 checkpoints |

`internal/game/levels_test.go` scripts a solution for every level (under 90s), checks naive attempts fail for the intended reason, and enforces the teaching order (each level adds at most one new mechanic). `internal/level/invariants_test.go` checks every embedded level (spawn/checkpoints stand on yellow ground, grey-only markers in open grey space, etc.).

Original characters and names only — no band names, logos, or lyrics anywhere.

## Architecture rules
- Game logic (`geom`, `world`, `level`, `physics`, `player`, `game`, `input`) must NOT import ebiten. Ebiten lives only in `internal/ebitenapp`.
- Input is abstracted as `input.Intent` + `input.Source`; logic tests feed `Intent` values directly.
- Package-by-feature under `internal/`. Idiomatic Go, SOLID, DRY, small files, clear names, no long comments or doc blocks. All code and identifiers in English.
- Table-driven tests for parsing, collision, jump timing, shift rules.

## Layout
- `cmd/game` — entry point (desktop and wasm).
- `internal/camera` — camera follow/clamp math.
- `internal/level` — grid, `TileSize` (10px), header + ASCII parser, embedded levels in `levels/*.txt`.
- `internal/physics` — axis-separated AABB vs tile grid. Sides/top of the map are solid, bottom is open (falling out restarts).
- `internal/player` — movement, gravity, coyote time and jump buffer (`tuning.go`).
- `internal/game` — `Session`: active world, shift outcomes and cost, fragments, watchers, spikes, death freeze, feedback timers, restart hint, level progression.
- `internal/ebitenapp` — keyboard, palettes, scene, cached ghost layer (tiles + spikes, per level), fragments (bobbing diamonds), spikes, watchers, checkpoints, world view with camera (`worldview.go`), HUD bar, title/end screens (`screens.go`), built-in 3x5 bitmap font with integer scaling (`font.go`), dev keys (`dev_on.go`/`dev_off.go`). Logical resolution 320x180, 60 TPS.
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
- `C` checkpoint: either layer; must stand on yellow ground (respawn is in yellow).
- Spawn must stand on yellow ground. 17 rows fit the view exactly; taller/wider levels scroll.
- Physics reference for design: run 1.4 px/tick, jump ≈ 26px high / ≈ 40px long, so 1-tile corridors block jumping and gaps of 5+ tiles can't be jumped.

## Controls
A/D or arrows move · Space/W/Up jump · Left/Right Shift swap worlds · R restart.

## Commands
`make run` · `make run-dev` (N/P level skip) · `make test` · `make lint` (go vet, also with `-tags dev`, + gofmt) · `make build-wasm` · `make serve` (web/ on :8080, override with `PORT=`). Go version pinned in `.tool-versions` (asdf).

## Roadmap
1. Core loop: movement, two worlds, shift, ghost — **done**
2. Yellow resource: fragments in the grey world, asymmetric shift cost, HUD, hatched ghost — **done**
3. Hazards: spikes, grey-world watchers, death/respawn, 2-level progression, compact HUD — **done**
4. Camera, checkpoints, six-level sequence, title and end screens, dev level skip — **done**
5. Art and audio pass
6. Stretch: beat-synced shifting
7. itch.io release polish

Only implement the step currently being worked on.
