# Story B: Action Registry & Contextual Actions

## Summary
Replace hardcoded global keybindings (e.g. `Ctrl+D` terminate visible everywhere)
with a config-driven action registry. Each table declares its available actions with
an optional `enabled_when` expression evaluated against live row data. The footer
hints and a new `Ctrl+X` action palette show only applicable, correctly-enabled actions.

---

## Background / Context

Currently `Ctrl+D` is bound globally as "terminate instance" regardless of which table
is active. This causes:
- "Terminate" appearing on process-definition and other views where it makes no sense.
- No way to declare "delete definition only when instanceCount == 0".
- No visual feedback distinguishing enabled vs disabled actions.

The action registry externalises all of this into `o8n-cfg.yaml`. Each `TableDef`
gains an `actions` list. The UI reads these at runtime and renders/dispatches
accordingly.

---

## Config Schema Addition

```yaml
# o8n-cfg.yaml — per-table actions block
- name: process-definition
  actions:
    - key: "ctrl+d"
      label: "Delete"
      scope: single          # single | multi
      confirm: true
      enabled_when: "instanceCount == 0"   # optional; empty = always enabled
      api_action:
        method: DELETE
        path: /process-definition/{id}

- name: process-instance
  actions:
    - key: "ctrl+d"
      label: "Terminate"
      scope: multi
      confirm: true
      enabled_when: "state != COMPLETED"
      api_action:
        method: DELETE
        path: /process-instance/{id}
    - key: "ctrl+s"
      label: "Suspend"
      scope: multi
      confirm: false
      enabled_when: "state == ACTIVE"
      api_action:
        method: PUT
        path: /process-instance/{id}/suspended
        body: '{"suspended": true}'
```

---

## Tasks

### T1 — Config structs for ActionDef
- [ ] `internal/config/config.go` — Add structs:
  ```go
  type ActionDef struct {
      Key         string       `yaml:"key"`
      Label       string       `yaml:"label"`
      Scope       string       `yaml:"scope"`        // "single" | "multi"
      Confirm     bool         `yaml:"confirm"`
      EnabledWhen string       `yaml:"enabled_when,omitempty"`
      ApiAction   ApiActionDef `yaml:"api_action"`
  }
  type ApiActionDef struct {
      Method string `yaml:"method"`
      Path   string `yaml:"path"`
      Body   string `yaml:"body,omitempty"`
  }
  ```
  Add `Actions []ActionDef` to `TableDef`.
- [ ] `o8n-cfg.yaml` — Add action declarations for:
  - `process-definition`: Delete (ctrl+d, enabled_when instanceCount == 0)
  - `process-instance`: Terminate (ctrl+d), Suspend (ctrl+s), Activate (ctrl+a — no, Ctrl+A is select-all; use ctrl+r for resume)
  - `job`: Suspend, Activate, Delete
  - `job-definition`: Suspend, Activate, Delete
- [ ] Test: config round-trip parses `enabled_when` and `api_action` fields.

### T2 — enabled_when expression evaluator
- [ ] `internal/app/action.go` (new file) — Implement `EvalEnabledWhen(expr string, row map[string]interface{}) bool`:
  - Empty expr → `true`.
  - Parse: `<field> <op> <value>` where op ∈ `{==, !=, >, <, >=, <=}`.
  - Type coercion: if both sides parse as float64, use numeric comparison; else string comparison.
  - Unknown field → `false` (safe default: disabled).
- [ ] Test suite for evaluator covering: always-enabled (empty), string equality, string inequality, numeric `==`, `!=`, `>`, `<`, `>=`, `<=`, unknown field.

### T3 — Remove hardcoded Ctrl+D global terminate
- [ ] `internal/app/update.go` — Remove the existing `case "ctrl+d":` global terminate
  handler. Replace with config-driven dispatch (T5 below).
- [ ] Ensure existing `ModalConfirmDelete` confirm flow is reachable from new dispatch.
- [ ] Test: pressing `ctrl+d` on a view with no actions defined does nothing.

### T4 — Ctrl+X action palette popup
- [ ] `internal/app/model.go` — Add `popupModeAction popupMode = 4`.
- [ ] `internal/app/update.go` — `case "ctrl+x":` opens popup in `popupModeAction`:
  - Items: `[key] Label` for each action in current table's `Actions` that matches
    `scope: single` (multi-select not active) or `scope: multi` (when selection active).
  - Disabled actions listed with `(disabled)` suffix and dimmed.
  - Cursor starts at first enabled action.
- [ ] `internal/app/view.go` — Render action palette: title = `"actions"`, list shows
  enabled actions normally, disabled actions with `Subtle` skin style.
- [ ] Esc closes palette.
- [ ] Enter dispatches selected action (T5).
- [ ] Test: `ctrl+x` in `process-instance` view opens popup with Terminate and Suspend.
- [ ] Test: Terminate shows as disabled when `state == COMPLETED` row is selected.

### T5 — Config-driven action dispatch
- [ ] `internal/app/commands.go` — Add `executeActionCmd(action ActionDef, ids []string) tea.Cmd`:
  - Performs `method` request to `path` (substituting `{id}`).
  - For multi-ID: fires N requests concurrently, collects errors.
  - Returns `actionResultMsg{succeeded int, failed int, errors []string}`.
- [ ] `internal/app/update.go` — Handle `actionResultMsg`:
  - Show success count in footer.
  - Show first error if any failures.
  - Trigger re-fetch of current table.
- [ ] Test: `executeActionCmd` builds correct URL for path `"/process-instance/{id}"` with id `"abc-123"`.

### T6 — Contextual footer hints
- [ ] `internal/app/view.go` — Replace hardcoded footer hint entries for `ctrl+d` etc.
  with dynamically built hints from `currentTable.Actions`, filtered by:
  - `scope == single` (or multi when selection active)
  - Enabled state (`EvalEnabledWhen`)
  - Visibility priority already in hint system
- [ ] Disabled actions: render dimmed (use `Subtle` role).
- [ ] When pressing a disabled action key directly: flash footer message
  `"<Label> not available: <enabled_when field> = <value>"`.
- [ ] Test: hint bar for `process-definition` shows Delete (dimmed if instanceCount > 0).
- [ ] Test: hint bar for `process-instance` does NOT show Delete (wrong table).

---

## Acceptance Criteria
- AC1: `enabled_when` evaluator handles all 6 operators, string and numeric.
- AC2: No hardcoded global `ctrl+d` terminate — fully config-driven.
- AC3: `Ctrl+X` opens action palette showing only current table's actions.
- AC4: Disabled actions appear dimmed; pressing their key shows explanatory footer msg.
- AC5: Configured actions are dispatched via API and result shown in footer.
- AC6: Footer hints show only current table's actions.

---

## Dev Agent Record
_To be filled during implementation._

## File List
_To be filled during implementation._
