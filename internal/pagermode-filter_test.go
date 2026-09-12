package internal

import (
	"testing"

	"github.com/walles/twin"
	"gotest.tools/v3/assert"
)

// Ref: https://github.com/walles/moor/issues/466
func TestPagerModeFilter_UpArrowShowsMostRecentSearchHistoryEntry(t *testing.T) {
	pager := createThreeLinesPager(t)
	pager.searchHistory.entries = []string{"b", "d"}

	filterMode := NewPagerModeFilter(pager, pager.scrollPosition)
	filterMode.onKey(twin.KeyUp)

	assert.Equal(t, "d", filterMode.inputBox.text)
}

func TestPagerModeFilter_DownArrowStepsBackTowardsNewestEntry(t *testing.T) {
	pager := createThreeLinesPager(t)
	pager.searchHistory.entries = []string{"b", "d"}

	filterMode := NewPagerModeFilter(pager, pager.scrollPosition)
	filterMode.onKey(twin.KeyUp)
	filterMode.onKey(twin.KeyUp)
	filterMode.onKey(twin.KeyDown)

	assert.Equal(t, "d", filterMode.inputBox.text)
}

func TestPagerModeFilter_EnterPersistsTypedTextToSearchHistory(t *testing.T) {
	pager := createThreeLinesPager(t)
	pager.searchHistory = &SearchHistory{} // No file backing, keep this test disk-free

	filterMode := NewPagerModeFilter(pager, pager.scrollPosition)
	filterMode.inputBox.setText("abc")
	filterMode.onKey(twin.KeyEnter)

	assert.DeepEqual(t, []string{"abc"}, pager.searchHistory.entries)
}

// Search and filter modes each get their own HistoryNavigator, but both wrap
// the same *SearchHistory on the pager, so an entry committed from one mode
// must be visible when navigating history from the other.
//
// Ref: https://github.com/walles/moor/issues/466
func TestPagerModeFilter_SharesSearchHistoryWithSearchMode(t *testing.T) {
	pager := createThreeLinesPager(t)
	pager.searchHistory = &SearchHistory{} // No file backing, keep this test disk-free

	filterMode := NewPagerModeFilter(pager, pager.scrollPosition)
	filterMode.inputBox.setText("from filter")
	filterMode.onKey(twin.KeyEnter)

	searchMode := NewPagerModeSearch(pager, SearchDirectionForward, pager.scrollPosition)
	searchMode.onKey(twin.KeyUp)

	assert.Equal(t, "from filter", searchMode.inputBox.text)
}
