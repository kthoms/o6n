# Story 3.2: Drill-Down Navigation & Breadcrumb

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an **operator**,
I want to drill down from a parent resource into related child resources and navigate back via Escape or breadcrumb,
so that I can traverse resource hierarchies (e.g., process definition → instances → variables) without losing my place.

## Acceptance Criteria

1. **Given** the current resource type has drilldown rules configured in `o8n-cfg.yaml`
   **When** the operator presses Enter on a row
   **Then** `prepareStateTransition(TransitionDrillDown)` is called, the current view is pushed onto the navigation stack, and the child resource type loads filtered to the selected parent.

2. **Given** the operator has drilled down one or more levels
   **When** the operator presses Escape
   **Then** `prepareStateTransition(TransitionPop)` is called and the previous view is fully restored (rows, cursor, columns, filters, breadcrumb).

3. **Given** the operator has drilled down multiple levels and the breadcrumb shows level labels
   **When** the operator selects a specific breadcrumb level (via mouse or future key binding)
   **Then** the application pops to that level and restores the view state captured at that level using `prepareStateTransition(TransitionPop)` or a series of pops.

## Tasks / Subtasks

- [ ] Audit and Refine `executeDrilldown` in `internal/app/nav.go` (AC: 1)
  - [ ] Verify `prepareStateTransition(TransitionDrillDown)` is called correctly.
  - [ ] Ensure `m.navigationStack` correctly captures the `viewState` (rows, cursor, filters, etc.) before the transition.
  - [ ] Verify `m.genericParams` is correctly populated with the drill-down filter.
  - [ ] Ensure child columns are pre-set to avoid "column flash" from parent.
- [ ] Audit and Refine `navigateToBreadcrumb` and Escape handling (AC: 2, 3)
  - [ ] Verify `prepareStateTransition(TransitionPop)` is used when popping.
  - [ ] Ensure `navigationStack` is correctly truncated and the restored state is applied to the model.
  - [ ] Verify cursor position and search filters are correctly restored.
- [ ] Add breadcrumb level selection logic (AC: 3)
  - [ ] Ensure `navigateToBreadcrumb(idx)` correctly handles stack restoration for any previous level.
- [ ] Tests and Validation (AC: 1, 2, 3)
  - [ ] Create `internal/app/main_drilldown_nav_test.go`
  - [ ] Test DrillDown: verify stack push and state clearing (except stack).
  - [ ] Test Pop (Escape): verify stack pop and full state restoration.
  - [ ] Test Jump: verify jumping to an early breadcrumb level restores the correct state.
- [ ] Documentation (AC: 1, 2, 3)
  - [ ] Ensure `specification.md` accurately describes the `TransitionDrillDown` and `TransitionPop` behavior.

## Dev Notes

- **State Transition Contract:** This story is the primary test for `TransitionDrillDown` and `TransitionPop`. Refer to `internal/app/nav.go` and `architecture.md` for the state clearing/restoration rules.
- **View State Persistence:** The `viewState` struct in `model.go` must be exhaustive enough to allow a perfect restoration of the parent view.
- **Navigation Stack:** The stack is the source of truth for the "Back" operation. Ensure it doesn't leak memory or grow boundlessly (though depth is usually shallow).
- **Wait for Data:** When drilling down, the UI should show a loading indicator while the child data is fetched.

### Project Structure Notes

- All logic resides in `internal/app/nav.go` (methods on `model`).
- Coordination with `update.go` for key handling (Enter and Escape).

### References

- Epic 3 Story 3.2: `_bmad/planning-artifacts/epics.md`
- State Transition Contract: `internal/app/nav.go`
- ViewState struct: `internal/app/model.go`

## Dev Agent Record

### Agent Model Used

Gemini 2.0 Flash

### Debug Log References

### Completion Notes List

### File List
