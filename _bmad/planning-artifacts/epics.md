---
stepsCompleted: ["step-01-validate-prerequisites", "step-02-design-epics", "step-03-create-stories", "step-04-final-validation"]
status: 'complete'
completedAt: '2026-03-03'
inputDocuments:
  - "_bmad/planning-artifacts/prd.md"
  - "_bmad/planning-artifacts/architecture.md"
---

# o8n - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for o8n, decomposing the requirements from the PRD and Architecture into implementable stories.

## Requirements Inventory

### Functional Requirements

FR1: Operator can navigate to any of the 35 configured resource types using the context switcher (`:` key)
FR2: Operator can browse a paginated table of resources in the current context
FR3: Operator can drill down from a parent resource to related child resources as configured in `o8n-cfg.yaml`
FR4: Operator can navigate back through the drill-down history level by level using Escape
FR5: Operator can jump directly to a specific level in the breadcrumb trail
FR6: Operator can execute any action configured for the current resource type on the selected row
FR7: Operator is prompted to confirm destructive actions before they are executed
FR8: Operator receives visible success or error feedback in the footer after an action completes
FR9: Operator can retry a failed job associated with an incident
FR10: Operator can set an annotation on an incident
FR11: All modal dialog types render with identical border style, padding, and button placement — produced by a shared modal factory, not per-type layout code
FR12: Operator can dismiss any modal by pressing Escape
FR13: Operator can confirm any modal by pressing Enter on the confirm action
FR14: Operator can interact with edit dialogs that validate input by type (string, integer, boolean, JSON)
FR15: Operator can see the primary available actions for the current view in the footer without opening a help screen
FR16: Operator can open a context-sensitive action menu via the Space key showing all available actions for the selected row
FR17: Operator can view the full key binding reference via the `?` key
FR18: Operator can switch between configured environments at any time
FR19: Operator can switch to any resource context using the `:` context switcher without leaving stale state
FR20: All navigation transitions (environment switch, context switch, drill-down, breadcrumb jump) clear prior view state completely
FR21: Application restores the last active context and environment on startup
FR22: Operator can claim an unassigned task
FR23: Operator can unclaim a task
FR24: Operator can complete a claimed task via a dialog that displays input variables read-only and allows editing output (form) variables
FR25: Task completion dialog supports variable types: string, integer, boolean
FR26: Operator can configure 2 or more named environments with distinct API URLs, credentials, and accent colors
FR27: Application reads resource types, columns, actions, and drilldown rules from `o8n-cfg.yaml` at startup
FR28: Contributor can add a new standard resource type by editing `o8n-cfg.yaml` without modifying Go source code
FR29: Operator can inspect process variables associated with a process instance
FR30: Operator can edit a process variable value inline with type validation
FR31: Operator can copy the selected resource row as YAML to the system clipboard
FR32: Operator can filter the current resource table by entering a search term
FR33: Operator can clear the active search filter and return to the full table
FR34: Operator can toggle auto-refresh to continuously update the current table view
FR35: Application renders without overflow or truncation of critical information at 120×20 terminal size
FR36: Application adapts column visibility and hint display when the terminal is narrower than 120 columns
FR37: Application handles terminal resize events without corrupting the layout
FR38: Operator can switch between available color skins
FR39: Operator can toggle vim-style key bindings in-session

### NonFunctional Requirements

NFR1: The UI responds to any key press within 100ms — no perceptible input lag
NFR2: API calls that take longer than 500ms surface a visible loading indicator; the UI does not appear frozen
NFR3: All API calls are asynchronous — network operations never block the application's event loop or UI rendering
NFR4: The application must not panic or crash on any malformed, partial, or unexpected API response
NFR5: All API errors surface as user-visible footer messages — no silent failures
NFR6: After a network timeout or connection failure, the application continues to accept user input and displays an error in the footer — no restart required
NFR7: Footer error and success messages auto-clear after 5 seconds and do not block further interaction
NFR8: Credentials (username, password) must never appear in log files, debug output, or clipboard operations
NFR9: `o8n-env.yaml` must be git-ignored and maintained at `chmod 600` file permissions at all times
NFR10: Application renders at 120×20 in VSCode integrated terminal and IntelliJ IDEA terminal without layout corruption, text overflow, or missing primary content
NFR11: Application operates in standard POSIX terminals (xterm, iTerm2, Alacritty) without modification — no missing key bindings, rendering artifacts, or color failures
NFR12: Application handles terminal resize events without corrupting the rendered layout
NFR13: Adding a standard resource type (table + columns + actions + drilldowns) requires only edits to `o8n-cfg.yaml` — no Go source changes
NFR14: The OpenAPI client in `internal/operaton/` remains auto-generated — no manual edits permitted; regenerated via `.devenv/scripts/generate-api-client.sh`
NFR15: The modal system is config-driven — new modal types are supported through the modal factory without hardcoded per-type Go logic

