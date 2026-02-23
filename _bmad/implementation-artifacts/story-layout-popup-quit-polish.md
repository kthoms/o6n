# Story A: Layout Polish, Popup Scroll Unification & Quit Ghost Fix

## Summary
Three independent but related polish items: (1) unify the skin-picker popup to use the
same generic scroll mechanism as context/search popups, (2) fix content pane not
resizing when the search popup closes, and (3) fix the quit confirmation dialog
remaining visible after the app has already issued `tea.Quit`.

---

## Background / Context

- Popup scroll offset (`popup.offset`) was introduced for context/search modes but the
  skin picker still uses its own linear walk and ignores `popup.offset`.
- `computePaneHeight()` subtracts popup height when the popup opens but nothing
  recomputes pane height when the search popup (popupModeSearch) closes — the table
  stays at its reduced height.
- `contextPopupHeight()` still contains stale "+1 for '…N more' line" logic that was
  removed from the view. Also doesn't handle `popupModeSearch` (counts rootContexts
  instead of filtered table rows).
- `tea.Quit` is returned from the quit-confirm dialog handler but `m.quitting` is never
  set, so `View()` renders one more frame containing the dialog overlay.

---

## Tasks

### T1 — Generic popup scroll for skin picker
- [ ] `internal/app/update.go` — Replace skin-picker `up`/`down` handlers with the
  same `offset`/`cursor` logic used by context/search mode.
  - `down`: if `cursor >= offset+maxShow` → `offset++`
  - `up`: if `cursor < offset` → `offset--`
  - Reset `popup.offset = 0` when skin popup opens (`Ctrl+T`).
- [ ] `internal/app/view.go` — Skin popup list rendering already uses the generic
  scroll window (`offset` → `offset+maxShow`). Verify no special-casing remains.
- [ ] Test: pressing `down` past the 8th skin increments `popup.offset`.
- [ ] Test: pressing `up` when `cursor == offset` decrements `popup.offset`.

### T2 — Content pane resizes when search popup closes
- [ ] `internal/app/nav.go` — Fix `contextPopupHeight()`:
  - Remove the `+1` for the removed "…N more" line; `shown = min(matchCount, maxShow)`.
  - For `popupModeSearch`: count visible filtered table rows (capped at `maxShow`) as
    the list height, not `rootContexts`.
- [ ] `internal/app/update.go` — Add `m.paneHeight = m.computePaneHeight();
  m.table.SetHeight(m.paneHeight - 1)` to every code path that closes the search popup:
  - Esc handler (`popupModeSearch` block)
  - Enter handler (`popupModeSearch` block)
- [ ] Test: after closing search popup (Esc), `m.paneHeight` equals `computePaneHeight()`
  called with `popup.mode == popupModeNone`.
- [ ] Test: after closing search popup (Enter/lock), same assertion.

### T3 — Quit dialog ghost fix
- [ ] `internal/app/model.go` — Add `quitting bool` field.
- [ ] `internal/app/update.go` — Set `m.quitting = true` immediately before every
  `return m, tea.Quit` call.
- [ ] `internal/app/view.go` — At the top of `View()`, if `m.quitting` return `""`.
- [ ] Test: after sending confirm-quit key (`ctrl+c` on quit dialog), model has
  `quitting == true`.
- [ ] Test: `View()` on a model with `quitting == true` returns empty string.

### T4 — Full-width selected row highlight
- [ ] `internal/app/model.go` — In `applyStyle()`, set the table `Selected` style width
  to `m.paneWidth` so the highlight spans the full pane, not just the content cells.
  Use `tStyles.Selected = tStyles.Selected.Width(m.paneWidth)`.
- [ ] `internal/app/update.go` / `tea.WindowSizeMsg` handler — After recomputing
  `m.paneWidth`, call `m.applyStyle()` to propagate the new width into the selection
  style.
- [ ] Test: selected style width equals `m.paneWidth` after a window resize message.

---

## Acceptance Criteria
- AC1: Skin popup scrolls identically to context popup — offset tracks cursor.
- AC2: After closing search popup, table height equals full available pane height.
- AC3: `View()` returns empty string when quitting — no ghost dialog.
- AC4: Selected row highlight spans the full pane width.

---

## Dev Agent Record
_To be filled during implementation._

## File List
_To be filled during implementation._
