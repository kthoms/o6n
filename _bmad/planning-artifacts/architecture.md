---
stepsCompleted: ["step-01-init", "step-02-context", "step-03-starter", "step-04-decisions", "step-05-patterns", "step-06-structure", "step-07-validation", "step-08-complete"]
lastStep: 8
status: 'complete'
completedAt: '2026-03-03'
inputDocuments:
  - "_bmad/planning-artifacts/prd.md"
  - "_bmad/planning-artifacts/prd-validation-report.md"
  - "CLAUDE.md"
  - "README.md"
  - "specification.md"
workflowType: 'architecture'
project_name: 'o6n'
user_name: 'Karsten'
date: '2026-03-03'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**
39 FRs across 10 capability areas. The sprint-critical FRs are the modal system (FR11-14),
action discoverability (FR15-17), state management (FR18-21), and rendering (FR35-37). FR28
and NFR13/15 define the config-driven architecture mandate — new resources and modal types
must not require Go source changes.

**Non-Functional Requirements:**
- Performance: 100ms UI response ceiling; 500ms threshold for loading indicator; strictly
  async API calls — no blocking the event loop (NFR1-3)
- Reliability: Panic-safe render with defer/recover; all errors surface as footer messages,
  never silently swallowed; 5s auto-clear (NFR4-7)
- Security: Credentials isolated to git-ignored `o6n-env.yaml` at 0600 perms; never appear
  in logs, debug output, or clipboard (NFR8-9)
- Terminal Compatibility: 120×20 primary target in VSCode and IntelliJ IDEA; column
  hide-order and priority-based hint system handle narrower widths; resize events must not
  corrupt layout (NFR10-12)
- Maintainability: Config-driven resource types (no Go changes), auto-generated API client
  (never edit), config-driven modal factory (NFR13-15)

**Scale & Complexity:**
- Primary domain: TUI/CLI — Go 1.24+, Bubble Tea ecosystem (Bubbles, Lipgloss)
- Complexity level: Medium (brownfield, rich feature set, active specification)
- Estimated architectural components: 8 internal packages + config layer + skin system

### Technical Constraints & Dependencies

- **Bubble Tea (Elm-inspired):** All state mutations are pure; all async operations return
  `tea.Cmd`; no goroutines may write to model directly — event loop is single-threaded
- **Generated API client:** `internal/operaton/` is regenerated from
  `resources/operaton-rest-api.json` — never edited manually
- **Config authority:** `o6n-cfg.yaml` is the single source of truth for resource types,
  columns, actions, drilldowns; architecture must not duplicate this authority in Go
- **Static binary:** `CGO_ENABLED=0`, `GO_TAGS=netgo` — no dynamic linking
- **Single binary, no subcommands:** All behavior is runtime-configured, not compile-time
- **specification.md is authoritative:** All architectural decisions must be documentable in
  and compatible with the existing specification structure

### Cross-Cutting Concerns Identified

1. **State Transition Contract** — `prepareStateTransition` must be called on ALL navigation
   transitions (context switch, environment switch, drill-down, breadcrumb jump, Esc);
   state leakage is the primary defect category
2. **Async/Concurrency Safety** — All API calls via `tea.Cmd`; shared cache access via
   `sync.RWMutex`; no ad-hoc goroutine writes to model
3. **Config-Driven Design** — Modal factory, resource definitions, action routing; any
   per-type hardcoded logic is an architecture violation
4. **Error Isolation** — `defer/recover` in view render functions; `errMsg` propagation
   through Update; footer as the sole error display surface
5. **Responsive Layout** — Column `hide_order`, hint priority system (1-9), width
   thresholds; layout decisions must use the established sizing contract
6. **Credential Security** — `o6n-env.yaml` must never appear in git, logs, or clipboard;
   Basic Auth is the only supported mechanism
7. **Theming** — 26 semantic color roles (including `env_name`); `env_name` per skin
   governs the environment label color in the fixed top-right header — the primary
   environment identity signal; `ui_color` per environment overrides border accent only
   (secondary accent); skin changes must not require restarts

## Starter Template Evaluation

### Primary Technology Domain

TUI/CLI tool (brownfield) — established Go codebase

### Starter Options Considered

This is a brownfield project. No starter template selection was required; the technology
stack is fully established and documented in `specification.md`. The evaluation below
documents the existing foundation for AI agent consistency.

### Established Foundation: o6n Codebase

**Rationale:** Existing codebase with active development, comprehensive specification,
and stable architecture. No migration or starter overlay warranted.

**Language & Runtime:** Go 1.24+. Static binary build with `CGO_ENABLED=0` and
`GO_TAGS=netgo`. Single binary, no subcommands, no dynamic linking.

