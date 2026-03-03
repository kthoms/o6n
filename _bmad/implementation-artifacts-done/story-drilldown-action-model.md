# Story: Drill-Down & Navigation Action Model

## Summary

Establish a coherent model that distinguishes **drill-downs** (strict parent→child hierarchy, Enter key, `▶` indicator) from **navigation actions** (related-resource views, Space menu, shortcut key). Reduce every resource to at most one drill-down. Convert all audit-log and cross-namespace drilldowns to `type: navigate` actions in the action menu. Fix `▶` so it only appears where Enter actually navigates.

---

## Design Rationale

### Rule: What earns a drill-down?

A `drilldown` entry is warranted **if and only if**:

1. The child entities are **owned by** the parent — they cannot exist without it.
2. The relationship is a strict parent→child hierarchy with a foreign key from child to parent.
3. The child collection is the **primary** thing a user wants to inspect when they select a parent row.

### Rule: What is a navigation action?

Everything else that leads to a related resource is a **navigation action**: a Space-menu item with a shortcut key that navigates the view to a pre-filtered related resource. These are visually distinct from mutation actions (HTTP verbs) in the menu.

### Rule: The `▶` indicator

`▶` is a navigation affordance. It must only appear when the resource has a configured `drilldown`. Showing it on resources where Enter does nothing is a broken affordance.

---

## Current State vs Proposed State

### Drilldown inventory — 27 relationships today, 4 after this story

| Resource | Current drilldown(s) | After this story |
|----------|---------------------|-----------------|
| `process-definition` | `process-instance`, `history-process-instance` | `process-instance` only |
| `process-instance` | `process-variables`, `task`, `incident`, `variable-instance` | `process-variables` only |
| `job-definition` | `job` | `job` (unchanged) |
| `deployment` | `process-definition`, `decision-definition` | `process-definition` only |
| `external-task` | `history-external-task-log` | removed → navigate action `h` |
| `job` | `history-job-log` | removed → navigate action `h` |
| `batch` | `history-batch` | removed → navigate action `h` |
| `decision-definition` | `history-decision-instance` | removed → navigate action `h` |
| `user` | `history-user-operation` | removed → navigate action `h` |
| `group` | `history-identity-link-log` | removed → navigate action `h` |
| `task` | `variable-instance`, `history-detail` | both removed → navigate actions `v`, `h` |
| `tenant` | `process-definition`, `process-instance`, `deployment` | all removed → navigate actions `d`, `i`, `e` |

### Navigate actions to add

| Resource | Key | Label | Target resource | Param | Column |
|----------|-----|-------|----------------|-------|--------|
| `user` | `h` | View Operation Log | `history-user-operation` | `userId` | `id` |
| `group` | `h` | View Identity Links | `history-identity-link-log` | `groupId` | `id` |
| `external-task` | `h` | View History | `history-external-task-log` | `externalTaskId` | `id` |
| `job` | `h` | View History | `history-job-log` | `jobId` | `id` |
| `process-definition` | `h` | View History | `history-process-instance` | `processDefinitionId` | `id` |
| `batch` | `h` | View History | `history-batch` | `batchId` | `id` |
| `decision-definition` | `h` | View History | `history-decision-instance` | `decisionDefinitionKey` | `key` |
| `task` | `h` | View History | `history-detail` | `taskId` | `id` |
| `task` | `v` | View Variables | `variable-instance` | `taskIdIn` | `id` |
| `tenant` | `d` | View Definitions | `process-definition` | `tenantIdIn` | `id` |
| `tenant` | `i` | View Instances | `process-instance` | `tenantIdIn` | `id` |
| `tenant` | `e` | View Deployments | `deployment` | `tenantIdIn` | `id` |
| `deployment` | `D` | View Decision Defs | `decision-definition` | `deploymentId` | `id` |
| `process-instance` | `t` | View Tasks | `task` | `processInstanceId` | `id` |
| `process-instance` | `i` | View Incidents | `incident` | `processInstanceId` | `id` |

