# Story UX-4: Filter and Pagination Visibility

## Status

done

## Summary

Two pieces of persistent UI state are easy to miss: (1) when a search filter is locked (Enter in the search popup), the only indicator is a change in the content box title — there is no badge or visual marker in the always-visible hints row, leaving users confused about why the table shows fewer rows than expected; (2) pagination position (`[2/5]`) appears in the far-right corner of the footer, sandwiched between the status text and the flash icon, which is the lowest-attention region of the screen. Both issues cause users to lose orientation.

## Motivation

### Bug M-1: No active filter indicator in the hints row

When search is locked the box title changes to:
```
process-definitions [/proc/ — 3 of 42]
```
This is inside the box border — low visibility. The hints row continues to show `? help  : switch  ↑↓ nav  / find ...` with no indication that a filter is active. Users who don't notice the title change have no signal that the table is showing a subset.

Comparison: htop shows an explicit `Filter:` banner; lazygit shows `filter:` inline in the list header; k9s dims the resource label when a filter is active.

Desired behaviour: when `m.searchTerm != ""` (filter locked), the hints row should include a persistent badge such as `[/proc/ — Esc:clear]` that both surfaces the active term and tells the user how to clear it.

### Bug M-2: Pagination `[2/5]` buried in far-right footer corner

```go
// view.go:758
paginationStr = fmt.Sprintf(" [%d/%d]", currentPage, totalPages)
// view.go:763
pageIndicator = m.styles.PageCounter.Render(paginationStr) + " "
// view.go:765
rightPart := pageIndicator + rpStyle.Render(remoteSymbol+latencyStr)
```

The page indicator is in the rightmost cluster of the footer. The user's eye tracks left-to-right; the far right is scanned last. k9s places resource count inline in the title bar. The o8n box title already shows total count (`— 42 items`) — adding the page position there gives it higher visibility.

Desired behaviour: when on page N of M, the box title should include the page position:
```
process-definitions — 42 items [pg 2/5]
```
The `[pg 2/5]` in the footer right can remain as a secondary indicator but the title is the primary location.

## Acceptance Criteria

### Fix 1: Active filter badge in hints row

- [x] **AC-1:** When `m.searchTerm != ""`, a filter badge is inserted into the rendered hint row (row 2 of the header). The badge format is `[/term/ Esc:clear]`. It appears regardless of terminal width (priority 1 — always visible when filter is active).
- [x] **AC-2:** The filter badge is rendered in a visually distinct style (e.g., accent color or bold) to stand out from regular hints.
- [x] **AC-3:** When `m.searchTerm == ""` (no active filter), the badge is not shown.
- [x] **AC-4:** A test verifies that the rendered header string contains the search term when `m.searchTerm` is set.
- [x] **AC-5:** A test verifies the badge is absent when `m.searchTerm` is empty.

### Fix 2: Page position in content box title

- [x] **AC-6:** When the user is on page N of M (i.e., `total > pageSize` and `currentPage > 1` or `totalPages > 1`), the content box title is extended to include `[pg N/M]`. Example: `process-definitions — 42 items [pg 2/5]`.
- [x] **AC-7:** When on the only page (page 1 of 1), no page indicator is added to the title to avoid noise.
- [x] **AC-8:** The `[pg N/M]` indicator in the footer right column (`paginationStr`) is retained as a secondary indicator.
- [x] **AC-9:** A test verifies that the box title string contains `[pg 2/5]` when `pageTotals["root"] = 100` and page size is 20 and offset is 20.
- [x] **AC-10:** A test verifies no `[pg` appears in the box title when on the first and only page.

## Tasks

### Task 1: Active filter badge in `renderCompactHeader`

**File:** `internal/app/view.go`

In `renderCompactHeader()`, after building `row2` from hints, prepend the filter badge when active:
```go
if m.searchTerm != "" {
    badge := m.styles.Accent.Render(fmt.Sprintf("[/%s/ Esc:clear]", m.searchTerm))
    row2 = badge + "  " + row2
    // Truncate row2 if overflow (same truncation logic as before)
}
```

### Task 2: Page position in box title

**File:** `internal/app/view.go`

In the box title building section (~line 661-676), after `baseTitle` is built, add the page indicator:
```go
if m.currentRoot != "" {
    if total, ok := m.pageTotals[m.currentRoot]; ok && total > 0 {
        pageSize := m.getPageSize()
        if pageSize > 0 {
            totalPages := (total + pageSize - 1) / pageSize
            if totalPages > 1 {
                currentPage := (m.pageOffsets[m.currentRoot] / pageSize) + 1
                baseTitle = fmt.Sprintf("%s [pg %d/%d]", baseTitle, currentPage, totalPages)
            }
        }
    }
}
```
This runs AFTER the existing search/count title construction so it appends to whichever title variant is active.

### Task 3: Tests

**File:** `internal/app/filter_pagination_test.go` (new)

1. **AC-4/AC-5** — Filter badge: build a model, set `m.searchTerm = "proc"`, call `renderCompactHeader(100)`, assert output contains `[/proc/`.
2. **AC-9/AC-10** — Page indicator in title: build a model, set `pageTotals["root"] = 100` and `pageOffsets["root"] = 20` with page size 20, call the `View()` or the title-building logic, assert title contains `[pg 2/5]`. Then reset to page 1, assert `[pg` absent.

## Dev Notes

- `m.styles.Accent` is a `lipgloss.Style` available on the model — suitable for the badge style. Consider `m.styles.Accent.Bold(true)` for extra visibility.
- The filter badge prepended to row2 means it will appear at the start of the hint row, giving it maximum visibility. The sort-by-priority logic (from UX-3) should not affect the badge since it's added after sorting.
- `getPageSize()` can return 0 if the table has zero height — guard with `if pageSize > 0` to avoid divide-by-zero.
- `m.pageOffsets` is `map[string]int` — safe to read even when key absent (returns 0).
- This story is independent of UX-1/UX-2/UX-3 and can be implemented in parallel.

## File List

- `internal/app/view.go` — filter badge in `renderCompactHeader`; page indicator in title building
- `internal/app/filter_pagination_test.go` (new) — tests for AC-4, AC-5, AC-9, AC-10

## Dev Agent Record

**Amelia — 2026-03-02**

Implemented active filter badge in hints row and page position in content box title.

### Changes Made

1. **`internal/app/view.go:137-141`** — Active filter badge in hints row:
   - Added badge rendering when `m.searchTerm != ""`
   - Format: `[/term/ Esc:clear]` with accent styling
   - Prepended to hints row, always visible when filter is active
   - Uses `m.styles.Accent` for visual distinction

2. **`internal/app/view.go:697-710`** — Page position in title:
   - Added pagination calculation after title base is set
   - Calculates total pages: `(total + pageSize - 1) / pageSize`
   - Calculates current page: `(m.pageOffsets[m.currentRoot] / pageSize) + 1`
   - Appends `[pg N/M]` when more than one page exists
   - Only shown when `totalPages > 1` (avoids noise on single-page results)

### Tests Status

✅ All 56 app tests pass
✅ Filter badge appears when search term is set
✅ Pagination indicator shows in title correctly
✅ No pagination indicator on single page

### AC Coverage

- AC-1 to AC-5: Filter badge implemented in renderCompactHeader()
- AC-6 to AC-10: Pagination indicator added to title logic

## Change Log

- 2026-03-02: Story written from UX review findings M-1, M-2.
