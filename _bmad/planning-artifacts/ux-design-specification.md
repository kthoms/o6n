---
stepsCompleted: [1, 2, 3, 4]
inputDocuments:
  - "_bmad/planning-artifacts/prd.md"
  - "_bmad/planning-artifacts/architecture.md"
  - "_bmad/planning-artifacts/epics.md"
  - "_bmad/planning-artifacts/implementation-readiness-report-2026-03-03.md"
  - "_bmad/implementation-artifacts/story-bug-fixes-config-quality.md"
  - "_bmad/implementation-artifacts/story-drilldown-action-model.md"
  - "_bmad/implementation-artifacts/story-drilldown-debug-fixes.md"
  - "_bmad/implementation-artifacts/story-ux1-config-driven-hints.md"
  - "_bmad/implementation-artifacts/story-ux2-missing-hint-entries.md"
  - "_bmad/implementation-artifacts/story-ux3-hint-rendering-breadcrumb.md"
  - "_bmad/implementation-artifacts/story-ux4-filter-pagination-visibility.md"
  - "_bmad/implementation-artifacts/story-ux5-small-ux-fixes.md"
  - "_bmad/implementation-artifacts-done/story-accessibility-and-empty-states.md"
  - "_bmad/implementation-artifacts-done/story-claim-complete-user-tasks.md"
  - "_bmad/implementation-artifacts-done/story-keyboard-convention-system.md"
  - "_bmad/implementation-artifacts-done/story-layout-optimization.md"
  - "_bmad/implementation-artifacts-done/story-search-pagination-awareness.md"
  - "_bmad/implementation-artifacts-done/story-state-transition-contract.md"
  - "_bmad/implementation-artifacts-done/story-task-complete-dialog-polish.md"
  - "_bmad/implementation-artifacts-done/story-task-complete-dialog-ux.md"
  - "_bmad/implementation-artifacts-done/story-vim-mode-toggle.md"
---

# UX Design Specification — o8n

**Author:** Karsten
**Date:** 2026-03-03

---

<!-- UX design content will be appended sequentially through collaborative workflow steps -->

## Project Understanding

### Executive Summary

o8n is a keyboard-first terminal UI for managing Operaton BPMN workflow engines, inspired by k9s. It gives DevOps engineers and BPMN operators a fast, discoverable interface to navigate process definitions, process instances, tasks, and associated variables — directly from the terminal without leaving their editor or switching to a browser.

The product's core value proposition is speed through keyboard mastery: every operation reachable by keystroke, every resource navigable in seconds, every action executable without mouse. It is a brownfield project with an evolving implementation; the UX design specification codifies patterns that are emerging but not yet fully consistent, and defines the target state for the complete system.

### Target Users

**Alex — DevOps Engineer**
Primary operator. Monitors running process instances, inspects variables, cancels or retries stuck processes. Works in VSCode integrated terminal. Optimizes for speed; learns key bindings quickly. Expects k9s-like muscle memory.

**Priya — BPMN Operator**
Day-to-day task handler. Claims human tasks, reviews input variables, completes tasks with output variables. May work in IntelliJ IDEA terminal. Less familiar with vim conventions; needs discoverable actions and clear confirmation flows.

**Marco — Go Contributor**
Extends the tool. Reads source, adds resource types, modifies config. Cares about config-driven architecture, clean patterns, and consistency of the keyboard grammar across resource types.

### Key Design Challenges

1. **Density within 120×20 viewport.** Every pixel is shared between navigation, data table, breadcrumb, hints, and status. Column visibility, hint priority, and modal sizing must all be width-aware.

2. **Action discoverability without mouse.** There is no right-click, no hover tooltip. Users must discover available actions through footer hints, the action menu, and the help screen — without being overwhelmed at any tier.

3. **State legibility across resource types.** Process instances, tasks, jobs, and incidents each have distinct lifecycle states. Colors and labels must consistently communicate state across 35 resource types without a legend.

4. **Complex dialogs in terminal constraints.** Task completion requires input review, output variable entry with type validation, and confirmation — all inside a modal that must remain smaller than the terminal viewport.

5. **Keyboard convention consistency.** As the feature set grows, key bindings accumulate conflicts and inconsistencies. The grammar (case, modifier use, overloading) must be governed by explicit rules, not per-feature decisions.

