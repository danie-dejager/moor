package internal

import (
	log "github.com/sirupsen/logrus"
	"github.com/walles/twin"
)

type SearchDirection bool

const (
	SearchDirectionForward  SearchDirection = false
	SearchDirectionBackward SearchDirection = true
)

type PagerModeSearch struct {
	pager                 *Pager
	initialScrollPosition scrollPosition // Pager position before search started
	direction             SearchDirection
	inputBox              *InputBox
	history               *HistoryNavigator
}

func NewPagerModeSearch(p *Pager, direction SearchDirection, initialScrollPosition scrollPosition) *PagerModeSearch {
	m := &PagerModeSearch{
		pager:                 p,
		initialScrollPosition: initialScrollPosition,
		direction:             direction,
		history:               NewHistoryNavigator(p.searchHistory),
	}
	m.inputBox = &InputBox{
		accept: INPUTBOX_ACCEPT_ALL,
		onTextChanged: func(text string) {
			m.pager.search.For(text)

			switch m.direction {
			case SearchDirectionBackward:
				m.pager.scrollToSearchHitsBackwards()
			case SearchDirectionForward:
				m.pager.scrollToSearchHits()
			}
		},
	}
	return m
}

func (m PagerModeSearch) drawFooter(_ string, _ string, _ string) {
	prompt := "Search: "
	if m.direction == SearchDirectionBackward {
		prompt = "Search backwards: "
	}
	m.inputBox.draw(m.pager.screen, "Type to search, 'ENTER' submits, 'ESC' cancels, '↑↓' navigate history", prompt)
}

// Exit search mode, skip back to where we started
func (m *PagerModeSearch) abort() {
	m.history.Commit(m.inputBox.text)
	m.pager.mode = PagerModeViewing{pager: m.pager}
	m.pager.scrollPosition = m.initialScrollPosition
	m.pager.setTargetLine(nil) // Viewing doesn't need all lines
}

func (m *PagerModeSearch) onKey(key twin.KeyCode) {
	if m.inputBox.handleKey(key) {
		m.history.TextEdited(m.inputBox.text)
		return
	}

	switch key {
	case twin.KeyEnter:
		m.history.Commit(m.inputBox.text)
		m.pager.mode = PagerModeViewing{pager: m.pager}
		m.pager.setTargetLine(nil) // Viewing doesn't need all lines

	case twin.KeyEscape:
		m.abort()

	case twin.KeyPgUp, twin.KeyPgDown:
		m.history.Commit(m.inputBox.text)
		m.pager.mode = PagerModeViewing{pager: m.pager}
		m.pager.mode.onKey(key)
		m.pager.setTargetLine(nil) // Viewing doesn't need all lines

	case twin.KeyUp:
		if text, ok := m.history.Move(-1); ok {
			m.inputBox.setText(text)
		}

	case twin.KeyDown:
		if text, ok := m.history.Move(1); ok {
			m.inputBox.setText(text)
		}

	default:
		log.Debugf("Unhandled search key event %v", key)
	}
}

func (m *PagerModeSearch) onRune(char rune) {
	// Handle ctrl-c
	if char == '\x03' {
		m.abort()
		return
	}

	m.inputBox.handleRune(char)
	m.history.TextEdited(m.inputBox.text)
}
