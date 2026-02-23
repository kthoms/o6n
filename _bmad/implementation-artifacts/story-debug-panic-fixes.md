# Story: Debug, Panic Recovery & Context Switch Fixes

## Status
review

## Summary
Five targeted fixes: (1) context-switch popup highlight foreground-only, (2) "engines"
undefined-table shows stale data, (3) errors not logged with stack + REST URL,
(4) index-out-of-bounds panics crash the app, (5) debug output needs command logging
and per-error named screen dumps linked in the log.

---

## Acceptance Criteria

- AC1: Selected row in context-switch popup shows foreground highlight (arrow + text); no background fill.
- AC2: Switching to a context with no TableDef in config shows an error footer message and clears the table; no stale rows.
- AC3: Every `errMsg` is written to `debug/o8n.log` with: error text, REST URL that caused it, and stack trace.
- AC4: A runtime panic in `Update()` is caught by a `recover()`, logs the stack to `debug/o8n.log`, and shows a footer error; the app does not exit.
- AC5: On `errMsg` a unique screen dump is saved to `debug/screen-<id>.txt` and the filename is referenced in the log entry.
- AC6: When `--debug` is active, each significant command dispatch is logged (fetch, terminate, delete, save, resize).

---

## Tasks

### T1 — Popup foreground highlight
- [x] `internal/app/view.go` — In the popup item loop, for the selected row (`globalIdx == m.popup.cursor`),
  render the full line (`cursor + " " + rc`) with `lipgloss.NewStyle().Foreground(col(m.skin, "borderFocus"))`.
  No background color change.
- [x] Test: popup selected row render contains ANSI foreground code; non-selected rows do not.

### T2 — Clear stale rows on error + undefined-table guard
- [x] `internal/app/update.go` — In `case errMsg:` handler: call `m.table.SetRows([]table.Row{})` to clear
  stale data before setting the footer error.
- [x] `internal/app/table.go` — `loadRootContexts`: after building the `roots` slice, filter to only
  include entries that have a matching `TableDef` in `m.config` (pass config as parameter, or filter
  in `model.go` after loading). Keep a `config *config.AppConfig` parameter. Alternatively, filter
  in `model.go` after `loadRootContexts` returns.
- [x] `internal/app/model.go` — After `m.rootContexts = loadRootContexts(...)`, filter:
  ```go
  filtered := m.rootContexts[:0]
  for _, rc := range m.rootContexts {
      if m.findTableDef(rc) != nil {
          filtered = append(filtered, rc)
      }
  }
  m.rootContexts = filtered
  ```
- [x] Test: `errMsg` clears table rows.
- [x] Test: `rootContexts` after model init contains only table names present in config.

### T3 — Error logging: stack trace + REST URL
- [x] `internal/app/commands.go` — In `fetchGenericCmd`, wrap all `errMsg{err}` returns with the URL
  included: `errMsg{fmt.Errorf("GET %s: %w", urlStr, err)}`.
- [x] `internal/app/update.go` — In `case errMsg:` handler, always log (not just debug mode).
  In debug mode, include full stack trace.
- [x] Test: `errMsg` handler logs error text (redirect `log` output to buffer in test).

### T4 — Panic recovery in Update()
- [x] `internal/app/update.go` — Named return + defer/recover at top of Update().
  Logs stack, saves screen dump, sets footer error without crashing.
- [x] `internal/app/view.go` — Package-level `lastRenderedView string`; set in View().
  Replaced `debugCh` goroutine.
- [x] `internal/app/model.go` — Removed `debugCh chan string` field and goroutine.
- [x] Test: Update handles common messages without panic.

### T5 — Debug command logging + named screen dumps on error
- [x] `internal/app/update.go` — `errMsg` saves `debug/screen-<id>.txt`, logs filename.
- [x] `internal/app/commands.go` — `fetchForRoot` logs `[cmd] fetch root=X` when debugEnabled.
- [x] `internal/app/update.go` — `WindowSizeMsg` logs `[resize] WxH` when debugEnabled.
- [x] `internal/app/view.go` — Writes `debug/last-screen.txt` directly (no goroutine).
- [x] Test: `errMsg` creates `debug/screen-*.txt` file.

---

## Dev Notes

- `runtime/debug` import: `"runtime/debug"` — `debug.Stack()` returns `[]byte`.
- Package-level `var lastRenderedView string` in `view.go` avoids value receiver mutation problem.
- The `Update()` recover defer: since `Update` returns `(tea.Model, tea.Cmd)` and uses value receiver,
  use named returns to allow defer to mutate the return value.
- `loadRootContexts` filtering: keep the function signature unchanged; do the filtering in `model.go`
  after the call, using `m.findTableDef(rc) != nil`.
- Log always (both debug and non-debug mode) for errors — this is the user's primary complaint.
- Screen dump ID: `fmt.Sprintf("%d", time.Now().UnixNano())` — monotonically increasing, unique per process.

---

## Dev Agent Record

### Implementation Notes
- T1: popup selected row styled with `lipgloss.NewStyle().Foreground(col(m.skin, "borderFocus")).Bold(true)` — foreground only, no background fill.
- T2: `errMsg` handler calls `m.table.SetRows([]table.Row{})` before setting footer. `rootContexts` filtered in `newModel()` to only include entries with a matching `TableDef`. This eliminates "engines", "batchs", "authorizations" etc. that come from the OpenAPI spec but have no config.
- T3: `fetchGenericCmd` wraps all errors with `fmt.Errorf("GET %s: %w", urlStr, err)`. `errMsg` handler always logs via `log.Printf`; stack trace added when `debugEnabled`.
- T4: `Update()` changed to named returns `(retModel tea.Model, retCmd tea.Cmd)`. `defer/recover` at top captures panics, logs stack + screen dump, sets footer error — app never exits on panic.
- T5: `debugCh chan string` and goroutine removed. Package-level `lastRenderedView string` in `view.go` set on every `View()` call. `errMsg` saves `debug/screen-<nanotime>.txt`. `fetchForRoot` logs `[cmd] fetch root=X`. `WindowSizeMsg` logs `[resize] WxH`. `view.go` writes `debug/last-screen.txt` directly when debugEnabled.

## File List
- `internal/app/view.go` — popup foreground highlight; `lastRenderedView` pkg var; direct last-screen.txt write; removed `debugCh` send
- `internal/app/update.go` — named returns + recover; errMsg clears rows, logs error, saves screen dump; resize logging; added `os`, `runtime/debug` imports
- `internal/app/model.go` — removed `debugCh` field and goroutine; rootContexts filtered to configured tables only
- `internal/app/commands.go` — fetchForRoot debug log; fetchGenericCmd errors include URL; added `log` import
- `internal/app/debug_fixes_test.go` — new test file covering all 5 tasks
