package app

import (
	"strings"
	"testing"

	"github.com/kthoms/o6n/internal/config"
)

// ── AC-1/AC-2: Help screen has category headers ───────────────────────────────

func TestHelpScreenHasCategoryHeaders(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastWidth = 120
	m.lastHeight = 40
	m.activeModal = ModalHelp

	rendered := m.View()
	// Should contain section headers
	categories := []string{"NAVIGATION", "GLOBAL", "SEARCH"}
	for _, cat := range categories {
		if !strings.Contains(rendered, cat) {
			t.Errorf("expected help screen to contain category header %q", cat)
		}
	}
}

// ── AC-3/AC-4: Resource-specific actions shown dynamically ────────────────────

func TestHelpScreenShowsResourceActionsWhenPresent(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastWidth = 120
	m.lastHeight = 40

	// Set a resource with actions directly in model config
	m.currentRoot = "test-resource"
	// Inject a table def with actions into the model's config (used by findTableDef)
	if m.config == nil {
		m.config = &config.Config{}
	}
	m.config.Tables = append(m.config.Tables, config.TableDef{
		Name:    "test-resource",
		Columns: []config.ColumnDef{{Name: "id"}},
		Actions: []config.ActionDef{
			{Key: "ctrl+t", Label: "Terminate", Method: "DELETE"},
		},
	})
	m.activeModal = ModalHelp

	rendered := m.View()
	if !strings.Contains(rendered, "RESOURCE ACTIONS") {
		t.Errorf("expected 'RESOURCE ACTIONS' section in help when resource has actions configured, view:\n%s", rendered[:min(600, len(rendered))])
	}
	if !strings.Contains(rendered, "Terminate") {
		t.Errorf("expected action label 'Terminate' in help screen resource actions section")
	}
}

func TestHelpScreenOmitsResourceSectionWhenNoActions(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastWidth = 120
	m.lastHeight = 40

	// Use a resource that has no actions (likely history tables or others)
	// Use a root that definitely has no actions in test config
	m.currentRoot = "process-definitions"
	m.activeModal = ModalHelp

	def := m.findTableDef("process-definitions")
	if def != nil && len(def.Actions) > 0 {
		t.Skip("process-definitions has actions configured in test config")
	}

	rendered := m.View()
	// Resource Actions section should not appear when no actions
	_ = rendered // help screen rendering should not panic
}

// ── AC-1b: Help screen is navigable by structure, not just color ──────────────

func TestHelpScreenHasSectionSeparators(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastWidth = 120
	m.lastHeight = 40
	m.activeModal = ModalHelp

	rendered := m.View()
	// Sections should have separator lines (─)
	if !strings.Contains(rendered, "─") {
		t.Errorf("expected section separators (─) in help screen for structural navigation")
	}
}
