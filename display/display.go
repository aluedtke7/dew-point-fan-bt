package display

// Display defines an interface for interacting with a display device, supporting operations like printing and clearing text.
// Backlight enables or disables the display backlight.
// Clear clears all content from the display.
// ClearLine clears content from a specific line on the display.
// Close closes the display and releases any associated resources.
// GetCharsPerLine retrieves the number of characters the display supports per line.
// GetMinMaxRowNum retrieves the minimum and maximum row indices supported by the display.
// PrintLine prints text to a specific line on the display, with an optional scrolling feature.
type Display interface {
	Backlight(on bool)
	Clear()
	ClearLine(ofs int)
	Close()
	GetCharsPerLine() int
	GetMinMaxRowNum() (int, int)
	PrintLine(line int, text string, scroll bool)
}