**UI Framework:** Charmbracelet ecosystem — Bubble Tea (Model/Update/View event loop),
Bubbles (reusable TUI components), Lipgloss (layout and styling). All state mutations are
pure; all async operations return `tea.Cmd`.

**API Client:** Auto-generated from `resources/operaton-rest-api.json` using
OpenAPI Generator (Docker). Lives in `internal/operaton/` — never edited manually.
Regenerated via `.devenv/scripts/generate-api-client.sh`.

**Configuration:** Three-file YAML split:
- `o6n-env.yaml` — credentials, URLs, accent colors (git-ignored, 0600)
- `o6n-cfg.yaml` — tables, columns, actions, drilldowns (version-controlled)
- `o6n-stat.yaml` — runtime state (auto-managed, git-ignored)

**Theming:** 35 built-in skins as YAML files in `skins/`. 26 semantic color roles —
no hardcoded color values. `env_name` semantic role governs the environment label color
in the fixed top-right header (primary environment signal). Environment `ui_color`
overrides border accent only (secondary accent).

**Testing:** Go stdlib `testing` + `httptest.NewServer()` for API tests; model state
assertion pattern for UI tests; temp files with `defer os.Remove()` for config tests.

**Code Organization:**
```
internal/
  app/           — TUI logic (model, update, view, nav, commands, skin, styles, table, edit)
  client/        — HTTP client wrapper around generated API client
  config/        — Config structs and loaders (LoadSplitConfig, LoadEnvConfig, Save)
  dao/           — DAO interfaces (DAO, HierarchicalDAO, ReadOnlyDAO)
  validation/    — Input validation (bool/int/float/json/text/auto)
  contentassist/ — Thread-safe suggestion cache (sync.RWMutex)
  operaton/      — Auto-generated OpenAPI client (do not edit)
skins/           — 35 YAML color theme files
resources/       — OpenAPI specification (operaton-rest-api.json)
```

**Development Workflow:**
```bash
make test       # Clear cache and run all tests
make cover      # HTML coverage report
make build      # Build binary to execs/o6n
go vet ./...    # Static analysis
gofmt -w .      # Format code
```

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (Block Implementation):**
- Modal factory pattern — defines how all modal rendering is structured
- State transition contract — defines when and how navigation state is cleared
- Footer hint rendering model — defines how views declare their available actions

**Important Decisions (Shape Architecture):**
- Space action dialog as a factory modal type (`ModalActionMenu`)
- Credential and config file ownership boundaries

**Deferred Decisions (Post-MVP):**
- Mouse support (growth feature per PRD)
- Scriptable/pipe output mode (growth feature per PRD)
- Smart completions beyond user suggestions (growth feature per PRD)

---

### Data Architecture

**Data Source:** Operaton REST API only. No local database, no persistent cache.
All application data is fetched on demand via `tea.Cmd` and discarded on navigation
transitions. The only local persistence is `o6n-stat.yaml` (runtime state: active env,
skin, last nav position) and `o6n-env.yaml` (credentials).

**In-Memory Cache:** `internal/contentassist` — thread-safe (`sync.RWMutex`) suggestion
cache for user name completions. Populated from API responses. No TTL — refreshed on data
loads. This pattern is reusable for future completion types (variable names, process keys).

**Config as Schema:** `o6n-cfg.yaml` is the authoritative schema for resource types,
columns, actions, and drilldown relationships. It is read-only at runtime — never written
by the application. Architecture decisions that would require runtime writes to this file
are rejected.

---

### Authentication & Security

**Authentication:** HTTP Basic Auth only — `operaton.ContextBasicAuth` context injected
per-request from the active environment's credentials. No token caching, no session
management, no OAuth.

**Credential Isolation:** Credentials live exclusively in `o6n-env.yaml` (git-ignored,
`chmod 600`). They must never appear in: log files, debug output (`./debug/`), clipboard
operations (`J`/`Ctrl+J` copies row data as JSON, not credentials), or any serialized state.

**API Security:** All API calls go through `internal/client/` wrapping the generated
`internal/operaton/` client. Direct use of `internal/operaton/` from `internal/app/` is
an architecture violation — client wrapping is the single entry point.

---

### API & Communication Patterns

**Client Generation:** `internal/operaton/` is auto-generated from
`resources/operaton-rest-api.json`. It must never be manually edited. When the Operaton
API spec evolves, regenerate via `.devenv/scripts/generate-api-client.sh`.

**Async Contract:** Every API call is a `tea.Cmd` returning a typed message
(`dataLoadedMsg`, `errMsg`, etc.). No goroutines may write to the model directly. The
Bubble Tea event loop is the sole writer to model state.

