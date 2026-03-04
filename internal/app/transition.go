package app

import "github.com/charmbracelet/bubbles/table"

// TransitionType identifies the category of navigation state transition.
// Every navigation path in internal/app/ MUST call prepareStateTransition before
// modifying view state. Any navigation code that does not call prepareStateTransition is a bug.
type TransitionType int

const (
	// TransitionFull performs a complete state reset: clears activeModal, footerError,
	// search, sort, cursor, navigationStack, and identity params.
	// Use for: environment switch, context switch, breadcrumb jump to root.
	TransitionFull TransitionType = iota

	// TransitionDrillDown captures the current viewState snapshot onto the navigationStack
	// (push-before-clear), then clears non-stack view state for the incoming child view.
	// Use for: drill-down navigation (Enter/→).
	TransitionDrillDown

	// TransitionPop pops the top viewState from navigationStack and restores all captured
	// fields (viewMode, breadcrumb, contentHeader, selectedKeys, tableRows, tableColumns,
	// tableCursor, genericParams, rowData). Performs no clearing.
	// Use for: Esc (back) and breadcrumb jump to non-root level (after caller truncates stack).
	TransitionPop
)

// prepareStateTransition is the single mandatory gate for all navigation changes.
// Call it before modifying any view state in a navigation handler.
func (m *model) prepareStateTransition(t TransitionType) {
	switch t {
	case TransitionFull:
		// Clear all view state for a fresh navigation context.
		m.activeModal = ModalNone
		m.footerError = ""
		m.footerStatusKind = footerStatusNone
		m.sortColumn = -1
		m.sortAscending = true
		m.searchTerm = ""
		m.searchMode = false
		m.searchInput.Blur()
		m.originalRows = nil
		m.filteredRows = nil
		m.navigationStack = nil
		m.genericParams = make(map[string]string)
		m.selectedDefinitionKey = ""
		m.selectedInstanceID = ""
		m.table.SetCursor(0)
		if m.popup.mode != popupModeNone {
			m.popup.mode = popupModeNone
			m.popup.input = ""
			m.popup.cursor = -1
			m.popup.offset = 0
		}

	case TransitionDrillDown:
		// Push current viewState BEFORE clearing (push-before-clear ordering ensures
		// the parent's cursor and column state are preserved in the snapshot).
		cols := m.table.Columns()
		var rows []table.Row
		if len(cols) > 0 {
			rows = normalizeRows(append([]table.Row{}, m.table.Rows()...), len(cols))
		} else {
			rows = append([]table.Row{}, m.table.Rows()...)
		}
		snapshot := viewState{
			viewMode:              m.viewMode,
			breadcrumb:            append([]string{}, m.breadcrumb...),
			contentHeader:         m.contentHeader,
			selectedDefinitionKey: m.selectedDefinitionKey,
			selectedInstanceID:    m.selectedInstanceID,
			tableRows:             rows,
			tableCursor:           m.table.Cursor(),
			cachedDefinitions:     m.cachedDefinitions,
			tableColumns:          append([]table.Column{}, cols...),
			genericParams:         m.genericParams,
			rowData:               append([]map[string]interface{}{}, m.rowData...),
		}
		m.navigationStack = append(m.navigationStack, snapshot)
		// Clear non-stack fields for the incoming child view.
		m.activeModal = ModalNone
		m.footerError = ""
		m.footerStatusKind = footerStatusNone
		m.sortColumn = -1
		m.sortAscending = true
		m.searchTerm = ""
		m.searchMode = false
		m.searchInput.Blur()
		m.originalRows = nil
		m.filteredRows = nil
		if m.popup.mode == popupModeSearch {
			m.popup.mode = popupModeNone
			m.popup.input = ""
			m.popup.cursor = -1
			m.popup.offset = 0
		}

	case TransitionPop:
		if len(m.navigationStack) == 0 {
			return
		}
		top := m.navigationStack[len(m.navigationStack)-1]
		m.navigationStack = m.navigationStack[:len(m.navigationStack)-1]
		// Restore all viewState fields — no clearing.
		m.viewMode = top.viewMode
		m.breadcrumb = top.breadcrumb
		m.contentHeader = top.contentHeader
		m.selectedDefinitionKey = top.selectedDefinitionKey
		m.selectedInstanceID = top.selectedInstanceID
		m.cachedDefinitions = top.cachedDefinitions
		m.genericParams = top.genericParams
		m.rowData = top.rowData
		// Restore table widget state: columns first, then rows, then cursor.
		if len(top.tableColumns) > 0 {
			m.table.SetRows(normalizeRows(nil, len(top.tableColumns)))
			m.table.SetColumns(top.tableColumns)
		}
		cols := m.table.Columns()
		if len(cols) > 0 {
			m.table.SetRows(normalizeRows(top.tableRows, len(cols)))
		} else {
			m.table.SetRows(top.tableRows)
		}
		m.table.SetCursor(top.tableCursor)
	}
}

// clampCursorAfterRowRemoval ensures the table cursor stays within bounds after
// rows are removed (e.g. after terminate, delete). Safe to call always.
func (m *model) clampCursorAfterRowRemoval() {
	rows := m.table.Rows()
	cursor := m.table.Cursor()
	if len(rows) == 0 {
		m.table.SetCursor(0)
		return
	}
	if cursor >= len(rows) {
		m.table.SetCursor(len(rows) - 1)
	}
}

// clearSearch resets client-side search state without touching sort or nav.
// Used when search is explicitly cancelled (Esc in search popup).
func (m *model) clearSearch() {
	m.searchTerm = ""
	if m.originalRows != nil {
		m.table.SetRows(m.originalRows)
	}
	m.originalRows = nil
	m.filteredRows = nil
}

// clearSort resets sort state and removes sort indicators from column headers.
func (m *model) clearSort() {
	m.sortColumn = -1
	m.sortAscending = true
	// rebuild columns without sort indicators
	cols := m.table.Columns()
	cleaned := make([]table.Column, len(cols))
	for i, c := range cols {
		title := c.Title
		if len(title) > 2 && (title[len(title)-2:] == " ^" || title[len(title)-2:] == " v") {
			title = title[:len(title)-2]
		}
		cleaned[i] = table.Column{Title: title, Width: c.Width}
	}
	m.table.SetColumns(cleaned)
}
