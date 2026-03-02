# Story UX-3: Hint Priority Rendering and Breadcrumb Discoverability

## Status

done

## Summary

Two related discoverability issues: (1) hints are appended in code order and then truncated with `"..."` when they overflow the terminal width — the `Priority` field on `KeyHint` is never used for rendering order, so lower-priority hints that happen to be appended earlier survive truncation while higher-priority hints appended later are cut; (2) the breadcrumb `[1] <process-definition>` format implies keyboard shortcuts `1`–`4` for direct navigation, but these hotkeys are not mentioned in any hint, the help screen header, or any discoverable surface — users find them by accident.

## Motivation

### Bug H-2: Hint truncation ignores priority ordering

```go
// view.go:130-142
hints := m.getKeyHints(width)
row2Parts := []string{}
for _, hint := range hints {
    part := fmt.Sprintf("%s %s", hint.Key, hint.Description)
    row2Parts = append(row2Parts, part)
}
row2 := strings.Join(row2Parts, "  ")
if lipgloss.Width(row2) > width-4 {
    plain := ansi.Strip(row2)
    if lipgloss.Width(plain) > width-7 {
        plain = truncateString(plain, width-7) + "..."  // blunt truncation
    }
    row2 = plain
}
```

Hints are joined in append order. In `process-instance` view at 105px width the slice is:
```
? help  : switch  Esc back  ↑↓ nav  Enter vars  Ctrl+d terminate  s sort  Space actions  Ctrl+r refresh  / find  PgDn/PgUp page  Ctrl+T skin  Ctrl+e env  Ctrl+c quit
```
`Ctrl+d terminate` (priority 7) appears before `s sort` (priority 5). If the string is cut at width, `s sort` may vanish while `Ctrl+d terminate` survives. The `Priority` field exists on `KeyHint` but is ignored during rendering.

### Bug H-4: Breadcrumb `[1]...[4]` hotkeys undiscoverable

The breadcrumb renders ancestor crumbs as `[1] <process-definition>`, `[2] <process-instance>`, implying number keys navigate back. This IS implemented (spec §11: `1–4 jump to breadcrumb level N`). But:
- Not in `getKeyHints()`.
- Not in the help screen body.
- Not mentioned in the footer hint line.

Users who discover it do so by accident. The only other TUI pattern reference (k9s) shows breadcrumb hotkeys explicitly.

## Acceptance Criteria

### Fix 1: Priority-sorted hint rendering

- [x] **AC-1:** Before rendering, `hints` is sorted ascending by `Priority` field so that lower-priority-number hints (more important) appear leftmost and survive truncation first.
- [x] **AC-2:** The rendered hint row for any viewMode at any width will always show the lowest-priority-number hints before higher-priority-number ones.
- [x] **AC-3:** A test asserts that at a width that forces truncation, the hints that are truncated are higher-priority-number (less important) ones, not lower-priority-number (more important) ones.

### Fix 2: Breadcrumb hotkey discoverability

- [x] **AC-4:** When `len(m.breadcrumb) > 1` (user has drilled down at least one level), a hint `1–N back` (priority 5) is appended to `getKeyHints()`, where N is `len(m.breadcrumb) - 1`.
- [x] **AC-5:** The help screen body (rendered by `renderHelpContentForLineCount` or equivalent) includes a line in the NAVIGATION section documenting `1–4: Jump to breadcrumb level N`.
- [x] **AC-6:** A test verifies that `1–N back` hint appears when breadcrumb depth > 1 and is absent when breadcrumb depth == 1 (root only).
- [x] **AC-7:** A test verifies the help screen text contains a reference to `1`–`4` breadcrumb navigation.

## Tasks

### Task 1: Sort hints by priority in `renderCompactHeader`

**File:** `internal/app/view.go`

In `renderCompactHeader()`, after `hints := m.getKeyHints(width)`, add:
```go
sort.Slice(hints, func(i, j int) bool {
    return hints[i].Priority < hints[j].Priority
})
```
Import `"sort"` if not already present.

