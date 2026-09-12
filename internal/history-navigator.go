package internal

// HistoryNavigator lets an InputBox-driven prompt (search, filter, ...) cycle
// up and down through a shared SearchHistory without losing whatever text the
// user is in the middle of typing.
type HistoryNavigator struct {
	history       *SearchHistory
	index         int
	lastTypedText string
}

// NewHistoryNavigator returns a HistoryNavigator positioned past the end of
// history's entries, i.e. it will show whatever the user types rather than a
// history entry until Move is called.
func NewHistoryNavigator(history *SearchHistory) *HistoryNavigator {
	return &HistoryNavigator{
		history: history,
		index:   len(history.entries),
	}
}

// TextEdited must be called whenever the user changes the input box text by
// hand, so that arrowing away and back finds that text again.
func (n *HistoryNavigator) TextEdited(text string) {
	n.index = len(n.history.entries)
	n.lastTypedText = text
}

// Move steps the history index by delta (negative towards older entries,
// positive towards newer), clamping at both ends. It returns the text that
// should now be shown, and whether there was anything to move to at all (an
// empty history has nothing to navigate).
func (n *HistoryNavigator) Move(delta int) (string, bool) {
	if len(n.history.entries) == 0 {
		return "", false
	}

	n.index += delta
	if n.index < 0 {
		n.index = 0
	}
	if n.index > len(n.history.entries) {
		n.index = len(n.history.entries)
	}

	if n.index == len(n.history.entries) {
		return n.lastTypedText, true
	}

	return n.history.entries[n.index], true
}

// Commit adds text as a new history entry. Call this when the user submits
// or aborts the prompt.
func (n *HistoryNavigator) Commit(text string) {
	n.history.addEntry(text)
}