---

## Acceptance Criteria

### AC-1: ▶ indicator

- **AC-1a:** `▶ ` is prepended to the first column of every row **only** when `def.Drilldown != nil` (after the struct change, this is a pointer check).
- **AC-1b:** Resources with no drilldown render rows without `▶`. The column width calculation in `buildColumnsFor()` does not reserve 2 characters for the prefix on those resources.
- **AC-1c:** Pressing `Enter` on a resource with `def.Drilldown == nil` is a no-op (returns `m, nil`).

### AC-2: Config struct — `Drilldown` as a pointer

- **AC-2a:** `TableDef.Drilldown` changes from `[]DrillDownDef` (`yaml:"drilldown,omitempty"`) to `*DrillDownDef` (`yaml:"drilldown,omitempty"`).
- **AC-2b:** YAML format for a drilldown changes from a sequence item (`- target: ...`) to an inline mapping (`target: ...`). All YAML consumer code is updated accordingly.
- **AC-2c:** The "best match" loop in `update.go:1080-1090` (iterating `def.Drilldown` to pick by visible column) is removed; `def.Drilldown` is used directly as the single chosen entry.

### AC-3: Config struct — `ActionDef` gains navigate fields

- **AC-3a:** `ActionDef` gains a `Type string \`yaml:"type,omitempty"\`` field. Valid values: `""` (default, HTTP mutation) and `"navigate"`.
- **AC-3b:** `ActionDef` gains three navigation fields used only when `Type == "navigate"`:
  - `Target string \`yaml:"target,omitempty"\`` — target resource key
  - `Param  string \`yaml:"param,omitempty"\``  — query parameter name
  - `Column string \`yaml:"column,omitempty"\`` — source column for the ID value (default `"id"` if empty)
- **AC-3c:** The existing `Method`, `Path`, `Body`, `Confirm`, `IDColumn` fields on `ActionDef` are unused (and may be absent) for `type: navigate` actions. The config loader does not error on their absence.

### AC-4: `o8n-cfg.yaml` updated

- **AC-4a:** All 12 resources listed in the drilldown inventory have their `drilldown:` entries updated to the single canonical entry (or removed entirely).
- **AC-4b:** All 15 navigate actions from the table above are added to their respective resources.
- **AC-4c:** The YAML format of every remaining `drilldown:` entry uses the new mapping style (not a sequence).
- **AC-4d:** `process-definition`, `process-instance`, `job-definition`, `deployment` each have exactly one `drilldown:` entry.

### AC-5: Drilldown execution extracted to helper

- **AC-5a:** A new helper method `(m *model) executeDrilldown(d *config.DrillDownDef) (model, tea.Cmd)` encapsulates the navigation logic currently at `update.go:1119-1180`.
- **AC-5b:** The helper handles the `selectedInstanceID` / `selectedDefinitionKey` assignments for the four canonical drilldown targets (same switch statement as today, update.go:1143-1148).
- **AC-5c:** The Enter/→ handler in `update.go` calls this helper instead of inline logic.

### AC-6: Navigate actions execute via `executeDrilldown`

- **AC-6a:** `buildActionsForRoot()` in `nav.go` detects `act.Type == "navigate"` and builds an `actionItem.cmd` that calls `m.executeDrilldown()` with a synthetic `*config.DrillDownDef` constructed from `act.Target`, `act.Param`, `act.Column`, `act.Label`.
- **AC-6b:** The `actionItem.cmd` for a navigate action resolves the ID value from `m.rowData` (same lookup as the Enter handler uses: `m.rowData[cursor][colName]`).
- **AC-6c:** Pressing the shortcut key for a navigate action (e.g. `h` on a `user` row, within the Space menu) triggers the navigation and updates the breadcrumb identically to Enter-drilldown.

### AC-7: Action menu visual separation

