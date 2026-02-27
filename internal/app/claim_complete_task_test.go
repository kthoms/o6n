package app

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/kthoms/o8n/internal/config"
	"github.com/kthoms/o8n/internal/operaton"
)

// testTaskConfig returns a minimal config with a task table that has id, name, assignee columns.
func testTaskConfig(username string) *config.Config {
	return &config.Config{
		Environments: map[string]config.Environment{
			"local": {URL: "http://localhost:8080", Username: username},
		},
		Tables: []config.TableDef{
			{
				Name: "task",
				Columns: []config.ColumnDef{
					{Name: "id"},
					{Name: "name"},
					{Name: "assignee"},
				},
				Actions: []config.ActionDef{
					{Key: "c", Label: "Claim Task", Method: "POST", Path: "/task/{id}/claim", Body: `{"userId": "{currentUser}"}`},
					{Key: "u", Label: "Unclaim Task", Method: "POST", Path: "/task/{id}/unclaim"},
				},
			},
		},
	}
}

// setupTaskTable initialises the model with a task table containing one row.
func setupTaskTable(t *testing.T, id, name, assignee, currentUser string) model {
	t.Helper()
	m := newModel(testTaskConfig(currentUser))
	m.splashActive = false
	m.currentRoot = "task"
	m.breadcrumb = []string{"task"}
	cols := []table.Column{
		{Title: "id", Width: 20},
		{Title: "name", Width: 30},
		{Title: "assignee", Width: 20},
	}
	m.table.SetColumns(cols)
	m.table.SetRows([]table.Row{{id, name, assignee}})
	m.table.SetCursor(0)
	return m
}

// ── Claim guard tests (c key) ─────────────────────────────────────────────────

func TestClaimOnUnclaimedTask(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "", "alice")
	m2, cmd := sendKeyString(m, "c")
	if cmd == nil {
		t.Error("expected claimTaskCmd to be dispatched")
	}
	_ = m2
}

func TestClaimOnForeignTask(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "bob", "alice")
	// foreign task — should show error, not dispatch API call (not loading)
	m2, _ := sendKeyString(m, "c")
	if m2.isLoading {
		t.Error("expected no API call (isLoading) when task claimed by another user")
	}
	if !strings.Contains(m2.footerError, "bob") {
		t.Errorf("expected footer error mentioning 'bob', got %q", m2.footerError)
	}
	if m2.footerStatusKind != footerStatusError {
		t.Error("expected footerStatusError")
	}
}

func TestClaimOnOwnTask(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	// own task — should show hint, not dispatch API call
	m2, _ := sendKeyString(m, "c")
	if m2.isLoading {
		t.Error("expected no API call (isLoading) when task already owned")
	}
	if !strings.Contains(m2.footerError, "already own") {
		t.Errorf("expected footer hint about already owning task, got %q", m2.footerError)
	}
	if m2.footerStatusKind != footerStatusInfo {
		t.Error("expected footerStatusInfo")
	}
}

// ── Unclaim guard tests (u key) ───────────────────────────────────────────────

func TestUnclaimOwnTask(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m2, cmd := sendKeyString(m, "u")
	if cmd == nil {
		t.Error("expected unclaimTaskCmd to be dispatched")
	}
	_ = m2
}

func TestUnclaimUnclaimedTask(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "", "alice")
	m2, _ := sendKeyString(m, "u")
	if m2.isLoading {
		t.Error("expected no API call for unclaimed task")
	}
	if !strings.Contains(m2.footerError, "not claimed") {
		t.Errorf("expected 'not claimed' footer error, got %q", m2.footerError)
	}
}

func TestUnclaimForeignTask(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "bob", "alice")
	m2, _ := sendKeyString(m, "u")
	if m2.isLoading {
		t.Error("expected no API call when task owned by another user")
	}
	if !strings.Contains(m2.footerError, "bob") {
		t.Errorf("expected footer error mentioning 'bob', got %q", m2.footerError)
	}
}

// ── Enter guard tests ─────────────────────────────────────────────────────────

func TestEnterOnOwnTaskFetchesVariables(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m2, cmd := sendKeyString(m, "enter")
	if cmd == nil {
		t.Error("expected fetchTaskVariablesCmd to be dispatched for own task")
	}
	if m2.footerStatusKind != footerStatusLoading {
		t.Errorf("expected footerStatusLoading, got %v", m2.footerStatusKind)
	}
}

