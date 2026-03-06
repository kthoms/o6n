# Story: Drilldown and Debug Fixes

## Status

done

## Summary

Four targeted bug fixes: (1) drilldown from `process-definition` opens the wrong table because the `target:` field is missing from every `drilldown:` block in `o6n-cfg.yaml`, causing `findTableDef("")` to match the first config table (`external-task`) and producing a blank or incorrect view; (2) REST calls are not logged in debug mode because `fetchGenericCmd` builds the URL inside a goroutine without emitting a log line; (3) screen captures in `debug/` are unreadable because ANSI escape sequences are written verbatim instead of being stripped; (4) context-switching to a table with a different column count panics because `SetColumns` is called while stale rows remain, causing `renderRow` to access an out-of-bounds index.

## Motivation

### Bug 1: Wrong table shown after drilldown from process-definition

`nav.go:executeDrilldown` uses `d.Target` as the key for both `findTableDef` and `fetchGenericCmd`. Every `drilldown:` block in `o6n-cfg.yaml` is missing the `target:` YAML field:

- `process-definition` drilldown → should `target: process-instance`
- `process-instance` drilldown → should `target: process-variables`
- `job-definition` drilldown → should `target: job`
- `deployment` drilldown → should `target: process-definition`

When `target` is `""`, `findTableDef("")` hits the last-resort prefix-match with `base=""`, which matches every table and returns `&config.Tables[0]` — the `external-task` table. The screen then shows TOPICNAME/WORKERID columns with "No results" and the API call goes to `base/` (empty path) which errors or returns unexpected data.

### Bug 2: REST calls not logged

`fetchGenericCmd` (`commands.go:322`) builds `urlStr` inside the goroutine closure but never logs it. The debug log only records the root name (`[cmd] fetch root=...`) but not the actual URL being fetched. When debugging API issues there is no way to see what URL was called or what response was received.

### Bug 3: Screen captures contain ANSI escape sequences

`view.go:870` writes `baseView` directly to `debug/last-screen.txt`. `update.go:27` and `update.go:1743` write `lastRenderedView` directly to panic/error screen files. All three writes include raw ANSI escape sequences (colors, bold, box-drawing escape codes) making the files unreadable in a text editor or `cat`.

### Bug 4: Context switch panics when column count changes

`update.go:1009–1016` (context switch handler) calls `m.table.SetColumns(cols)` before clearing the existing rows. `SetColumns` internally calls `UpdateViewport` → `renderRow(cursor)` which iterates over the current rows' cells and accesses `m.cols[i]`. If a row has more cells than the new column count, this panics with `index out of range [N] with length N` (and symmetrically for fewer cells).

Observed: switching to `batchs` (4 columns) while table held 3-cell rows from a previous view caused `panic: runtime error: index out of range [3] with length 3` at `update.go:1013`.

The correct fix is a three-step sequence: (1) `SetRows([]table.Row{})` clears rows — UpdateViewport loop runs 0 iterations, safe; (2) `SetColumns(cols)` sets new column count — UpdateViewport loop still runs 0 iterations, safe; (3) `SetRows(normalizeRows(nil, len(cols)))` sets placeholder rows matching the new column count — safe because cells and columns match.

## Acceptance Criteria

### Fix 1: Drilldown target fields

- [x] **AC-1:** `o6n-cfg.yaml` — `process-definition` drilldown gains `target: process-instance`. Pressing Enter on a process-definition row navigates to and fetches the process-instance table filtered by `processDefinitionId`.
- [x] **AC-2:** `o6n-cfg.yaml` — `process-instance` drilldown gains `target: process-variables`. Pressing Enter on a process-instance row navigates to and fetches the `process-variables` table filtered by `processInstanceId`.
- [x] **AC-3:** `o6n-cfg.yaml` — `job-definition` drilldown gains `target: job`. Pressing Enter on a job-definition row navigates to the job table filtered by `jobDefinitionId`.
- [x] **AC-4:** `o6n-cfg.yaml` — `deployment` drilldown gains `target: process-definition`. Pressing Enter on a deployment row navigates to the process-definition table filtered by `deploymentId`.
- [x] **AC-5:** A test verifies that after dispatching Enter on a `process-definition` table, the model's `currentRoot` is `"process-instance"` and `genericParams` contains the expected `processDefinitionId` key.
- [x] **AC-6:** A test verifies that after dispatching Enter on a `process-instance` table, the model's `currentRoot` is `"process-variables"` and `genericParams` contains the expected `processInstanceId` key.