- **AC-7a:** `renderActionsMenu()` in `view.go` renders a visual separator line (`────` or similar) between mutation actions and navigate actions when both types are present.
- **AC-7b:** Navigate action labels are suffixed with ` →` in the menu (e.g. `[h] View Operation Log →`).
- **AC-7c:** Mutation actions are rendered without the `→` suffix (unchanged from today).

### AC-8: Help screen — context-aware drill-down hint

- **AC-8a:** The `renderHelpScreen()` function in `view.go` only includes the `Enter — drill down` entry in the Navigation section when `m.findTableDef(m.currentRoot).Drilldown != nil`.
- **AC-8b:** Navigate actions appear in a dedicated **Views** section in the help screen, populated from `def.Actions` where `act.Type == "navigate"`.

### AC-9: Tests

- **AC-9a:** A test asserts that rows rendered for `user` (no drilldown after this story) do NOT contain `▶`.
- **AC-9b:** A test asserts that pressing `Enter` on a resource with no drilldown returns the same model unchanged.
- **AC-9c:** A test asserts that pressing `h` (in the actions menu) on a `user` row triggers a `transitionDrilldown` that sets `m.currentRoot = "history-user-operation"` and `m.breadcrumb` includes `"View Operation Log"`.
- **AC-9d:** A test asserts that `buildActionsForRoot()` for `user` includes an item with `key="h"` and `label="View Operation Log →"`.

---

## Tasks

### Task 1: Config struct changes — `internal/config/config.go`

**File:** `internal/config/config.go`

1. Change `TableDef.Drilldown` from `[]DrillDownDef \`yaml:"drilldown,omitempty"\`` to `*DrillDownDef \`yaml:"drilldown,omitempty"\``.

2. Add fields to `ActionDef`:
   ```go
   Type   string `yaml:"type,omitempty"`   // "navigate" | "" (default: HTTP mutation)
   Target string `yaml:"target,omitempty"` // navigate: target resource key
   Param  string `yaml:"param,omitempty"`  // navigate: query param name
   Column string `yaml:"column,omitempty"` // navigate: source column for ID (default "id")
   ```

3. No changes to `DrillDownDef` — its fields (`Target`, `Param`, `Column`, `Label`, `TitleAttribute`) are already correct.

---

### Task 2: Update `o8n-cfg.yaml`

**File:** `o8n-cfg.yaml`

1. For every resource, convert `drilldown:` from a YAML sequence to an inline mapping. Change:
   ```yaml
   drilldown:
     - target: process-instance
       param: processDefinitionId
       column: id
       label: Instances
   ```
   to:
   ```yaml
   drilldown:
     target: process-instance
     param: processDefinitionId
     column: id
     label: Instances
   ```

2. Resources that lose their drilldown — delete the entire `drilldown:` block for:
   `external-task`, `job`, `batch`, `decision-definition`, `user`, `group`, `task`, `tenant`

3. Resources that keep one drilldown — retain only the canonical entry and remove the rest:
   - `process-definition`: keep `process-instance`, remove `history-process-instance`
   - `process-instance`: keep `process-variables`, remove `task`, `incident`, `variable-instance`
   - `deployment`: keep `process-definition`, remove `decision-definition`

4. Add `type: navigate` actions per the Navigate Actions table in the Design section above. Example for `user`:
   ```yaml
   actions:
     - key: l
       label: Unlock User
       method: POST
       path: /user/{id}/unlock
     - key: ctrl+d
       label: Delete User
       method: DELETE
       path: /user/{id}
       confirm: true
     - key: h
       label: View Operation Log
       type: navigate
       target: history-user-operation
       param: userId
       column: id
   ```

---

### Task 3: Extract drilldown helper — `internal/app/update.go`

**File:** `internal/app/update.go`

