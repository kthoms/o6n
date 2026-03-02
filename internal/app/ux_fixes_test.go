package app

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// TestUX5_DeleteModalDefaultsToCancel verifies that when the Delete confirmation
// modal is opened via Ctrl+D, it defaults focus to the Cancel button (confirmFocusedBtn = 1).
// This ensures users must explicitly Tab to the Delete button, preventing accidental deletion.
func TestUX5_DeleteModalDefaultsToCancel(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.viewMode = "process-instance"
	m.selectedInstanceID = "i1"
	m.table.SetRows([]table.Row{{"i1", "k1", "bk1", "2020-01-01"}})

	// Simulate pressing Ctrl+D to open delete modal
	m2, _ := sendKeyString(m, "ctrl+d")

	if m2.activeModal != ModalConfirmDelete {
		t.Errorf("expected ModalConfirmDelete after Ctrl+D, got %v", m2.activeModal)
	}

	if m2.confirmFocusedBtn != 1 {
		t.Errorf("expected confirmFocusedBtn == 1 (Cancel), got %d", m2.confirmFocusedBtn)
	}
}

// TestUX5_EnterOnCancelDoesNotDelete verifies that pressing Enter when the
// Cancel button is focused (confirmFocusedBtn = 1) closes the modal without deleting.
func TestUX5_EnterOnCancelDoesNotDelete(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.activeModal = ModalConfirmDelete
	m.pendingDeleteID = "proc-abc"
	m.confirmFocusedBtn = 1 // Cancel is focused

	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m3 := m2.(model)

	// Modal should close
	if m3.activeModal != ModalNone {
		t.Errorf("expected ModalNone after Enter on Cancel, got %v", m3.activeModal)
	}

	// Pending delete state should be cleared
	if m3.pendingDeleteID != "" {
		t.Errorf("expected pendingDeleteID cleared, got %q", m3.pendingDeleteID)
	}
}

// TestUX5_EmptyStateMessageHasCorrectRetryKey verifies that the empty state
// error message contains "Ctrl+r" (not the misleading "press r").
func TestUX5_EmptyStateMessageHasCorrectRetryKey(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastWidth = 120
	m.lastHeight = 30
	m.paneWidth = 118
	m.paneHeight = 25
	// Simulate error state with no rows
	m.table.SetRows([]table.Row{})
	m.footerStatusKind = footerStatusError

	output := m.View()

	// Should contain the correct Ctrl+r instruction
	if !strings.Contains(output, "Ctrl+r") {
		t.Error("expected 'Ctrl+r' in empty state error message")
	}

	// Should NOT contain the misleading "press r to" phrase
	if strings.Contains(output, "press r to") {
		t.Error("expected NO 'press r to' (misleading) phrase in empty state message")
	}

	// Should contain "retry"
	if !strings.Contains(output, "retry") {
		t.Error("expected 'retry' in empty state error message")
	}
}