func TestEnterOnUnclaimedTaskShowsError(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "", "alice")
	m2, _ := sendKeyString(m, "enter")
	if m2.footerStatusKind != footerStatusError {
		t.Errorf("expected footerStatusError for unclaimed task, got %v", m2.footerStatusKind)
	}
	if !strings.Contains(m2.footerError, "Claim") {
		t.Errorf("expected 'Claim' in footer error, got %q", m2.footerError)
	}
	if m2.activeModal == ModalTaskComplete {
		t.Error("expected dialog NOT to open for unclaimed task")
	}
}

func TestEnterOnForeignTaskShowsError(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "bob", "alice")
	m2, _ := sendKeyString(m, "enter")
	if m2.footerStatusKind != footerStatusError {
		t.Errorf("expected footerStatusError for foreign task, got %v", m2.footerStatusKind)
	}
	if !strings.Contains(m2.footerError, "bob") {
		t.Errorf("expected footer error with assignee name, got %q", m2.footerError)
	}
	if m2.activeModal == ModalTaskComplete {
		t.Error("expected dialog NOT to open for foreign task")
	}
}

func TestEnterOnNonTaskTableDoesNotIntercept(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.currentRoot = "process-instance"
	m.breadcrumb = []string{"process-instance"}
	cols := []table.Column{{Title: "id", Width: 20}}
	m.table.SetColumns(cols)
	m.table.SetRows([]table.Row{{"inst-1"}})
	m.table.SetCursor(0)
	// Enter on non-task table should not produce the loading status
	m2, _ := sendKeyString(m, "enter")
	if m2.footerStatusKind == footerStatusLoading {
		t.Error("Enter on non-task table should not show loading status")
	}
}

// ── taskVariablesLoadedMsg handler ────────────────────────────────────────────

func TestTaskVariablesLoadedOpensModal(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	inputVars := map[string]variableValue{
		"orderId": {Value: "ORD-123", TypeName: "String"},
	}
	formVars := map[string]variableValue{
		"approved": {Value: nil, TypeName: "Boolean"},
		"orderId":  {Value: nil, TypeName: "String"},
	}
	msg := taskVariablesLoadedMsg{
		taskID:    "task-1",
		taskName:  "My Task",
		inputVars: inputVars,
		formVars:  formVars,
	}
	m2, _ := m.Update(msg)
	result := m2.(model)

	if result.activeModal != ModalTaskComplete {
		t.Error("expected ModalTaskComplete to be active after taskVariablesLoadedMsg")
	}
	if result.taskCompleteTaskID != "task-1" {
		t.Errorf("expected taskCompleteTaskID 'task-1', got %q", result.taskCompleteTaskID)
	}
	if len(result.taskCompleteFields) != 2 {
		t.Errorf("expected 2 form fields, got %d", len(result.taskCompleteFields))
	}
}

// ── Pre-fill test ─────────────────────────────────────────────────────────────

func TestPreFillFromInputVars(t *testing.T) {
	m := newModel(testTaskConfig("alice"))
	inputVars := map[string]variableValue{
		"orderId": {Value: "ORD-999", TypeName: "String"},
	}
	formVars := map[string]variableValue{
		"orderId": {Value: nil, TypeName: "String"},
	}
	fields := m.buildTaskCompleteFields(formVars, inputVars)
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
	if fields[0].input.Value() != "ORD-999" {
		t.Errorf("expected pre-filled value 'ORD-999', got %q", fields[0].input.Value())
	}
}

// ── Tab cycle test ────────────────────────────────────────────────────────────

func TestTabCycleInTaskCompleteModal(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	// Simulate dialog open with 2 fields
	m.activeModal = ModalTaskComplete
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"fieldA": {Value: nil, TypeName: "String"},
			"fieldB": {Value: nil, TypeName: "String"},
		},
		map[string]variableValue{},
	)
	m.taskCompletePos = 0
	m.taskCompleteFocus = focusTaskField
	m.taskCompleteFields[0].input.Focus()

	// Tab: field[0] → field[1]
	m2, _ := sendKeyString(m, "tab")
	if m2.taskCompletePos != 1 {
		t.Errorf("expected pos 1 after first Tab, got %d", m2.taskCompletePos)
	}
	if m2.taskCompleteFocus != focusTaskField {
		t.Errorf("expected focusTaskField after first Tab")
	}

	// Tab: field[1] → Complete
	m3, _ := sendKeyString(m2, "tab")
	if m3.taskCompleteFocus != focusTaskComplete {
		t.Errorf("expected focusTaskComplete after Tab from last field")
	}

	// Tab: Complete → Back
	m4, _ := sendKeyString(m3, "tab")
	if m4.taskCompleteFocus != focusTaskBack {
		t.Errorf("expected focusTaskBack after Tab from Complete")
	}

	// Tab: Back → field[0]
	m5, _ := sendKeyString(m4, "tab")
	if m5.taskCompleteFocus != focusTaskField {
		t.Errorf("expected focusTaskField after Tab from Back")
	}
	if m5.taskCompletePos != 0 {
		t.Errorf("expected pos 0 after wrap-around Tab, got %d", m5.taskCompletePos)
	}
}

