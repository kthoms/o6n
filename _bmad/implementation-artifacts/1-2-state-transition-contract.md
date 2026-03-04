# Story 1.2: State Transition Contract

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **developer contributing to o8n**,
I want a `TransitionType` enum and `prepareStateTransition(t TransitionType)` function enforced across all navigation paths in `internal/app/nav.go` and `internal/app/update.go`,
So that every navigation action uses a single, auditable gate that eliminates state leakage bugs.

## Acceptance Criteria

1. **Given** `TransitionType` is defined with constants: `TransitionFull`, `TransitionDrillDown`, `TransitionPop`
   **When** any navigation action is taken (environment switch, context switch, drill-down, Esc, breadcrumb jump)
   **Then** `prepareStateTransition(transitionType)` is called before any view state is modified
   **And** `TransitionFull` clears: `activeModal`, `footerError`, `searchQuery`, `searchActive`, `sortColumn`, `sortDirection`, `tableCursor`, `navigationStack`
   **And** `TransitionDrillDown` pushes current `viewState` onto `navigationStack`, then clears non-stack state fields
   **And** `TransitionPop` pops `viewState` from `navigationStack` and restores all captured state — no field clearing
   **And** a code review of all nav paths confirms no navigation code bypasses `prepareStateTransition`

2. **Given** `internal/app/main_transition_test.go` is created with tests covering all three `TransitionType` cases
   **When** `make cover` is run
   **Then** line coverage on the transition logic in `transition.go` is ≥ 80%

## Tasks / Subtasks

- [ ] Task 1: Define exported `TransitionType` enum in `transition.go` (AC: 1)
  - [ ] Add `type TransitionType int` with constants `TransitionFull`, `TransitionDrillDown`, `TransitionPop` (exported, PascalCase)
  - [ ] Remove (or deprecate) the unexported `transitionScope` type and its 5 constants after all callers are migrated

- [ ] Task 2: Refactor `prepareStateTransition` to accept `TransitionType` (AC: 1)
  - [ ] Change signature from `(scope transitionScope, depth ...int)` to `(t TransitionType)` (pointer receiver stays)
  - [ ] `TransitionFull` implementation: clears `activeModal`, `footerError`, `searchTerm`, `searchMode`, `sortColumn`, `sortAscending`, `searchInput.Blur()`, `originalRows`, `filteredRows`, `navigationStack`, `table.SetCursor(0)`, `genericParams`, `selectedDefinitionKey`, `selectedInstanceID`, popup state reset
  - [ ] `TransitionDrillDown` implementation: (a) capture current viewState snapshot BEFORE clearing, (b) push snapshot to `navigationStack`, (c) clear `activeModal`, `footerError`, `searchTerm`, `searchMode`, `sortColumn`, `sortAscending`, `searchInput.Blur()`, `originalRows`, `filteredRows`, popup reset
  - [ ] `TransitionPop` implementation: (a) pop last entry from `navigationStack`, (b) restore all viewState fields to model (viewMode, breadcrumb, contentHeader, selectedDefinitionKey, selectedInstanceID, cachedDefinitions, genericParams, rowData), (c) restore table widget state (`SetRows`, `SetColumns`, `SetCursor`), (d) no clearing

- [ ] Task 3: Update `executeDrilldown` in `nav.go` (AC: 1)
  - [ ] Replace `m.prepareStateTransition(transitionDrilldown)` with `m.prepareStateTransition(TransitionDrillDown)`
  - [ ] Remove the manual `navigationStack = append(...)` push block — `TransitionDrillDown` now pushes internally
  - [ ] Keep all other drilldown logic (target resolution, breadcrumb update, column pre-set, fetch command)

- [ ] Task 4: Update `navigateToBreadcrumb` in `nav.go` (AC: 1)
  - [ ] Truncate `m.navigationStack` to target depth `idx` BEFORE calling `prepareStateTransition`
  - [ ] Replace `m.prepareStateTransition(transitionBreadcrumb, idx)` with `m.prepareStateTransition(TransitionPop)` after stack truncation
  - [ ] Keep breadcrumb slice truncation and fetch logic

