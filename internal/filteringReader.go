package internal

import (
	"fmt"
	"math"
	"regexp"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/walles/moor/v2/internal/linemetadata"
	"github.com/walles/moor/v2/internal/reader"
)

// Filters lines based on the search query from the pager.

type FilteringReader struct {
	BackingReader reader.Reader

	// This is a reference to a reference so that we can track changes to the
	// original pattern, including if it is set to nil.
	FilterPattern **regexp.Regexp

	// Protects filteredLinesCache, unfilteredLineCountWhenCaching, and
	// filterPatternWhenCaching.
	lock sync.Mutex

	// nil means no filtering has happened yet
	filteredLinesCache *[]*reader.NumberedLine

	// This is what the reader's line count was when we filtered. If the
	// reader's current line count doesn't match, then our cache needs to be
	// rebuilt.
	unfilteredLineCountWhenCaching int

	// This is the pattern that was used when we cached the lines. If it
	// doesn't match the current pattern, then our cache needs to be rebuilt.
	filterPatternWhenCaching *regexp.Regexp
}

// Please hold the lock when calling this method.
func (f *FilteringReader) rebuildCache() {
	t0 := time.Now()

	cache := make([]*reader.NumberedLine, 0)
	filterPattern := *f.FilterPattern

	// Mark cache base conditions
	f.unfilteredLineCountWhenCaching = f.BackingReader.GetLineCount()
	f.filterPatternWhenCaching = filterPattern

	// Repopulate the cache
	allBaseLines := f.BackingReader.GetLines(linemetadata.Index{}, math.MaxInt)
	resultIndex := 0
	for _, line := range allBaseLines.Lines {
		if filterPattern != nil && len(filterPattern.String()) > 0 && !filterPattern.MatchString(line.Line.Plain(&line.Index)) {
			// We have a pattern but it doesn't match
			continue
		}

		cache = append(cache, &reader.NumberedLine{
			Line:   line.Line,
			Index:  linemetadata.IndexFromZeroBased(resultIndex),
			Number: line.Number,
		})
		resultIndex++
	}

	f.filteredLinesCache = &cache

	log.Debugf("Filtered out %d/%d lines in %s",
		len(allBaseLines.Lines)-len(cache), len(allBaseLines.Lines), time.Since(t0))
}

func (f *FilteringReader) getAllLines() []*reader.NumberedLine {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.filteredLinesCache == nil {
		f.rebuildCache()
		return *f.filteredLinesCache
	}

	if f.unfilteredLineCountWhenCaching != f.BackingReader.GetLineCount() {
		f.rebuildCache()
		return *f.filteredLinesCache
	}

	var currentFilterPattern string
	if *f.FilterPattern != nil {
		currentFilterPattern = (*f.FilterPattern).String()
	}
	var cacheFilterPattern string
	if f.filterPatternWhenCaching != nil {
		cacheFilterPattern = f.filterPatternWhenCaching.String()
	}
	if currentFilterPattern != cacheFilterPattern {
		f.rebuildCache()
		return *f.filteredLinesCache
	}

	return *f.filteredLinesCache
}

func (f *FilteringReader) shouldPassThrough() bool {
	f.lock.Lock()
	defer f.lock.Unlock()

	if *f.FilterPattern == nil || len((*f.FilterPattern).String()) == 0 {
		// Cache is not needed
		f.filteredLinesCache = nil

		// No filtering, so pass through all
		return true
	}

	return false
}

func (f *FilteringReader) GetLineCount() int {
	if f.shouldPassThrough() {
		return f.BackingReader.GetLineCount()
	}

	return len(f.getAllLines())
}

func (f *FilteringReader) ShouldShowLineCount() bool {
	panic("Unexpected call to FilteringReader.ShouldShowLineCount()")
}

func (f *FilteringReader) GetLine(index linemetadata.Index) *reader.NumberedLine {
	if f.shouldPassThrough() {
		return f.BackingReader.GetLine(index)
	}

	allLines := f.getAllLines()
	if index.Index() < 0 || index.Index() >= len(allLines) {
		return nil
	}
	return allLines[index.Index()]
}

func (f *FilteringReader) GetLines(firstLine linemetadata.Index, wantedLineCount int) *reader.InputLines {
	if f.shouldPassThrough() {
		return f.BackingReader.GetLines(firstLine, wantedLineCount)
	}

	acceptedLines := f.getAllLines()

	if len(acceptedLines) == 0 || wantedLineCount == 0 {
		return &reader.InputLines{
			StatusText: f.createStatus(nil),
		}
	}

	lastLine := firstLine.NonWrappingAdd(wantedLineCount - 1)

	// Prevent reading past the end of the available lines
	maxLineIndex := *linemetadata.IndexFromLength(len(acceptedLines))
	if lastLine.IsAfter(maxLineIndex) {
		lastLine = maxLineIndex

		// If one line was requested, then first and last should be exactly the
		// same, and we would get there by adding zero.
		firstLine = lastLine.NonWrappingAdd(1 - wantedLineCount)

		return f.GetLines(firstLine, firstLine.CountLinesTo(lastLine))
	}

	return &reader.InputLines{
		Lines:      acceptedLines[firstLine.Index() : firstLine.Index()+wantedLineCount],
		StatusText: f.createStatus(&lastLine),
	}
}

// In the general case, this will return a text like this:
// "Filtered: 1234/5678 lines  22%"
func (f *FilteringReader) createStatus(lastLine *linemetadata.Index) string {
	baseCount := f.BackingReader.GetLineCount()
	if baseCount == 0 {
		return "Filtered: No input lines"
	}

	baseCountString := "/" + linemetadata.IndexFromLength(baseCount).Format()
	if !f.BackingReader.ShouldShowLineCount() {
		baseCountString = ""
	}

	if lastLine == nil {
		// 100% because we're showing all 0 lines
		return "Filtered: 0" + baseCountString + " lines  100%"
	}

	acceptedCount := f.GetLineCount()
	acceptedCountString := linemetadata.IndexFromLength(acceptedCount).Format()

	percent := int(math.Floor(100 * float64(lastLine.Index()+1) / float64(acceptedCount)))

	lineString := "line"
	if (len(baseCountString) > 0 && baseCount != 1) || (len(baseCountString) == 0 && acceptedCount != 1) {
		lineString += "s"
	}

	return fmt.Sprintf("Filtered: %s%s %s  %d%%",
		acceptedCountString, baseCountString, lineString, percent)
}
