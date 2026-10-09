package internal

import (
	"os"
	"testing"

	"github.com/walles/moor/v2/internal/linemetadata"
	"github.com/walles/moor/v2/internal/reader"
	"github.com/walles/twin"
	"gotest.tools/v3/assert"
)

func TestErrUnlessExecutable_yes(t *testing.T) {
	// Find our own executable
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	// Check that it's executable
	err = errUnlessExecutable(executable)
	if err != nil {
		t.Fatal(err)
	}
}

func TestErrUnlessExecutable_no(t *testing.T) {
	textFile := "pagermode-viewing_test.go"
	if _, err := os.Stat(textFile); os.IsNotExist(err) {
		t.Fatal("Test setup failed, text file not found: " + textFile)
	}

	err := errUnlessExecutable(textFile)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// Counts how many times lines are fetched from the wrapped reader
type lineReadCountingReader struct {
	reader.Reader
	lineReads int
}

func (r *lineReadCountingReader) GetLine(index linemetadata.Index) *reader.NumberedLine {
	r.lineReads++
	return r.Reader.GetLine(index)
}

func (r *lineReadCountingReader) GetLines(firstLine linemetadata.Index, wantedLineCount int) reader.InputLines {
	r.lineReads++
	return r.Reader.GetLines(firstLine, wantedLineCount)
}

func (r *lineReadCountingReader) GetLinesPreallocated(firstLine linemetadata.Index, resultLines *[]reader.NumberedLine) reader.Status {
	r.lineReads++
	return r.Reader.GetLinesPreallocated(firstLine, resultLines)
}

// Inertial scrolling can queue hundreds of scroll-down events, and the screen
// is redrawn only after all of them have been handled. So handling one must not
// touch the input lines, all such work waits for the next redraw.
func TestViewingScrollDownIsCheap(t *testing.T) {
	r := reader.NewFromTextForTesting("test", "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n")

	pager := NewPager(r)
	pager.screen = twin.NewFakeScreen(20, 5)
	countingReader := &lineReadCountingReader{Reader: r}
	pager.filteringReader.BackingReader = countingReader

	mode := PagerModeViewing{pager: pager}
	mode.onKey(twin.KeyDown)

	assert.Equal(t, 0, countingReader.lineReads)
}

// Scrolling down to the last line starts following the end of the input.
func TestViewingScrollDownToEndFollows(t *testing.T) {
	r := reader.NewFromTextForTesting("test", "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n")

	pager := NewPager(r)
	pager.screen = twin.NewFakeScreen(20, 5)
	assert.Assert(t, pager.TargetLine == nil)

	// One line per press is guaranteed to get us to the end
	mode := PagerModeViewing{pager: pager}
	for range r.GetLineCount() {
		mode.onKey(twin.KeyDown)
	}
	pager.redraw("")

	assert.Assert(t, pager.TargetLine != nil)
	assert.Equal(t, linemetadata.IndexMax(), *pager.TargetLine)
}

// Choosing a target line after scrolling down to the end keeps that target.
func TestViewingScrollDownToEndThenTargetLine(t *testing.T) {
	r := reader.NewFromTextForTesting("test", "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n")

	pager := NewPager(r)
	pager.screen = twin.NewFakeScreen(20, 5)

	// One line per press is guaranteed to get us to the end
	mode := PagerModeViewing{pager: pager}
	for range r.GetLineCount() {
		mode.onKey(twin.KeyDown)
	}
	target := linemetadata.IndexFromOneBased(2)
	pager.setTargetLine(&target)
	pager.redraw("")

	assert.Assert(t, pager.TargetLine != nil)
	assert.Equal(t, target, *pager.TargetLine)
}

func TestViewingFooter_WithSpinner(t *testing.T) {
	r := reader.NewFromTextForTesting("", "text")

	pager := NewPager(r)
	pager.ShowStatusBar = true

	// Attach a fake screen large enough to render the footer
	screen := twin.NewFakeScreen(80, 10)
	pager.screen = screen

	// Drive the footer rendering directly via PagerModeViewing
	mode := PagerModeViewing{pager: pager}
	spinner := "<->"
	mode.drawFooter("", "1 line  100%", spinner)

	footer := rowToString(screen.GetRow(9))

	// Quotes are stripped in rendering; expect plain keys
	expectedHelp := "Press ESC / q to exit, / to search, & to filter, h for help"
	expected := "1 line  100%  " + spinner + "  " + expectedHelp

	assert.Equal(t, expected, footer)
}
