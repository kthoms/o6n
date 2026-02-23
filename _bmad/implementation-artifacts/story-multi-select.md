# Story C: Multi-Select & Batch Actions

## Summary
Add Space-to-toggle row selection, Ctrl+A select-all, full-width selection highlight
using a distinct skin role, selection count in the footer, and batch action dispatch
through the action registry (depends on Story B).

**Depends on:** Story B (action registry, `ActionDef.Scope`, `executeActionCmd`).

---

## Background / Context

Power users need to terminate/suspend/delete multiple rows at once without repeating
the operation per row. The selection model should feel familiar (Space in ranger/
lazygit/k9s), be visually distinct from the cursor row, and integrate cleanly with
Story B's action registry — actions with `scope: multi` are available when a selection
is active.

---

## Tasks

### T1 — Selection state on model
- [ ] `internal/app/model.go` — Add:
  ```go
  selectedRows map[int]bool  // keyed by table row index; nil = no selection
  ```
  Initialise to `nil` (no allocation until first Space press).
- [ ] Selection is reset (`selectedRows = nil`) whenever:
  - A new data load completes (`genericLoadedMsg` handler)
  - Navigation stack is pushed (drilldown) or popped (Esc back)
  - Root context is switched

### T2 — Space toggles selection
- [ ] `internal/app/update.go` — `case " ":` (space):
  - Guard: only when `popup.mode == popupModeNone`, `activeModal == ModalNone`,
    `searchMode == false`, table has rows.
  - If `selectedRows == nil` → initialise `selectedRows = make(map[int]bool)`.
  - Toggle `selectedRows[m.table.Cursor()]`.
  - If the map becomes empty after toggle → set `selectedRows = nil`.
  - Move cursor down by 1 after toggle (ergonomic: keeps selection flow).
- [ ] Test: space on row 0 adds it to selectedRows; space again removes it.
- [ ] Test: selectedRows is nil after toggling the only selected row off.
- [ ] Test: cursor advances after space.

### T3 — Ctrl+A selects / deselects all
- [ ] `internal/app/update.go` — `case "ctrl+a":`:
  - Guard: popup closed, no modal.
  - If all visible rows are selected → clear (`selectedRows = nil`).
  - Else → select all (`selectedRows[i] = true` for i in 0..len(rows)-1).
- [ ] Test: ctrl+a on 3 rows produces selectedRows with 3 entries.
- [ ] Test: ctrl+a again when all selected clears selectedRows to nil.

### T4 — Full-width selection highlight (skin role)
- [ ] `internal/app/skin.go` — Add `rowSelected` role to the `Colors` struct and to all
  skin YAML files (default: slightly lighter background than cursor, distinct from
  `borderFocus`).
- [ ] `internal/app/styles.go` — Add `RowSelected lipgloss.Style` to `StyleSet`;
  build from `rowSelected` skin role in `buildStyleSet`.
- [ ] `internal/app/table.go` — In `colorizeRows` (or a post-pass in `setTableRowsSorted`),
  prepend a full-width selection indicator to rows in `selectedRows`:
  - Use a left-margin marker character (e.g. `✓ `) styled with `RowSelected`.
  - Or: apply `RowSelected` background to the entire rendered row string.
  - The highlight must span full pane width (not just cell content).
- [ ] `internal/app/model.go` — `applyStyle()` already sets table `Selected` style;
  leave cursor highlight unchanged. Selection highlight is a separate visual layer.
- [ ] Test: `colorizeRows` output for a selected row index contains the selection marker.

### T5 — Selection count in footer
- [ ] `internal/app/view.go` — When `len(selectedRows) > 0`, append
  `[N selected]` (styled with `AccentPrimary`) to the footer hint row 2.
- [ ] Test: view output contains `[2 selected]` when two rows are selected.

### T6 — Batch action dispatch via Ctrl+X
- [ ] `internal/app/update.go` — When `popupModeAction` Enter is triggered with a
  `scope: multi` action:
  - Collect IDs: for each index in `selectedRows`, look up `rowData[idx]["id"]`.
  - Pass ID slice to `executeActionCmd(action, ids)`.
  - Clear `selectedRows = nil` after dispatch.
- [ ] When `scope: single` and selection is active, apply to cursor row only (ignore
  selection).
- [ ] Test: batch dispatch collects correct IDs from selectedRows.
- [ ] Test: selectedRows cleared after dispatch.

### T7 — Clear selection on data events
- [ ] `internal/app/update.go` — In `genericLoadedMsg` handler: `m.selectedRows = nil`.
- [ ] In navigation push (drilldown): `m.selectedRows = nil`.
- [ ] In navigation pop (Esc): `m.selectedRows = nil`.
- [ ] In context switch: `m.selectedRows = nil`.
- [ ] Test: selection cleared after genericLoadedMsg.

---

## Acceptance Criteria
- AC1: Space toggles row selection; cursor advances.
- AC2: Ctrl+A selects all rows; second Ctrl+A deselects all.
- AC3: Selected rows have a visually distinct full-width highlight.
- AC4: Footer shows `[N selected]` count when selection is active.
- AC5: Ctrl+X with multi-scope action dispatches to all selected rows.
- AC6: Selection is cleared on data load, drilldown, and context switch.

---

## Dev Agent Record
_To be filled during implementation._

## File List
_To be filled during implementation._