- [ ] Task 5: Update Esc/back and env/context-switch handlers in `update.go` (AC: 1)
  - [ ] Replace `prepareStateTransition(transitionBack)` with `prepareStateTransition(TransitionPop)` — remove manual stack pop (now internal to `TransitionPop`)
  - [ ] Replace `prepareStateTransition(transitionEnvSwitch)` with `prepareStateTransition(TransitionFull)` for all environment-switching paths
  - [ ] Replace `prepareStateTransition(transitionContextSwitch)` with `prepareStateTransition(TransitionFull)` for all context-switching paths (`:` key)

- [ ] Task 6: Audit all nav paths to confirm no bypasses (AC: 1)
  - [ ] `grep -rn "prepareStateTransition\|navigationStack\|breadcrumb\|viewMode\|currentRoot"` across `update.go` and `nav.go`
  - [ ] Confirm every code path that changes `viewMode`, `currentRoot`, `breadcrumb`, or `navigationStack` calls `prepareStateTransition`
  - [ ] Document any edge cases found

- [ ] Task 7: Write tests in `main_transition_test.go` (AC: 2)
  - [ ] Test `TransitionFull`: model with active modal, search, sort, cursor, and non-empty navigationStack → all fields cleared after call
  - [ ] Test `TransitionDrillDown`: model with cursor=5, sort on column 1 → call transition → viewState pushed to stack; stack length = 1; model has sort/search cleared; pushed viewState has correct cursor/columns/rows
  - [ ] Test `TransitionPop`: model with 1 entry on navigationStack → pop → all viewState fields restored; stack empty
  - [ ] Test `TransitionPop` with empty stack: safe, no panic
  - [ ] Test `navigateToBreadcrumb` integration: 2-level breadcrumb, jump to idx=0 → stack empty, state restored from idx=0 snapshot
  - [ ] `make cover` confirms ≥80% coverage on `transition.go`

- [ ] Task 8: Verify all tests pass (AC: all)
  - [ ] `make test` passes with zero regressions
  - [ ] `go vet ./...` passes
  - [ ] `gofmt -w .` produces no changes

## Dev Notes

### Context: Why This Story Exists

State leakage between navigation transitions is the primary defect category in o8n. The current `prepareStateTransition` in `transition.go` uses an unexported `transitionScope` enum with 5 constants (`transitionEnvSwitch`, `transitionContextSwitch`, `transitionDrilldown`, `transitionBack`, `transitionBreadcrumb`) and mixes two concerns: it clears some fields but leaves other cleanup (navStack push/pop) to callers. This split responsibility makes it hard to audit whether a nav path is fully correct.

Additionally, the function does NOT currently clear `activeModal` or `footerError` on any transition — meaning if a modal is open and the user triggers a navigation, the modal stays visible with stale context.

This story aligns the implementation with Architecture Decision 2: a typed, mandatory transition function that fully handles the push/pop/clear lifecycle for each transition category.

### Existing Code to Understand Before Starting

**`internal/app/transition.go`** (primary file to modify):

```go
// CURRENT — unexported enum
type transitionScope int
const (
    transitionEnvSwitch     transitionScope = iota
    transitionContextSwitch
    transitionDrilldown
    transitionBack
    transitionBreadcrumb
)

// CURRENT — pointer receiver, variadic depth param for breadcrumb
func (m *model) prepareStateTransition(scope transitionScope, depth ...int) {
    // always: clears sortColumn, sortAscending, searchTerm, searchMode, searchInput.Blur(),
    //         originalRows, filteredRows, popup reset
    // transitionEnvSwitch: also clears navigationStack, genericParams, selectedDefinitionKey, selectedInstanceID
    // transitionContextSwitch: also clears navigationStack, genericParams
    // transitionDrilldown: only sort/search (push handled by caller)
    // transitionBack: only sort/search (pop handled by caller)
    // transitionBreadcrumb: truncates navigationStack to depth[0]
}
```

