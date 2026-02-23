# Story: Config-Driven Skin System with Generic Command Palette

**Story Key**: skin-system
**Status**: in-progress
**Priority**: High

---

## Story

As a user of o8n,
I want a fully config-driven color/theme system with no hardcoded colors and a smooth runtime skin picker,
So that the UI is consistently beautiful, completely themeable, and I can switch skins instantly without restarting.

---

## Background & Design

### Design Author: Sally (UX Designer)
### Implementation: Amelia (Developer Agent)
### Approved by: Karsten

### Key Decisions

1. **Semantic color roles** — skin defines named roles (accent, danger, surface…), code uses only role names — never hex literals
2. **Full `StyleSet` struct** — replaces all `lipgloss.NewStyle()` + hardcoded colors in `styles.go` and `view.go`; rebuilt on skin switch
3. **`env.UIColor` removed entirely** — drop the field from struct + config + env YAML; no deprecation
4. **Generic command palette popup** — the existing context-switcher (`:`) becomes a reusable popup widget with `popupMode`; new `Ctrl+T` opens skin picker using the same widget
5. **Live preview on ↑↓** in skin picker; Esc reverts; Enter commits
6. **34 existing skin YAMLs updated** with new semantic role fields; new `skins/o8n-cyber.yaml` added as a showcase theme

### Semantic Color Roles (skin schema)

| Role | Purpose |
|---|---|
| `bg` | Terminal background |
| `fg` | Primary text |
| `fgMuted` | Secondary/dim text (hints, counters, muted) |
| `accent` | Primary accent (focus borders, keys, active elements) |
| `accentAlt` | Secondary accent (logo, breadcrumb active bg) |
| `surface` | Elevated surface (table header bg, modal bg) |
| `surfaceAlt` | Table row cursor highlight bg |
| `success` | OK / operational / saved confirmation |
| `warning` | Unknown / loading / caution |
| `danger` | Error / kill / unreachable |
| `info` | Informational messages |
| `borderFg` | Inactive border color |
| `borderFocus` | Active/focused border color |
| `crumbBg` | Breadcrumb bar background |
| `crumbFg` | Breadcrumb text |
| `crumbActiveBg` | Current breadcrumb segment background |
| `crumbActiveFg` | Current breadcrumb segment text |
| `jsonKey` | JSON detail view: key color |
| `jsonValue` | JSON detail view: string value color |
| `jsonNumber` | JSON detail view: number color |
| `jsonBool` | JSON detail view: boolean color |
| `btnPrimaryBg` | Save/confirm button background |
| `btnPrimaryFg` | Save/confirm button foreground |
| `btnSecondaryBg` | Cancel button background |
| `btnSecondaryFg` | Cancel button foreground |

### o8n-cyber Palette (new default showcase)

```yaml
bg:            "#0D1117"
fg:            "#E6EDF3"
fgMuted:       "#7D8590"
accent:        "#2F81F7"
accentAlt:     "#3DC9B0"
surface:       "#161B22"
surfaceAlt:    "#1F2937"
success:       "#3FB950"
warning:       "#D29922"
danger:        "#F85149"
info:          "#58A6FF"
borderFg:      "#30363D"
borderFocus:   "#2F81F7"
crumbBg:       "#161B22"
crumbFg:       "#7D8590"
crumbActiveBg: "#2F81F7"
crumbActiveFg: "#E6EDF3"
jsonKey:       "#2F81F7"
jsonValue:     "#3DC9B0"
jsonNumber:    "#D29922"
jsonBool:      "#F85149"
btnPrimaryBg:  "#238636"
btnPrimaryFg:  "#E6EDF3"
btnSecondaryBg:"#21262D"
btnSecondaryFg:"#7D8590"
```

---

## Acceptance Criteria

