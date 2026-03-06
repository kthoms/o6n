# Story 3.6: Process Variable Inspection, Editing & JSON View

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an **operator**,
I want to **inspect and edit process variables** associated with a process instance, and **view or copy any resource row as JSON**,
so that I can diagnose and correct stuck processes and extract resource data directly from the TUI.

## Acceptance Criteria

1. **Given** the operator drills down to the variables view for a process instance
   **When** the variables table loads
   **Then** all process variables for that instance are displayed with name, type, and value columns.

2. **Given** the operator selects a variable row and presses the **edit key (`e`)**
   **When** the edit dialog opens
   **Then** the current variable value is pre-populated and the type is displayed
   **And** input is validated against the variable type (string, integer, boolean, JSON)
   **And** on confirmation, the variable is updated via the API and the table refreshes.

3. **Given** the operator presses **`J`** on any table row
   **When** the JSON viewer opens
   **Then** `ModalJSONView` opens as an `OverlayLarge` modal with title `ResourceType: ID` and the row's full data formatted as JSON in a scrollable `viewport.Model`.
   **And** background content remains visible behind the modal.
   **And** a hint line is rendered at the bottom: `Ctrl+J Copy  Esc Close`.

4. **Given** `ModalJSONView` is open
   **When** the operator presses **`Ctrl+J`**
   **Then** the JSON content is copied to the system clipboard (`github.com/atotto/clipboard`) and the footer confirms: `✓ Copied to clipboard`.
   **And** the modal remains open.

5. **Given** the operator presses **`Ctrl+J`** directly on any table row (without opening the viewer)
   **When** the copy action executes
   **Then** the row's JSON is copied directly to the system clipboard and the footer confirms: `✓ Copied to clipboard`.
   **And** `ModalJSONView` does not open.

6. **Given** the operator opens the Action Menu (**`Ctrl+Space`**)
   **When** the menu renders
   **Then** `[J] View as JSON` and `[Ctrl+J] Copy JSON` are the final two items in the list.

## Tasks / Subtasks

- [ ] **Variable Inspection & Editing (AC: 1, 2)**
  - [ ] Audit `process-variables` TableDef in `o8n-cfg.yaml`: ensure `value` is `editable: true` with `input_type: auto`.
  - [ ] Update `internal/app/edit.go`: ensure CAMUNDA type validation (String, Integer, Boolean, Json) is strictly enforced.
  - [ ] Ensure `editSavedMsg` triggers a data re-fetch for the variables context to ensure server-side sync.
- [ ] **ModalJSONView Surgical Refactor (AC: 3, 6)**
  - [ ] Rename `ModalDetailView` → `ModalJSONView` across `model.go`, `modal.go`, and `view.go`.
  - [ ] Register `ModalJSONView` in `modal.go` with `OverlayLarge` size hint and established `HintLine`.
  - [ ] Refactor `renderJSONViewBody` (was `modalDetailViewBody`) to use `github.com/charmbracelet/bubbles/viewport` for performance.
  - [ ] Use the existing `syntaxHighlightJSON` function in `view.go` for all inner JSON rendering.
  - [ ] Verify `buildActionsForRoot` in `nav.go` correctly appends these actions to the `ModalActionMenu`.
- [ ] **JSON Integration & Feedback (AC: 4, 5)**
  - [ ] Map `J` and `Ctrl+J` keys in `update.go` for the main table view.
  - [ ] Implement `Ctrl+J` handler inside the `ModalJSONView` logic block in `update.go`.
  - [ ] Use `setFooterStatus` with `footerStatusSuccess` for "Copied to clipboard" feedback.
- [ ] **Verification (AC: all)**
  - [ ] Create `internal/app/main_json_view_test.go` covering `J` trigger, modal presence, and `Ctrl+J` clipboard logic (mocked).
  - [ ] Ensure variable editing remains functional and type-safe via `make test`.

## Dev Notes

### Architecture Compliance
- **Modal Factory:** `ModalJSONView` must use `renderModal(m, cfg)` in the view path.
- **Async API:** All variable updates must return a `tea.Cmd`.
- **Component Pattern:** Use `viewport.Model` for JSON scrolling; it is already established in the codebase for help/detail views.

### UI/UX & Formatting
- **Formatting:** Use `json.MarshalIndent` with 2 spaces for the viewer content.
- **Syntax Highlighting:** Ensure skin-based JSON semantic colors (`jsonKey`, `jsonValue`, etc.) are applied.

### Project Structure Notes
- **View Logic:** `internal/app/view.go`
- **Navigation/Actions:** `internal/app/nav.go`
- **Update/Keys:** `internal/app/update.go`

### References
- [Source: `_bmad/planning-artifacts/epics.md#Story 3.6`]
- [Source: `_bmad/planning-artifacts/ux-design-specification.md#Journey 4: Variable Inspection and JSON Export`]
- [Source: `internal/app/view.go#syntaxHighlightJSON`]

## Dev Agent Record

### Agent Model Used

Gemini 2.0 Flash

### Debug Log References

### Completion Notes List

### File List