### Additional Requirements

- **Modal Factory (arch-critical):** `ModalConfig` struct + `renderModal()` factory must be extracted into `internal/app/modal.go` (new file). New modal types registered via config — no changes to render switch path.
- **State Transition Contract (arch-critical):** `prepareStateTransition(TransitionType)` is the mandatory gate for all navigation. `TransitionFull`, `TransitionDrillDown`, `TransitionPop` types must be defined and enforced across all nav paths in `internal/app/nav.go`.
- **Footer Hint Push Model (arch-critical):** `Hint{Key, Label, MinWidth, Priority}` struct + per-view hint functions extracted into `internal/app/hints.go` (new file). Footer renderer is stateless.
- **ModalActionMenu:** Space action dialog implemented as a factory-registered modal type (`ModalActionMenu`) — reads config actions at render time, single-char shortcuts dispatch actions, visual separator before navigate actions.
- **Brownfield — no starter bootstrap required:** Existing codebase is the foundation. No template scaffolding needed.
- **`specification.md` update obligation:** Post-sprint, `specification.md` MUST be updated to document the modal factory, hint push model, and `TransitionType` enum. This is part of the definition of done for each sprint task.
- **Async contract:** All API calls via `tea.Cmd` — no goroutines writing to model directly. No blocking of Bubble Tea event loop.
- **Generated client:** `internal/operaton/` is auto-generated and must never be edited manually. Regenerated via `.devenv/scripts/generate-api-client.sh`.
- **Credential isolation:** `o8n-env.yaml` git-ignored, `chmod 600`. Credentials must never appear in logs, debug output (`./debug/`), or clipboard operations.
- **Test coverage target:** 80%+ line coverage on new sprint files (`modal.go`, `hints.go`). Verify with `make cover`.

### FR Coverage Map

FR1: Epic 3 — Navigate to any of 35 resource types via context switcher
FR2: Epic 3 — Browse paginated table of resources
FR3: Epic 3 — Drill down to child resources per config
FR4: Epic 3 — Navigate back through drill-down history (Escape)
FR5: Epic 3 — Jump directly to breadcrumb level
FR6: Epic 3 — Execute any configured action on selected row
FR7: Epic 3 — Confirm destructive actions before execution
FR8: Epic 3 — Success/error feedback in footer after action
FR9: Epic 3 — Retry failed job associated with incident
FR10: Epic 3 — Set annotation on incident
FR11: Epic 1 — All modals render from shared factory (identical border, padding, buttons)
FR12: Epic 1 — Dismiss any modal with Escape
FR13: Epic 1 — Confirm any modal with Enter
FR14: Epic 1 — Edit dialogs validate input by type (string, integer, boolean, JSON)
FR15: Epic 2 — Primary actions visible in footer without help screen
FR16: Epic 2 — Ctrl+Space opens context-sensitive action menu (Space reserved for row selection)
FR17: Epic 2 — Full key binding reference via `?` key
FR18: Epic 1 — Switch between configured environments at any time
FR19: Epic 1 — Switch resource context via `:` without stale state
FR20: Epic 1 — All navigation transitions clear prior view state completely
FR21: Epic 1 — Restore last active context and environment on startup
FR22: Epic 3 — Claim an unassigned task
FR23: Epic 3 — Unclaim a task
FR24: Epic 3 — Complete task via dialog (input read-only, output editable)
FR25: Epic 3 — Task completion dialog supports string, integer, boolean variable types
FR26: Epic 5 — Configure 2+ named environments with distinct URLs, credentials, accent colors
FR27: Epic 5 — Read resource types, columns, actions, drilldowns from `o8n-cfg.yaml` at startup
FR28: Epic 5 — Add standard resource type via `o8n-cfg.yaml` only (no Go changes)
FR29: Epic 3 — Inspect process variables associated with a process instance
FR30: Epic 3 — Edit process variable value inline with type validation
FR31: Epic 3 — Copy selected resource row as YAML to system clipboard
FR32: Epic 3 — Filter current resource table by search term
FR33: Epic 3 — Clear active search filter and return to full table
FR34: Epic 3 — Toggle auto-refresh for continuous table updates
FR35: Epic 4 — Render without overflow or truncation at 120×20
FR36: Epic 4 — Adapt column visibility and hint display below 120 columns
FR37: Epic 4 — Handle terminal resize events without layout corruption
FR38: Epic 4 — Switch between color skins
FR39: Epic 4 — Toggle vim-style key bindings in-session