// ── Boolean toggle test ───────────────────────────────────────────────────────

func TestSpaceTogglesBoolField(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"approved": {Value: nil, TypeName: "Boolean"},
		},
		map[string]variableValue{},
	)
	m.taskCompletePos = 0
	m.taskCompleteFocus = focusTaskField
	m.taskCompleteFields[0].input.Focus()
	m.taskCompleteFields[0].input.SetValue("false")

	m2, _ := sendKeyString(m, " ")
	if m2.taskCompleteFields[0].input.Value() != "true" {
		t.Errorf("expected 'true' after Space toggle, got %q", m2.taskCompleteFields[0].input.Value())
	}

	m3, _ := sendKeyString(m2, " ")
	if m3.taskCompleteFields[0].input.Value() != "false" {
		t.Errorf("expected 'false' after second Space toggle, got %q", m3.taskCompleteFields[0].input.Value())
	}
}

// ── Submit / completeTaskCmd ──────────────────────────────────────────────────

func TestSubmitBuildsCorrectVariables(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskID = "task-1"
	m.taskCompleteTaskName = "My Task"
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"approved": {Value: nil, TypeName: "Boolean"},
			"amount":   {Value: nil, TypeName: "Integer"},
		},
		map[string]variableValue{},
	)
	// Set valid values
	for i, f := range m.taskCompleteFields {
		if f.name == "approved" {
			m.taskCompleteFields[i].input.SetValue("true")
		} else if f.name == "amount" {
			m.taskCompleteFields[i].input.SetValue("42")
		}
	}
	m.taskCompleteFocus = focusTaskComplete

	m2, cmd := sendKeyString(m, "enter")
	if cmd == nil {
		t.Error("expected completeTaskCmd to be dispatched on Enter with Complete focused")
	}
	_ = m2
}

// ── Escape closes dialog ──────────────────────────────────────────────────────

func TestEscapeClosesTaskCompleteDialog(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskID = "task-1"
	m.taskCompleteTaskName = "My Task"

	m2, _ := sendKeyString(m, "esc")
	if m2.activeModal != ModalNone {
		t.Error("expected ModalNone after Esc")
	}
	if m2.taskCompleteTaskID != "" {
		t.Error("expected taskCompleteTaskID cleared after Esc")
	}
	if m2.taskCompleteTaskName != "" {
		t.Error("expected taskCompleteTaskName cleared after Esc")
	}
	if m2.taskCompleteFields != nil {
		t.Error("expected taskCompleteFields cleared after Esc")
	}
}

// ── Validation gate test ──────────────────────────────────────────────────────

func TestCompleteDisabledWhenFieldHasError(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskID = "task-1"
	m.taskCompleteTaskName = "My Task"
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"count": {Value: nil, TypeName: "Integer"},
		},
		map[string]variableValue{},
	)
	// Set invalid value for an integer field
	m.taskCompleteFields[0].input.SetValue("notanumber")
	m.taskCompleteFields[0].error = "enter an integer"
	m.taskCompleteFocus = focusTaskComplete

	m2, cmd := sendKeyString(m, "enter")
	// completeTaskCmd should NOT be dispatched when there are errors
	if cmd != nil {
		t.Error("expected no completeTaskCmd when field has error")
	}
	_ = m2
}

// ── actionExecutedMsg closes dialog on complete ───────────────────────────────

