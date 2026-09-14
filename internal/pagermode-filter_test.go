package internal

import (
	"testing"

	"github.com/walles/moor/v2/internal/linemetadata"
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

// Pressing ESC should return the view to wherever it was scrolled to before
// filtering started, the same way search's ESC does, even though filtering
// out non-matching lines may have moved the pager's scroll position while
// typing.
//
// Ref: https://github.com/walles/moor/issues/466
func TestPagerModeFilter_EscapeRestoresScrollPositionFromBeforeFiltering(t *testing.T) {
	pager := createThreeLinesPager(t)
	pager.scrollPosition = NewScrollPositionFromIndex(linemetadata.IndexFromZeroBased(4), "before filtering")
	initialIndex := pager.lineIndex().Index()

	filterMode := NewPagerModeFilter(pager, pager.scrollPosition)
	filterMode.onRune('a') // Filters down to only the line containing "a"

	// This mirrors what a real redraw does, and confirms filtering really did
	// move the scroll position as a side effect.
	assert.Assert(t, pager.lineIndex().Index() != initialIndex, "Filtering should have moved the scroll position")

	filterMode.onKey(twin.KeyEscape)

	assert.Equal(t, initialIndex, pager.lineIndex().Index())
}

// Ref: https://github.com/walles/moor/issues/466
func TestPagerModeFilter_PgDownPersistsTypedTextAndSwitchesToViewing(t *testing.T) {
	pager := createThreeLinesPager(t)
	pager.searchHistory = &SearchHistory{} // No file backing, keep this test disk-free

	filterMode := NewPagerModeFilter(pager, pager.scrollPosition)
	pager.mode = filterMode
	filterMode.inputBox.setText("abc")
	filterMode.onKey(twin.KeyPgDown)

	assert.DeepEqual(t, []string{"abc"}, pager.searchHistory.entries)
	assert.Equal(t, "Viewing", modeName(pager))
}