NFR1: Epic 3 — UI responds to any key press within 100ms
NFR2: Epic 3 — API calls >500ms surface loading indicator
NFR3: Epic 3 — All API calls asynchronous, never block event loop
NFR4: Epic 1 — No panic or crash on malformed API response
NFR5: Epic 3 — All API errors surface as footer messages
NFR6: Epic 3 — Application continues after network failure, no restart required
NFR7: Epic 3 — Footer messages auto-clear after 5 seconds
NFR8: Epic 5 — Credentials never appear in logs, debug output, or clipboard
NFR9: Epic 5 — `o8n-env.yaml` git-ignored, `chmod 600` at all times
NFR10: Epic 4 — Renders at 120×20 in VSCode and IntelliJ IDEA without corruption
NFR11: Epic 4 — Operates in standard POSIX terminals without modification
NFR12: Epic 4 — Handles terminal resize without layout corruption
NFR13: Epic 5 — Standard resource type addition requires only `o8n-cfg.yaml` edits
NFR14: Epic 5 — `internal/operaton/` remains auto-generated, never manually edited
NFR15: Epic 1 — Modal system is config-driven via factory, no per-type hardcoded logic

## Epic List

### Epic 1: Consistent & Reliable Modal System
Operators can interact with any modal type — confirm, delete, help, edit, sort, environment switch — with completely consistent UX: identical border styling, Esc always dismisses, Enter always confirms, all produced by a shared factory. Navigation between resources and environments never leaves stale state behind.
**FRs covered:** FR11, FR12, FR13, FR14, FR18, FR19, FR20, FR21
**NFRs covered:** NFR4, NFR15
**Architecture deliverables:** Modal factory extraction (`internal/app/modal.go`), state transition contract enforcement (`internal/app/nav.go`)

### Epic 2: Discoverable Actions
Operators can see the primary available actions for any view directly in the footer — no help screen required. `Ctrl+Space` opens a context-sensitive action menu showing all configured actions for the selected row.
**FRs covered:** FR15, FR16, FR17
**Architecture deliverables:** Footer hint push model (`internal/app/hints.go` + `view.go` wiring), `ModalActionMenu` registration in `internal/app/update.go`

### Epic 3: Core Operational Workflows
Operators can navigate all 35 resource types, execute configured actions with success/error feedback, drill down through resource hierarchies, handle incidents (retry jobs, set annotations), claim and complete tasks with typed form variables, inspect and edit process variables, and search/filter any table.
**FRs covered:** FR1, FR2, FR3, FR4, FR5, FR6, FR7, FR8, FR9, FR10, FR22, FR23, FR24, FR25, FR29, FR30, FR31, FR32, FR33, FR34
**NFRs covered:** NFR1, NFR2, NFR3, NFR5, NFR6, NFR7

### Epic 4: Stable Rendering & Visual Polish
The application renders correctly at 120×20 in VSCode and IntelliJ IDEA terminals, adapts gracefully to narrower viewports (column and hint visibility), handles terminal resize events without layout corruption, and gives operators access to color skins and vim-style key bindings.
**FRs covered:** FR35, FR36, FR37, FR38, FR39
**NFRs covered:** NFR10, NFR11, NFR12

### Epic 5: Configuration, Security & Documentation
Operators can configure 2+ named environments with distinct credentials and accent colors; contributors can add standard resource types via `o8n-cfg.yaml` alone. Credentials are never exposed in logs, debug output, or clipboard. `specification.md` accurately reflects post-sprint behavior.
**FRs covered:** FR26, FR27, FR28
**NFRs covered:** NFR8, NFR9, NFR13, NFR14
**Additional:** `specification.md` update is part of the definition of done for all sprint tasks

---

## Epic 1: Consistent & Reliable Modal System

Operators can interact with any modal type — confirm, delete, help, edit, sort, environment switch — with completely consistent UX: identical border styling, Esc always dismisses, Enter always confirms, all produced by a shared factory. Navigation between resources and environments never leaves stale state behind.

### Story 1.1: Modal Factory Foundation

As a **developer contributing to o8n**,
I want a `ModalConfig` struct and `renderModal()` factory function in `internal/app/modal.go`,
So that all modal types are rendered from a single, consistent code path with no per-type layout logic in the view render path.

**Acceptance Criteria:**