### Fix 2: HTTP request debug logging

- [x] **AC-7:** In `commands.go:fetchGenericCmd`, when `m.debugEnabled` is true, a `log.Printf("[http] GET %s", urlStr)` line is emitted just before `http.DefaultClient.Do(req)`.
- [x] **AC-8:** The same guard applies to the count sub-request: `log.Printf("[http] GET %s (count)", countURL)` emitted before the count `http.DefaultClient.Do`.
- [x] **AC-9:** No log output is emitted when `debugEnabled` is false (production default).

### Fix 3: ANSI-free screen captures

- [x] **AC-10:** `view.go:870` — `debug/last-screen.txt` is written with `ansi.Strip(baseView)` instead of `baseView`.
- [x] **AC-11:** `update.go:27` (panic handler) — panic screen file is written with `ansi.Strip(lastRenderedView)`.
- [x] **AC-12:** `update.go:1743` (error handler) — error screen file is written with `ansi.Strip(lastRenderedView)`.
- [x] **AC-13:** A test verifies that the string written to the screen file does not contain the ESC character (`\x1b`) after a simulated panic or error event.

### Fix 4: Context switch column-count panic

- [x] **AC-14:** `update.go:1013–1016` — corrected to three-step sequence: `SetRows([]table.Row{})` → `SetColumns(cols)` → `SetRows(normalizeRows(nil, len(cols)))`. This is safe in both directions (wide→narrow and narrow→wide) because rows are cleared before column count changes.
- [x] **AC-15:** A test sets up a model with rows having N cells, simulates a context switch to a table with M ≠ N columns, and asserts no panic and that the resulting table has M columns.

## Tasks

### Task 1: Add missing drilldown `target:` fields in config

**Files:** `o6n-cfg.yaml`

1. `process-definition` drilldown block: add `target: process-instance`
2. `process-instance` drilldown block: add `target: process-variables`
3. `job-definition` drilldown block: add `target: job`
4. `deployment` drilldown block: add `target: process-definition`

Run `wc -l o6n-cfg.yaml` before and after — line count may increase by 4. Verify no deletions with `git diff o6n-cfg.yaml`.

### Task 2: HTTP debug logging in fetchGenericCmd

**Files:** `internal/app/commands.go`

1. After `urlStr` is fully assembled (after paging params), add:
   ```go
   if m.debugEnabled {
       log.Printf("[http] GET %s", urlStr)
   }
   ```
2. After `countURL` is fully assembled, add:
   ```go
   if m.debugEnabled {
       log.Printf("[http] GET %s (count)", countURL)
   }
   ```

### Task 3: Strip ANSI from screen captures

**Files:** `internal/app/view.go`, `internal/app/update.go`

1. `view.go:870` — change `os.WriteFile(..., []byte(baseView), ...)` to `os.WriteFile(..., []byte(ansi.Strip(baseView)), ...)`; `ansi` package is already imported.
2. `update.go:27` — change `[]byte(lastRenderedView)` to `[]byte(ansi.Strip(lastRenderedView))`; import `"github.com/charmbracelet/x/ansi"` if not already present.
3. `update.go:1743` — same change as above.

### Task 4: Fix SetColumns/SetRows order in context switch

**Files:** `internal/app/update.go`

At the context switch handler (~line 1013–1016), apply the three-step clear-first sequence:
```go
m.table.SetRows([]table.Row{})             // 1. clear rows — UpdateViewport loop runs 0 iterations
m.table.SetColumns(cols)                    // 2. set new column count — still 0 rows, safe
m.table.SetRows(normalizeRows(nil, len(cols))) // 3. set placeholder rows matching new cell count
```