**Error Propagation:** All API errors produce an `errMsg`. `Update()` stores the error
text in `m.footerError` and schedules a `clearErrorMsg` after 5 seconds via `tea.Tick`.
There are no silent failures and no error dialogs — the footer is the sole error surface.

**Nullable Types:** Use `NullableString`, `NullableInt32`, `NullableBool` helpers from the
generated client for all optional API response fields. Never nil-dereference raw pointers
from API responses.

---

### TUI Architecture

#### Decision 1: Modal Factory Pattern

**Decision:** The modal system is config-driven via a `ModalConfig` struct. The view path
contains no `switch modalType` logic for rendering individual modal bodies. Instead, each
modal type registers a `ModalConfig` at initialization that specifies:
- Size hint: `OverlayCenter` (compact dialogs, ~50%×auto), `OverlayLarge` (~80%×80%, rich
  content), or `FullScreen` (immersive flows)
- Title string
- Body renderer: `func(m Model) string` — produces the modal's inner content
- Button layout: derived from modal config (confirm/cancel labels, positions)
- Border and padding: uniform across all types (enforced by factory, not per-type code)
- Hint line: `[]Hint` — required for all `OverlayLarge` and `FullScreen` modals; rendered
  at the modal bottom using the same `Hint{Key, Label, MinWidth, Priority}` system as the
  main footer. Must include at minimum `Esc Close` and the modal's primary action.

**Size class usage:**
- `OverlayCenter` — Operational modals: Edit, Sort, ConfirmDelete, ConfirmQuit, ModalActionMenu, FirstRunModal
- `OverlayLarge` — Contextual modals: ModalHelp, ModalDetailView, ModalJSONView. Preserves
  background context; operator retains spatial orientation while viewing rich content.
  Approximate dimensions: ~80% termWidth × ~80% termHeight (at 120×20 minimum: ~96×16).
- `FullScreen` — Immersive flows: TaskComplete dialog.

The factory function signature:
```go
func renderModal(m Model, cfg ModalConfig) string
```

New modal types are added by defining a `ModalConfig` and registering it — no changes to
the factory function or view render path required.

**Rationale:** Eliminates per-type hardcoded layout code. Ensures FR11 (identical border,
padding, button placement) is structurally enforced. Satisfies NFR15 (config-driven modal
system). Three size classes allow operators to retain background context for reference
modals (OverlayLarge) while keeping destructive confirmations compact (OverlayCenter).
Reduces cognitive load for contributors adding new modal types.

#### Decision 2: State Transition Contract

**Decision:** `prepareStateTransition` is the single mandatory gate for all navigation
changes. It accepts a transition type parameter to control the scope of resets:

```go
type TransitionType int
const (
    TransitionFull      TransitionType = iota // context/environment switch: full reset
    TransitionDrillDown                        // drill-down push: partial reset
    TransitionPop                              // Esc/breadcrumb pop: restore from stack
)
```

`TransitionFull` clears: `activeModal`, `footerError`, `searchQuery`, `searchActive`,
`sortColumn`, `sortDirection`, `tableCursor`, `navigationStack` (full reset).

`TransitionDrillDown` pushes current `viewState` onto `navigationStack`, then clears:
`activeModal`, `footerError`, `searchQuery`, `searchActive`, `sortColumn`, `sortDirection`,
`tableCursor` (cursor resets for child view).

`TransitionPop` pops `viewState` from `navigationStack` and restores all captured state
(rows, cursor, columns, filters, breadcrumb) — no clearing.

**Any navigation code that does not call `prepareStateTransition` is a bug.**

**Rationale:** State leakage between navigation transitions is the primary defect category.
A typed, mandatory transition function makes violations detectable in code review and
testable in isolation. Satisfies FR20 (all transitions clear prior state completely).

#### Decision 3: Footer Hint Rendering — Push Model

**Decision:** Each view handler returns a `[]Hint` slice declaring the hints it contributes.
The footer renderer is stateless — it receives the hint slice and terminal width, filters
by each hint's `minWidth` threshold, and renders the result.

```go
type Hint struct {
    Key      string
    Label    string
    MinWidth int  // terminal columns required; 0 = always show
    Priority int  // 1 = highest (shown first when space is tight)
}
```

View handlers produce hints at render time — hints are not stored in model state. The
footer renderer calls `currentViewHints(m)` which dispatches to the active view's hint
function. This makes hints independently testable: `viewHints(m)` → `[]Hint`, no
rendering required.

**Rationale:** Separates hint declaration (per-view concern) from hint display (footer
concern). Satisfies FR15 (primary actions visible in footer). Enables clean testing of
hint visibility logic without full render cycle.

#### Decision 4: Space Action Dialog — ModalActionMenu