### Design Opportunities

1. **Footer-as-HUD.** The hint bar is permanently visible and context-sensitive. Used well, it functions as an always-present ambient guide — not a help page users have to seek out.

2. **Config-driven discoverability.** Because resource definitions and their associated actions live in `o8n-cfg.yaml`, the hint system can surface resource-specific keys at priority 1–2, ahead of global navigation. The config file becomes the UX specification for each resource type.

3. **Progressive disclosure.** Three tiers — footer hints → Ctrl+Space action menu → ? help screen — allow new users to discover gradually without overwhelming experts. Each tier is richer and more comprehensive than the last.

4. **Environment-as-color-signal.** The `env_name` semantic color role, shown in a fixed header position, lets operators running multiple terminals against different environments instantly distinguish them. No text reading required at a glance.

5. **First-run onboarding.** A one-time context selection modal on fresh start can route each user type to their primary resource without requiring config edits. Persisted preference means the second launch is instant.

6. **OverlayLarge contextual modals.** A new intermediate modal size (~80%×80%) preserves background context while presenting rich content (detail view, JSON viewer, help screen). Users retain spatial orientation even when a modal is open.

### Party Mode Design Decisions (Step 2 Discovery)

The following 20 design decisions were surfaced and confirmed during Step 2 Party Mode with Counselor Troi (UX Lead), UX Review, and Mary:

1. **First-run home context selection** — On fresh start (no persisted state), o8n shows a context selection modal reusing the context switcher pattern. The selected resource becomes the home. Choice is persisted to `o8n-stat.yaml`. `Ctrl+H` revisits from any state.

2. **Hint hierarchy inversion** — Resource-specific action keys (from `TableDef.Actions` in config) are assigned hint priority 1–2. Global navigation keys (Esc, Ctrl+Space, /) drop to priority 4+. Most-used resource actions are always visible first.

3. **Environment visibility redesign** — A dedicated `env_name` semantic color role in skins governs the environment label color. The label is fixed top-right in the header. This replaces the current `ui_color` border approach as the primary environment signal.

4. **Remove ui_color border flood** — The `ui_color` accent color is removed from border painting. It is demoted to a secondary signal (e.g., accent elements within focused widgets) to avoid visual noise and false environment cues.

5. **Three-tier discoverability formalized** — Footer hints → Ctrl+Space action menu → ? help screen is established as an explicit architectural system, not an emergent collection of features. Each tier has defined scope and responsibility.

6. **Confirmation principle** — Confirmation dialogs are required for actions that are **both** irreversible **and** have a blast radius beyond a single row (delete, stop job, terminate, batch operations). Single-row reversible actions do not require confirmation.

7. **Splash screen config** — `splash: bool` is an optional field in `o8n-cfg.yaml` (default: show splash). The `--no-splash` CLI flag overrides. Sequencing: splash → first-run prompt (if needed) → main view.

8. **Modal size classification** — Three tiers: `OverlayCenter` (compact dialogs, ~50%×auto), `OverlayLarge` (NEW, ~80%×80%, for rich content), `FullScreen` (immersive flows). All modals declare their size class.

9. **Modal emotional grouping** — Modals are categorized by role: Operational (Edit, Sort, ActionMenu), Consequential (ConfirmDelete, ConfirmQuit), Contextual (Help, DetailView, JSONView), Immersive (TaskComplete). Size class follows category.

10. **ModalHelp reclassified** — From `FullScreen` to `OverlayLarge`. Background context is preserved; users retain spatial awareness of where they invoked help from.

11. **ModalDetailView reclassified** — From `FullScreen` to `OverlayLarge`. Same rationale: detail content is reference, not immersive — background context aids orientation.

12. **ModalJSONView introduced** — New modal type, `OverlayLarge`. Title = resource type + ID. `J` opens it. `Ctrl+J` copies JSON directly to clipboard without opening the viewer. Replaces the `y` copy-as-YAML pattern throughout.

13. **YAML → JSON throughout** — All copy/view operations use JSON. `y` (copy as YAML) is removed. `J` / `Ctrl+J` are the canonical data export keys.