**Given** `ModalConfig` is defined with: `sizeHint` (OverlayCenter / FullScreen), `title string`, `bodyRenderer func(m Model) string`, and button label fields
**When** `renderModal(m Model, cfg ModalConfig)` is called for any registered modal type
**Then** the rendered modal has identical border style, padding, and button placement regardless of which type it is
**And** the view render path contains no `switch modalType` statement for modal body content
**And** new modal types are added by defining a new `ModalConfig` — no changes to `renderModal()` required
**And** all existing modal types (confirm, delete, help, edit, sort, environment) are migrated to use the factory
**And** `internal/app/main_modal_test.go` is created covering the factory with ≥80% line coverage on `modal.go`

### Story 1.2: State Transition Contract

As a **developer contributing to o8n**,
I want a `TransitionType` enum and `prepareStateTransition(t TransitionType)` function enforced across all navigation paths in `internal/app/nav.go`,
So that every navigation action uses a single, auditable gate that eliminates state leakage bugs.

**Acceptance Criteria:**

**Given** `TransitionType` is defined with constants: `TransitionFull`, `TransitionDrillDown`, `TransitionPop`
**When** any navigation action is taken (environment switch, context switch, drill-down, Esc, breadcrumb jump)
**Then** `prepareStateTransition(transitionType)` is called before any view state is modified
**And** `TransitionFull` clears: `activeModal`, `footerError`, `searchQuery`, `searchActive`, `sortColumn`, `sortDirection`, `tableCursor`, `navigationStack`
**And** `TransitionDrillDown` pushes current `viewState` onto `navigationStack`, then clears non-stack state fields
**And** `TransitionPop` pops `viewState` from `navigationStack` and restores all captured state — no field clearing
**And** a code review of all nav paths confirms no navigation code bypasses `prepareStateTransition`

### Story 1.3: Environment & Context Switching Correctness

As an **operator**,
I want switching between environments and resource contexts to always start fresh with no stale data,
So that I can trust the view I see reflects exactly the environment and context I selected.

**Acceptance Criteria:**

**Given** the operator is in any resource context with an active filter, modal, or cursor position
**When** the operator switches to a different environment
**Then** `prepareStateTransition(TransitionFull)` is called, clearing all prior view state
**And** the new environment's accent color and API URL are active immediately
**And** the table reloads from the new environment with no rows, filters, or modal from the previous environment visible

**Given** the operator is in any resource context with active state
**When** the operator uses `:` to switch to a different resource type
**Then** `prepareStateTransition(TransitionFull)` is called before the new context loads
**And** the new context view has no residual search query, sort order, modal, or cursor from the previous context

### Story 1.4: Consistent Modal UX (Esc/Enter + Type Validation)

As an **operator**,
I want Esc to always dismiss any modal and Enter to always confirm, and edit dialogs to validate input by type,
So that modal interaction is fully predictable across every modal in the application.

**Acceptance Criteria:**

**Given** any modal is active (confirm, delete, help, edit, sort, environment)
**When** the operator presses Esc
**Then** the modal is dismissed and the underlying view is restored — no action is executed

**Given** any modal with a confirm action is active
**When** the operator presses Enter on the confirm button
**Then** the confirm action is executed

**Given** an edit dialog is open for a typed field (integer, boolean, JSON, string)
**When** the operator enters a value that does not match the declared type (e.g., "abc" for integer)
**Then** the dialog displays an inline validation error and does not submit the value
**And** when a valid value is entered and Enter pressed, the value is accepted and saved

### Story 1.5: Startup State Restoration & API Resilience

As an **operator**,
I want the application to restore my last context and environment on startup and never crash on unexpected API responses,
So that each session starts where I left off and operational incidents don't require tool restarts.

**Acceptance Criteria:**

**Given** `o8n-stat.yaml` contains the last active context and environment from a previous session
**When** the application starts
**Then** it navigates directly to the last active context in the last active environment and loads data automatically

**Given** no previous state file exists (first run)
**When** the application starts
**Then** it opens with a sensible default context without error

**Given** the API returns a malformed, partial, or empty JSON response for any resource type
**When** the application processes the response
**Then** an `errMsg` is produced and displayed in the footer
**And** the application remains fully interactive — no freeze, no panic, no restart required
**And** the footer error auto-clears after 5 seconds

---

## Epic 2: Discoverable Actions

Operators can see the primary available actions for any view directly in the footer — no help screen required. `Ctrl+Space` opens a context-sensitive action menu showing all configured actions for the selected row.

### Story 2.1: Footer Hint Push Model

As an **operator**,
I want the primary available actions for the current view visible in the footer at all times,
So that I can discover what I can do without opening the `?` help screen.

**Acceptance Criteria:**

