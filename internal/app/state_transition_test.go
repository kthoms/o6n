package app

import (
	"testing"

	table "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// ── Task 1: prepareStateTransition exists ──────────────────────────────────

func TestPrepareStateTransitionExists(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	// Should not panic — just verify all scopes can be called
	m.prepareStateTransition(transitionEnvSwitch)
	m.prepareStateTransition(transitionContextSwitch)
	m.prepareStateTransition(transitionDrilldown)
	m.prepareStateTransition(transitionBack)
	m.prepareStateTransition(transitionBreadcrumb, 0)
}

// ── Task 2: Environment switch clears all leaking state ────────────────────

func TestEnvSwitchClearsNavigationStack(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	// simulate a drilldown push
	m.navigationStack = []viewState{{viewMode: "process-definition"}}

	m.prepareStateTransition(transitionEnvSwitch)

	if m.navigationStack != nil {
		t.Errorf("expected navigationStack nil after env switch, got %v", m.navigationStack)
	}
}

func TestEnvSwitchClearsGenericParams(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.genericParams = map[string]string{"processInstanceId": "abc"}

	m.prepareStateTransition(transitionEnvSwitch)

	if len(m.genericParams) != 0 {
		t.Errorf("expected empty genericParams after env switch, got %v", m.genericParams)
	}
}

func TestEnvSwitchClearsSelectedKeys(t *testing.T) {
	m := newTestModel(t)
	m.selectedDefinitionKey = "my-process"
	m.selectedInstanceID = "inst-123"

	m.prepareStateTransition(transitionEnvSwitch)

	if m.selectedDefinitionKey != "" {
		t.Errorf("expected empty selectedDefinitionKey, got %q", m.selectedDefinitionKey)
	}
	if m.selectedInstanceID != "" {
		t.Errorf("expected empty selectedInstanceID, got %q", m.selectedInstanceID)
	}
}

// ── Task 3: Context switch clears sort state ───────────────────────────────

func TestContextSwitchClearsSortState(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.sortColumn = 3
	m.sortAscending = false

	m.prepareStateTransition(transitionContextSwitch)

	if m.sortColumn != -1 {
		t.Errorf("expected sortColumn -1 after context switch, got %d", m.sortColumn)
	}
	if !m.sortAscending {
		t.Errorf("expected sortAscending true after context switch")
	}
}

func TestContextSwitchClearsNavStack(t *testing.T) {
	m := newTestModel(t)
	m.navigationStack = []viewState{{viewMode: "process-definition"}, {viewMode: "process-instance"}}

	m.prepareStateTransition(transitionContextSwitch)

	if m.navigationStack != nil {
		t.Errorf("expected nil navigationStack after context switch, got %v", m.navigationStack)
	}
}

// ── Task 4: Breadcrumb navigation truncates navStack ──────────────────────

func TestBreadcrumbNavToRootClearsNavStack(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.navigationStack = []viewState{
		{viewMode: "a"},
		{viewMode: "b"},
		{viewMode: "c"},
	}

	m.prepareStateTransition(transitionBreadcrumb, 0)

	if m.navigationStack != nil {
		t.Errorf("expected nil navStack after breadcrumb(0), got %v", m.navigationStack)
	}
}

func TestBreadcrumbNavToDepth1TruncatesNavStack(t *testing.T) {
	m := newTestModel(t)
	m.navigationStack = []viewState{
		{viewMode: "a"}, // depth 0→1 transition saved here
		{viewMode: "b"}, // depth 1→2
		{viewMode: "c"}, // depth 2→3
	}

	m.prepareStateTransition(transitionBreadcrumb, 1)

	if len(m.navigationStack) != 1 {
		t.Errorf("expected navStack len 1 after breadcrumb(1), got %d", len(m.navigationStack))
	}
}

// ── Task 5: All transitions clear sort and search ─────────────────────────

func TestAllTransitionsClearSortAndSearch(t *testing.T) {
	scopes := []transitionScope{
		transitionEnvSwitch,
		transitionContextSwitch,
		transitionDrilldown,
		transitionBack,
	}
	for _, scope := range scopes {
		m := newTestModel(t)
		m.sortColumn = 2
		m.sortAscending = false
		m.searchTerm = "myterm"
		m.originalRows = []table.Row{{"id1"}}

		m.prepareStateTransition(scope)

		if m.sortColumn != -1 {
			t.Errorf("scope %d: expected sortColumn -1, got %d", scope, m.sortColumn)
		}
		if !m.sortAscending {
			t.Errorf("scope %d: expected sortAscending true", scope)
		}
		if m.searchTerm != "" {
			t.Errorf("scope %d: expected empty searchTerm, got %q", scope, m.searchTerm)
		}
		if m.originalRows != nil {
			t.Errorf("scope %d: expected nil originalRows", scope)
		}
	}
}

// ── Task 6: Cursor bounds after row deletion ──────────────────────────────

func TestCursorBoundsAfterTerminate(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false

	// Set up table with 3 rows, cursor at row 2 (last)
	cols := []table.Column{{Title: "ID", Width: 10}}
	rows := []table.Row{{"row1"}, {"row2"}, {"row3"}}
	m.table.SetColumns(cols)
	m.table.SetRows(rows)
	m.table.SetCursor(2)

	// Terminate the row at index 2 — after removal only 2 rows remain
	m2raw, _ := m.Update(terminatedMsg{id: "row3"})
	m2 := m2raw.(model)

	remaining := m2.table.Rows()
	cursor := m2.table.Cursor()
	if cursor >= len(remaining) && len(remaining) > 0 {
		t.Errorf("cursor %d out of bounds for %d rows after terminate", cursor, len(remaining))
	}
}

func TestCursorBoundsAfterDeleteOnlyRow(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false

	cols := []table.Column{{Title: "ID", Width: 10}}
	rows := []table.Row{{"only-row"}}
	m.table.SetColumns(cols)
	m.table.SetRows(rows)
	m.table.SetCursor(0)

	m2raw, _ := m.Update(terminatedMsg{id: "only-row"})
	m2 := m2raw.(model)

	cursor := m2.table.Cursor()
	if cursor > 0 {
		t.Errorf("expected cursor <= 0 after deleting only row, got %d", cursor)
	}
}

// ── Integration: Esc after env switch doesn't restore stale state ─────────

func TestEscAfterEnvSwitchDoesNotRestoreOldStack(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	// Simulate drilled-in state in env A
	m.navigationStack = []viewState{{viewMode: "process-definition"}}
	m.genericParams = map[string]string{"processInstanceId": "abc"}

	// Simulate env switch (what prepareStateTransition(envSwitch) does)
	m.prepareStateTransition(transitionEnvSwitch)

	// Now press Esc — should NOT pop the old stack
	m2raw, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m2 := m2raw.(model)

	if len(m2.navigationStack) > 0 {
		t.Errorf("expected empty navStack after Esc post-env-switch, got %v", m2.navigationStack)
	}
}