**Decision:** The `Ctrl+Space` action dialog is `ModalActionMenu`, a factory-registered
modal type. It reads the current table's `actions` slice from the resolved config at render
time. Its body renderer produces a list of action entries:
- Mutation actions listed first (HTTP verbs)
- Visual separator before the first `type: navigate` action
- Navigate actions shown with `→` suffix
- `[J] View as JSON` / `[Ctrl+J] Copy JSON` appended as the last item always

Keyboard handling for `ModalActionMenu`: single-character shortcuts dispatch directly to
the corresponding action. `Esc` closes the menu. The modal uses `OverlayCenter` size hint.

`Space` (without Ctrl) is reserved for future row selection — it must not trigger the
action menu.

**Rationale:** Gives the action dialog full visual consistency with other modals (border,
Esc/Enter). Keeps action definitions in `o6n-cfg.yaml` — the factory reads config, it does
not duplicate it. Satisfies FR16 (context-sensitive action menu via `Ctrl+Space`).

---

### Infrastructure & Deployment

**Build:** `make build` produces a static binary at `execs/o6n` with `CGO_ENABLED=0` and
`GO_TAGS=netgo`. No Docker, no containers, no cloud deployment — o6n is a local developer
tool distributed as a binary.

**Distribution:** Direct binary download or community package manager (e.g., Homebrew).
No server-side infrastructure required or planned for MVP.

**Debug Mode:** `--debug` flag creates `./debug/` with `o6n.log` (errors/debug messages),
`last-screen.txt` (last rendered frame), and `screen-{timestamp}.txt` (panic dumps). Debug
output must never contain credentials.

---

### Decision Impact Analysis

**Implementation Sequence (sprint order):**
1. Modal factory extraction — highest risk, foundational for all modal work
2. State transition contract enforcement — audit all nav code paths
3. Footer hint push model — wire per-view hint functions to footer renderer
4. `ModalActionMenu` registration — builds on factory once factory is stable
5. Rendering validation at 120×20 — validates the above changes in target terminals

**Cross-Component Dependencies:**
- Modal factory (TUI arch) → requires `internal/app/` refactor; no config or client changes
- State transition contract → affects all navigation paths in `internal/app/update.go` and nav files
- Footer hints → affects `internal/app/view.go` and all per-view render functions
- `ModalActionMenu` → depends on modal factory; reads `o6n-cfg.yaml` action config at render time

## Implementation Patterns & Consistency Rules

### Naming Patterns

**Go Identifiers:**
- Exported types, functions, constants: `PascalCase` — e.g., `ModalConfig`, `ViewState`
- Unexported identifiers: `camelCase` — e.g., `footerError`, `navigationStack`
- Acronyms follow Go convention: `URL`, `ID`, `API` not `Url`, `Id`, `Api`
- Test files: `main_<feature>_test.go` co-located with tested code — not in a `tests/` dir

**Message Types:** Suffix with `Msg`, past-tense noun — describes what happened:
- ✅ `dataLoadedMsg`, `errMsg`, `splashDoneMsg`, `terminatedMsg`
- ❌ `loadDataMsg`, `errorMsg`, `splashMsg`

**Command Functions:** Suffix with `Cmd`, imperative verb — describes what to do:
- ✅ `fetchDataCmd`, `terminateInstanceCmd`, `setVariableCmd`
- ❌ `dataFetchCmd`, `instanceTerminator`

**Modal Types:** Prefix with `Modal`, PascalCase noun:
- ✅ `ModalNone`, `ModalConfirmDelete`, `ModalActionMenu`
- ❌ `DeleteModal`, `ConfirmModal`

**Config Keys (YAML):** `snake_case` for all YAML keys in `o6n-cfg.yaml`,
`o6n-env.yaml`, and skin files:
- ✅ `api_path`, `hide_order`, `id_column`
- ❌ `apiPath`, `hideOrder`

---

### Structure Patterns

**Package Organization:** New feature code belongs in `internal/app/`. Use descriptive
file names that reflect the concern, not the layer:
- ✅ `internal/app/modal.go`, `internal/app/hints.go`, `internal/app/nav.go`
- ❌ `internal/app/helpers.go`, `internal/app/utils.go` (vague, accumulate debt)

**Off-Limits Paths:** Never add files to or modify `internal/operaton/` — regenerate
only via `.devenv/scripts/generate-api-client.sh`. No exceptions.

**New Top-Level Packages:** Do not create new packages under `internal/` without
architectural justification. Features that operate only on model state belong in
`internal/app/`. Only genuinely reusable, domain-independent logic warrants a new package.