**Given** `Hint{Key, Label, MinWidth, Priority}` is defined in `internal/app/hints.go`
**When** the footer renderer is called
**Then** it receives a `[]Hint` slice from `currentViewHints(m)` — the footer renderer itself contains no hint declarations
**And** hints with `MinWidth: 0` are always visible regardless of terminal width
**And** hints with `MinWidth > 0` are hidden when the terminal is narrower than that threshold
**And** when multiple hints must be dropped due to width, higher `Priority` integer hints are dropped first (lower int = higher priority)
**And** `internal/app/main_hints_test.go` is created with tests for hint visibility filtering logic (≥80% coverage on `hints.go`)

### Story 2.2: Per-View Hint Functions

As an **operator**,
I want each view to declare its own available actions as hints in the footer,
So that the footer always reflects the actions relevant to my current context.

**Acceptance Criteria:**

**Given** the operator is in the main table view
**When** the footer renders
**Then** primary actions for that view (e.g., Enter to drill down, `Ctrl+Space` for actions, `/` to filter, `r` for refresh, `?` for help) appear as key hints in the footer

**Given** the operator is in the context switcher, a modal, or the search input
**When** the footer renders
**Then** the hints reflect the actions available in that specific sub-view, not the table's actions

**Given** the terminal width is reduced below a hint's `MinWidth` threshold
**When** the footer renders
**Then** that hint is omitted from the footer without breaking the layout of remaining hints

### Story 2.3: Action Dialog via Ctrl+Space (ModalActionMenu)

As an **operator**,
I want pressing `Ctrl+Space` on any table row to open a context-sensitive action menu showing all configured actions for that resource,
So that I can discover and execute any available action without memorising every key binding, while Space remains available for row selection.

**Acceptance Criteria:**

**Given** `ModalActionMenu` is registered as a factory modal type using `ModalConfig`
**When** the operator presses `Ctrl+Space` on a table row
**Then** `ModalActionMenu` opens as a centered overlay listing all configured actions for the current resource type from `o8n-cfg.yaml`
**And** mutation actions (HTTP verbs) are listed first
**And** a visual separator appears before the first `type: navigate` action
**And** navigate actions display a `→` suffix
**And** `[y] View as YAML` is always the last item

**Given** `ModalActionMenu` is open
**When** the operator presses a single-character shortcut matching an action's configured key
**Then** the action is dispatched immediately without requiring cursor movement or Enter
**And** `Up`/`Down` moves the cursor; Enter dispatches the highlighted action
**And** Esc closes the menu without executing any action

**Note:** `Space` (without Ctrl) is reserved for row selection — it must not trigger the action menu.

### Story 2.4: Full Key Binding Reference

As an **operator**,
I want to open the full key binding reference via the `?` key,
So that I can discover any binding I've forgotten without leaving the application.

**Acceptance Criteria:**

**Given** the operator is in any view
**When** the operator presses `?`
**Then** the help modal opens as a full-screen overlay displaying all key bindings organized by category

**Given** the help modal is open
**When** the operator presses Esc or `?` again
**Then** the help modal closes and the operator returns to the previous view unchanged
**And** the help modal is rendered via the modal factory (`ModalConfig` with `FullScreen` size hint)

---

## Epic 3: Core Operational Workflows

Operators can navigate all 35 resource types, execute configured actions with success/error feedback, drill down through resource hierarchies, handle incidents (retry jobs, set annotations), claim and complete tasks with typed form variables, inspect and edit process variables, and search/filter any table.

### Story 3.1: Resource Navigation & Context Switcher

As an **operator**,
I want to navigate to any of the 35 configured resource types using the context switcher,
So that I can reach any operational view in seconds from anywhere in the application.

**Acceptance Criteria:**

**Given** 35 resource types are defined in `o8n-cfg.yaml`
**When** the operator presses `:` to open the context switcher
**Then** a searchable list of all configured resource types is displayed

**Given** the context switcher is open
**When** the operator types a partial name and selects a match
**Then** `prepareStateTransition(TransitionFull)` is called and the selected resource type loads with a fresh, paginated table
**And** the breadcrumb shows the new context name

**Given** the selected resource type returns data from the API
**When** the table loads
**Then** the correct columns for that resource type (as defined in `o8n-cfg.yaml`) are displayed and the first row is selected

### Story 3.2: Drill-Down Navigation & Breadcrumb

As an **operator**,
I want to drill down from a parent resource into related child resources and navigate back via Escape or breadcrumb,
So that I can traverse resource hierarchies (e.g., process definition → instances → variables) without losing my place.

**Acceptance Criteria:**