This replaces the original two-step code (SetColumns then SetRows) plus removes the duplicate second SetRows call. The key insight: `renderRow` in `bubbles/table@v1.0.0` iterates over row cells and indexes into `m.cols`, so any mismatch between row cell count and column count causes a panic — in either direction (wider→narrower or narrower→wider). Clearing rows first avoids all cross-size rendering.

### Task 5: Tests

**Files:** `internal/app/drilldown_debug_fixes_test.go` (new)

1. **AC-5** — Drilldown from process-definition: build model with `process-definition` table rows, dispatch Enter key, assert `m.currentRoot == "process-instance"` and `m.genericParams["processDefinitionId"] != ""`.
2. **AC-6** — Drilldown from process-instance: build model with `process-instance` table rows, dispatch Enter key, assert `m.currentRoot == "process-variables"` and `m.genericParams["processInstanceId"] != ""`.
3. **AC-13** — Screen capture is ANSI-free: after an `errMsg` is dispatched, read the written screen file and assert `!strings.Contains(content, "\x1b")`.
4. **AC-15** — Context switch no panic: set up model with 3-cell rows, trigger context switch to a table with 4 columns, assert no panic and `len(m.table.Columns()) == 4`.

## Dev Notes

- `ansi` package (`github.com/charmbracelet/x/ansi`) is already imported in `view.go` and `nav.go` — check import in `update.go`.
- The `lastRenderedView` variable is package-level in `view.go`, accessible from `update.go` without qualification.
- `findTableDef("")` returning `&config.Tables[0]` is a by-product of the last-resort prefix-match loop. Fixing the missing `target:` fields is the correct fix; the prefix-match fallback is intentional for singular/plural name aliasing.
- After adding `target:` to the drilldown blocks, the `switch d.Target` in `executeDrilldown` will correctly set `m.selectedDefinitionKey` for `"process-instance"` and `m.selectedInstanceID` for `"process-variables"`, which are used by the edit/save flows.
- The `renderRow` function in `bubbles/table@v1.0.0` iterates over `m.rows[r]` cells and indexes `m.cols[i]`, so it panics whenever row cell count > column count. The three-step clear pattern handles both wide→narrow and narrow→wide switches safely.

## File List

- `o6n-cfg.yaml` — added `target:` to four drilldown blocks
- `internal/app/commands.go` — added HTTP debug log lines
- `internal/app/view.go` — strip ANSI from last-screen.txt write
- `internal/app/update.go` — strip ANSI from panic/error screen writes; fix SetColumns/SetRows order with clear-first pattern
- `internal/app/drilldown_debug_fixes_test.go` (new) — tests for AC-5, AC-6, AC-13, AC-15

## Dev Agent Record

Implemented by dev agent 2026-03-02. Code-reviewed and fixed by dev agent 2026-03-02.

All 5 tasks completed. Review fixed 3 issues:
- H-1: Added AC-7/AC-8/AC-9 unit tests (`TestFetchGenericCmd_DebugLogsHTTP`, `TestFetchGenericCmd_NoDebugNoLog`) — HTTP debug logging now has full test coverage.
- M-1: Updated Task 4 description in story body to show actual three-step implementation.
- M-2: Added `defer` restore of `lastRenderedView` in `TestScreenDumpIsAnsiFree` to prevent test contamination.

Key insight on Task 4: `bubbles/table@v1.0.0:renderRow` iterates row cells and indexes `m.cols[i]` — ANY mismatch panics. The three-step clear-first pattern is safe in both directions (wider→narrower and narrower→wider).

Full test suite passes: `go test ./...` — all packages green.

## Change Log

- 2026-03-02: Story written from user-reported bugs with screen + log evidence.
- 2026-03-02: All tasks implemented and tested. Status → review.
- 2026-03-02: Code review completed. Fixed H-1 (missing AC-7/8/9 tests), M-1 (Task 4 doc), M-2 (test isolation). Status → done.