**Test Co-Location:** Test files live next to the code they test:
- ✅ `internal/app/main_modal_test.go` alongside `internal/app/modal.go`
- ❌ `tests/modal_test.go`

---

### Format Patterns

**Bubble Tea Model Mutations:** All model mutations happen exclusively in `Update()`.
`View()` is a pure function — it reads model state and returns a string. No side effects
in `View()`. No model mutations in command functions.

```go
// ✅ Correct
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case dataLoadedMsg:
        m.rows = msg.rows
        return m, nil
    }
}

// ❌ Wrong — mutation in View
func (m Model) View() string {
    m.lastRendered = time.Now() // NEVER
    ...
}
```

**API Response Handling:** Always use `Nullable*` helpers for optional API fields.
Never nil-dereference raw pointers from generated client types:
```go
// ✅ Correct
name := NullableStringValue(item.GetName())

// ❌ Wrong
name := *item.Name
```

**YAML Config Keys:** All action path templates use `{id}`, `{name}`, `{value}`,
`{type}`, `{parentId}` as placeholders. Never invent new placeholder names without
updating the path template resolver.

---

### Communication Patterns (Bubble Tea Events)

**Command Pattern:** Every async operation follows this shape:
```go
func fetchSomethingCmd(client *client.Client, param string) tea.Cmd {
    return func() tea.Msg {
        result, err := client.GetSomething(param)
        if err != nil {
            return errMsg{err: err}
        }
        return somethingLoadedMsg{data: result}
    }
}
```
No goroutines. No channels. No direct model writes. Always returns `tea.Msg`.

**Error Events:** All errors produce `errMsg{err: error}`. The `Update()` handler for
`errMsg` stores to `m.footerError` and schedules auto-clear. Never use `log.Fatal`,
`os.Exit`, or bare `panic` in business logic.

**State Transitions:** All navigation changes call `prepareStateTransition(transitionType)`
before modifying view state. No exceptions. An agent that navigates without calling this
function is introducing a state leakage bug.

---

### Process Patterns

**Error Handling:**
- Business logic errors → `errMsg` → footer display
- Render panics → `defer/recover` in view functions → debug log + graceful fallback
- Config load errors → logged + user-visible fatal message at startup (before TUI starts)
- Never display errors in modal dialogs — footer is the sole error surface

**Loading States:**
- Footer center column: spinner text during active requests (e.g., `Loading…`)
- Footer right column: `⚡` flash (200ms) on each API call
- No loading indicators in the content box or table
- No skeleton screens or placeholder rows

**Modal Registration:**
- New modal type = new `ModalConfig` struct + registration
- Do NOT add `case ModalXxx:` to any `switch` statement in view render path
- Modal body renderer is a `func(m Model) string` — it may read model state but
  must not modify it

**Hint Registration:**
- New key hint = add `Hint{Key, Label, MinWidth, Priority}` to the relevant
  view's hint function return value
- Do NOT append hints inside the footer renderer function itself

**Skin/Color Usage:**
- Always use semantic color roles via the active skin's `lipgloss.Style`
- Never hardcode hex values or ANSI codes in Go source
- `env_name` semantic role: applied to the environment label in the fixed top-right
  header position — this is the primary environment identity signal. Set per-skin;
  configured in `o6n-env.yaml` per environment to distinguish prod/staging/dev.
- Environment `ui_color` is applied only to border accent — never to text content colors.
  It is a secondary accent, not the primary environment signal.

---

### Enforcement Guidelines

**All AI Agents MUST:**
- Call `prepareStateTransition` before any navigation state change
- Return `errMsg` for all async operation failures — no silent drops
- Keep `View()` free of side effects and model mutations
- Use the modal factory for all new modal types — no per-type view switch additions
- Register new hints in the view's hint function — not in the footer renderer
- Never edit `internal/operaton/` — regenerate via script only
- Use semantic color roles — never hardcode colors
- Co-locate test files with the code they test

**Anti-Patterns (Reject in Code Review):**
- Goroutines writing directly to model fields
- `case ModalXxx:` added to view render switch for a new modal type
- `m.footerError = "..."` set directly instead of via `errMsg`
- Hardcoded `#RRGGBB` or ANSI escape codes for colors
- Navigation code that bypasses `prepareStateTransition`
- Any manual edit to files under `internal/operaton/`
- Test files in a `tests/` directory rather than co-located

## Project Structure & Boundaries

### Complete Project Directory Structure