**Given** the current resource type has drilldown rules configured in `o8n-cfg.yaml`
**When** the operator presses Enter on a row
**Then** `prepareStateTransition(TransitionDrillDown)` is called, the current view is pushed onto the navigation stack, and the child resource type loads filtered to the selected parent

**Given** the operator has drilled down one or more levels
**When** the operator presses Escape
**Then** `prepareStateTransition(TransitionPop)` is called and the previous view is fully restored (rows, cursor, columns, filters, breadcrumb)

**Given** the operator has drilled down multiple levels and the breadcrumb shows level labels
**When** the operator selects a specific breadcrumb level
**Then** the application pops to that level and restores the view state captured at that level

### Story 3.3: Action Execution with Feedback

As an **operator**,
I want to execute any configured action on a selected row and receive clear success or error feedback,
So that I always know whether an action succeeded or failed without guessing.

**Acceptance Criteria:**

**Given** the operator selects a row and presses the configured action key
**When** the action executes successfully
**Then** a success message appears in the footer (e.g., `✓ Job retried`)
**And** the table refreshes to reflect the updated state

**Given** a destructive action key is pressed (e.g., delete)
**When** the operator presses the key
**Then** a confirmation modal is shown before the action executes
**And** confirming executes the action; Esc cancels without any side effect

**Given** an action's API call fails
**When** the error is received
**Then** an error message appears in the footer and auto-clears after 5 seconds
**And** no silent failure occurs

### Story 3.4: Incident Operations (Retry & Annotate)

As an **operator** (Alex persona),
I want to retry a failed job and set an annotation on an incident from within the TUI,
So that I can resolve incidents without switching to a browser or writing curl commands.

**Acceptance Criteria:**

**Given** the operator is viewing the incidents resource type and selects a row
**When** the operator executes the configured Retry action
**Then** the corresponding job retry API call is made and the footer shows `✓ Job retried`
**And** the incident table refreshes

**Given** the operator is viewing an incident row
**When** the operator executes the configured Annotate action
**Then** an edit dialog opens for the annotation text field
**And** on confirmation, the annotation is set via the API and the footer confirms success

**Given** the operator drills down from an incident to its process instance
**When** Enter is pressed on the incident row
**Then** the process instance view loads filtered to the incident's process instance ID

### Story 3.5: Task Claim, Unclaim & Complete

As an **operator** (Priya persona),
I want to claim, unclaim, and complete user tasks via the TUI including a form variable dialog,
So that I can process my task queue without switching to Operaton Cockpit.

**Acceptance Criteria:**

**Given** the operator is viewing the task resource type and selects an unassigned task
**When** the operator executes the Claim action
**Then** the task is claimed via the API, the assignee column updates immediately, and the footer confirms

**Given** the operator has a claimed task selected
**When** the operator executes the Unclaim action
**Then** the task is unclaimed via the API and the assignee column clears

**Given** the operator has a claimed task selected and executes the Complete action
**When** the task completion dialog opens
**Then** input variables with no corresponding output variable are displayed read-only
**And** output (form) variables are displayed as editable fields with type validation (string, integer, boolean)
**And** output variables whose name matches an input variable are pre-populated with that input variable's value — the input variable is not shown separately as a duplicate read-only row
**And** output variables with no matching input variable are presented as empty editable fields
**And** pressing Enter on `[Complete]` submits the completion and the task disappears from the table
**And** the footer confirms: `✓ Completed: [task name]`

### Story 3.6: Process Variable Inspection & Editing

As an **operator**,
I want to inspect and edit process variables associated with a process instance,
So that I can diagnose and correct stuck processes directly from the TUI.

**Acceptance Criteria:**

**Given** the operator drills down to the variables view for a process instance
**When** the variables table loads
**Then** all process variables for that instance are displayed with name, type, and value columns

**Given** the operator selects a variable row and presses the edit key
**When** the edit dialog opens
**Then** the current variable value is pre-populated and the type is displayed
**And** input is validated against the variable type (string, integer, boolean, JSON)
**And** on confirmation, the variable is updated via the API and the table refreshes

**Given** the operator presses `y` on any table row
**When** the copy action executes
**Then** the row is copied as YAML to the system clipboard and the footer confirms

### Story 3.7: Search, Filter & Auto-Refresh

As an **operator**,
I want to filter the current resource table by search term and toggle auto-refresh,
So that I can efficiently find specific resources and monitor live state.

**Acceptance Criteria:**

**Given** the operator is viewing any resource table
**When** the operator presses `/` and types a search term
**Then** the table is filtered to rows matching the term in real time
**And** the footer or header indicates the active filter

