package display

// Display defines an interface for interacting with a display device, supporting operations like printing and clearing text.
type Display interface {

	// Backlight toggles the display backlight on or off based on the provided boolean value.
	Backlight(on bool)

	// Clear removes all content from the display.
	Clear()

	// ClearLine clears the content of a specific line on the display, identified by the provided line offset.
	ClearLine(ofs int)

	// Close releases resources associated with the display and properly shuts down the interface.
	Close()

	// GetCharsPerLine retrieves the number of characters that can fit on a single line of the display.
	GetCharsPerLine() int

	// GetMinMaxRowNum retrieves the minimum and maximum row indices supported by the display.
	GetMinMaxRowNum() (int, int)

	// PrintLine outputs the provided text to a specific line on the display and optionally enables scrolling.
	PrintLine(line int, text string, scroll bool)
}
