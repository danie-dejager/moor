package twin

// Used for testing.
//
// Try GetRow() after some SetCell() calls to see what you got.
type FakeScreen struct {
	width  int
	height int
	cells  [][]StyledRune
}

func NewFakeScreen(width int, height int) *FakeScreen {
	rows := make([][]StyledRune, height)
	for i := 0; i < height; i++ {
		rows[i] = make([]StyledRune, width)
	}

	return &FakeScreen{
		width:  width,
		height: height,
		cells:  rows,
	}
}

func (screen *FakeScreen) Close() {
	// This method intentionally left blank
}

func (screen *FakeScreen) Clear() {
	// This method's contents has been copied from UnixScreen.Clear()

	empty := NewStyledRune(' ', StyleDefault)

	width, height := screen.Size()
	for row := 0; row < height; row++ {
		for column := 0; column < width; column++ {
			screen.cells[row][column] = empty
		}
	}
}

func (screen *FakeScreen) SetCell(column int, row int, styledRune StyledRune) int {
	// This method's contents has been copied from UnixScreen.Clear()

	if column < 0 {
		return styledRune.Width()
	}
	if row < 0 {
		return styledRune.Width()
	}

	width, height := screen.Size()
	if column >= width {
		return styledRune.Width()
	}
	if row >= height {
		return styledRune.Width()
	}

	if column+styledRune.Width() > width {
		// This cell is too wide for the screen, write a space instead
		screen.cells[row][column] = NewStyledRune(' ', styledRune.Style)
		return styledRune.Width()
	}

	screen.cells[row][column] = styledRune

	return styledRune.Width()
}

func (screen *FakeScreen) Show() {
	// This method intentionally left blank
}

func (screen *FakeScreen) ShowNLines(int) {
	// This method intentionally left blank
}

func (screen *FakeScreen) Size() (width int, height int) {
	return screen.width, screen.height
}

func (screen *FakeScreen) RequestTerminalBackgroundColor() {
	// This method intentionally left blank
}

func (screen *FakeScreen) ShowCursorAt(_ int, _ int) {
	// This method intentionally left blank
}

func (screen *FakeScreen) Events() chan Event {
	// TODO: Do better here if or when this becomes a problem
	return nil
}

func (screen *FakeScreen) GetRow(row int) []StyledRune {
	return withoutHiddenRunes(screen.cells[row])
}