**Given** an active search filter is set
**When** the operator presses Escape or the clear-filter key
**Then** the filter is cleared and the full table is restored

**Given** the operator presses `r` to toggle auto-refresh
**When** auto-refresh is enabled
**Then** the table reloads on the configured interval and a visual indicator (pulsing badge) appears in the footer
**And** pressing `r` again disables auto-refresh and the indicator disappears

**Given** any API call takes longer than 500ms
**When** the request is in flight
**Then** a loading indicator appears in the footer center column and the UI remains fully interactive

---

## Epic 4: Stable Rendering & Visual Polish

The application renders correctly at 120×20 in VSCode and IntelliJ IDEA terminals, adapts gracefully to narrower viewports (column and hint visibility), handles terminal resize events without layout corruption, and gives operators access to color skins and vim-style key bindings.

### Story 4.1: 120×20 Rendering Validation

As an **operator**,
I want the application to render without overflow or truncation of critical information at the 120×20 minimum viewport,
So that o8n is fully usable in VSCode's integrated terminal and IntelliJ IDEA without manual resizing.

**Acceptance Criteria:**

**Given** the terminal is sized to exactly 120 columns × 20 rows
**When** the main table view renders in VSCode integrated terminal and IntelliJ IDEA terminal
**Then** the header, table body, and footer are all visible with no overflow, truncation of critical content, or layout corruption
**And** at least the primary columns for each resource type are visible at this width
**And** the breadcrumb, environment indicator, and footer hints all render within their allocated rows without wrapping into the table area

**Given** 120×20 rendering is validated
**When** the test is documented
**Then** a test or checklist entry covers each of the two target terminals (VSCode, IntelliJ IDEA)

### Story 4.2: Responsive Column & Hint Visibility

As an **operator**,
I want the application to gracefully adapt when the terminal is narrower than 120 columns,
So that the most important information remains visible even in constrained terminal widths.

**Acceptance Criteria:**

**Given** the terminal is narrower than 120 columns
**When** the table renders
**Then** columns are hidden in `hide_order` sequence (lowest-priority columns hidden first, as defined in `o8n-cfg.yaml`)
**And** the table never renders with truncated column headers or overflowing cell content

**Given** the terminal is narrower than a hint's `MinWidth` threshold
**When** the footer renders
**Then** that hint is omitted cleanly; remaining hints are not shifted or broken
**And** hints with `MinWidth: 0` (always-show) remain visible at any width

**Given** the terminal is very narrow (below ~80 columns)
**When** the application renders
**Then** at minimum the row cursor, resource name column, and footer error/status are visible — the UI does not crash or produce garbled output

### Story 4.3: Terminal Resize Handling

As an **operator**,
I want the application to adapt cleanly when I resize my terminal window,
So that layout corruption doesn't interrupt my workflow.

**Acceptance Criteria:**

**Given** the application is running and displaying any view
**When** the operator resizes the terminal window (larger or smaller)
**Then** Bubble Tea's `tea.WindowSizeMsg` is handled and the layout reflows correctly within one render cycle
**And** no text overflows, no borders break, and no content from a previous size persists as artifacts
**And** the table, header, and footer proportions are recalculated based on the new terminal dimensions

**Given** the terminal is resized to below the 120×20 minimum
**When** the application renders
**Then** it degrades gracefully per Story 4.2 — no panic, no corruption

### Story 4.4: Color Skins

As an **operator**,
I want to switch between the available color skins to match my terminal theme or differentiate environments,
So that I can customise the visual experience and use color as an environment signal.

**Acceptance Criteria:**

**Given** 35 built-in skin files exist in `skins/`
**When** the operator switches skins using the skin selection key
**Then** the active skin is applied immediately without a restart
**And** all UI elements use semantic color roles from the new skin — no hardcoded colors remain
**And** the active skin name is persisted to `o8n-stat.yaml` and restored on next startup

**Given** `ui_color` is set in `o8n-env.yaml` for the active environment
**When** the skin is applied
**Then** `ui_color` overrides the border accent and footer breadcrumb background regardless of which skin is active

### Story 4.5: Vim-Style Key Bindings Toggle

As an **operator** (Marco persona),
I want to toggle vim-style key bindings on and off in-session without restarting,
So that keyboard-native operators can use familiar `j`/`k`/`g`/`G` navigation alongside the default bindings.

**Acceptance Criteria:**

**Given** the application is running
**When** the operator presses `V` to toggle vim mode
**Then** vim-style bindings (`j`/`k` for up/down, `g`/`G` for top/bottom) become active immediately
**And** the footer or status area indicates that vim mode is on
**And** pressing `V` again restores the default bindings