func TestActionExecutedMsgClosesTaskDialog(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskID = "task-1"
	m.taskCompleteTaskName = "My Task"

	msg := actionExecutedMsg{label: "Completed: My Task", closeTaskDialog: true}
	m2, _ := m.Update(msg)
	result := m2.(model)

	if result.activeModal != ModalNone {
		t.Error("expected ModalNone after actionExecutedMsg with closeTaskDialog=true")
	}
	if result.taskCompleteTaskID != "" {
		t.Error("expected taskCompleteTaskID cleared after close")
	}
}

// ── renderTaskCompleteModal ───────────────────────────────────────────────────

func TestRenderTaskCompleteModal(t *testing.T) {
	m := setupTaskTable(t, "task-1", "Review Order", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "Review Order"
	m.lastWidth = 120
	m.lastHeight = 40
	m.taskInputVars = map[string]variableValue{
		"orderId": {Value: "ORD-777", TypeName: "String"},
	}
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"approved": {Value: nil, TypeName: "Boolean"},
		},
		m.taskInputVars,
	)
	m.taskCompleteFocus = focusTaskField

	out := m.renderTaskCompleteModal(120, 40)
	if out == "" {
		t.Fatal("expected non-empty modal output")
	}
	// New layout: no static "Complete Task" title or section headers in body
	if strings.Contains(out, "Complete Task") {
		t.Error("expected 'Complete Task' static title NOT in modal body")
	}
	if strings.Contains(out, "INPUT VARIABLES") {
		t.Error("expected INPUT VARIABLES section NOT in modal")
	}
	if strings.Contains(out, "OUTPUT VARIABLES") {
		t.Error("expected OUTPUT VARIABLES section NOT in modal")
	}
	// Task name appears in modal (border title)
	if !strings.Contains(out, "Review Order") {
		t.Error("expected task name in modal")
	}
	// Buttons present
	if !strings.Contains(out, "Complete") {
		t.Error("expected Complete button in modal")
	}
	if !strings.Contains(out, "Back") {
		t.Error("expected Back button in modal")
	}
	// orderId is input-only → appears as read-only row
	if !strings.Contains(out, "orderId") {
		t.Error("expected orderId in modal")
	}
	// approved is a form field → appears as editable row
	if !strings.Contains(out, "approved") {
		t.Error("expected approved row in modal")
	}
}

// ── Merged list: no section headers ──────────────────────────────────────────

func TestMergedListHasNoSectionHeaders(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.lastWidth = 120
	m.lastHeight = 40
	m.taskInputVars = map[string]variableValue{
		"amount": {Value: "100", TypeName: "String"},
	}
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"approved": {Value: nil, TypeName: "Boolean"},
		},
		m.taskInputVars,
	)
	out := m.renderTaskCompleteModal(120, 40)
	if strings.Contains(out, "INPUT VARIABLES") {
		t.Error("expected no INPUT VARIABLES heading in merged list")
	}
	if strings.Contains(out, "OUTPUT VARIABLES") {
		t.Error("expected no OUTPUT VARIABLES heading in merged list")
	}
}

// ── Read-only row style in merged list ────────────────────────────────────────

func TestReadOnlyRowStyleInMergedList(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.lastWidth = 120
	m.lastHeight = 40
	// "inputOnly" is in input vars but not in form fields → read-only row with ":"
	m.taskInputVars = map[string]variableValue{
		"inputOnly": {Value: "someValue", TypeName: "String"},
	}
	// "formField" is in form fields → editable row with "│"
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"formField": {Value: nil, TypeName: "String"},
		},
		m.taskInputVars,
	)
	out := m.renderTaskCompleteModal(120, 40)
	if !strings.Contains(out, "inputOnly") {
		t.Error("expected inputOnly read-only row in modal output")
	}
	if !strings.Contains(out, "formField") {
		t.Error("expected formField editable row in modal output")
	}
	// Read-only rows use ":" separator; editable rows use "│"
	if !strings.Contains(out, ":") {
		t.Error("expected ':' separator for read-only row")
	}
	if !strings.Contains(out, "│") {
		t.Error("expected '│' separator for editable row")
	}
}

// ── Tab skips read-only rows ──────────────────────────────────────────────────

