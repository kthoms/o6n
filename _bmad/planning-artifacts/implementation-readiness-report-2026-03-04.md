---
stepsCompleted: ["step-01-document-discovery", "step-02-prd-analysis", "step-03-epic-coverage-validation", "step-04-ux-alignment", "step-05-epic-quality-review", "step-06-final-assessment"]
inputDocuments:
  - "_bmad/planning-artifacts/prd.md"
  - "_bmad/planning-artifacts/architecture.md"
  - "_bmad/planning-artifacts/ux-design-specification.md"
  - "_bmad/planning-artifacts/epics.md"
date: '2026-03-04'
project: 'o6n'
runNumber: 2
---

# Implementation Readiness Assessment Report

**Date:** 2026-03-04 (Run 2 — re-run after UX design completion)
**Project:** o6n

---

## Document Inventory

### PRD Documents
- `_bmad/planning-artifacts/prd.md` — whole document, status: complete (all 11 steps, validated 5/5 Excellent, 2026-03-03)

### Architecture Documents
- `_bmad/planning-artifacts/architecture.md` — whole document, status: complete (all 8 steps, 2026-03-03)

### UX Design Documents
- `_bmad/planning-artifacts/ux-design-specification.md` — whole document, status: complete (all 14 steps, 2026-03-04)

### Epics & Stories Documents
- `_bmad/planning-artifacts/epics.md` — whole document, status: complete (5 epics / 25 stories, 2026-03-03)

### Supporting Documents (not assessed)
- `_bmad/planning-artifacts/prd-validation-report.md` — PRD validation report (reference only)
- `_bmad/planning-artifacts/implementation-readiness-report-2026-03-03.md` — prior IR run (superseded)

---

## PRD Analysis

### Functional Requirements

FR1: Operator can navigate to any of the 35 configured resource types using the context switcher (`:` key)
FR2: Operator can browse a paginated table of resources in the current context
FR3: Operator can drill down from a parent resource to related child resources as configured in `o6n-cfg.yaml`
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
FR16: Operator can open a context-sensitive action menu **via the Space key** showing all available actions for the selected row ⚠️
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
FR27: Application reads resource types, columns, actions, and drilldown rules from `o6n-cfg.yaml` at startup
FR28: Contributor can add a new standard resource type by editing `o6n-cfg.yaml` without modifying Go source code
FR29: Operator can inspect process variables associated with a process instance
FR30: Operator can edit a process variable value inline with type validation
FR31: Operator can copy the selected resource row **as YAML** to the system clipboard ⚠️
FR32: Operator can filter the current resource table by entering a search term
FR33: Operator can clear the active search filter and return to the full table
FR34: Operator can toggle auto-refresh to continuously update the current table view
FR35: Application renders without overflow or truncation of critical information at 120×20 terminal size
FR36: Application adapts column visibility and hint display when the terminal is narrower than 120 columns
FR37: Application handles terminal resize events without corrupting the layout
FR38: Operator can switch between available color skins
FR39: Operator can toggle vim-style key bindings in-session

**Total FRs: 39**
⚠️ FR16: Still says "Space key" — must change to "Ctrl+Space" (B-1)
⚠️ FR31: Still says "copy as YAML" — must change to "copy as JSON (J / Ctrl+J)" (B-2)

### Non-Functional Requirements

NFR1–NFR15: All 15 NFRs present and unmodified. All passed previous validation (5/5). No changes detected.

**Total NFRs: 15**

### Additional Requirements

- Application Interface > Output Formats: "The `y` key copies the selected row as YAML to the clipboard" ⚠️ Also needs updating to `J` / `Ctrl+J` JSON (related to B-2)
- All other interface, config, and scripting requirements: present and correct.

### PRD Completeness Assessment

PRD is structurally complete — 39 FRs, 15 NFRs, 8 sections, validated at 5/5. Two targeted text corrections remain outstanding (B-1, B-2). No structural defects.

---

## Epic Coverage Validation

### FR Coverage Matrix