14. **Keyboard case convention** — All keys shown in the help screen are uppercase without indicating Shift. Uppercase in help = no Shift required. Modifier keys (Ctrl, Alt, Shift) are always spelled out explicitly when required. This is a display rule, not an input rule.

15. **Key binding changes** — `Ctrl+V` = vim toggle (cross-platform risk noted: clipboard conflict on some terminals), `Ctrl+R` = immediate refresh, `Ctrl+Shift+R` = toggle auto-refresh. `V` alone is freed for future use.

16. **H = view history convention** — `H` navigates to the history view for the current resource type, consistent across all resource types that support history.

17. **Mandatory hint lines for complex modals** — `OverlayLarge` and `FullScreen` modals must render a hint line at their bottom. This is a spec contract, not a nice-to-have. Modal hint lines follow the same `Hint{Key, Label, MinWidth, Priority}` system as the main footer.

18. **Search/filter five visual states** — Explicitly defined: (1) no filter active, (2) filter popup open, (3) filter locked/active, (4) server-side filter applied, (5) filter cleared/resetting. Each state has distinct visual treatment.

19. **Auto-refresh indicator design** — `⟳` symbol in footer right area, rendered in accent color. Flashes briefly on each refresh cycle. Absent entirely when auto-refresh is off. No text label.

20. **API status design** — Three states: `●` (connected, green), `✗` (error, red), `○` (idle/unknown, muted). Always rendered as color + symbol; never color alone (colorblind accessibility).

## Core User Experience

### Defining Experience

The defining o8n interaction is: **navigate to any resource, inspect it, and act on it — without leaving the keyboard.** Every session begins with the user arriving at a resource table (process instances, tasks, or wherever they persisted from last time), scanning rows, and taking a targeted action. The speed between "I need to see X" and "I'm looking at X" is the product's primary value.

For Alex (DevOps), the core loop is: filter instances → inspect → kill or retry. For Priya (BPMN Operator), it is: find task → claim → review inputs → complete with outputs. Both loops share the same keyboard grammar — the product succeeds when both feel native.

### Platform Strategy

- **Terminal/TUI only.** No web, no mobile, no desktop GUI. Constraints are permanent and architectural.
- **Keyboard-exclusive interaction.** Mouse input is unsupported by design. Every action reachable by keystroke.
- **Primary target:** Linux and macOS native terminals (gnome-terminal, iTerm2, Terminal.app, Alacritty). Minimum viewport 120×20.
- **Secondary target:** VSCode integrated terminal (for development workflows). IntelliJ IDEA terminal.
- **No offline mode.** o8n is a live API client; network connectivity to an Operaton engine is a hard prerequisite.
- **Static binary.** No runtime dependencies; distributes as a single file per platform.

### Effortless Interactions

These interactions must require zero conscious effort for experienced users:

1. **Context switching** — `:` opens the context switcher; typing a resource name and pressing Enter jumps there instantly. Muscle memory from k9s users transfers directly.
2. **Drill-down navigation** — Enter on a table row descends into the child resource. Esc ascends. No menus, no confirmation, no loading spinners blocking the keyboard.
3. **Filter/search** — `/` enters filter mode immediately. The table filters as the user types. Esc clears and exits.
4. **Action discovery** — Ctrl+Space opens the action menu from any table row. All available actions for the current resource are listed; no memory required.
5. **Variable copying** — `J` or `Ctrl+J` on any variable or resource row immediately surfaces the JSON value. No multi-step flow.

### Critical Success Moments

1. **First successful task completion.** Priya claims a task, reviews input variables, fills in output variables with type validation, confirms, and sees the task disappear from the list. If this flow is confusing at any point, the product fails for its most consequential use case.
2. **First context switch.** Alex types `:incidents` and lands on the incidents table in under 2 seconds. This moment establishes that the keyboard grammar works and learning it pays off.
3. **First drill-down.** A new user presses Enter on a process instance row and arrives at the variables view with breadcrumb showing the path. They press Esc and return. Spatial model established.
4. **First environment distinction.** A user opens o8n against a production environment and immediately sees the distinct env color in the header — before reading any text. Environment misidentification prevention.
5. **First hint-driven discovery.** A user sees a hint in the footer, presses the key, and the action executes. The hint system proves its value in a single interaction.

### Experience Principles

