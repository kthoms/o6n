package app

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kthoms/o8n/internal/config"
)

func TestUserRowsNoDrillPrefix(t *testing.T) {
	envCfg, err := config.LoadEnvConfig("o8n-env.yaml")
	if err != nil {
		t.Fatalf("load env config: %v", err)
	}
	appCfg, err := config.LoadAppConfig("o8n-cfg.yaml")
	if err != nil {
		t.Fatalf("load app config: %v", err)
	}
	m := newModelEnvApp(envCfg, appCfg, "")
	m.breadcrumb = []string{"user"}

	item := map[string]interface{}{"id": "u1", "firstName": "Alice", "lastName": "Smith"}
	cols := m.buildColumnsFor("user", m.paneWidth-4)
	r := make(table.Row, len(cols))
	for i, col := range cols {
		key := strings.ToLower(col.Title)
		val := ""
		if v, ok := item[key]; ok {
			val = fmt.Sprintf("%v", v)
		}
		r[i] = val
	}
	def := m.findTableDef("user")
	hasDrilldown := def != nil && def.Drilldown != nil
	if !hasDrilldown && strings.HasPrefix(r[0], "▶ ") {
		t.Fatalf("expected no drill prefix for user, got %q", r[0])
	}
}

func TestEnterNoopOnUser(t *testing.T) {
	envCfg, err := config.LoadEnvConfig("o8n-env.yaml")
	if err != nil {
		t.Fatalf("load env config: %v", err)
	}
	appCfg, err := config.LoadAppConfig("o8n-cfg.yaml")
	if err != nil {
		t.Fatalf("load app config: %v", err)
	}
	m := newModelEnvApp(envCfg, appCfg, "")
	m.breadcrumb = []string{"user"}
	m.table.SetRows([]table.Row{{"u1", "Alice"}})
	m.table.SetCursor(0)

	origRoot := m.currentRoot
	origStack := len(m.navigationStack)

	ret, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newM := ret.(model)
	if newM.currentRoot != origRoot {
		t.Fatalf("expected currentRoot unchanged, got %q", newM.currentRoot)
	}
	if len(newM.navigationStack) != origStack {
		t.Fatalf("expected navigation stack unchanged, got len=%d", len(newM.navigationStack))
	}
}

func TestUserNavigateActionH(t *testing.T) {
	envCfg, err := config.LoadEnvConfig("o8n-env.yaml")
	if err != nil {
		t.Fatalf("load env config: %v", err)
	}
	appCfg, err := config.LoadAppConfig("o8n-cfg.yaml")
	if err != nil {
		t.Fatalf("load app config: %v", err)
	}
	m := newModelEnvApp(envCfg, appCfg, "")
	m.breadcrumb = []string{"user"}
	m.table.SetRows([]table.Row{{"u1", "Alice"}})
	m.table.SetCursor(0)

	// Build actions and ensure 'h' navigate exists
	items := m.buildActionsForRoot()
	found := false
	for _, it := range items {
		if it.key == "h" && strings.HasSuffix(it.label, " →") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected navigate action 'h' for user")
	}

	// Open actions menu and simulate pressing 'h'
	m.activeModal = ModalActionMenu
	m.actionsMenuItems = items
	ret, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	newM := ret.(model)
	if newM.currentRoot != "history-user-operation" {
		t.Fatalf("expected currentRoot to become history-user-operation, got %q", newM.currentRoot)
	}
	last := newM.breadcrumb[len(newM.breadcrumb)-1]
	if last != "View Operation Log" {
		t.Fatalf("expected breadcrumb to include 'View Operation Log', got %q", last)
	}
}

func TestBuildActionsForRootIncludesNavigateLabel(t *testing.T) {
	envCfg, err := config.LoadEnvConfig("o8n-env.yaml")
	if err != nil {
		t.Fatalf("load env config: %v", err)
	}
	appCfg, err := config.LoadAppConfig("o8n-cfg.yaml")
	if err != nil {
		t.Fatalf("load app config: %v", err)
	}
	m := newModelEnvApp(envCfg, appCfg, "")
	m.breadcrumb = []string{"user"}
	items := m.buildActionsForRoot()
	found := false
	for _, it := range items {
		if it.key == "h" && it.label == "View Operation Log →" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected action item with key=h and label='View Operation Log →'")
	}
}
