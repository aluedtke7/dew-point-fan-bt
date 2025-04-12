//go:build !linux && !arm

package display

import (
	"github.com/d2r2/go-logger"
	"strings"
)

var lgTerm = logger.NewPackageLogger("term", logger.InfoLevel)

// TerminalDisplay simulates a display device by printing to the terminal.
type TerminalDisplay struct {
	backlight      bool     // Represents the backlight status
	rows           []string // Simulates rows in the display
	charsPerLine   int      // Number of characters supported per line
	minRow, maxRow int      // Range of row indices supported by the display
}

// New creates a new TerminalDisplay
func New(_ bool, _ int, _ int) (*TerminalDisplay, error) {
	const charsPerLine = 20
	const numRows = 4
	return &TerminalDisplay{
		backlight:    false,
		rows:         make([]string, numRows),
		charsPerLine: charsPerLine,
		minRow:       0,
		maxRow:       numRows - 1,
	}, nil
}

// Backlight enables or disables the display backlight.
func (d *TerminalDisplay) Backlight(on bool) {
	d.backlight = on
	lgTerm.Infof("Backlight set to: %t", on)
}

// Clear clears all content from the display.
func (d *TerminalDisplay) Clear() {
	for i := range d.rows {
		d.rows[i] = strings.Repeat(" ", d.charsPerLine)
	}
	lgTerm.Info("Display cleared")
}

// ClearLine clears content from a specific line on the display.
func (d *TerminalDisplay) ClearLine(ofs int) {
	if ofs < d.minRow || ofs > d.maxRow {
		lgTerm.Warnf("Line %d is out of range [%d-%d]", ofs, d.minRow, d.maxRow)
		return
	}
	d.rows[ofs] = strings.Repeat(" ", d.charsPerLine)
	lgTerm.Infof("Line %d cleared", ofs)
}

// Close closes the display (no-op for this implementation).
func (d *TerminalDisplay) Close() {
	lgTerm.Info("Display closed")
}

// GetCharsPerLine retrieves the number of characters the display supports per line.
func (d *TerminalDisplay) GetCharsPerLine() int {
	return d.charsPerLine
}

// GetMinMaxRowNum retrieves the minimum and maximum row indices supported by the display.
func (d *TerminalDisplay) GetMinMaxRowNum() (int, int) {
	return d.minRow, d.maxRow
}

// PrintLine prints text to a specific line on the display, with an optional scrolling feature.
// The scrolling is not implemented.
func (d *TerminalDisplay) PrintLine(line int, text string, scroll bool) {
	if line < d.minRow || line > d.maxRow {
		lgTerm.Warnf("Line %d is out of range [%d-%d]\n", line, d.minRow, d.maxRow)
		return
	}

	// Truncate or pad the text to fit into the line
	if len(text) > d.charsPerLine {
		if scroll {
			text = text[len(text)-d.charsPerLine:]
		} else {
			text = text[:d.charsPerLine]
		}
	} else {
		text = text + strings.Repeat(" ", d.charsPerLine-len(text))
	}

	d.rows[line] = text
	lgTerm.Infof("Line %d: %s", line, text)
}