| FR | PRD Requirement (summary) | Epic Coverage | Status |
|---|---|---|---|
| FR1 | Navigate to any of 35 resource types via `:` | Epic 3 — Story 3.1 | ✅ Covered |
| FR2 | Browse paginated table | Epic 3 — Story 3.1 | ✅ Covered |
| FR3 | Drill down to child resources per config | Epic 3 — Story 3.2 | ✅ Covered |
| FR4 | Navigate back via Escape | Epic 3 — Story 3.2 | ✅ Covered |
| FR5 | Jump to breadcrumb level | Epic 3 — Story 3.2 | ✅ Covered |
| FR6 | Execute configured action on selected row | Epic 3 — Story 3.3 | ✅ Covered |
| FR7 | Confirm destructive actions | Epic 3 — Story 3.3 | ✅ Covered |
| FR8 | Success/error feedback in footer | Epic 3 — Story 3.3 | ✅ Covered |
| FR9 | Retry failed job (incident) | Epic 3 — Story 3.4 | ✅ Covered |
| FR10 | Set annotation on incident | Epic 3 — Story 3.4 | ✅ Covered |
| FR11 | All modals from shared factory | Epic 1 — Story 1.1 | ✅ Covered |
| FR12 | Esc dismisses any modal | Epic 1 — Story 1.4 | ✅ Covered |
| FR13 | Enter confirms any modal | Epic 1 — Story 1.4 | ✅ Covered |
| FR14 | Edit dialogs with type validation | Epic 1 — Story 1.4 | ✅ Covered |
| FR15 | Primary actions in footer | Epic 2 — Stories 2.1/2.2 | ✅ Covered |
| FR16 | Action menu trigger key | Epic 2 — Story 2.3 | ⚠️ Epics say Ctrl+Space ✓; PRD still says Space (B-1) |
| FR17 | Full key binding reference via `?` | Epic 2 — Story 2.4 | ✅ Covered |
| FR18–FR21 | Env/context switching, state restoration | Epic 1 — Stories 1.3/1.5 | ✅ Covered |
| FR22–FR25 | Task claim/unclaim/complete with typed vars | Epic 3 — Story 3.5 | ✅ Covered |
| FR26–FR28 | Env config, cfg.yaml-driven resources | Epic 5 — Stories 5.1/5.3 | ✅ Covered |
| FR29–FR31 | Variable inspection, editing, row copy | Epic 3 — Story 3.6 | ⚠️ FR31 YAML→JSON discrepancy (B-2/B-8) |
| FR32–FR34 | Search, filter, auto-refresh | Epic 3 — Story 3.7 | ✅ Covered |
| FR35–FR39 | Rendering, skins, vim mode | Epic 4 — Stories 4.1–4.5 | ✅ Covered |

### NFR Coverage

All 15 NFRs covered across Epics 1–5. No gaps.

### Coverage Statistics

- Total PRD FRs: 39
- FRs covered in epics: 39
- Coverage percentage: **100%**
- NFRs covered: 15/15 (100%)
- Key/format discrepancies requiring resolution: 2 (FR16 key binding, FR31 format)

---

## UX Alignment Assessment

### UX Document Status

Found and complete: `_bmad/planning-artifacts/ux-design-specification.md` — 14 steps, completed 2026-03-04.

### PRD ↔ UX Alignment

| Item | PRD | UX Spec | Status |
|---|---|---|---|
| FR16: Action menu trigger | `Space` key | `Ctrl+Space` (Space reserved for row selection) | 🔴 Conflict — PRD needs B-1 fix |
| FR31: Row export format | "copy as YAML" | `J` = JSON viewer, `Ctrl+J` = copy JSON; `y`/YAML removed | 🔴 Conflict — PRD needs B-2 fix |
| FR21: Startup state restoration | Restore last context/env | FirstRunModal on fresh start; `Ctrl+H` revisits | ✅ Extended (backward compatible) |
| FR34: Auto-refresh toggle | Toggle auto-refresh | `Ctrl+Shift+R` key specified; `⟳` indicator | ✅ Aligned |
| FR39: Vim mode toggle | Toggle in-session | `Ctrl+Shift+V` key specified; arrow keys remain primary | ✅ Aligned |
| NFR10: 120×20 target | VSCode + IntelliJ IDEA primary | 120×20 minimum; VSCode/IntelliJ = secondary | ⚠️ Priority framing differs (non-blocking) |

### Architecture ↔ UX Alignment