```
o6n/
├── main.go                          # Entry point — calls internal/app.Run()
├── go.mod                           # Module: Go 1.24+
├── go.sum
├── Makefile                         # build, test, cover targets
├── CLAUDE.md                        # AI agent context and conventions
├── README.md                        # User-facing documentation
├── specification.md                 # Authoritative technical spec (keep updated)
│
├── o6n-env.yaml                     # [git-ignored] Credentials, URLs, ui_color
├── o6n-env.yaml.example             # Template for first-time setup
├── o6n-cfg.yaml                     # [version-controlled] Resource definitions
├── o6n-stat.yml                     # [git-ignored, auto-generated] Runtime state
│
├── internal/
│   ├── app/                         # All TUI logic — primary sprint work area
│   │   ├── model.go                 # Model struct, ModalType consts, ViewState
│   │   ├── update.go                # Update() — all message handlers
│   │   ├── view.go                  # View() — pure render; dispatches to per-view renderers
│   │   ├── nav.go                   # Navigation: drill-down, Esc, breadcrumb, context switch
│   │   ├── commands.go              # tea.Cmd factories for all API operations
│   │   ├── table.go                 # Table rendering, column sizing, responsive layout
│   │   ├── edit.go                  # Edit modal logic: field cycling, validation, save
│   │   ├── skin.go                  # Skin loading, lipgloss style construction
│   │   ├── styles.go                # Semantic style accessors (border, fg, accent, etc.)
│   │   ├── util.go                  # Shared helpers (overlayCenter, truncate, etc.)
│   │   │
│   │   ├── modal.go                 # ★ SPRINT: ModalConfig struct + factory renderModal()
│   │   ├── hints.go                 # ★ SPRINT: Hint struct + per-view hint functions
│   │   │
│   │   ├── main_test.go             # Core model/update test helpers
│   │   ├── main_modal_test.go       # ★ SPRINT: Modal factory tests
│   │   ├── main_hints_test.go       # ★ SPRINT: Hint push model tests
│   │   └── main_<feature>_test.go  # Pattern for all feature tests
│   │
│   ├── client/                      # HTTP client wrapper (do not bypass — use this)
│   │   ├── client.go                # Client struct, auth context, request execution
│   │   └── client_test.go
│   │
│   ├── config/                      # Config structs and loaders
│   │   ├── config.go                # Config, EnvConfig, AppConfig, StatConfig structs
│   │   ├── loader.go                # LoadSplitConfig, LoadEnvConfig, LoadAppConfig, Save
│   │   └── config_test.go
│   │
│   ├── dao/                         # Data access interfaces
│   │   └── dao.go                   # DAO, HierarchicalDAO, ReadOnlyDAO interfaces
│   │
│   ├── validation/                  # Input validation for edit dialogs
│   │   ├── validation.go            # ValidateBool, ValidateInt, ValidateFloat, ValidateJSON
│   │   └── validation_test.go
│   │
│   ├── contentassist/               # Thread-safe suggestion cache
│   │   ├── contentassist.go         # SetUserCache, SuggestUsers (sync.RWMutex)
│   │   └── contentassist_test.go
│   │
│   └── operaton/                    # ★ AUTO-GENERATED — never edit manually
│       └── ...                      # Regenerate: .devenv/scripts/generate-api-client.sh
│
├── skins/                           # 35 YAML color theme files
│   ├── dracula.yaml
│   ├── nord.yaml
│   └── ...
│
├── resources/
│   └── operaton-rest-api.json       # OpenAPI spec — source for client generation
│
├── execs/
│   └── o6n                          # Build output (git-ignored)
│
├── debug/                           # [git-ignored] Created by --debug flag
│   ├── o6n.log
│   ├── last-screen.txt
│   └── screen-{timestamp}.txt
│
└── .devenv/
    └── scripts/
        └── generate-api-client.sh   # Regenerates internal/operaton/ via Docker
```

### Architectural Boundaries

**Package Ownership:**