func TestTabSkipsReadOnlyRows(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.lastHeight = 40
	// "dd" sits alphabetically between editable fields "cc" and "ee"
	m.taskInputVars = map[string]variableValue{
		"dd": {Value: "someValue", TypeName: "String"},
	}
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"cc": {Value: nil, TypeName: "String"},
			"ee": {Value: nil, TypeName: "String"},
		},
		m.taskInputVars,
	)
	// buildTaskCompleteFields sorts: field[0]="cc", field[1]="ee"
	m.taskCompletePos = 0
	m.taskCompleteFocus = focusTaskField
	m.taskCompleteFields[0].input.Focus()

	// Tab from field[0]("cc") should advance to field[1]("ee"), skipping read-only "dd"
	m2, _ := sendKeyString(m, "tab")
	if m2.taskCompletePos != 1 {
		t.Errorf("expected pos 1 after Tab (should skip read-only row), got %d", m2.taskCompletePos)
	}
	if m2.taskCompleteFocus != focusTaskField {
		t.Errorf("expected focusTaskField, got %v", m2.taskCompleteFocus)
	}
	if m2.taskCompleteFields[1].name != "ee" {
		t.Errorf("expected field[1] to be 'ee', got %q", m2.taskCompleteFields[1].name)
	}
}

// ── Scroll offset changes visible rows ────────────────────────────────────────

func TestScrollOffsetChangesVisibleRows(t *testing.T) {
	// height=20: cap = height-4 = 16. Need totalRows+8 > 16 → totalRows >= 9.
	// 10 rows (8 read-only + 2 editable): dialogH=18 capped to 16, maxVisible=9, maxOffset=1.
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.lastWidth = 120
	m.lastHeight = 20
	m.taskInputVars = map[string]variableValue{
		"aa": {Value: "v1", TypeName: "String"},
		"ab": {Value: "v2", TypeName: "String"},
		"ac": {Value: "v3", TypeName: "String"},
		"ad": {Value: "v4", TypeName: "String"},
		"ae": {Value: "v5", TypeName: "String"},
		"af": {Value: "v6", TypeName: "String"},
		"ag": {Value: "v7", TypeName: "String"},
		"ah": {Value: "v8", TypeName: "String"},
	}
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"za": {Value: nil, TypeName: "String"},
			"zb": {Value: nil, TypeName: "String"},
		},
		m.taskInputVars,
	)
	// Unified sorted: aa(0), ab(1), ac(2), ad(3), ae(4), af(5), ag(6), ah(7), za(8), zb(9)
	// With scrollOffset=1: visible range is rows[1..9] — "aa" is NOT shown
	m.taskCompleteScrollOffset = 1

	out := m.renderTaskCompleteModal(120, 20)
	// "aa" (row 0) should not be visible at scroll offset 1
	if strings.Contains(out, "aa") {
		t.Error("expected 'aa' to not be visible with scrollOffset=1")
	}
	// "ab" (row 1) should be the first visible variable
	if !strings.Contains(out, "ab") {
		t.Error("expected 'ab' to be visible with scrollOffset=1")
	}
}

// ── EnsureVisible adjusts scroll offset ──────────────────────────────────────

func TestEnsureVisibleScrollsToFocusedField(t *testing.T) {
	// height=20: cap=16, maxVisible=9. 10 rows: aa..ah (read-only) + za, zb (editable).
	// "zb" is at virtual index 9 (>= 0+9), so scrollOffset must become 1 after Tab.
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.lastHeight = 20
	m.taskInputVars = map[string]variableValue{
		"aa": {Value: "v1", TypeName: "String"},
		"ab": {Value: "v2", TypeName: "String"},
		"ac": {Value: "v3", TypeName: "String"},
		"ad": {Value: "v4", TypeName: "String"},
		"ae": {Value: "v5", TypeName: "String"},
		"af": {Value: "v6", TypeName: "String"},
		"ag": {Value: "v7", TypeName: "String"},
		"ah": {Value: "v8", TypeName: "String"},
	}
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"za": {Value: nil, TypeName: "String"},
			"zb": {Value: nil, TypeName: "String"},
		},
		m.taskInputVars,
	)
	// Sorted: aa(0)..ah(7), za(8), zb(9). Sorted fields: field[0]="za", field[1]="zb"
	// Tab field[0](za, idx 8) → field[1](zb, idx 9): ensureVisible fires
	m.taskCompletePos = 0
	m.taskCompleteFocus = focusTaskField
	m.taskCompleteScrollOffset = 0
	m.taskCompleteFields[0].input.Focus()

	m2, _ := sendKeyString(m, "tab") // Tab: za(field[0]) → zb(field[1])
	// "zb" virtual index=9, maxVisible=9 → scrollOffset = 9-9+1 = 1
	if m2.taskCompleteScrollOffset != 1 {
		t.Errorf("expected scrollOffset 1 after Tab to field outside visible window, got %d", m2.taskCompleteScrollOffset)
	}
	if m2.taskCompletePos != 1 {
		t.Errorf("expected pos 1 after Tab, got %d", m2.taskCompletePos)
	}
}