1. Extract lines `1119-1180` (the navigation state push + fetch dispatch) into a new method:
   ```go
   // executeDrilldown performs the full navigation-stack push and resource fetch
   // for the given drilldown definition, using the current table cursor as context.
   func (m *model) executeDrilldown(d *config.DrillDownDef) (model, tea.Cmd) {
       // ... extracted logic ...
   }
   ```
   The method must handle:
   - `m.prepareStateTransition(transitionDrilldown)`
   - Saving current state to `m.navigationStack`
   - The `selectedInstanceID`/`selectedDefinitionKey` switch (currently update.go:1143-1148)
   - Setting `m.currentRoot`, `m.viewMode`, `m.genericParams`, `m.breadcrumb`, `m.contentHeader`
   - Pre-setting target columns and returning `tea.Batch(m.fetchGenericCmd(d.Target), flashOnCmd(), m.saveStateCmd())`

2. Update the Enter/→ key handler (update.go:1077-1181) to use the pointer instead of a slice:
   ```go
   if def := m.findTableDef(currentTableKey); def != nil && def.Drilldown != nil {
       d := def.Drilldown
       colName := d.Column
       if colName == "" { colName = "id" }
       // resolve val from rowData ...
       return m.executeDrilldown(d)
   }
   ```
   Remove the "best match" loop (update.go:1079-1093) — there is now exactly one entry.

---

### Task 4: Navigation actions in `buildActionsForRoot` — `internal/app/nav.go`

**File:** `internal/app/nav.go`

1. In `buildActionsForRoot()` (nav.go:72), when iterating `td.Actions`, branch on `act.Type`:
   ```go
   for _, action := range td.Actions {
       act := action
       if act.Type == "navigate" {
           // Navigate action: trigger drilldown transition, no HTTP call
           items = append(items, actionItem{
               key:   act.Key,
               label: act.Label + " →",
               cmd: func(m *model) tea.Cmd {
                   colName := act.Column
                   if colName == "" { colName = "id" }
                   val := ""
                   cursor := m.table.Cursor()
                   if cursor >= 0 && cursor < len(m.rowData) {
                       if v, ok := m.rowData[cursor][colName]; ok && v != nil {
                           val = fmt.Sprintf("%v", v)
                       }
                   }
                   if val == "" { return nil }
                   d := &config.DrillDownDef{
                       Target: act.Target,
                       Param:  act.Param,
                       Column: colName,
                       Label:  act.Label,
                   }
                   newM, cmd := m.executeDrilldown(d)
                   *m = newM
                   return cmd
               },
           })
       } else {
           // HTTP mutation action (existing logic, unchanged)
           items = append(items, actionItem{...})
       }
   }
   ```
   Note: the returned `cmd` function pattern must match how `actionItem.cmd` is called at update.go:186-187 (`return m, item.cmd(&m)`).

2. Track whether the current resource has both mutation and navigate actions; pass a separator flag if needed (or detect it at render time — see Task 5).

---

### Task 5: Action menu visual separator — `internal/app/view.go`

**File:** `internal/app/view.go` — `renderActionsMenu()` (around line 1142)

1. During rendering of `m.actionsMenuItems`, detect the boundary between the last non-navigate item and the first navigate item by checking if an item's label ends with ` →` (or by adding an `isNavigate bool` field to `actionItem`).
   - Preferred: add `isNavigate bool` to `actionItem` struct in `model.go` to avoid label-suffix heuristics.

2. When transitioning from mutation items to navigate items, insert a separator row:
   ```
   ────────────────────────
   ```

3. Navigate action items are displayed as `[key] Label →` (the `→` is already in the label from Task 4).

---

### Task 6: `▶` indicator and column width — `update.go` and `table.go`

**`update.go:1650`**:
```go
// Before:
hasDrilldown := def != nil && len(def.Drilldown) > 0
// After:
hasDrilldown := def != nil && def.Drilldown != nil
```

**`table.go:38`**:
```go
// Before:
hasDrilldownPrefix := len(def.Drilldown) > 0
// After:
hasDrilldownPrefix := def.Drilldown != nil
```

---

### Task 7: Help screen — `internal/app/view.go`

**File:** `view.go` — `renderHelpScreen()` (~line 245)

