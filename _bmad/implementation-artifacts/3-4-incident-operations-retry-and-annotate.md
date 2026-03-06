# Story 3.4: Incident Operations (Retry & Annotate)

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As an **operator** (Alex persona),
I want to retry a failed job and set an annotation on an incident from within the TUI,
so that I can resolve incidents without switching to a browser or writing curl commands.

## Acceptance Criteria

1. **Given** the operator is viewing the `incident` resource table and selects a row
   **When** the operator executes the configured Retry action (key `r`)
   **Then** the corresponding job retry API call (`PUT /job/{jobId}/retries`) is made
   **And** the footer shows `✓ Job retried`
   **And** the incident table refreshes automatically.

2. **Given** the operator is viewing an incident row
   **When** the operator presses `a` (Annotate) or `e` (Edit)
   **Then** the `ModalEdit` factory modal opens for the `annotation` field.
   **And** on saving, the annotation is persisted via `PUT /incident/{id}/annotation`
   **And** the footer confirms success.

3. **Given** the operator drills down from an incident row
   **When** the operator presses `Enter`
   **Then** `prepareStateTransition(TransitionDrillDown)` is called
   **And** the application navigates to the `process-instance` view filtered to the incident's `processInstanceId`.

4. **Given** the operator is in the incident table view
   **When** the footer renders
   **Then** the hints `r Retry` and `a Annotate` are visible (if space allows).

## Tasks / Subtasks

- [ ] Update `o8n-cfg.yaml` for the `incident` table (AC: 1, 2, 3)
  - [ ] Add `jobId` column definition (`visible: false`).
  - [ ] Add `annotation` column definition (`editable: true`).
  - [ ] Add `edit_action` for the table:
    - `method: PUT`
    - `path: /incident/{id}/annotation`
    - `body_template: '{"annotation": "{value}"}'`
    - `name_column: id`
  - [ ] Add `actions`:
    - `key: r`, `label: Retry`, `method: PUT`, `path: /job/{jobId}/retries`, `body: '{"retries": 1}'`, `id_column: jobId`
    - `key: a`, `label: Annotate` (shortcut for editing the annotation column)
  - [ ] Ensure `drilldown` target is `process-instance` with `param: id` and `column: processInstanceId`.

- [ ] Improve Action ID Resolution in `internal/app/nav.go` (AC: 1)
  - [ ] Update `resolveActionID(action config.ActionDef)` to check `m.rowData` for hidden columns if the `IDColumn` is not found in the visible table columns.

- [ ] Enhance Hint System in `internal/app/hints.go` (AC: 4)
  - [ ] Update `tableViewHints(m model)` to dynamically append hints from `TableDef.Actions`.
  - [ ] Ensure these resource-specific hints have high priority (e.g., `Priority: 4`).

- [ ] Tests and Validation (AC: 1, 2, 3, 4)
  - [ ] Create `internal/app/main_incident_ops_test.go`.
  - [ ] Test Retry: verify `jobId` resolution and API command emission.
  - [ ] Test Annotate: verify `ModalEdit` opening and `edit_action` command emission.
  - [ ] Test Drilldown: verify navigation to `process-instance` with correct filter param.
  - [ ] Test Hints: verify `r` and `a` appear in footer for incident table.

## Dev Notes

- **Hidden Column Resolution**: The `jobId` is required for the retry action but shouldn't clutter the main table. Using `rowData` lookup in `resolveActionID` is critical.
- **Edit Action Pattern**: Camunda's annotation endpoint is non-standard for PUT (it takes a single JSON field). The `edit_action` pattern in `o8n-cfg.yaml` is designed for this.
- **Success Feedback**: Story 3.3 established the `actionExecutedMsg` pattern. Use it to trigger the refresh and success message.

### Project Structure Notes

- All implementation remains within `internal/app/` and `o8n-cfg.yaml`.
- Respect `prepareStateTransition` for the drilldown path.

### References

- [Source: `_bmad/planning-artifacts/epics.md#Story 3.4`]
- [Source: `o8n-cfg.yaml`] — `incident` table section
- [Source: `internal/app/nav.go`] — `resolveActionID` and `executeDrilldown`
- [Source: `internal/app/hints.go`] — `tableViewHints`

## Dev Agent Record

### Agent Model Used

Gemini 2.0 Flash

### Debug Log References

None.

### Completion Notes List

### File List