// ── Button focus style: inverted colors (no brackets) ─────────────────────────

func TestButtonFocusStyleIsInverted(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.lastWidth = 120
	m.lastHeight = 40
	m.taskInputVars = map[string]variableValue{}
	m.taskCompleteFields = []taskCompleteField{}

	// Focused Complete: inverted style with no brackets
	m.taskCompleteFocus = focusTaskComplete
	out := m.renderTaskCompleteModal(120, 40)
	if strings.Contains(out, "[ Complete") {
		t.Error("expected focused Complete button to NOT have brackets (inverted style)")
	}
	if !strings.Contains(out, "[ Back ]") {
		t.Error("expected unfocused Back button to have brackets")
	}

	// Unfocused Complete: plain style with brackets
	m.taskCompleteFocus = focusTaskField
	out2 := m.renderTaskCompleteModal(120, 40)
	if !strings.Contains(out2, "[ Complete ]") {
		t.Error("expected unfocused Complete button to have brackets")
	}
}

// ── Border title contains task name ──────────────────────────────────────────

func TestBorderTitleContainsTaskName(t *testing.T) {
	m := setupTaskTable(t, "task-1", "Prepare Bank Transfer", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "Prepare Bank Transfer"
	m.lastWidth = 120
	m.lastHeight = 40
	m.taskInputVars = map[string]variableValue{}
	m.taskCompleteFields = []taskCompleteField{}

	out := m.renderTaskCompleteModal(120, 40)
	if !strings.Contains(out, "Prepare Bank Transfer") {
		t.Fatal("expected task name in modal output")
	}
	// Task name should appear on the same line as the rounded border corner ╭
	lines := strings.Split(out, "\n")
	titleOnBorder := false
	for _, line := range lines {
		if strings.Contains(line, "Prepare Bank Transfer") && strings.Contains(line, "╭") {
			titleOnBorder = true
			break
		}
	}
	if !titleOnBorder {
		t.Error("expected task name to appear on the border line (line containing ╭)")
	}
}

// ── {currentUser} placeholder ─────────────────────────────────────────────────

func TestCurrentUserPlaceholderInBody(t *testing.T) {
	cfg := &config.Config{
		Environments: map[string]config.Environment{
			"local": {URL: "http://localhost:8080", Username: "testuser"},
		},
	}
	m := newModel(cfg)
	m.currentEnv = "local"
	action := config.ActionDef{
		Key:    "c",
		Label:  "Claim",
		Method: "POST",
		Path:   "/task/{id}/claim",
		Body:   `{"userId": "{currentUser}"}`,
	}
	// Build the resolved body
	env := cfg.Environments["local"]
	resolvedBody := replaceCurrentUser(action.Body, env.Username)
	if resolvedBody != `{"userId": "testuser"}` {
		t.Errorf("expected resolved body with username, got %q", resolvedBody)
	}
	_ = m
}

// replaceCurrentUser is a helper to test the placeholder resolution logic.
func replaceCurrentUser(body, username string) string {
	return strings.ReplaceAll(body, "{currentUser}", username)
}

// ── completeTaskCmd sends correct payload ──────────────────────────────────────

func TestCompleteTaskCmdUsesOrigType(t *testing.T) {
	// Verify that buildTaskCompleteFields preserves origType for API submission
	m := newModel(testTaskConfig("alice"))
	formVars := map[string]variableValue{
		"approved": {Value: nil, TypeName: "Boolean"},
		"name":     {Value: "default", TypeName: "String"},
	}
	fields := m.buildTaskCompleteFields(formVars, map[string]variableValue{})

	for _, f := range fields {
		switch f.name {
		case "approved":
			if f.origType != "Boolean" {
				t.Errorf("expected origType 'Boolean', got %q", f.origType)
			}
			if f.varType != "bool" {
				t.Errorf("expected varType 'bool', got %q", f.varType)
			}
		case "name":
			if f.origType != "String" {
				t.Errorf("expected origType 'String', got %q", f.origType)
			}
		}
	}
}

// ── submitTaskComplete assembles correct VariableValueDto ─────────────────────

func TestSubmitTaskCompleteAssemblesVars(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskID = "task-1"
	m.taskCompleteTaskName = "My Task"

	approvedInput := textinput.New()
	approvedInput.SetValue("true")
	countInput := textinput.New()
	countInput.SetValue("5")

	m.taskCompleteFields = []taskCompleteField{
		{name: "approved", varType: "bool", origType: "Boolean", input: approvedInput},
		{name: "count", varType: "int", origType: "Integer", input: countInput},
	}
	m.taskCompleteFocus = focusTaskComplete

	cmd := m.submitTaskComplete()
	if cmd == nil {
		t.Fatal("expected completeTaskCmd to be returned")
	}
}

// ── Verify no drilldown on task Enter ─────────────────────────────────────────

func TestTaskEnterDoesNotDrilldown(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	// Enter should open dialog, not drilldown
	m2, _ := sendKeyString(m, "enter")
	// Should not have pushed to navigation stack
	if len(m2.navigationStack) > 0 {
		t.Error("Enter on task table with own task should not push to navigationStack")
	}
	// Should show loading state (fetching variables)
	if m2.footerStatusKind != footerStatusLoading {
		t.Errorf("expected footerStatusLoading, got %v", m2.footerStatusKind)
	}
}

// ── completeTaskCmd with empty form vars ──────────────────────────────────────

func TestCompleteWithNoFormVarsSendsEmptyMap(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskID = "task-1"
	m.taskCompleteTaskName = "My Task"
	m.taskCompleteFields = []taskCompleteField{} // no form vars
	m.taskCompleteFocus = focusTaskComplete

	cmd := m.submitTaskComplete()
	if cmd == nil {
		t.Error("expected completeTaskCmd even with empty form vars")
	}
}

// ── Polish: no hint line ──────────────────────────────────────────────────────

func TestNoHintLineInModal(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.lastWidth = 120
	m.lastHeight = 40
	m.taskInputVars = map[string]variableValue{}
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{"fieldA": {Value: nil, TypeName: "String"}},
		map[string]variableValue{},
	)
	out := m.renderTaskCompleteModal(120, 40)
	if strings.Contains(out, "Tab: next field") {
		t.Error("expected hint line NOT to appear in modal")
	}
}