| UX Decision | Architecture Status | Gap |
|---|---|---|
| Three modal size classes: `OverlayCenter` / `OverlayLarge` (NEW) / `FullScreen` | Architecture Decision 1 documents `OverlayCenter` + `FullScreen` only | 🔴 **Blocking (B-3)** — `OverlayLarge` missing from `ModalConfig.sizeHint` enum |
| `env_name` semantic color role (fixed top-right header, primary environment signal) | Architecture references `ui_color` as border accent only; no `env_name` role | 🔴 **Blocking (B-4)** — `env_name` role must be added; `ui_color` demoted |
| `ModalActionMenu` trigger: `Ctrl+Space` | Architecture Decision 4 still states "via Space" in rationale | 🔴 **Blocking** — Architecture must align trigger key to `Ctrl+Space` |
| `ModalActionMenu` last item: `[J] View as JSON` / `[Ctrl+J] Copy JSON` | Architecture Decision 4 says "`[y] View as JSON`" — wrong key, partially correct format | 🔴 **Blocking** — Architecture must update to `J`/`Ctrl+J` pattern |
| JSONView modal (OverlayLarge) — `J` opens, `Ctrl+J` copies | Not in architecture decisions | 🟠 **High** — New modal type; compatible with factory but not documented |
| FirstRunModal (OverlayCenter) — `Ctrl+H` revisit | Story 1.5 covers startup restoration but not first-run context selection | 🟠 **High (B-6)** — Story 1.5 first-run AC insufficient |
| Mandatory hint lines for OverlayLarge + FullScreen modals | Architecture/Story 1.1 do not specify hint contract for modals | 🟠 **High (B-5)** — `ModalConfig` must include hint line contract |
| `Ctrl+Shift+V` / `Ctrl+Shift+R` intercept risk in VSCode/IntelliJ | Not addressed in architecture | 🟡 **Advisory** — Testing required; fallback bindings may be needed |

### Warnings

🔴 **PRD requires 2 text corrections (B-1, B-2):**
1. FR16: "Space key" → "`Ctrl+Space`"
2. FR31: "copy as YAML" → "copy as JSON (`J` / `Ctrl+J`)" + update Output Formats section

🔴 **Architecture requires 4 corrections (B-3, B-4 + 2 new):**
1. Decision 1: Add `OverlayLarge` as third `sizeHint` enum value
2. Decision 4: Update trigger from "Space" to "Ctrl+Space" in rationale + last item from `[y] View as JSON` to `[J] View as JSON` / `[Ctrl+J] Copy JSON`
3. Cross-cutting concerns: Add `env_name` semantic color role; demote `ui_color` to secondary accent
4. Process patterns: Add modal hint line contract (required for OverlayLarge + FullScreen)

🟠 **Epics require 5 story corrections (B-5 through B-9):**
All identified and detailed in Epic Quality Review below.

---

## Epic Quality Review

### Epic Structure Validation

#### User Value Check

| Epic | Title | User Outcome | Pass? |
|---|---|---|---|
| 1 | Consistent & Reliable Modal System | Operators trust every modal to behave predictably | ✅ Pass |
| 2 | Discoverable Actions | Operators find available actions without help screens | ✅ Pass |
| 3 | Core Operational Workflows | Operators complete full operational tasks in the TUI | ✅ Pass |
| 4 | Stable Rendering & Visual Polish | Operators use the tool in any terminal without corruption | ✅ Pass |
| 5 | Configuration, Security & Documentation | Operators configure environments; contributors extend safely | ✅ Pass |

No technical milestone epics. All 5 epics deliver user/contributor outcomes.

#### Epic Independence Check

| Epic | Depends On | Forward Dependency? | Verdict |
|---|---|---|---|
| Epic 1 | None (brownfield baseline) | — | ✅ Standalone |
| Epic 2 | Epic 1 (modal factory for ModalActionMenu) | No — prior epic | ✅ Acceptable |
| Epic 3 | Epic 1 (state transition, modal factory) | No — prior epic | ✅ Acceptable |
| Epic 4 | None (rendering layer independent) | — | ✅ Standalone |
| Epic 5 | None (config/security/docs layer) | — | ✅ Standalone |

No circular or forward dependencies.

### Story Quality Assessment

All 25 stories use Given/When/Then BDD format. Stories are appropriately sized for brownfield sprint. No trivially-small or oversized stories.

#### 🔴 Critical Violations

None — no technical epics, no forward-breaking dependencies.

#### 🟠 Major Issues

**Issue QR-1 (B-7): Story 2.3 — Wrong action menu last item**

Story 2.3 AC states: `"[y] View as YAML" is always the last item` in `ModalActionMenu`.

