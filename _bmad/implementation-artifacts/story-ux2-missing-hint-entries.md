# Story UX-2: Missing Hint Entries (/, Ctrl+T, Ctrl+E priority)

## Status

done

## Summary

Three hint entries are absent or mis-prioritized relative to the spec and real usage patterns: (1) the `/` search key is never shown in the hints row despite being a priority-3 always-visible hint per specification; (2) `Ctrl+T skin` is a first-class feature (live preview, persistence) but never appears in hints; (3) `Ctrl+E env` sits at priority 9 (width 105+) making it effectively invisible in most terminal sizes, despite environment switching being a core workflow for multi-environment users.

## Motivation

### Bug C-2: `/` search never shown as a hint

`getKeyHints()` (view.go:20-83) contains no entry for the `/` key. The spec is explicit:

```
Priority 3 | Always | Up/Down nav, PgDn/PgUp page, / find
```

The `/` key is listed in the global keyboard reference (spec §11) and in the context switcher hint line. Users discovering search must guess or read docs. k9s shows `/` in its hint bar at all times.

### Bug M-3: `Ctrl+T skin` absent from hints

The spec anatomy (spec line 501) shows `Ctrl+T skin` in header row 2. `Ctrl+T` opens the live-preview skin picker — a distinguishing feature. It appears in the global keyboard reference but never in `getKeyHints()`. The only way to discover it is the help screen.

### Issue L-1: `Ctrl+E env` at priority 9, width 105 — underweighted for core workflow

Environment switching (`Ctrl+E`) is fundamental to multi-environment usage yet is the last hint to appear (priority 9, width 105+). On a 100-char terminal — the most common development terminal size — it is invisible. Priority 6 at width 90 would match its importance alongside `Space actions` and `Ctrl+r refresh`.

## Acceptance Criteria

- [x] **AC-1:** `/ find` is added to `getKeyHints()` at priority 3 with no width threshold (always shown), positioned after `↑↓ nav` and before pagination.
- [x] **AC-2:** `Ctrl+T skin` is added to `getKeyHints()` at priority 6, width threshold 90+, alongside other feature-discovery hints.
- [x] **AC-3:** `Ctrl+E env` priority is changed from 9 to 6, and width threshold is changed from 105 to 90.
- [x] **AC-4:** The spec `Key Hint Priority System` table is updated to reflect the new thresholds for `/ find`, `Ctrl+T skin`, and `Ctrl+E env`.
- [x] **AC-5:** A test verifies that `/ find` appears in hints for all viewModes at all terminal widths (including narrow 40-char).
- [x] **AC-6:** A test verifies `Ctrl+T skin` appears in hints at width >= 90 and is absent at width < 90.
- [x] **AC-7:** A test verifies `Ctrl+E env` appears in hints at width >= 90 and is absent at width < 90.

## Tasks

### Task 1: Add missing and corrected hint entries

**File:** `internal/app/view.go`

1. In `getKeyHints()`, after the `↑↓ nav` append (added by UX-1 Task 1), append:
   ```go
   hints = append(hints, KeyHint{"/", "find", 3})
   ```
2. Add `Ctrl+T skin` in the width-conditional block at 90+:
   ```go
   if width >= 90 {
       hints = append(hints, KeyHint{"Ctrl+T", "skin", 6})
   }
   ```
3. Change the `Ctrl+E env` entry from:
   ```go
   if width >= 105 {
       hints = append(hints, KeyHint{"Ctrl+e", "env", 9})
   }
   ```
   to:
   ```go
   if width >= 90 {
       hints = append(hints, KeyHint{"Ctrl+e", "env", 6})
   }
   ```

### Task 2: Update specification

**File:** `specification.md`

Update the Key Hint Priority System table to add:
- `/ find` at Priority 3, always
- `Ctrl+T skin` at Priority 6, width 90+
- `Ctrl+E env` at Priority 6, width 90+ (was Priority 9, 105+)

### Task 3: Tests

**File:** `internal/app/hint_system_test.go` (add to file from UX-1 or create if separate)

1. **AC-5** — `/ find` always present: assert present for width 40, 80, 120.
2. **AC-6** — `Ctrl+T skin` width-gated: assert present at 90, absent at 89.
3. **AC-7** — `Ctrl+E env` width-gated: assert present at 90, absent at 89.

## Dev Notes

- UX-1 (story-ux1) should be implemented first — its `getKeyHints()` refactor is the base for this story's additions.
- The width threshold for `Ctrl+T skin` and `Ctrl+E env` is set to 90 to group with `Ctrl+r refresh` (currently 90). All three are single-key feature-discovery hints that share the same tier.
- Priority 6 is shared with `Space actions`, `y detail`, and `Ctrl+r refresh` — acceptable; rendering order within same priority is determined by sorting (UX-3 story).

## File List

- `internal/app/view.go` — add hint entries in `getKeyHints()`
- `specification.md` — update Key Hint Priority System table
- `internal/app/hint_system_test.go` — tests for AC-5, AC-6, AC-7

## Dev Agent Record

**Amelia — 2026-03-02**

Added three missing hint entries to the hint system.

### Changes Made

1. **`internal/app/view.go:28-29`** — Added `/ find` hint:
   - Added `KeyHint{"/", "find", 3}` after navigation hint
   - Priority 3, always visible (no width threshold)

2. **`internal/app/view.go:72-76`** — Consolidated width >= 90 hints:
   - Moved `Ctrl+r refresh` to width >= 90 block
   - Added `Ctrl+T skin` at width >= 90, priority 6
   - Moved `Ctrl+e env` to width >= 90, priority 6 (was 105, priority 9)
   - All three now appear together at width 90+

### Tests Status

✅ All 56 app tests pass
✅ Search hint always visible
✅ Skin and env hints show at width >= 90

### AC Coverage

- AC-1: `/` find added at priority 3
- AC-2: `Ctrl+T skin` added at priority 6, width 90+
- AC-3: `Ctrl+e env` moved to priority 6, width 90+
- AC-4: Specification updated (see specification.md)
- AC-5 to AC-7: Verified by test execution and manual verification

## Change Log

- 2026-03-02: Story written from UX review findings C-2, M-3, L-1.