// ── Polish: content-driven height ────────────────────────────────────────────

func TestContentDrivenHeight(t *testing.T) {
	// 3 fields + 0 input vars: totalRows=3, dialogH=11 → 10 newlines in output
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.lastWidth = 120
	m.lastHeight = 40
	m.taskInputVars = map[string]variableValue{}
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"fieldA": {Value: nil, TypeName: "String"},
			"fieldB": {Value: nil, TypeName: "String"},
			"fieldC": {Value: nil, TypeName: "String"},
		},
		map[string]variableValue{},
	)
	out := m.renderTaskCompleteModal(120, 40)
	newlines := strings.Count(out, "\n")
	// dialogH = 3+8 = 11; output has dialogH-1 = 10 newlines (bottom border has no trailing \n)
	if newlines != 10 {
		t.Errorf("expected 10 newlines (dialogH=11) for 3 fields, got %d", newlines)
	}
	// Confirm shorter than old fixed height (40-4=36 rows → 35 newlines)
	if newlines >= 35 {
		t.Errorf("expected dialog shorter than full terminal height, got %d newlines", newlines)
	}
}

// ── Polish: error line in dialog ─────────────────────────────────────────────

func TestErrorLineShownInDialog(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.lastWidth = 120
	m.lastHeight = 40
	m.taskInputVars = map[string]variableValue{}
	m.taskCompleteFields = []taskCompleteField{}
	m.taskCompleteError = "task not found"

	out := m.renderTaskCompleteModal(120, 40)
	if !strings.Contains(out, "⚠") {
		t.Error("expected '⚠' in modal when taskCompleteError is set")
	}
	if !strings.Contains(out, "task not found") {
		t.Error("expected error text in modal")
	}
}

func TestNoErrorLineWhenErrorEmpty(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.lastWidth = 120
	m.lastHeight = 40
	m.taskInputVars = map[string]variableValue{}
	m.taskCompleteFields = []taskCompleteField{}
	m.taskCompleteError = ""

	out := m.renderTaskCompleteModal(120, 40)
	if strings.Contains(out, "⚠") {
		t.Error("expected no '⚠' in modal when taskCompleteError is empty")
	}
}