UX spec Decision (point 12/13): `ModalJSONView` replaces `y`/YAML. The last item must be `[J] View as JSON` with `Ctrl+J` copying directly.

**Required fix:** Replace `[y] View as YAML` → `[J] View as JSON` / `[Ctrl+J] Copy JSON` in Story 2.3.

---

**Issue QR-2 (B-8): Story 3.6 — Wrong copy key and format**

Story 3.6 AC states: `"Given the operator presses 'y' on any table row / When the copy action executes / Then the row is copied as YAML to the system clipboard"`

UX spec: `y` key removed; `J` opens JSON viewer; `Ctrl+J` copies JSON.

**Required fix:** Replace YAML clipboard AC with JSON viewer/copy ACs in Story 3.6.

---

**Issue QR-3 (B-9): Story 5.2 — Wrong copy key and format**

Story 5.2 AC states: `"Given the operator presses 'y' to copy a row as YAML to the clipboard / When the clipboard content is inspected / Then no credential fields... are present"`

UX spec: `y` removed; `J`/`Ctrl+J` are the copy keys.

**Required fix:** Update Story 5.2 to reference `J`/`Ctrl+J` JSON copy (security test: verify no credentials in JSON clipboard content).

---

**Issue QR-4 (B-5): Story 1.1 — `ModalConfig.sizeHint` enum incomplete**

Story 1.1 AC defines: `sizeHint (OverlayCenter / FullScreen)` only.

UX spec introduces `OverlayLarge` (~80%×80%) as a mandatory third size class used by: ModalHelp (reclassified from FullScreen), ModalDetailView (reclassified from FullScreen), ModalJSONView (new). Additionally, the UX spec mandates that all `OverlayLarge` and `FullScreen` modals render a `HintLine` at the bottom.

**Required fix:** Extend Story 1.1 AC to:
1. Add `OverlayLarge` to `sizeHint` enum: `sizeHint (OverlayCenter / OverlayLarge / FullScreen)`
2. Add mandatory hint line contract: `ModalConfig` must include a `HintLine []Hint` field; all `OverlayLarge` and `FullScreen` instances must populate it

---

**Issue QR-5 (B-6): Story 1.5 — First-run AC insufficient**

Story 1.5 AC states: `"Given no previous state file exists (first run) / When the application starts / Then it opens with a sensible default context without error"`

UX spec requires a `FirstRunModal` that prompts the operator to select their home context. No default context is meaningful without operator input. `Ctrl+H` opens this modal in subsequent sessions to revisit the home context choice.

**Required fix:** Replace the "sensible default" AC with:
- First-run: `FirstRunModal` (OverlayCenter) opens, prompting operator to select home context from the configured resource types
- On selection: home context is persisted to `o6n-stat.yaml`; app navigates to it
- `Ctrl+H`: re-opens the `FirstRunModal` to allow home context change at any time

#### 🟡 Minor Concerns

**QR-6: Story 4.4 skin count discrepancy** — Story says 35 skins; architecture says 35; CLAUDE.md says 36. Verify and correct to match actual `skins/` directory count.

**QR-7: Story 4.3 cross-story reference** — AC contains `"degrades gracefully per Story 4.2"`. Non-breaking (4.2 precedes 4.3 in sequence), but prefer self-contained AC language.

**QR-8: Story 2.3 Space note** — Note clarifying Space is reserved for "future" row selection should explicitly mark this as post-sprint intent, not current-sprint scope.

### Best Practices Compliance Summary

| Check | Status |
|---|---|
| All epics deliver user value | ✅ Pass |
| All epics function independently | ✅ Pass |
| No technical-milestone epics | ✅ Pass |
| No forward epic dependencies | ✅ Pass |
| Brownfield project correctly identified | ✅ Pass |
| Stories appropriately sized | ✅ Pass |
| Given/When/Then AC format | ✅ Pass |
| FR traceability maintained | ✅ Pass |
| JSON/key consistency (B-7, B-8, B-9) | ❌ 3 stories still use `y`/YAML |
| ModalConfig sizeHint complete (B-5) | ❌ OverlayLarge missing + no hint line contract |
| First-run AC matches UX spec (B-6) | ❌ "Sensible default" insufficient |

---

## Summary and Recommendations

### Overall Readiness Status

**⚠️ NEEDS WORK — 9 blocking items remain unresolved. Not ready for Sprint Planning.**