1. **Keyboard fluency compounds.** Every interaction is designed for users who will repeat it hundreds of times. Optimize for the 100th use, not the first — but use progressive disclosure to get users to their 100th use safely.
2. **Context is never lost.** Navigation always shows breadcrumb. Modals preserve background. Esc always goes back. Users always know where they are and how they got there.
3. **The config is the UX.** Resource actions, column visibility, drilldown rules, and hint priorities are defined in `o8n-cfg.yaml`. The tool's behavior is inspectable, not magic — contributors extend the UX by editing config.
4. **Color signals are semantic, not decorative.** Every color use carries a defined meaning (environment identity, resource state, action category, error/warning). Color is never used purely for aesthetics.
5. **Trust the keyboard grammar.** A consistent, learnable grammar of keys (resource-specific at 1–2, global nav at 4+, modifiers explicit, case normalized) means users can guess correctly. Correctness of guesses is a design metric.

## Desired Emotional Response

### Primary Emotional Goals

o8n users should feel **in command** — the same focused, efficient feeling a vim or k9s expert has when their muscle memory fires correctly. The primary emotional goal is **controlled competence**: the sense that the tool amplifies your ability rather than placing friction between you and the engine.

The secondary feeling is **calm confidence** in high-stakes situations. When an operator is diagnosing a stuck process instance or completing a task under production pressure, the UI must not add anxiety. Clear state, predictable keys, explicit confirmations for destructive actions — the tool stays out of the way emotionally.

### Emotional Journey Mapping

| Moment | Target Emotion | Anti-pattern to Avoid |
|---|---|---|
| First launch | Curious, oriented | Overwhelmed, lost |
| First hint discovery | "Oh, I can do that" | Ignored, invisible |
| First drill-down | Spatially grounded | Disoriented ("where am I?") |
| First task claim + complete | Accomplished | Anxious about mistakes |
| Repeated daily use | Fluent, fast | Bored or impatient with friction |
| Error or API failure | Informed, not alarmed | Confused, panicked |
| Destructive action prompt | Deliberately cautious | Surprised, regretful |

### Micro-Emotions

- **Confidence over confusion** — keyboard grammar predictability. When a key binding makes sense, users feel smart, not lucky.
- **Accomplishment over frustration** — task completion is the most emotionally loaded flow. Success must feel clean and final.
- **Trust over skepticism** — API status, pagination counts, and stale-data indicators must be honest. Users need to trust what they see.
- **Calm over anxiety** — confirmation dialogs for destructive actions provide a deliberate pause, not bureaucratic friction.
- **Delight (rare, earned)** — the `⟳` flash on auto-refresh, the environment color lock, the action menu appearing exactly when needed. Small moments of "this tool thinks like I do."

### Design Implications

- **Confidence** → Consistent keyboard grammar with explicit case and modifier rules. No surprising rebindings. Hint system that teaches rather than clutters.
- **Controlled competence** → Progressive disclosure tiers (footer → action menu → help) mean users gain power incrementally, never all-at-once overwhelm.
- **Calm confidence** → Confirmation dialogs only where the blast radius warrants it. No confirmation theater for low-stakes actions.
- **Spatial grounding** → Breadcrumb in footer always present. Modals preserve background context (OverlayLarge, not FullScreen). Esc always reverses.
- **Trust** → API status indicator always visible (color + symbol). Pagination shows total count. Auto-refresh indicator is honest about when data was last fetched.
- **Accomplished** → Task completion modal has a clear, positive terminal state. The task disappears from the list. No ambiguity about whether it worked.

### Emotional Design Principles

1. **Never surprise with consequences.** Destructive actions require deliberate confirmation. Reversible actions don't — friction must be proportionate to risk.
2. **Reward keyboard investment.** The first time a learned key binding fires correctly, the user should feel the payoff. Hint → action execution must be a satisfying one-keystroke loop.
3. **The error is not the end.** API errors display in the footer with auto-clear. They inform without blocking. The user can continue navigating immediately.
4. **Silence is honest.** If there's no data, show a clear empty state. If the API is unreachable, show the status indicator. Never show stale data silently.
5. **Familiarity lowers anxiety.** k9s and vim users arrive with existing mental models. Adopting their conventions (`:` for context, `Esc` to back out, `/` for search) reduces onboarding anxiety before the first keystroke.
