# Story 3.7: Search, Filter & Auto-Refresh

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an **operator**,
I want to **filter the current resource table by search term** and **toggle auto-refresh**,
so that I can efficiently find specific resources and monitor live state.

## Acceptance Criteria

1. **Given** the operator is viewing any resource table
   **When** the operator presses **`/`** and types a search term
   **Then** the table is filtered to rows matching the term in real time.
   **And** the `FilterBar` (UX specified) indicates the active filter state via a dedicated `renderFilterBar(m model)` function.

2. **Given** an active search filter is set
   **When** the operator presses **Escape** or the clear-filter key
   **Then** the filter is cleared and the full table is restored from `m.originalRows` without a re-fetch.

3. **Given** the operator presses **`Ctrl+Shift+R`** to toggle auto-refresh
   **When** auto-refresh is enabled
   **Then** the table reloads on the configured interval (`5s`) and a visual indicator (**`⟳`**) appears in the footer right area.
   **And** the indicator flashes in sync with `m.flashActive`.
   **And** pressing `Ctrl+Shift+R` again disables auto-refresh and the indicator disappears.

4. **Given** any API call is in flight
   **When** the request takes longer than 500ms
   **Then** a loading indicator (spinner) appears in the footer center column and the UI remains fully interactive.

5. **Given** the application is connected to an environment
   **When** the footer renders
   **Then** the API status indicator (**`●`** green, **`✗`** red, or **`○`** muted) is visible in the footer right column.

## Tasks / Subtasks

- [ ] **Implement FilterBar & Search Logic (AC: 1, 2)**
  - [ ] Add `renderFilterBar(m model) string` to `view.go` implementing the 5 UX states (Hidden, Active input, Locked/applied, Server-side active, Clearing).
  - [ ] Update `Update()` key handler for **`/`** to enter `searchMode` and capture `m.originalRows`.
  - [ ] Ensure `prepareStateTransition` properly clears search/filter state on navigation.
- [ ] **Implement Auto-Refresh & Indicator (AC: 3)**
  - [ ] Update key binding for auto-refresh toggle to **`Ctrl+Shift+R`**.
  - [ ] Implement periodic refresh command using `tea.Tick`.
  - [ ] Add the **`⟳`** badge to `renderFooter` in `view.go`, ensuring it uses the `m.flashActive` state for its flash effect.
- [ ] **Implement API Status & Loading Debounce (AC: 4, 5)**
  - [ ] Implement **`●`/`✗`/`○`** indicators in `renderFooter` based on `m.envStatus`.
  - [ ] Implement a delayed loading message or tick to ensure the spinner only appears for requests > 500ms.
  - [ ] Ensure the loading spinner in the footer center column is tied to `m.isLoading`.
- [ ] **Hint System Integration (AC: 1, 3)**
  - [ ] Add hints for `/ Filter` and `Ctrl+Shift+R Refresh` to `tableViewHints` in `internal/app/hints.go` with high priority.
- [ ] **Verify & Test (AC: all)**
  - [ ] Create `internal/app/main_filter_refresh_test.go` covering search mode entry, live filtering from `originalRows`, and auto-refresh toggling.
  - [ ] Verify footer indicators appear/disappear correctly based on model state.

## Dev Notes

### Architecture Compliance
- **State Transition:** Search/filter clearing MUST use `prepareStateTransition` or be integrated into its logic to prevent state leakage.
- **Async Pattern:** Auto-refresh MUST be implemented via `tea.Tick` and `tea.Cmd`.
- **Footer Hint Push Model:** Update `hints.go` to include the search and refresh keys.

### UI/UX Standards
- **Filter States:**
    1. No filter: `FilterBar` hidden.
    2. Input active: `FilterBar` appears with prompt.
    3. Locked: Filter text shown in locked style.
    4. Server-side: Indicator for API-level filtering.
    5. Clearing: Transition state.
- **Status Symbols:** Use `●` (operational), `✗` (unreachable), `○` (unknown).

### Project Structure Notes
- **View Logic:** `internal/app/view.go` (footer and FilterBar).
- **Update Logic:** `internal/app/update.go` (key handlers and tick handling).
- **Transition Logic:** `internal/app/transition.go`.

### References
- [Source: `_bmad/planning-artifacts/epics.md#Story 3.7`]
- [Source: `_bmad/planning-artifacts/ux-design-specification.md#Search and Filter Patterns`]
- [Source: `_bmad/planning-artifacts/ux-design-specification.md#Auto-refresh indicator design`]

## Dev Agent Record

### Agent Model Used

Gemini 2.0 Flash

### Debug Log References

### Completion Notes List

### File List