- **AC1**: `Skin` struct in `skin.go` maps all 25 semantic color roles via a flat `Colors` struct; `skin.Color(role)` returns the hex string (empty string = terminal default)
- **AC2**: `StyleSet` struct in `styles.go` contains all lipgloss styles used by the app; `buildStyleSet(skin *Skin) StyleSet` constructs it; model holds `m.styles StyleSet`; **zero** hardcoded `lipgloss.Color("...")` calls remain outside `buildStyleSet`
- **AC3**: `env.UIColor` field removed from `Environment` struct, `o8n-env.yaml.example`, and all rendering code; no references remain
- **AC4**: The context popup state is refactored to a generic `popup` struct in `model.go` with `mode popupMode` (None/Context/Skin); all 4 `root*` fields replaced
- **AC5**: `Ctrl+T` opens the skin picker popup using the generic popup; ↑↓ live-previews the skin; Esc reverts to skin before popup opened; Enter commits and saves to `o8n-stat.yml`
- **AC6**: `:` key still opens context/resource picker — behaviour unchanged, now using generic popup
- **AC7**: All 34 existing skin YAML files updated with the new semantic role fields (backward-compatible: old `o8n.body.fgColor` etc. fields are REMOVED; new flat `colors:` section used); `skins/o8n-cyber.yaml` added
- **AC8**: `stock.yaml` updated with all 25 roles; serves as canonical reference
- **AC9**: README documents `Ctrl+T` skin picker; specification.md documents full skin schema
- **AC10**: All existing tests pass; new tests cover popup mode switching, skin preview/revert/commit

---

## Tasks/Subtasks

### Task 1: Extend Skin struct + Color() helper
- [ ] 1.1 Rewrite `Skin` struct in `skin.go`: new flat `Colors` struct with 25 role fields (string, yaml tags)
- [ ] 1.2 Add `Color(role string) string` method — returns role value or `""` for terminal default
- [ ] 1.3 Add `loadSkin(name string) (*Skin, error)` — unchanged interface, reads new schema
- [ ] 1.4 Add backward-compat migration: if new `colors:` section is empty, attempt to map old `o8n.body.fgColor` etc. to roles (for skins not yet updated)
- [ ] 1.5 Write tests: `TestSkinColorLookup`, `TestSkinMissingRoleReturnsEmpty`, `TestLoadSkinFile`

### Task 2: StyleSet — replace all hardcoded styles
- [ ] 2.1 Define `StyleSet` struct in `styles.go` — one field per distinct style in the app (~30 styles)
- [ ] 2.2 Write `buildStyleSet(skin *Skin) StyleSet` — all `lipgloss.NewStyle()` calls with skin colors live here
- [ ] 2.3 Add `m.styles StyleSet` to model; call `buildStyleSet` in `newModelEnvApp` and `applyStyle`
- [ ] 2.4 Replace all hardcoded `lipgloss.Color("...")` in `view.go` with `m.styles.X`
- [ ] 2.5 Replace all hardcoded `lipgloss.Color("...")` in `model.go` `applyStyle()` with `buildStyleSet` calls
- [ ] 2.6 Remove all package-level `var` styles from `styles.go` (they become fields of `StyleSet`)
- [ ] 2.7 Write tests: `TestBuildStyleSetUsesAccentForBorderFocus`, `TestBuildStyleSetNonEmpty`

### Task 3: Remove env.UIColor
- [ ] 3.1 Delete `UIColor string` field from `Environment` struct in `internal/config/config.go`
- [ ] 3.2 Remove `ui_color` from `o8n-env.yaml.example`
- [ ] 3.3 Remove all `env.UIColor` references in `view.go` (~12 occurrences) — replace with `m.styles` lookups
- [ ] 3.4 Remove `color := ...` / `if env.UIColor` blocks from `view.go`
- [ ] 3.5 Verify `go vet ./...` and `go build ./...` clean
- [ ] 3.6 Update `internal/config/config_test.go` — remove any UIColor assertions