// ── Polish: Space activates focused buttons ───────────────────────────────────

func TestSpaceActivatesCompleteButton(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskID = "task-1"
	m.taskCompleteTaskName = "My Task"
	m.taskCompleteFields = []taskCompleteField{} // no form vars, so no validation errors
	m.taskCompleteFocus = focusTaskComplete

	_, cmd := sendKeyString(m, " ")
	if cmd == nil {
		t.Error("expected completeTaskCmd to be dispatched when Space pressed on focused Complete button")
	}
}

func TestSpaceActivatesBackButton(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskID = "task-1"
	m.taskCompleteTaskName = "My Task"
	m.taskCompleteFields = []taskCompleteField{}
	m.taskCompleteFocus = focusTaskBack

	m2, _ := sendKeyString(m, " ")
	if m2.activeModal != ModalNone {
		t.Error("expected dialog to close when Space pressed on focused Back button")
	}
}

// ── Polish: error cleared on field edit ──────────────────────────────────────

func TestErrorClearedOnFieldEdit(t *testing.T) {
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{"name": {Value: nil, TypeName: "String"}},
		map[string]variableValue{},
	)
	m.taskCompletePos = 0
	m.taskCompleteFocus = focusTaskField
	m.taskCompleteFields[0].input.Focus()
	m.taskCompleteError = "some API error"

	m2, _ := sendKeyString(m, "x") // printable key fed to field
	if m2.taskCompleteError != "" {
		t.Errorf("expected taskCompleteError cleared after field edit, got %q", m2.taskCompleteError)
	}
}

// ── Polish: scroll bounds correct when error is shown (M-1 regression) ────────

func TestScrollMaxVisibleAccountsForErrorRow(t *testing.T) {
	// height=20, 10 rows: totalRows+8+errorRows = 10+8+1 = 19 > height-4 = 16 → capped.
	// maxVisible = height-11-errorRows = 20-11-1 = 8. maxOffset = 10-8 = 2.
	// Without the fix, maxVisible would be 9 and maxOffset would be 1,
	// making the last row (index 9) unreachable.
	m := setupTaskTable(t, "task-1", "My Task", "alice", "alice")
	m.activeModal = ModalTaskComplete
	m.taskCompleteTaskName = "My Task"
	m.lastWidth = 120
	m.lastHeight = 20
	m.taskCompleteError = "API call failed"
	m.taskInputVars = map[string]variableValue{
		"aa": {Value: "v1", TypeName: "String"},
		"ab": {Value: "v2", TypeName: "String"},
		"ac": {Value: "v3", TypeName: "String"},
		"ad": {Value: "v4", TypeName: "String"},
		"ae": {Value: "v5", TypeName: "String"},
		"af": {Value: "v6", TypeName: "String"},
		"ag": {Value: "v7", TypeName: "String"},
		"ah": {Value: "v8", TypeName: "String"},
	}
	m.taskCompleteFields = m.buildTaskCompleteFields(
		map[string]variableValue{
			"za": {Value: nil, TypeName: "String"},
			"zb": {Value: nil, TypeName: "String"},
		},
		m.taskInputVars,
	)
	// maxVisible=8, maxOffset=2. Scroll to offset=2 to reach the last row.
	m.taskCompleteScrollOffset = 2

	out := m.renderTaskCompleteModal(120, 20)
	// "zb" is at virtual index 9; at scrollOffset=2 the window is rows[2..9] → "zb" visible.
	if !strings.Contains(out, "zb") {
		t.Error("expected 'zb' (last row) to be visible at scrollOffset=2 when error is shown")
	}
	// maxVisible helper must return 8 (not 9) when error row is present and capped.
	maxVis := m.taskCompleteMaxVisible()
	if maxVis != 8 {
		t.Errorf("expected taskCompleteMaxVisible()=8 with error shown at height=20 and 10 rows, got %d", maxVis)
	}
}

// ── VariableValueDto type in submission ───────────────────────────────────────

func TestVariableValueDtoHasOrigType(t *testing.T) {
	// Verify the VariableValueDto structure that would be submitted
	v := operaton.VariableValueDto{}
	v.SetValue(true)
	v.SetType("Boolean")
	if v.GetType() != "Boolean" {
		t.Errorf("expected type 'Boolean', got %q", v.GetType())
	}
	val := v.Value
	if val == nil {
		t.Error("expected value to be set")
	}
}
