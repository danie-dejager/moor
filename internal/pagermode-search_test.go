package internal

import (
	"testing"

	"github.com/walles/twin"
	"gotest.tools/v3/assert"
)

// Ref: https://github.com/walles/moor/issues/466
func TestPagerModeSearch_PgDownPersistsTypedTextAndSwitchesToViewing(t *testing.T) {
	pager := createThreeLinesPager(t)
	pager.searchHistory = &SearchHistory{} // No file backing, keep this test disk-free

	searchMode := NewPagerModeSearch(pager, SearchDirectionForward, pager.scrollPosition)
	pager.mode = searchMode
	searchMode.inputBox.setText("abc")
	searchMode.onKey(twin.KeyPgDown)

	assert.DeepEqual(t, []string{"abc"}, pager.searchHistory.entries)
	assert.Equal(t, "Viewing", modeName(pager))
}