This is an identical finding to the previous IR run (2026-03-03). No blocking items have been applied to any document since that run. The core planning structure is sound: 100% FR/NFR coverage, 5 well-formed user-value epics, 25 properly sized BDD stories, no forward dependencies. All 9 blockers are targeted, surgical text corrections.

### Critical Issues Requiring Immediate Action (Blocking)

| # | Document | What's Wrong | Fix Required |
|---|---|---|---|
| B-1 | PRD FR16 | "via the Space key" | Change to "via `Ctrl+Space`" |
| B-2 | PRD FR31 + Output Formats | "copy as YAML" / "y key" | Change to "copy as JSON (`J` / `Ctrl+J`)" in FR31 and Output Formats section |
| B-3 | Architecture Decision 1 | `sizeHint` only has OverlayCenter / FullScreen | Add `OverlayLarge` (~80%×80%) as third size class |
| B-4 | Architecture (Theming + Cross-cutting) | `ui_color` is primary env signal; no `env_name` role | Add `env_name` semantic color role; demote `ui_color` to secondary accent |
| B-3b | Architecture Decision 4 | Trigger still says "via Space"; last item is `[y] View as JSON` (wrong key) | Change trigger to "Ctrl+Space"; change last item to `[J] View as JSON` / `[Ctrl+J] Copy JSON` |
| B-5 | Story 1.1 | `ModalConfig.sizeHint` missing `OverlayLarge`; no hint line contract | Add `OverlayLarge` to enum; add `HintLine []Hint` contract for OverlayLarge + FullScreen |
| B-6 | Story 1.5 | First-run: "sensible default context" | Replace with FirstRunModal flow: prompt operator to select home context; persist; `Ctrl+H` revisit |
| B-7 | Story 2.3 | `[y] View as YAML` as last action menu item | Replace: `[J] View as JSON` + `[Ctrl+J] Copy JSON` |
| B-8 | Story 3.6 | `y` key YAML copy | Replace: `J` opens JSON viewer; `Ctrl+J` copies JSON |
| B-9 | Story 5.2 | `y` key YAML clipboard security test | Replace: credential test uses `J`/`Ctrl+J` JSON copy |

> **Note on B-3b:** Decision 4 in architecture.md has a partial inconsistency not captured in the original IR: it already says "View as JSON" (not YAML) but still uses `y` as the key and "Space" as the trigger. These must both be corrected as part of the B-3/B-4 architecture pass.

### Recommended Actions (High Priority, Non-Blocking)

| # | Document | Issue | Action |
|---|---|---|---|
| R-1 | Story 3.7 | Filter bar 5-state visual design not in ACs | Add ACs for the 5 filter bar states (idle, active-pending, active, active-no-results, clearing) per UX spec |
| R-2 | Architecture / Epics | History view (`H` key) in UX spec but no story or architecture entry | Add story or explicitly mark post-sprint with rationale |

### Advisory Items (Low Priority)

| # | Document | Issue |
|---|---|---|
| A-1 | Architecture | `Ctrl+Shift+V` / `Ctrl+Shift+R` intercept risk in VSCode/IntelliJ not addressed |
| A-2 | Story 4.4 | Skin count: story says 35 — verify against actual `skins/` directory |
| A-3 | Story 2.3 | Space reservation note needs to clarify future intent vs current sprint scope |

### Recommended Next Steps

1. **Fix B-1/B-2** — Edit `_bmad/planning-artifacts/prd.md` (FR16 key + FR31 format + Output Formats section)
2. **Fix B-3/B-4/B-3b** — Edit `_bmad/planning-artifacts/architecture.md` (Decision 1 sizeHint enum + Decision 4 trigger/key + env_name role + hint line contract)
3. **Fix B-5 through B-9** — Edit `_bmad/planning-artifacts/epics.md` (5 story AC corrections)
4. **Address R-1** — Add filter bar 5-state ACs to Story 3.7 in the same epics pass
5. **Run Sprint Planning** — `/bmad-bmm-sprint-planning` once all 9 blocking items are resolved

### Final Note

This assessment (Run 2) confirms that all 9 blocking items from the previous IR run remain open in the source documents. No document was updated between IR runs. The architectural structure and story quality are otherwise solid. Correction effort is low — all 9 items are targeted text changes, not structural redesigns.

**Assessor:** Lt. Commander Data (Architect Agent) — BMad IR Workflow — 2026-03-04
**Status at completion:** NEEDS WORK (9 blocking items — unchanged from Run 1)
