# Story 3.1: Resource Navigation & Context Switcher

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an **operator**,
I want to navigate to any of the 35 configured resource types using the context switcher (key `:`),
so that I can reach any operational view in seconds from anywhere in the application.

## Acceptance Criteria

1. **Given** 35 resource types are defined in `o8n-cfg.yaml`
   **When** the operator presses `:` to open the context switcher
   **Then** a searchable list of all configured resource types is displayed in an `OverlayCenter` modal.

2. **Given** the context switcher is open
   **When** the operator types a partial name and selects a match
   **Then** `prepareStateTransition(TransitionFull)` is called, clearing all prior view state
   **And** the selected resource type loads with a fresh, paginated table
   **And** the breadcrumb shows the new context name.

3. **Given** the selected resource type returns data from the API
   **When** the table loads
   **Then** the correct columns for that resource type (as defined in `o8n-cfg.yaml`) are displayed and the first row is selected.

## Tasks / Subtasks

- [ ] Migrate `popupModeContext` legacy logic to `ModalContextSwitcher` factory modal (AC: 1, 2)
  - [ ] Add `ModalContextSwitcher` to `ModalType` in `internal/app/model.go`
  - [ ] Update `update.go` to set `m.activeModal = ModalContextSwitcher` on `:` instead of setting `popupModeContext`
  - [ ] Remove `popupModeContext` from `popupMode` enum and clean up associated legacy logic in `model.go`, `update.go`, `view.go`, `nav.go`
- [ ] Implement `ModalContextSwitcher` registration in `modal.go` (AC: 1)
  - [ ] Register with `SizeHint: OverlayCenter`
  - [ ] Register body renderer `renderContextSwitcherBody`
  - [ ] Define `HintLine` showing `↑↓ Nav`, `Enter Select`, `Esc Close`
- [ ] Create `renderContextSwitcherBody` in `internal/app/view.go` (AC: 1)
  - [ ] Reuse `m.rootContexts` (filtered by `m.popup.input`) as items
  - [ ] Use `lipgloss.Place` to render a centered list with current skin styles
- [ ] Handle selection and navigation in `internal/app/update.go` (AC: 2, 3)
  - [ ] In `ModalContextSwitcher` key handler, call `prepareStateTransition(TransitionFull)` on `Enter`
  - [ ] Update `m.currentRoot`, `m.breadcrumb`, and `m.viewMode` to the selected resource
  - [ ] Dispatch `m.fetchForRoot(selected)` cmd
  - [ ] Ensure `m.table.SetCursor(0)` is called for the new context
- [ ] Tests and Validation (AC: 1, 2, 3)
  - [ ] Create `internal/app/main_context_switcher_test.go`
  - [ ] Test `:` trigger, partial name search, selection, and transition state clearing
- [ ] Documentation (AC: 3)
  - [ ] Update `specification.md` modal keyboard behavior table with `ModalContextSwitcher`

## Dev Notes

- **Architecture Compliance:** Must use the `ModalConfig` and `renderModal()` factory established in Story 1.1. Do NOT add a `case ModalContextSwitcher` to the main `View` render switch.
- **State Transition Contract:** Mandatory use of `prepareStateTransition(TransitionFull)` from Story 1.2 to ensure zero state leakage between resource contexts.
- **Config-Driven:** Resource list is already filtered to only those with `TableDef` in `o8n-cfg.yaml` within `newModel`.
- **Legacy Cleanup:** The current `:` implementation in `update.go` (line ~708) and `view.go` (line ~630) is a legacy "popup" system. Migrating it to the modal factory is a key architectural goal of this story.

### Project Structure Notes

- All changes in `internal/app/`.
- No changes to `internal/operaton/` or `internal/client/`.
- Test file `main_context_switcher_test.go` co-located with implementation.

### References

- Epic 3 section: `_bmad/planning-artifacts/epics.md`
- Modal Factory: `internal/app/modal.go`
- State Transition Contract: `internal/app/nav.go` (check `prepareStateTransition`)
- Resource definitions: `o8n-cfg.yaml`

## Dev Agent Record

### Agent Model Used

Gemini 2.0 Flash

### Debug Log References

### Completion Notes List

### File List