**`internal/app/nav.go`** — Two primary callsites:

1. `executeDrilldown()` (line ~383): calls `prepareStateTransition(transitionDrilldown)` then manually builds `currentStateDrill viewState` and appends to `m.navigationStack`. After migration to `TransitionDrillDown`, the manual push block must be **removed** — the function now handles it internally.

2. `navigateToBreadcrumb()` (line ~359): calls `prepareStateTransition(transitionBreadcrumb, idx)`. After migration, the stack truncation logic moves to the caller; the call becomes `prepareStateTransition(TransitionPop)`.

**`internal/app/update.go`** — Additional callsites (search in the file):
- Esc handler (back navigation): calls `prepareStateTransition(transitionBack)` — update to `TransitionPop`
- Environment switch: calls `prepareStateTransition(transitionEnvSwitch)` — update to `TransitionFull`
- Context switch (`:` key): calls `prepareStateTransition(transitionContextSwitch)` — update to `TransitionFull`

**`internal/app/model.go`** — Key model fields:

```go
// Model field names (NOT the AC spec field names — see mapping table below)
type viewState struct {
    viewMode              string
    breadcrumb            []string
    contentHeader         string
    selectedDefinitionKey string
    selectedInstanceID    string
    tableRows             []table.Row
    tableCursor           int
    cachedDefinitions     []config.ProcessDefinition
    tableColumns          []table.Column
    genericParams         map[string]string
    rowData               []map[string]interface{}
}

type model struct {
    activeModal     ModalType          // NOT cleared by current prepareStateTransition
    footerError     string             // NOT cleared by current prepareStateTransition
    searchTerm      string             // AC calls this "searchQuery"
    searchMode      bool               // AC calls this "searchActive"
    sortColumn      int                // -1 means no sort
    sortAscending   bool               // AC calls this "sortDirection"; reset = true
    navigationStack []viewState
    // ... table.Model accessed via m.table.Cursor(), m.table.SetCursor(0)
}
```

### Critical: AC Field Name → Actual Go Field Mapping

The Acceptance Criteria uses architecture-spec names; the actual model uses different names:

| AC Spec Name | Actual Go Field | Clear Value |
|---|---|---|
| `searchQuery` | `m.searchTerm` | `""` |
| `searchActive` | `m.searchMode` | `false` |
| `sortColumn` | `m.sortColumn` | `-1` |
| `sortDirection` | `m.sortAscending` | `true` (ascending is default) |
| `tableCursor` | `m.table.SetCursor(0)` | via table widget method |
| `activeModal` | `m.activeModal` | `ModalNone` |
| `footerError` | `m.footerError` | `""` |
| `navigationStack` | `m.navigationStack` | `nil` |

Additionally, the existing code also clears these (keep doing so):
- `m.searchInput.Blur()` — defocuses the text input widget
- `m.originalRows = nil` — cached pre-filter rows
- `m.filteredRows = nil` — cached filtered rows
- Popup state reset (if popup.mode == popupModeSearch): `m.popup.mode = popupModeNone`, `m.popup.input = ""`, `m.popup.cursor = -1`, `m.popup.offset = 0`

### Critical: `TransitionDrillDown` Push-Before-Clear Ordering

The push MUST happen BEFORE any clearing:

```go
case TransitionDrillDown:
    // 1. Capture current viewState BEFORE clearing (includes active sort, cursor, rows)
    cols := m.table.Columns()
    rows := normalizeRows(append([]table.Row{}, m.table.Rows()...), len(cols))
    snapshot := viewState{
        viewMode:              m.viewMode,
        breadcrumb:            append([]string{}, m.breadcrumb...),
        contentHeader:         m.contentHeader,
        selectedDefinitionKey: m.selectedDefinitionKey,
        selectedInstanceID:    m.selectedInstanceID,
        tableRows:             rows,
        tableCursor:           m.table.Cursor(), // captured BEFORE clearing
        cachedDefinitions:     m.cachedDefinitions,
        tableColumns:          append([]table.Column{}, cols...),
        genericParams:         m.genericParams,
        rowData:               append([]map[string]interface{}{}, m.rowData...),
    }
    // 2. Push snapshot (parent state, with sort/cursor intact)
    m.navigationStack = append(m.navigationStack, snapshot)
    // 3. THEN clear (for the new child view)
    m.activeModal = ModalNone
    m.footerError = ""
    m.sortColumn = -1
    m.sortAscending = true
    m.searchTerm = ""
    m.searchMode = false
    m.searchInput.Blur()
    m.originalRows = nil
    m.filteredRows = nil
    if m.popup.mode == popupModeSearch { ... reset popup ... }
```

Note: `tableCursor` is NOT cleared by `TransitionDrillDown`. The caller (`executeDrilldown`) calls `m.table.SetCursor(0)` for the child view — keep that call in nav.go.

### Critical: `TransitionPop` Restore Logic

```go
case TransitionPop:
    if len(m.navigationStack) == 0 {
        return // safe no-op when stack is empty
    }
    // Pop top entry
    top := m.navigationStack[len(m.navigationStack)-1]
    m.navigationStack = m.navigationStack[:len(m.navigationStack)-1]
    // Restore all viewState fields
    m.viewMode = top.viewMode
    m.breadcrumb = top.breadcrumb
    m.contentHeader = top.contentHeader
    m.selectedDefinitionKey = top.selectedDefinitionKey
    m.selectedInstanceID = top.selectedInstanceID
    m.cachedDefinitions = top.cachedDefinitions
    m.genericParams = top.genericParams
    m.rowData = top.rowData
    // Restore table widget state
    if len(top.tableColumns) > 0 {
        m.table.SetRows(normalizeRows(nil, len(top.tableColumns)))
        m.table.SetColumns(top.tableColumns)
    }
    m.table.SetRows(top.tableRows)
    m.table.SetCursor(top.tableCursor)
    // NO clearing — restored state is the correct parent state
```

### Breadcrumb Navigation: Stack Truncate Then Pop

```go
// navigateToBreadcrumb — updated pattern
func (m *model) navigateToBreadcrumb(idx int) tea.Cmd {
    if idx < 0 || idx >= len(m.breadcrumb) {
        m.footerError = "Invalid breadcrumb index"
        return nil
    }
    // 1. Truncate navStack to target depth (entries above idx are discarded)
    if idx <= 0 {
        m.navigationStack = nil
    } else if idx < len(m.navigationStack) {
        m.navigationStack = m.navigationStack[:idx]
    }
    // 2. Call TransitionPop to restore state from new top of stack
    m.prepareStateTransition(TransitionPop)
    // 3. Truncate breadcrumb to target (still needed — viewState.breadcrumb restores the parent's slice)
    m.breadcrumb = append([]string{}, m.breadcrumb[:idx+1]...)
    // ... rest of fetch logic
}
```

**Note:** When idx == 0 (jump to root), the stack becomes nil, so `TransitionPop` is a no-op. The caller must still set `m.currentRoot`, `m.viewMode`, etc. from the breadcrumb target. This is a special case: jumping to the root level doesn't restore from stack (there's nothing to restore from); instead the caller reconstructs root state directly. For idx > 0, there IS a stack entry to restore from.

**Revised navigateToBreadcrumb logic:** Only call `TransitionPop` if the stack is non-empty after truncation (i.e., idx > 0 and the stack has entries). For idx == 0, call `TransitionFull` instead (full reset for root navigation).

```go
if idx == 0 {
    m.prepareStateTransition(TransitionFull)
} else {
    // truncate stack then pop
    if idx < len(m.navigationStack) {
        m.navigationStack = m.navigationStack[:idx]
    }
    m.prepareStateTransition(TransitionPop)
}
```

