# Two Realities

Entry for GitHub Game Off 2026 (theme: **SHIFT**, deadline **Dec 1, 2026**). Go + Ebitengine v2, shipped as WebAssembly on itch.io.

## Concept
2D puzzle-platformer. The player exists in two overlapping realities:
- **Yellow** world: warm, free, open (yellows, warm whites).
- **Grey** world: cold, watched, oppressive (greys, muted reds).

Shift swaps the active reality. Only the active world's tiles collide. The inactive world is drawn as a ghost (dotted outlines in that world's tile color, clearly "not solid here"; watchers and the Warden's eye too) so the player can plan.

### Art direction (step 5, option B)
- 320x180 and 10px tiles stay; art never changes geometry or logic timing (the shift wipe is visual only).
- Four colors per world (`palette.go`): yellow = cream background, amber tiles, dark brown ink (player/text), orange exit; grey = graphite background, slate tiles, light ink, exit in fragment yellow.
- Shift wipe (`wipe.go`): when the active world changes (not on respawn or level change), the previous frame slides off to the right behind a yellow/black hazard band for `wipeTicks` (9 ticks). Purely visual; logic already switched. Shifts land on the beat in quantized/auto, so the wipe does too.
- Glow (`glow.go`): additive radial glow (one pre-rendered gradient drawn with `BlendLighter`, no shader) on the warm things while grey is active: fragments, secret, open exit, active checkpoint flag, lit and current Warden lamps; a slow pulse. Skipped in yellow, where it would be invisible on cream.
- Grey-only backdrop (`backdrop.go`): two layers of tower silhouettes (parallax 0.25 and 0.5, seeded per level, anchored to the level bottom) and a screen-locked vignette, both drawn into the canvas before the tiles, so they shade the background only and never dim tiles, the player or hazards.
- Player feel (`playerfx.go`, visual only): a 2px eye on the facing side (from horizontal movement), and on landings faster than `landingSpeed` a short squash (wider and lower, feet fixed) plus three dust specks in the tile color. No trail, no rounded corners. Resets on respawn and level change.
- Pieces and HUD: watchers are `enemyColor` blocks with a 4px red visor on the facing side; the exit is a portal (frame in the world's exit color, light door, 4px taller than its tile, dim while locked); the HUD bar takes the world's background and its panels are dark (`panel`/`panelInk`) with a 1px frame in the tile color (the world chip framed in the shown world's tile color); the transmission box uses the same panel colors.
- Shared accents: fragment yellow (in grey it only appears on fragments, exit, lamps, secret and lit checkpoints), `dangerColor` red only for danger (active spikes, watcher visor and sight, the Warden's beam/rail, blocked flash), `enemyColor` near-black for watcher and Warden bodies.
- The inactive world is a ghost: dotted outlines for tiles, watchers and the eye, outlined spikes, all in that world's tile color (or `enemyColor` for enemies), so it never reads as solid or as active danger.

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

### Shift modes and rhythm
- Per-level header `shift: free|quantized|auto` (default `free`). Levels 1-8 are free.
  - **free**: Shift resolves immediately (rules above).
  - **quantized**: Shift toggles a queued request (pressing again cancels) that resolves on the next beat through `TryShift`, so cost and both blocked outcomes apply at resolution. The HUD world chip blinks the target world while queued.
  - **auto**: the world swaps by itself every `auto_every: N` beats (required in auto, rejected elsewhere), counted from spawn/respawn. Player Shift is ignored, swaps are free, `Y` is a parse error. A swap that leaves the player inside a solid tile kills (`DeathCrushed`). A 3-2-1 countdown is drawn above the player.
- `internal/rhythm` (pure): `Clock` interface (`Beat()`, `Phase()` 0..1), `BPM = 100` (600ms beats), `At(elapsed, bpm)`, and `Manual` for tests. The session only reacts to beat index changes; logic never counts frames for rhythm.
- The ebiten adapter implements `Clock` from the audio player's position minus an audio latency offset (`defaultAudioLatency` in `latency_js.go` / `latency_other.go`), so music and gameplay never drift. The player sets the offset in the pause menu ("AUDIO OFFSET", `AudioOffsetStep` 10ms, clamped to ±`AudioOffsetLimit` 500ms, positive = beats later for laggy output such as Bluetooth); it lives in `Flow`, is saved as `audio_offset_ms`, and a dot next to the item pulses on the beat for calibration.
- Placeholder music is synthesized in `synth.go`: one 4-beat bar (pitch-swept sine kick on beats, noise hi-hat on off-beats) looped by an endless reader into a float32 player. It starts on the title Shift press (browser autoplay rule). **M** mutes (volume 0, the clock keeps running). The HUD bottom border pulses on every beat.

### Transmissions
- `T` marker (either layer) + text sections after the layers: `[transmissions]` (English, required when there are `T` markers) and optional `[transmissions.pt]`, each with one line per `T` in reading order (yellow layer first, then grey, row by row). Count mismatch in any language is a parse error; a missing Portuguese section falls back to English. Text is uppercased and must be drawable by `internal/font` (max `MaxTransmissionLength` runes).
- Touching a `T` while its layer is active opens it (once per attempt: death and R reset it). Typewriter box placed away from the player (top of the view when the player is in the lower half, bottom otherwise; it keeps its side while open and only flips if the player walks under it, `textbox.go`); it fades after `TransmissionTicks` (4s) or when the player is farther than `TransmissionRange` (12 tiles). Heard markers are drawn dimmed.
- Tone: notes left by someone who tried to escape before. Placeholder lines, to be rewritten.

### Speedrun timer and records
- Level timer starts on the first acting input (move, jump, shift) and stops at the exit; deaths and R do not reset it. The run timer starts on the first input of level 1 and stops at the last exit. All timing is in ticks (60 TPS), shown as `M:SS.CC`.
- `internal/save`: `Storage` interface (`Load() (Data, error)`, `Save(Data) error`), `Data` = best ticks per level name, best run, secrets found per level, audio offset, language. Adapters: `LocalStorage` (wasm), `File` (desktop, `os.UserConfigDir()/two-realities/save.json`), `Memory` (tests). `save.Default()` picks per platform; load errors fall back to empty data.
- Every browser storage key goes through `save.Key(name)` (prefix `two-realities:`; the save lives in `two-realities:save`). The pre-0.3.1 key `two-realities/save` is migrated on load.
- On wasm, `Default()` wraps `LocalStorage` in `save.Fallback`: a probe write on load detects blocked/full storage; `ErrUnavailable` on load or any failed write switches to in-memory storage for the session (`Degraded()`), and no error or panic ever reaches the game. Corrupt JSON starts empty but keeps localStorage. The title screen shows "PROGRESS WON'T BE SAVED IN THIS BROWSER" until the first run starts (`Flow.StorageNotice`).
- JS access to `localStorage` uses `Reflect.get` through `Value.Call`: `wasm_exec.js` wraps `Call` in try/catch but not `Get`, so a throwing getter (blocked third-party storage) would freeze the Go runtime.
- `game.Records` wraps the storage. Level bests save on each exit unless assist was on during that level; the run best saves only for runs started at level 1, without dev skips and without assist (end screen shows `PRACTICE` / `ASSIST`).
- **T** toggles the HUD timer (hidden by default). The level-complete banner (`TIME ... BEST ...` / `NEW BEST`) always shows for 2.5s.

### Assist mode and pause
- **P, Enter or Esc** open the pause menu (`ScreenPaused`): resume, assist on/off, mute on/off, audio offset (left/right or -/= while selected), language, back to title. Arrows or W/S move, Enter/Space/Shift confirm (confirm wins over pause, so Enter picks the item), P/Esc resume. The session is frozen while paused.
- Auto-pause: `internal/platform.Signals` turns window blur, `visibilitychange` to hidden, fullscreen exit (`fullscreenchange`, or a resize from screen size, since an iframe made fullscreen by the parent page may not get the event) and desktop focus loss (`ebiten.IsFocused` edge) into one `Intent.FocusLost`. `Flow` pauses only from `ScreenPlaying` and never resumes by itself. Listeners live in `platform/browser_js.go`; the keyboard adapter polls the signals.
- Assist: watcher alert needs `2*AlertTicks`; an auto swap into a solid tile pushes the player to the nearest free spot (up, then sideways, up to 3 tiles) instead of killing. The assist choice persists across runs within the session (`Flow`).
- Mute lives in `Flow` (M or the menu); the adapter only applies the volume.

### Secrets
- `S` marker, at most one per level, either layer; drawn only while its layer is active (never in the ghost). Collecting it adds a fragment and marks the secret; it is saved when the level is finished with it (checkpoints snapshot it). The end screen shows secrets found out of the total.
- Every level uses the same hidden shape: a 2-tile ledge in yellow 2 tiles above the floor, and a 2-tile ledge in the secret's layer 2 tiles higher and 1 tile to the right, with the secret on top. Reaching it needs a mid-air shift at the top of the jump (queued on the beat in L7, timed to the auto flip in L8). `secrets_levels_test.go` grabs every secret.

### Languages
- English and Brazilian Portuguese. `internal/i18n` (pure): `Lang` (`EN`, `PT`), typed `Key`s, one table per language, `lang.T(key, args...)`; a test checks every key exists in both and every string is drawable by the font.
- `Flow` owns the language: a saved choice (`language` in save data) wins, otherwise `UseLocale(platform.Locale())` detects it (`navigator.language` on wasm, `LC_ALL`/`LC_MESSAGES`/`LANG` on desktop; anything starting with `pt` is Portuguese). **L** on the title screen or the pause menu "LANGUAGE" item switches and saves.
- The game name stays "TWO REALITIES" in both languages. Code, identifiers, commits and docs stay in English.
- `internal/font` (pure) holds the 3x5 glyphs, including Á À Â Ã É Ê Í Ó Ô Õ Ú Ç with a one-row accent above (cedilla below); it is the single authority on drawable text (`font.Missing`).

### Boss: The Warden (O Carcereiro)
- Level header `boss: warden` (then `shift`/`auto_every` are rejected: each phase sets its own mode). Single-screen 32x17 arena, `13_warden`, the last level.
- **Eye** rides the `O` rail under the ceiling (grey layer, one unbroken row). Its **beam** drops straight down and stops at the first grey solid tile under either edge of the beam (`BeamWidth` 6px): grey platforms are cover, the yellow world is immune. Instant death (`DeathBeam`) while grey is active and the player overlaps it.
- **3 lamps** (`L`, on yellow ground, open in grey), lit in order bottom to top (then left to right) by touching the current one **while grey is active**. Each lit lamp ends a phase, opens that phase's gate, acts as a checkpoint (death keeps lit lamps; R restarts the fight) and gives `RestBeats` (2) without beam. The exit stays locked (drawn dim) until the third lamp.
- **Gates** `1`/`2` (either layer): solid in both worlds until lamp 1/2 is lit. The session keeps its own copy of the layers with closed gates (`Session.layer`).
- Phases (`WardenPhases` in `game/warden.go`), each one taking more control of the shift:
  1. free shift; searchlight **sweeps** continuously at 80 px/beat, ping-pong along the rail, starting at the end farthest from the player.
  2. quantized shift; sweep at 100 px/beat (queue the yellow before the light arrives: queued shifts resolve before the beam on the same beat).
  3. auto shift every 2 beats; the beam **tracks**: aims at the player's column, dotted telegraph for 2 beats, fires for half a beat.
- Assist: sweep ×0.6, track telegraph +1 beat. Music adds a layer per phase (bass, double hats, snare) by swapping the looped bar only at the loop boundary, so the beat clock never jumps.
- `warden_level_test.go` solves the fight on all 36 beat phases with a planner that predicts the light and only crosses open stretches it can clear; rushing the bridge and standing still in phase 3 die to the beam.

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
- Dev build (`make run-dev`, `-tags dev`): `]` next level, `[` previous level.

### Level list (one new idea per level)
| # | File | Teaches | Needs |
|---|------|---------|-------|
| 1 | `01_first_shift` | move, jump a step, first shift over a grey-only bridge | no fragments, never return |
| 2 | `02_fragments` | fragments + asymmetric cost: grab one on the bridge, pay it to pass a grey-only wall | 1 fragment |
| 3 | `03_spikes` | spikes per world: alternate worlds 4 times | 2 fragments, exactly |
| 4 | `04_watcher` | first watcher: hide on the yellow pocket while it passes, return behind it | 1 fragment |
| 5 | `05_planning` | watcher + spikes with exactly 2 fragments: wasting one means R | 2 fragments, exactly |
| 6 | `06_finale` | everything + checkpoints: corridor on top, drop, spikes walking left at the bottom | 3 fragments, 2 checkpoints |
| 7 | `07_crossfire` | mastery (Act II): two watchers split by a grey-only pillar crossed through yellow; three hiding places, three fragments | 3 fragments, exactly; 1 checkpoint |
| 8 | `08_descent` | mastery (Act II): vertical drop; shift to grey before landing on the watcher's floor, follow it and jump over it when it turns, shift mid-fall, alternate spikes at the bottom | 2 fragments |
| 9 | `09_on_the_beat` | quantized shift: queue the hide on the yellow pocket before the watcher turns back; reacting on sight gets caught when the beat is far | 1 fragment |
| 10 | `10_metronome` | auto shift every 4 beats: wait on shared ground, cross single-world floors within one phase, don't stand under the grey block | no fragments |
| 11 | `11_offbeat` | mastery (Act III): alternating spike run in quantized mode (queue on the safe tiles), then the corridor watcher with an early queue | 3 fragments |
| 12 | `12_clockwork` | mastery (Act III): auto every 2 beats (the boss's final rhythm); wait on crates above the watcher's sight, cross the yellow spikes in a grey window behind its back, mind the grey blocks | no fragments |
| 13 | `13_warden` | boss: light 3 lamps in grey while the eye's beam sweeps (free, then quantized) and finally tracks you under auto shift | 1 fragment (2 available) |

Acts: I = 1-4 (one mechanic each), II = 5-8 (combining, ends with two mastery levels), III = 9-13 (rhythm modes, two mastery levels, the Warden). Mastery levels add no new mechanic; `mastery_levels_test.go` solves 11/12 on every beat phase and checks their naive attempts.

`internal/game/levels_test.go` scripts a solution for every level (under 90s), checks naive attempts fail for the intended reason, and enforces the teaching order (each level adds at most one new mechanic; a shift mode counts as one). `rhythm_levels_test.go` drives L9/L10 with a manual clock: L9's solution works on every beat phase, reacting on sight fails on some. `internal/level/invariants_test.go` checks every embedded level (spawn/checkpoints stand on yellow ground, grey-only markers in open grey space, etc.).

Original characters and names only — no band names, logos, or lyrics anywhere.

## Architecture rules
- Game logic (`geom`, `world`, `level`, `physics`, `player`, `game`, `input`, `rhythm`, `save`, `platform`, `i18n`, `font`) must NOT import ebiten. Ebiten lives only in `internal/ebitenapp`.
- Input is abstracted as `input.Intent` + `input.Source`; logic tests feed `Intent` values directly.
- Package-by-feature under `internal/`. Idiomatic Go, SOLID, DRY, small files, clear names, no long comments or doc blocks. All code and identifiers in English.
- Table-driven tests for parsing, collision, jump timing, shift rules.

## Layout
- `cmd/game` — entry point (desktop and wasm).
- `internal/camera` — camera follow/clamp math.
- `internal/rhythm` — beat clock interface, BPM, beat/phase math.
- `internal/i18n` — languages, string tables, locale detection.
- `internal/font` — 3x5 glyph table with accents, drawable-text check.
- `internal/save` — `Storage` interface, save `Data`, key prefix, `Fallback`, localStorage/file/memory adapters.
- `internal/platform` — focus/visibility/fullscreen signals turned into a pause request; browser listeners via `syscall/js`.
- `internal/level` — grid, `TileSize` (10px), header + ASCII parser, embedded levels in `levels/*.txt`.
- `internal/physics` — axis-separated AABB vs tile grid. Sides/top of the map are solid, bottom is open (falling out restarts).
- `internal/player` — movement, gravity, coyote time and jump buffer (`tuning.go`).
- `internal/game` — `Session`: active world, shift outcomes and cost, shift modes (`beats.go`), transmissions, secrets, timer and `Records`, assist, pause menu in `Flow`, fragments, watchers, spikes, death freeze, feedback timers, restart hint, level progression, the Warden boss (`warden.go`, `beam.go`).
- `internal/ebitenapp` — keyboard, palettes, scene, cached ghost layer (tiles + spikes, per level), fragments (bobbing diamonds), spikes, watchers, checkpoints, world view with camera (`worldview.go`), HUD bar, title/end screens (`screens.go`), synthesized music + audio beat clock (`synth.go`, `music.go`), beat pulse/queued blink/countdown (`beatfx.go`), transmission box (`textbox.go`), Warden eye/beam/lamps/gates and HUD pips (`warden.go`), pause menu (`pause.go`), signals, font drawing with integer scaling over `internal/font` glyphs (`font.go`), dev keys (`dev_on.go`/`dev_off.go`). Logical resolution 320x180, 60 TPS.
- Level tests in `internal/game/levels_test.go` script a full solution for each level and prove the intended mechanic is required.

## Level format
```
fragments: 0
shift: free

[yellow]
#P....#
[grey]
#..Y.E#
```
- Header: `key: value` lines before the first section. `fragments: N` (non-negative, starting count) is required; `shift: free|quantized|auto` is optional; `auto_every: N` (positive) is required in auto and rejected otherwise. Unknown or duplicate keys are errors.
- `#` solid, `.` empty, `P` spawn and `E` exit exactly once each across both layers.
- `Y` fragment and `W` watcher: grey layer only (in yellow they are parse errors).
- `^` spike: either layer.
- Both layers must be the same size.
- `C` checkpoint: either layer; must stand on yellow ground (respawn is in yellow).
- `T` transmission: either layer, paired with `[transmissions]` lines.
- `S` secret: either layer, at most one per level, allowed in auto levels. Boss levels have none.
- Boss levels (`boss: warden`): `L` lamp ×3, `O` rail (grey only), `1`/`2` gate tiles; these markers are errors in other levels.
- Spawn must stand on yellow ground. 17 rows fit the view exactly; taller/wider levels scroll.
- Physics reference for design: run 1.4 px/tick, jump ≈ 26px high / ≈ 40px long, so 1-tile corridors block jumping and gaps of 5+ tiles can't be jumped.

## Controls
A/D or arrows move · Space/W/Up jump · Left/Right Shift swap worlds · R restart · M mute · T timer · L language (title) · P/Enter/Esc pause (arrows/W/S + Enter/Space/Shift in the menu).

## Commands
`make run` · `make run-dev` (`]`/`[` level skip) · `make test` · `make lint` (go vet, also with `-tags dev`, + gofmt) · `make build-wasm` · `make package` (itch upload folder in `dist/html5`) · `make serve` (web/ on :8080, override with `PORT=`). Go version pinned in `.tool-versions` (asdf).

## Release (CI + itch.io)
- **CI** (`.github/workflows/ci.yml`): on push to `main` and on PRs → setup Go from `go.mod` (module cache on), `make lint`, `make test`, `make build-wasm`.
- **Deploy** (`.github/workflows/deploy.yml`): on tags `v*` or manual `workflow_dispatch` → `make test`, `make package`, install butler from `broth.itch.zone`, `butler push dist/html5 carvalhocaio/two-realities:html5 --userversion <tag | dev-<short sha>>`.
  - Auth: repository secret `BUTLER_API_KEY`, passed only as env to the push step (butler reads it). Never echo it or write it to disk.
  - The itch target lives in `ITCH_TARGET` at the top of `deploy.yml`; change it there if the itch user or project slug differs.
- **Local package**: `make package` builds `dist/html5/` (`index.html` at the root, `game.wasm`, `wasm_exec.js`), the exact folder uploaded, and prints raw/gzip `.wasm` size (≈13.5 MB / 3.2 MB at v0.1.0).
- **Cutting a release**: push `main`, then `git tag -a vX.Y.Z -m "vX.Y.Z"` and `git push origin vX.Y.Z`. The build lands in the `html5` channel of the itch project (draft until the page is published).
- **Prerequisites**: itch project `two-realities` of kind HTML (draft is fine) and the `BUTLER_API_KEY` secret (`gh secret set BUTLER_API_KEY`).
- `web/index.html` is itch-friendly: full-iframe canvas, no scrollbars, grabs focus on load/click/window focus (Ebitengine does not focus its canvas inside iframes), forwards keys that land outside the canvas so the first press counts, and blocks Space/arrows/Shift from scrolling the parent page. Verified in headless Chrome with a cross-origin iframe: Shift starts the game with zero clicks and the parent page does not scroll. Pixel scaling stays sharp because Ebitengine's default final screen uses nearest (integer) or `FilterPixelated` (fractional) filtering.

## Roadmap
1. Core loop: movement, two worlds, shift, ghost — **done**
2. Yellow resource: fragments in the grey world, asymmetric shift cost, HUD, hatched ghost — **done**
3. Hazards: spikes, grey-world watchers, death/respawn, 2-level progression, compact HUD — **done**
4. Camera, checkpoints, six-level sequence, title and end screens, dev level skip — **done**
4.5. CI + itch.io html5 deploy via butler, v0.1.0 live as draft — **done**
5. Art and audio pass — 5a art (option B: palettes, dotted ghost, shift wipe, glow, grey backdrop, player feel, pieces/HUD) **done**; 5b music pending
6. Beat-synced shifting: rhythm clock, quantized/auto modes, placeholder synth, L7-L8 — **done**
6.5. Atmosphere and replay: transmissions, speedrun timer + records, assist mode + pause menu, secrets — **done**
6.6. itch iframe fixes: P/Enter pause, auto-pause on blur/hidden/fullscreen exit, prefixed keys, storage fallback + notice — **done**
6.7. Audio offset in the pause menu and Brazilian Portuguese localization (UI + transmissions) — **done**
6.8. The Warden boss: three-phase single-screen fight (sweep, quantized sweep, tracking beam under auto shift), lamps, gates, layered music — **done**
6.9. Four mastery levels (07, 08, 11, 12) and the three-act order — **done**
7. itch.io release polish

Only implement the step currently being worked on.
