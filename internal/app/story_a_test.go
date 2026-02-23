package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// ── T1: Skin popup uses generic offset/cursor scroll ─────────────────────────

func TestSkinPopupDownPastMaxShowScrollsOffset(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	// Provide 12 skins so we exceed maxShow=8
	m.availableSkins = []string{"s0", "s1", "s2", "s3", "s4", "s5", "s6", "s7", "s8", "s9", "s10", "s11"}
	m.popup.mode = popupModeSkin
	m.popup.cursor = 7 // at last visible item (0-indexed)
	m.popup.offset = 0

	m2, _ := sendKeyString(m, "down")

	if m2.popup.cursor != 8 {
		t.Errorf("expected cursor=8 after down from 7, got %d", m2.popup.cursor)
	}
	if m2.popup.offset == 0 {
		t.Errorf("expected popup.offset > 0 after scrolling past maxShow, got %d", m2.popup.offset)
	}
}

func TestSkinPopupUpAtOffsetBoundaryScrollsOffset(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.availableSkins = []string{"s0", "s1", "s2", "s3", "s4", "s5", "s6", "s7", "s8", "s9"}
	m.popup.mode = popupModeSkin
	m.popup.cursor = 3 // cursor == offset → at top of visible window
	m.popup.offset = 3

	m2, _ := sendKeyString(m, "up")

	if m2.popup.cursor != 2 {
		t.Errorf("expected cursor=2 after up, got %d", m2.popup.cursor)
	}
	if m2.popup.offset != 2 {
		t.Errorf("expected offset=2 after scrolling up, got %d", m2.popup.offset)
	}
}

func TestSkinPopupOffsetResetOnOpen(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.availableSkins = []string{"s0", "s1"}
	m.popup.offset = 5 // stale from previous session

	m2, _ := sendKeyString(m, "ctrl+t")

	if m2.popup.offset != 0 {
		t.Errorf("expected popup.offset=0 on skin popup open, got %d", m2.popup.offset)
	}
}

// ── T2: Pane height restores when search popup closes ────────────────────────

func TestSearchPopupEscRestoresPaneHeight(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastHeight = 40
	m.lastWidth = 120

	// Open search popup
	m2, _ := sendKeyString(m, "/")
	if m2.popup.mode != popupModeSearch {
		t.Fatal("expected popup.mode == popupModeSearch")
	}
	reducedHeight := m2.paneHeight

	// Close with Esc
	m3, _ := sendKeyString(m2, "esc")

	if m3.popup.mode != popupModeNone {
		t.Fatal("expected popup closed after Esc")
	}
	// pane height should have been recomputed (should be >= reduced height)
	if m3.paneHeight <= reducedHeight {
		t.Errorf("expected paneHeight to grow after closing search popup; was %d, now %d", reducedHeight, m3.paneHeight)
	}
	// verify it matches what computePaneHeight would give with no popup
	expected := m3.computePaneHeight()
	if m3.paneHeight != expected {
		t.Errorf("expected paneHeight=%d (from computePaneHeight), got %d", expected, m3.paneHeight)
	}
}

func TestSearchPopupEnterRestoresPaneHeight(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastHeight = 40
	m.lastWidth = 120

	m2, _ := sendKeyString(m, "/")
	reducedHeight := m2.paneHeight

	// Close with Enter (lock filter)
	m3, _ := sendKeyString(m2, "enter")

	if m3.popup.mode != popupModeNone {
		t.Fatal("expected popup closed after Enter")
	}
	if m3.paneHeight <= reducedHeight {
		t.Errorf("expected paneHeight to grow after Enter lock; was %d, now %d", reducedHeight, m3.paneHeight)
	}
	expected := m3.computePaneHeight()
	if m3.paneHeight != expected {
		t.Errorf("expected paneHeight=%d, got %d", expected, m3.paneHeight)
	}
}

// ── T3: Quit dialog ghost fix ─────────────────────────────────────────────────

func TestQuitConfirmSetsQuittingFlag(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.activeModal = ModalConfirmQuit
	m.confirmFocusedBtn = 0 // confirm button focused

	m2, cmd := sendKeyString(m, "enter")

	if !m2.quitting {
		t.Error("expected m.quitting=true after confirming quit")
	}
	if cmd == nil {
		t.Error("expected tea.Quit command returned")
	}
}

func TestQuitCtrlCConfirmSetsQuittingFlag(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.activeModal = ModalConfirmQuit

	m2, _ := sendKeyString(m, "ctrl+c")

	if !m2.quitting {
		t.Error("expected m.quitting=true after ctrl+c on quit dialog")
	}
}

func TestViewReturnsEmptyWhenQuitting(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false
	m.lastWidth = 120
	m.lastHeight = 40
	m.quitting = true

	out := m.View()

	if out != "" {
		t.Errorf("expected empty View() when quitting, got %q (len %d)", out[:min(len(out), 40)], len(out))
	}
}

// ── T4: Full-width selected row highlight ─────────────────────────────────────

func TestPaneWidthSetOnWindowResize(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false

	resizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	m2raw, _ := m.Update(resizeMsg)
	m2 := m2raw.(model)

	if m2.paneWidth <= 0 {
		t.Errorf("expected paneWidth > 0 after resize, got %d", m2.paneWidth)
	}
}

func TestPaneWidthUpdatesOnSubsequentResize(t *testing.T) {
	m := newTestModel(t)
	m.splashActive = false

	m1raw, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m1 := m1raw.(model)
	w1 := m1.paneWidth

	m2raw, _ := m1.Update(tea.WindowSizeMsg{Width: 150, Height: 30})
	m2 := m2raw.(model)
	w2 := m2.paneWidth

	if w1 == w2 {
		t.Errorf("expected paneWidth to change on resize: w1=%d w2=%d", w1, w2)
	}
	if w2 <= w1 {
		t.Errorf("expected paneWidth to increase from %d to some value > that for width 150, got %d", w1, w2)
	}
}
