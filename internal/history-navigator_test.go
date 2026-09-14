package internal

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestHistoryNavigator_MoveCyclesThroughEntriesNewestFirst(t *testing.T) {
	history := &SearchHistory{entries: []string{"a", "b", "c"}}
	navigator := NewHistoryNavigator(history)

	text, ok := navigator.Move(-1)
	assert.Assert(t, ok)
	assert.Equal(t, "c", text)

	text, ok = navigator.Move(-1)
	assert.Assert(t, ok)
	assert.Equal(t, "b", text)

	text, ok = navigator.Move(-1)
	assert.Assert(t, ok)
	assert.Equal(t, "a", text)
}

func TestHistoryNavigator_MoveClampsAtOldestEntry(t *testing.T) {
	history := &SearchHistory{entries: []string{"a", "b"}}
	navigator := NewHistoryNavigator(history)

	navigator.Move(-1)
	navigator.Move(-1)
	text, ok := navigator.Move(-1)

	assert.Assert(t, ok)
	assert.Equal(t, "a", text)
}

func TestHistoryNavigator_MoveClampsAtNewestEntry(t *testing.T) {
	history := &SearchHistory{entries: []string{"a", "b"}}
	navigator := NewHistoryNavigator(history)

	navigator.Move(-1)
	navigator.Move(1)
	text, ok := navigator.Move(1)

	assert.Assert(t, ok)
	assert.Equal(t, "", text)
}

func TestHistoryNavigator_MoveOnEmptyHistoryDoesNothing(t *testing.T) {
	history := &SearchHistory{}
	navigator := NewHistoryNavigator(history)

	_, ok := navigator.Move(-1)

	assert.Assert(t, !ok)
}

func TestHistoryNavigator_TypedTextSurvivesUpThenDown(t *testing.T) {
	history := &SearchHistory{entries: []string{"a", "b"}}
	navigator := NewHistoryNavigator(history)

	navigator.TextEdited("typing this")

	text, ok := navigator.Move(-1)
	assert.Assert(t, ok)
	assert.Equal(t, "b", text)

	text, ok = navigator.Move(1)
	assert.Assert(t, ok)
	assert.Equal(t, "typing this", text)
}

func TestHistoryNavigator_MoveStepsThroughEveryEntryBothWays(t *testing.T) {
	history := &SearchHistory{entries: []string{"a", "b", "c"}}
	navigator := NewHistoryNavigator(history)

	text, ok := navigator.Move(-1)
	assert.Assert(t, ok)
	assert.Equal(t, "c", text)

	text, ok = navigator.Move(-1)
	assert.Assert(t, ok)
	assert.Equal(t, "b", text)

	text, ok = navigator.Move(-1)
	assert.Assert(t, ok)
	assert.Equal(t, "a", text)

	text, ok = navigator.Move(1)
	assert.Assert(t, ok)
	assert.Equal(t, "b", text)

	text, ok = navigator.Move(1)
	assert.Assert(t, ok)
	assert.Equal(t, "c", text)

	text, ok = navigator.Move(1)
	assert.Assert(t, ok)
	assert.Equal(t, "", text)
}

func TestHistoryNavigator_CommitAddsEntryToHistory(t *testing.T) {
	history := &SearchHistory{entries: []string{"a"}}
	navigator := NewHistoryNavigator(history)

	navigator.Commit("b")

	assert.DeepEqual(t, []string{"a", "b"}, history.entries)
}