### Esc Handler in `update.go`

The current Esc/back handler (search for `transitionBack` in update.go) manually pops from the navStack and restores state. After this story, the pop and restore move into `TransitionPop`. The handler simplifies to:

```go
// Before (current pattern — schematic):
top := m.navigationStack[len(m.navigationStack)-1]
m.navigationStack = m.navigationStack[:len(m.navigationStack)-1]
m.viewMode = top.viewMode
// ... many manual field restores ...
m.prepareStateTransition(transitionBack)  // only cleared sort/search

// After (new pattern):
m.prepareStateTransition(TransitionPop)  // does pop + restore + nothing else
// m.currentRoot, m.breadcrumb etc. are already restored by TransitionPop
```

**Important**: If the handler had additional logic after the restore (like resetting `m.currentRoot` from the restored breadcrumb), that logic should still remain. `TransitionPop` restores the viewState fields but does NOT call `fetchForRoot` — callers still issue the fetch command if needed.

### `removeExistingTransitionScope` — Removal Strategy

The old `transitionScope` type and its 5 constants can be removed ONLY after all callers are migrated. Do this:

1. Migrate all callers first (Tasks 3-5)
2. Run `grep -rn "transitionScope\|transitionEnvSwitch\|transitionContextSwitch\|transitionDrilldown\|transitionBack\|transitionBreadcrumb"` — should show zero results
3. Delete the type and constants from `transition.go`

Do NOT use a deprecation alias — just delete them cleanly.

### Testing Requirements

Test file: `internal/app/main_transition_test.go` (co-located with `transition.go`).

Coverage target: ≥80% on `transition.go` per architecture convention.

Key test scenarios:

**`TransitionFull`:**
```go
func TestTransitionFull_ClearsAllFields(t *testing.T) {
    m := newTestModel(t)
    m.activeModal = ModalConfirmDelete
    m.footerError = "some error"
    m.searchTerm = "filter"
    m.searchMode = true
    m.sortColumn = 2
    m.sortAscending = false
    m.table.SetCursor(5)
    m.navigationStack = []viewState{{viewMode: "process-definition"}}

    m.prepareStateTransition(TransitionFull)

    if m.activeModal != ModalNone { t.Error("activeModal not cleared") }
    if m.footerError != "" { t.Error("footerError not cleared") }
    if m.searchTerm != "" { t.Error("searchTerm not cleared") }
    if m.searchMode { t.Error("searchMode not cleared") }
    if m.sortColumn != -1 { t.Error("sortColumn not cleared") }
    if !m.sortAscending { t.Error("sortAscending not reset to true") }
    if m.table.Cursor() != 0 { t.Error("tableCursor not reset to 0") }
    if len(m.navigationStack) != 0 { t.Error("navigationStack not cleared") }
}
```

**`TransitionDrillDown`:**
```go
func TestTransitionDrillDown_PushesStateAndClears(t *testing.T) {
    m := newTestModel(t)
    m.viewMode = "process-definition"
    m.table.SetCursor(3)
    // set up some rows/columns

    m.prepareStateTransition(TransitionDrillDown)

    if len(m.navigationStack) != 1 { t.Error("expected 1 entry pushed") }
    pushed := m.navigationStack[0]
    if pushed.viewMode != "process-definition" { t.Error("pushed wrong viewMode") }
    if pushed.tableCursor != 3 { t.Error("pushed wrong tableCursor") }
    if m.activeModal != ModalNone { t.Error("activeModal not cleared after drilldown") }
    if m.sortColumn != -1 { t.Error("sortColumn not cleared") }
}
```