| Package | Owns | May NOT |
|---|---|---|
| `internal/app/` | All TUI state, rendering, event handling, navigation | Call Operaton API directly — use `internal/client/` |
| `internal/client/` | HTTP execution, auth context injection | Hold TUI state or know about model fields |
| `internal/config/` | Config structs, file I/O, merging | Know about TUI or API response types |
| `internal/dao/` | Data access interface contracts | Implement API calls (that's `internal/client/`) |
| `internal/validation/` | Field-level input validation rules | Know about model, TUI, or API |
| `internal/contentassist/` | In-memory suggestion cache | Persist to disk or make API calls |
| `internal/operaton/` | Generated API types and HTTP wrappers | Be manually edited |

**Data Flow:**
```
User keypress
  → tea.KeyMsg
  → Update() dispatches to handler
  → handler calls xxxCmd() returning tea.Cmd
  → tea.Cmd executes async (via internal/client/)
  → returns typed Msg (dataLoadedMsg / errMsg)
  → Update() handles Msg, mutates model copy
  → View() called, reads model, renders to string
  → Terminal display updated
```

**Config Authority Boundary:**
- `o6n-cfg.yaml` is the sole authority for: resource type names, columns, actions,
  drilldown relationships, key bindings per resource
- Go code reads and interprets config — it never duplicates config knowledge
- If a behavior can be expressed in config, it must not be hardcoded in Go

### Requirements to Structure Mapping

**Sprint Work — Files Created or Significantly Modified:**

| File | Sprint Task | FRs Addressed |
|---|---|---|
| `internal/app/modal.go` | Create: `ModalConfig` struct + `renderModal()` factory | FR11, NFR15 |
| `internal/app/hints.go` | Create: `Hint` struct + per-view hint functions | FR15 |
| `internal/app/view.go` | Modify: wire hint push model to footer renderer | FR15 |
| `internal/app/nav.go` | Modify: enforce `prepareStateTransition` on all paths | FR20 |
| `internal/app/update.go` | Modify: register `ModalActionMenu` handler | FR16 |
| `internal/app/main_modal_test.go` | Create: modal factory tests | FR11 |
| `internal/app/main_hints_test.go` | Create: hint visibility tests | FR15 |

**Pre-existing Files (read-only during sprint unless fixing bugs):**
- `internal/client/`, `internal/config/`, `internal/validation/`, `internal/contentassist/`
- `o6n-cfg.yaml`, `o6n-env.yaml`, `skins/`

### Integration Points

**Internal Communication:**
- `internal/app/` → `internal/client/` via dependency injection (client passed to cmd factories)
- `internal/app/` → `internal/config/` at startup via `LoadSplitConfig()`
- `internal/app/` → `internal/validation/` in edit modal field save path
- `internal/app/` → `internal/contentassist/` for user input suggestions

**External Integration:**
- Operaton REST API via `internal/client/` wrapping `internal/operaton/`
- System clipboard via clipboard library (for `y` key copy-as-YAML)
- Terminal via Bubble Tea's `tea.Program` (ANSI/POSIX interface)

**Data Boundaries:**
- No inter-process communication — single binary, single process
- No database — all state in memory during session; `o6n-stat.yaml` on clean exit
- No network state persistence — each session starts fresh from the API

## Architecture Validation Results

### Coherence Validation ✅

**Decision Compatibility:** All technology choices are mutually compatible. Bubble Tea,
Bubbles, and Lipgloss are pure Go — no CGO requirements, compatible with static binary
build. Modal factory (`ModalConfig`) reads from `o6n-cfg.yaml` without duplicating its
authority. State transition contract operates as pure model mutations in `Update()` —
fully compatible with Bubble Tea's immutable model pattern.

**Pattern Consistency:** Naming conventions (`Msg`/`Cmd` suffixes, `Modal` prefix,
`snake_case` YAML) are consistent across all architecture sections. Push model for hints
aligns with `View()` purity requirement. `tea.Cmd` async pattern aligns with
`sync.RWMutex` usage in `contentassist` — no goroutine model writes in either.

**Structure Alignment:** Package layering (`internal/app/` → `internal/client/` →
`internal/operaton/`) has no circular dependencies. Sprint file targets (`modal.go`,
`hints.go`) are within `internal/app/` — correct package placement per boundaries.

### Requirements Coverage Validation ✅

**Functional Requirements:** All 39 FRs across 10 categories have architectural support.
Sprint-critical FRs (FR11-17, FR20) are mapped to specific new or modified files.
FR28 and NFR13/15 (config-driven extensibility) are enforced by the modal factory pattern
and `o6n-cfg.yaml` authority boundary.

**Non-Functional Requirements:** All 15 NFRs covered. Performance (NFR1-3): `tea.Cmd`
async contract and footer loading state. Reliability (NFR4-7): `defer/recover` in view
render + `errMsg` propagation + 5s auto-clear. Security (NFR8-9): credential isolation,
git-ignored env file, no credential logging. Terminal Compatibility (NFR10-12):
`hide_order` column visibility + priority-based hint system + resize event handling.
Maintainability (NFR13-15): config-driven resources, auto-generated client, modal factory.

### Implementation Readiness Validation ✅

**Decision Completeness:** All 4 sprint-critical architectural decisions documented with
rationale, Go type signatures, and behavioral contracts. Technology stack fully specified.

**Structure Completeness:** Complete project tree with file-level annotations. Sprint
files explicitly called out with ★ SPRINT markers. Package ownership table defines what
each package may and may not do.

**Pattern Completeness:** 7 naming pattern categories, 4 structure patterns, format
patterns with code examples (✅/❌), command/error/state-transition communication
patterns, 5 process patterns, enforcement guidelines with explicit anti-pattern list.

### Gap Analysis Results

**Important Gaps — Addressed Here:**

1. **`ModalActionMenu` cursor behaviour:** The action list within `ModalActionMenu`
   follows the existing sort modal cursor pattern — `Up`/`Down` moves the list cursor;
   single-character shortcuts (matching the action's `key` in config) dispatch the action
   directly without requiring cursor selection. `Enter` dispatches the currently
   highlighted action. `Esc` closes without action.

2. **`Hint.Priority` semantics:** Lower integer = higher priority = shown first when
   terminal width constrains the number of visible hints. `MinWidth: 0` means always
   visible regardless of terminal width. When two hints have the same priority, preserve
   declaration order. Footer renderer truncates from the end of the sorted list.

3. **`specification.md` update obligation:** Post-sprint, `specification.md` MUST be
   updated to document: (a) the modal factory and `ModalConfig` struct, (b) the hint
   push model and `Hint` struct, (c) the `TransitionType` enum and
   `prepareStateTransition` signature. This is a business success criterion from the PRD
   (documentation accuracy). The implementing agent must include spec updates as part of
   the definition of done for each sprint task.

**Nice-to-Have:** Target 80%+ line coverage on `modal.go` and `hints.go` (new sprint
files). Use `make cover` to verify after implementation.

### Architecture Completeness Checklist

**✅ Requirements Analysis**
- [x] Project context thoroughly analyzed (39 FRs, 15 NFRs, 10 FR categories)
- [x] Scale and complexity assessed (Medium / brownfield / Go TUI)
- [x] Technical constraints identified (Bubble Tea, generated client, config authority)
- [x] Cross-cutting concerns mapped (7 concerns documented)

**✅ Architectural Decisions**
- [x] Modal factory pattern — `ModalConfig` struct + `renderModal()` factory
- [x] State transition contract — `TransitionType` enum + `prepareStateTransition`
- [x] Footer hint push model — `Hint` struct + per-view hint functions
- [x] Space action dialog — `ModalActionMenu` factory modal type
- [x] Data flow, authentication, API client, error propagation all documented

**✅ Implementation Patterns**
- [x] Naming conventions — identifiers, messages, commands, modals, YAML keys
- [x] Structure patterns — package placement, test co-location, off-limits paths
- [x] Format patterns — model mutation rules, API response handling, code examples
- [x] Communication patterns — `tea.Cmd` shape, error events, state transitions
- [x] Process patterns — error handling, loading states, modal/hint/color registration
- [x] Enforcement guidelines — mandatory rules + anti-pattern list

**✅ Project Structure**
- [x] Complete directory structure with file-level annotations
- [x] Package ownership table (Owns / May NOT)
- [x] Data flow diagram (keypress → display)
- [x] Sprint file mapping (7 files: 2 new, 5 modified)
- [x] Integration points documented (internal + external)

### Architecture Readiness Assessment

**Overall Status: READY FOR IMPLEMENTATION**

**Confidence Level: High**

**Key Strengths:**
- Architecture is grounded in the existing brownfield codebase — no speculative patterns
- Sprint work is precisely scoped to 7 files with clear FR-to-file mapping
- All four sprint-critical decisions have concrete type signatures, not just concepts
- Anti-pattern list gives AI agents explicit rejection criteria for code review
- `specification.md` update obligation is part of the architecture, not an afterthought

**Areas for Future Enhancement (post-sprint):**
- Mouse support for table rows and modal buttons (Growth feature per PRD)
- Smart completions beyond user suggestions: variable names, process keys (Growth)
- CI/CD pipeline definition (not required for MVP single-developer sprint)
- Formal test coverage gates (currently advisory 80% guideline)

### Implementation Handoff

**AI Agent Guidelines:**
- Follow all architectural decisions exactly as documented — deviations require explicit
  architectural discussion, not unilateral implementation choices
- Use implementation patterns consistently — naming, structure, communication, process
- Respect package boundaries — `internal/app/` does not call `internal/operaton/` directly
- Refer to this document and `specification.md` for all architectural questions
- Update `specification.md` as part of done for each sprint task

**First Implementation Priority (sprint order):**
1. `internal/app/modal.go` — `ModalConfig` struct + `renderModal()` factory
2. `internal/app/main_modal_test.go` — factory tests (write tests first)
3. Enforce `prepareStateTransition` in `internal/app/nav.go`
4. `internal/app/hints.go` — `Hint` struct + per-view hint functions
5. `internal/app/main_hints_test.go` — hint visibility tests
6. Wire hints into `internal/app/view.go` footer renderer
7. Register `ModalActionMenu` in `internal/app/update.go`
8. Update `specification.md` to document new architectural elements