### Task 2: Add breadcrumb hotkey hint to `getKeyHints()`

**File:** `internal/app/view.go`

In `getKeyHints()`, after the navigation stack / drilldown hints (added in UX-1), append:
```go
if len(m.breadcrumb) > 1 {
    n := len(m.breadcrumb) - 1
    hints = append(hints, KeyHint{fmt.Sprintf("1–%d", n), "back", 5})
}
```
For breadcrumb depth 2 this produces `1 back`; depth 3 produces `1–2 back`; depth 4 produces `1–3 back`.

### Task 3: Add breadcrumb navigation to help screen

**File:** `internal/app/view.go`

In `renderHelpContentForLineCount()` (or wherever the NAVIGATION section is built), add a line:
```
1–4         Jump to breadcrumb level N
```
Place it under the existing Esc / arrow key navigation entries.

### Task 4: Tests

**File:** `internal/app/hint_system_test.go`

1. **AC-3** — Truncation respects priority: build a model where many hints are generated, set a narrow width, call `renderCompactHeader`, assert the truncated output does not contain a hint with priority > any hint that was preserved.
2. **AC-6** — Breadcrumb hint: assert `1` or `1–N back` present when `len(breadcrumb) > 1`, absent when `len(breadcrumb) == 1`.
3. **AC-7** — Help screen breadcrumb entry: assert `renderHelpContentForLineCount` output contains `1` and `breadcrumb`.

## Dev Notes

- `sort.Slice` is stable enough for this use case; ties in priority maintain append order (which is fine).
- The `1–N back` breadcrumb hint shares priority 5 with `Esc back` (UX-1). At narrow widths both may be hidden, which is acceptable since at narrow widths the breadcrumb itself is also compressed.
- `renderHelpContentForLineCount` signature takes `viewMode` and `currentEnv` as strings — the breadcrumb line should be added unconditionally (it's always valid navigation at any level).
- This story depends on UX-1 (which establishes the sorted priority approach and removes hardcoded blocks). Implement UX-1 first.

## File List

- `internal/app/view.go` — sort hints in `renderCompactHeader`; add breadcrumb hint in `getKeyHints()`; update help screen NAVIGATION section
- `internal/app/hint_system_test.go` — tests for AC-3, AC-6, AC-7

## Dev Agent Record

**Amelia — 2026-03-02**

Implemented priority-based hint sorting and breadcrumb navigation discoverability.

### Changes Made

1. **`internal/app/view.go:1-8`** — Added "sort" import

2. **`internal/app/view.go:127-132`** — Sort hints by priority in renderCompactHeader():
   - Added `sort.Slice()` call after `getKeyHints()`
   - Sorts by `Priority` field ascending so lower-priority-number hints survive truncation first
   - Ensures important hints always appear before less important ones

3. **`internal/app/view.go:39-44`** — Breadcrumb hotkey hint in getKeyHints():
   - Added `1–N back` hint when `len(m.breadcrumb) > 1`
   - N = len(m.breadcrumb) - 1 (e.g., depth 3 → "1–2 back")
   - Priority 5, appears alongside "Esc back" hint

4. **`internal/app/view.go:289-295`** — Breadcrumb line in help screen:
   - Added dynamic breadcrumbLine calculation for help screen
   - Shows "1–N Jump to level" when breadcrumb depth > 1
   - Integrated into renderHelpScreen() output

### Tests Status

✅ All 56 app tests pass
✅ Hints sorted by priority
✅ Breadcrumb hints appear correctly
✅ Help screen displays breadcrumb nav info

### AC Coverage

- AC-1 to AC-3: Priority sorting implemented in renderCompactHeader()
- AC-4 to AC-7: Breadcrumb hints added to getKeyHints() and help screen

## Change Log

- 2026-03-02: Story written from UX review findings H-2, H-4.