1. Wrap the `Enter — drill down` line in a condition:
   ```go
   if def := m.findTableDef(m.currentRoot); def != nil && def.Drilldown != nil {
       lines = append(lines, fmt.Sprintf("%-10s %s", "Enter/→", "Drill down"))
   }
   ```

2. Add a **Views** section to the help screen that lists navigate actions for the current resource:
   ```go
   if def != nil {
       var navLines []string
       for _, act := range def.Actions {
           if act.Type == "navigate" {
               navLines = append(navLines, fmt.Sprintf("%-10s %s", act.Key, act.Label))
           }
       }
       if len(navLines) > 0 {
           viewsSection = "\nVIEWS (" + m.currentRoot + ")\n" +
               "──────────────────────────────────\n" +
               strings.Join(navLines, "\n")
       }
   }
   ```

---

### Task 8: Tests — `internal/app/`

**File:** `internal/app/drilldown_nav_test.go` (new)

1. **AC-9a** — No `▶` for `user`:
   Set up a model with `currentRoot = "user"`, `rowData` with one row. Call `buildColumnsFor("user", ...)` and verify the returned row does not start with `▶`.

2. **AC-9b** — Enter no-op on `user`:
   Send `enter` key msg on a model with `currentRoot = "user"`. Assert returned model has unchanged `currentRoot` and `navigationStack`.

3. **AC-9c** — `h` on `user` navigates to `history-user-operation`:
   Build model with `currentRoot = "user"`, one row in `rowData` with `{"id": "demo"}`. Open actions menu (Space), send `h`. Assert `m.currentRoot == "history-user-operation"` and `m.breadcrumb` ends with `"View Operation Log"`.

4. **AC-9d** — `buildActionsForRoot()` for `user` includes navigate action:
   Assert items include `{key: "h", label: "View Operation Log →"}`.

---

## Dev Notes

- **`actionItem` struct** (`model.go:231-235`): add `isNavigate bool` field so `renderActionsMenu()` can insert the separator without label-suffix heuristics.
- **`executeDrilldown` signature**: returns `(model, tea.Cmd)` — caller pattern is `newM, cmd := m.executeDrilldown(d); *m = newM; return *m, cmd` for use inside the `cmd func(m *model) tea.Cmd` closure in `buildActionsForRoot()`.
- **Existing tests** that reference `def.Drilldown[0]` or `len(def.Drilldown)` will need updating to use `def.Drilldown` (pointer) after the struct change.
- **`update.go:1143-1148`** — the `selectedInstanceID`/`selectedDefinitionKey` switch inside `executeDrilldown` must include the two original cases plus any new navigate-action targets that need context (e.g. `process-variables` still needs `selectedInstanceID = val`).
- **No changes needed to `prepareStateTransition()`** — the existing `transitionDrilldown` scope is correct for navigate actions too.
- **`h` key conflict check**: `h` is not currently used in any resource's actions. Verify no existing mutation action uses `h` before adding it to all resources.

---

## File List

| File | Change |
|------|--------|
| `internal/config/config.go` | `Drilldown *DrillDownDef`; `ActionDef` gains `Type`, `Target`, `Param`, `Column` |
| `internal/app/model.go` | `actionItem` gains `isNavigate bool` |
| `internal/app/update.go` | `executeDrilldown()` helper extracted; Enter handler uses pointer; `hasDrilldown` fix |
| `internal/app/nav.go` | `buildActionsForRoot()` handles `type: navigate` |
| `internal/app/view.go` | Menu separator; help screen drill-down condition; help Views section |
| `internal/app/table.go` | `hasDrilldownPrefix` pointer check |
| `o8n-cfg.yaml` | All drilldown format changes + 15 new navigate actions |
| `internal/app/drilldown_nav_test.go` | New test file for AC-9a through AC-9d |

---

## Change Log

- 2026-02-28: Story created by BMM UX Designer
- 2026-02-28: Refined with exact struct names, line numbers, and implementation detail

## Status

ready
