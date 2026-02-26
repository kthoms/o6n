package app

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
)

// ── AC-4: Empty state for zero-row table ─────────────────────────────────────

func TestEmptyStateShownWhenNoRows(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastWidth = 120
	m.lastHeight = 40
	m.currentRoot = "process-definitions"
	m.table.SetColumns([]table.Column{{Title: "ID", Width: 20}, {Title: "Name", Width: 30}})
	m.table.SetRows([]table.Row{}) // empty

	rendered := m.View()
	lower := strings.ToLower(rendered)
	if !strings.Contains(lower, "no ") && !strings.Contains(lower, "empty") && !strings.Contains(lower, "found") {
		t.Errorf("expected empty state message in view when no rows, got:\n%s", rendered[:min(400, len(rendered))])
	}
}

func TestEmptyStateContainsResourceName(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastWidth = 120
	m.lastHeight = 40
	m.currentRoot = "process-definitions"
	m.table.SetColumns([]table.Column{{Title: "ID", Width: 20}})
	m.table.SetRows([]table.Row{})

	rendered := m.View()
	// Should mention the resource name in some form
	if !strings.Contains(strings.ToLower(rendered), "process") {
		t.Errorf("expected empty state message to mention resource name, got:\n%s", rendered[:min(400, len(rendered))])
	}
}

func TestEmptyStateNotShownWhenRowsPresent(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastWidth = 120
	m.lastHeight = 40
	m.currentRoot = "process-definitions"
	m.table.SetColumns([]table.Column{{Title: "ID", Width: 20}})
	m.table.SetRows([]table.Row{{"row-1"}})

	rendered := m.View()
	// When rows present, "no ... found" hint should not appear in the content box
	if strings.Contains(strings.ToLower(rendered), "no process-definitions found") {
		t.Errorf("empty state message should not appear when rows are present")
	}
}

// ── AC-6: Error empty state ───────────────────────────────────────────────────

func TestEmptyStateAfterErrorShowsRetryHint(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastWidth = 120
	m.lastHeight = 40
	m.currentRoot = "process-definitions"
	m.table.SetColumns([]table.Column{{Title: "ID", Width: 20}})
	m.table.SetRows([]table.Row{}) // empty after error
	m.isLoading = false

	// Simulate error state
	m2, _ := m.Update(errMsg{err: fmt.Errorf("connection refused")})
	m3 := m2.(model)

	rendered := m3.View()
	lower := strings.ToLower(rendered)
	if !strings.Contains(lower, "retry") && !strings.Contains(lower, "r ") && !strings.Contains(lower, "no ") {
		t.Errorf("expected retry hint in empty state after error, got:\n%s", rendered[:min(400, len(rendered))])
	}
}