### Task 4: Generic popup widget
- [ ] 4.1 Add `popupMode` type + constants (`popupModeNone`, `popupModeContext`, `popupModeSkin`) to `model.go`
- [ ] 4.2 Add `popup struct { mode popupMode; input string; cursor int; items []string; title string; hint string; previewSkin string }` to model
- [ ] 4.3 Remove `showRootPopup bool`, `rootPopupCursor int`, `rootInput string` fields — replace usages with `m.popup.*`
- [ ] 4.4 Refactor `view.go` popup rendering to use `m.popup` — title line uses `m.popup.title`, hint uses `m.popup.hint`
- [ ] 4.5 Refactor `update.go` popup key handlers to use `m.popup` — `:` sets `popupModeContext` + populates `m.popup.items` from `m.rootContexts`
- [ ] 4.6 Helper `openPopup(mode popupMode) tea.Cmd` sets up the popup and loads items
- [ ] 4.7 Write tests: `TestPopupOpenContext`, `TestPopupOpenSkin`, `TestPopupEscCloses`, `TestPopupFilterItems`

### Task 5: Runtime skin switch (Ctrl+T) with live preview
- [ ] 5.1 Add `listSkinsCmd()` → returns `skinsLoadedMsg{names []string}` (reads `skins/` directory)
- [ ] 5.2 Handle `skinsLoadedMsg` in `update.go` — stores to `m.availableSkins []string`
- [ ] 5.3 On `Ctrl+T` key: call `openPopup(popupModeSkin)` — sets `m.popup.previewSkin = m.activeSkin`
- [ ] 5.4 On popup cursor move in skin mode: load skin, call `buildStyleSet`, apply to `m.styles` (live preview) — store original skin name in `m.popup.previewSkin` for revert
- [ ] 5.5 On Esc in skin mode: reload `m.popup.previewSkin` skin, rebuild styles, close popup
- [ ] 5.6 On Enter in skin mode: commit — set `m.activeSkin`, save to `o8n-stat.yml` via `m.saveStateCmd()`
- [ ] 5.7 Load available skins on app startup (call `listSkinsCmd` in `Init()`)
- [ ] 5.8 Write tests: `TestSkinPickerOpens`, `TestSkinPickerEscReverts`, `TestSkinPickerEnterCommits`

### Task 6: Skin YAML files — new schema
- [ ] 6.1 Write `skins/o8n-cyber.yaml` with all 25 role fields (GitHub dark palette)
- [ ] 6.2 Update `skins/stock.yaml` with all 25 roles
- [ ] 6.3 Update all 34 remaining skin YAML files: add `colors:` section with best-effort role mapping from existing palette
- [ ] 6.4 Verify all skins load without error: `go test ./internal/app/... -run TestLoadSkinFile`
- [ ] 6.5 Add `skins/transparent.yaml` compatibility (all roles = `""` = terminal default)

### Task 7: Tests, README, docs
- [ ] 7.1 Run full test suite; fix any regressions
- [ ] 7.2 Add `Ctrl+T` to README key bindings table; document skin picker UX
- [ ] 7.3 Update `specification.md`: skin schema, StyleSet, popup modes, skin switching
- [ ] 7.4 Add skin picker to hints bar in footer (visible hint `Ctrl+T:theme`)
- [ ] 7.5 Commit: `feat: config-driven skin system with generic command palette`

---

## Dev Agent Record

_To be filled during implementation_

### Decisions Made
- `popupModeContext` retains `:` key trigger — zero behavior change for existing users
- `Ctrl+T` chosen for skin picker: T = Theme, not occupied, discoverable via footer hint
- Live preview: preview on cursor move (not on every keypress during typing — only when cursor moves in the list)
- Skin YAML new schema: flat `colors:` map replaces nested `o8n.*` hierarchy; `loadSkin` maps old format for any YAML that doesn't have `colors:` key yet (Task 1.4)
- `buildStyleSet` is the ONLY place that calls `lipgloss.NewStyle().Foreground/Background(...)` with color values

### Files Changed
_To be listed per task_

---

## Story DoD Checklist

- [ ] All ACs verified
- [ ] All tasks/subtasks marked [x]
- [ ] All tests pass (go test ./...)
- [ ] go vet clean
- [ ] No hardcoded color literals in production code outside `buildStyleSet`
- [ ] README updated
- [ ] specification.md updated
- [ ] Committed