**`TransitionPop`:**
```go
func TestTransitionPop_RestoresStateFromStack(t *testing.T) {
    m := newTestModel(t)
    m.navigationStack = []viewState{
        {viewMode: "process-definition", tableCursor: 7, breadcrumb: []string{"defs"}},
    }

    m.prepareStateTransition(TransitionPop)

    if len(m.navigationStack) != 0 { t.Error("stack not emptied") }
    if m.viewMode != "process-definition" { t.Error("viewMode not restored") }
    if m.table.Cursor() != 7 { t.Error("tableCursor not restored") }
}

func TestTransitionPop_EmptyStackIsNoop(t *testing.T) {
    m := newTestModel(t)
    m.navigationStack = nil
    // Must not panic
    m.prepareStateTransition(TransitionPop)
}
```

Use `newTestModel(t)` from `main_test.go` as the model factory — same pattern as `main_modal_test.go`.

### Project Structure Notes

**Primary file modified:** `internal/app/transition.go`
- Add `TransitionType` type + 3 exported constants
- Refactor `prepareStateTransition` body
- Remove old `transitionScope` type + 5 unexported constants (after callers migrated)

**Other files modified:**
- `internal/app/nav.go` — update `executeDrilldown` and `navigateToBreadcrumb` callsites
- `internal/app/update.go` — update Esc/back, env switch, context switch callsites

**New file:**
- `internal/app/main_transition_test.go` — test file for transition contract

**Off-limits:** `internal/operaton/` — never modify. `o8n-cfg.yaml` / `o8n-env.yaml` — no changes needed.

### Async / Bubble Tea Constraints

- `prepareStateTransition` is a pure synchronous method on `*model`. No `tea.Cmd`, no goroutines.
- All calls happen within `Update()` message handlers — this is correct Bubble Tea pattern (model mutations in Update only).
- `TransitionPop` sets rows/columns on `m.table` — these are synchronous widget operations, safe within Update.
- Callers may issue `tea.Cmd` (e.g., `fetchForRoot`) AFTER calling `prepareStateTransition`, as is the current pattern.

### Specification Update Obligation

Per the architecture definition of done: after implementation, `specification.md` MUST be updated to document:
- The `TransitionType` enum (`TransitionFull`, `TransitionDrillDown`, `TransitionPop`) and their behavioral contracts
- The `prepareStateTransition(t TransitionType)` function signature
- The push-before-clear ordering guarantee for `TransitionDrillDown`
- The field-clearing list for `TransitionFull`

### Learnings from Story 1.1

- The Senior Developer Review will apply fixes for any specification compliance gaps. Write the code as spec'd and document any tricky design choices in Completion Notes.
- `gofmt -w .` after modifying large files (especially if using heredocs or raw strings).
- The test file pattern `main_<feature>_test.go` co-located with the production file — confirmed by Story 1.1 (main_modal_test.go alongside modal.go).
- Coverage target ≥80% on the primary modified file — use `make cover` to verify.
- Leave `KeyHint` (model.go:131) and `getKeyHints()` untouched — Story 2.1 handles that rename.

### References

- [Source: `_bmad/planning-artifacts/architecture.md#Decision 2: State Transition Contract`] — TransitionType enum, field-clearing specs, "any navigation code that does not call prepareStateTransition is a bug"
- [Source: `_bmad/planning-artifacts/epics.md#Story 1.2`] — Acceptance criteria (BDD format)
- [Source: `internal/app/transition.go`] — Existing prepareStateTransition with transitionScope enum
- [Source: `internal/app/nav.go:383`] — executeDrilldown callsite (manual push to remove)
- [Source: `internal/app/nav.go:359`] — navigateToBreadcrumb callsite (transitionBreadcrumb to replace)
- [Source: `internal/app/model.go:239`] — viewState struct (all fields to capture/restore in DrillDown/Pop)
- [Source: `_bmad/implementation-artifacts/1-1-modal-factory-foundation.md`] — Previous story patterns (test structure, completion notes format)

## Dev Agent Record

### Agent Model Used

(to be filled by implementing agent)

### Debug Log References

### Completion Notes List

### File List

### Senior Developer Review (AI)

### Change Log