**Given** vim mode is active and the operator presses `j` or `k`
**When** the table has multiple rows
**Then** the cursor moves down or up respectively — identical to arrow key behavior

**Given** `--vim` flag is passed at startup
**When** the application initializes
**Then** vim mode is active from the first rendered frame

---

## Epic 5: Configuration, Security & Documentation

Operators can configure 2+ named environments with distinct credentials and accent colors; contributors can add standard resource types via `o8n-cfg.yaml` alone. Credentials are never exposed in logs, debug output, or clipboard. `specification.md` accurately reflects post-sprint behavior.

### Story 5.1: Multi-Environment Configuration

As an **operator**,
I want to configure 2 or more named environments with distinct API URLs, credentials, and accent colors,
So that I can operate across local, staging, and production without editing config between sessions.

**Acceptance Criteria:**

**Given** `o8n-env.yaml` contains 2 or more named environment entries, each with `name`, `api_url`, `username`, `password`, and `ui_color`
**When** the application starts
**Then** all configured environments are available in the environment switcher
**And** the active environment's `ui_color` is applied to the border accent and footer breadcrumb background

**Given** the operator switches to a different environment
**When** the switch completes
**Then** all subsequent API calls use the new environment's `api_url` and credentials
**And** the environment name and accent color update immediately in the UI

**Given** a single environment is configured
**When** the application starts
**Then** it loads normally — multi-environment is not required for the app to function

### Story 5.2: Credential Security

As an **operator**,
I want credentials to be isolated in a git-ignored, permission-restricted file and never appear in logs, debug output, or clipboard operations,
So that sensitive API credentials cannot leak into version control or observability tooling.

**Acceptance Criteria:**

**Given** `o8n-env.yaml` contains credentials
**When** the file is checked
**Then** it is listed in `.gitignore` and carries `chmod 600` file permissions

**Given** the application is running with `--debug` flag
**When** debug output is written to `./debug/o8n.log` and `./debug/last-screen.txt`
**Then** no credentials (username, password, API URL with embedded auth) appear in any debug file

**Given** the operator presses `y` to copy a row as YAML to the clipboard
**When** the clipboard content is inspected
**Then** no credential fields from `o8n-env.yaml` are present in the clipboard content

**Given** any API error is displayed in the footer
**When** the error message is rendered
**Then** it contains no credential values — only the error type and message from the API response

### Story 5.3: Config-Driven Resource Extensibility

As a **contributor** (Marco persona),
I want to add a new standard resource type by editing only `o8n-cfg.yaml`,
So that the community can extend o8n's resource coverage without Go source code changes.

**Acceptance Criteria:**

**Given** a contributor adds a new table entry to `o8n-cfg.yaml` with columns, actions, and drilldown rules
**When** the application is built and started
**Then** the new resource type appears in the context switcher (`:`)
**And** the table loads with the defined columns
**And** configured actions execute the correct API calls
**And** configured drilldown rules navigate to the child resource type
**And** no Go source file was modified to achieve this

**Given** `internal/operaton/` needs updating for a new API endpoint
**When** the contributor regenerates the client
**Then** `.devenv/scripts/generate-api-client.sh` produces the updated client without manual file edits to `internal/operaton/`

### Story 5.4: Specification & Documentation Accuracy

As a **contributor** (Marco persona),
I want `specification.md` to accurately reflect the post-sprint implementation,
So that I can understand the full system from documentation alone without reading source code.

**Acceptance Criteria:**

**Given** the modal factory pattern has been implemented
**When** `specification.md` is reviewed
**Then** it documents: the `ModalConfig` struct fields, the `renderModal()` factory signature, and the rule that new modal types are added via config registration — not via switch statements

**Given** the footer hint push model has been implemented
**When** `specification.md` is reviewed
**Then** it documents: the `Hint` struct fields, the `currentViewHints(m)` dispatch pattern, and `MinWidth`/`Priority` semantics

**Given** the state transition contract has been implemented
**When** `specification.md` is reviewed
**Then** it documents: the `TransitionType` enum values and the rule that `prepareStateTransition` is mandatory for all navigation changes

**Given** the `Ctrl+Space` action dialog has been implemented
**When** `specification.md` is reviewed
**Then** it documents: `ModalActionMenu`, its `Ctrl+Space` trigger, the action list ordering (mutations first, separator, navigate with `→`), and single-char shortcut dispatch

**Note:** `specification.md` updates are part of the definition of done for each sprint task throughout all epics — this story serves as the final verification pass.
