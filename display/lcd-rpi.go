//go:build linux && arm

package display

import (
	"time"

	device "github.com/d2r2/go-hd44780"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
)

const (
	numChars = 20
	numLines = 4
	cmdClear = iota
	cmdBacklightOn
	cmdBacklightOff
	cmdPrintln
)

var lgLcd = logger.NewPackageLogger("lcd", logger.InfoLevel)

// lcdData represents the internal structure for managing an LCD using an I2C bus and supporting
// scrolling functionality.
type lcdData struct {
	i2cBus       *i2c.I2C
	dev          *device.Lcd
	lines        [numLines]device.ShowOptions
	ticker       [numLines]*time.Ticker
	cmdChan      chan command
	scrollSpeed  int
	charsPerLine int
	initDelay    int
	retryCount   int
}

// command represents a structure for issuing commands to an LCD.
// cmd specifies the type of command to execute.
// lineNum indicates the target line number for the command.
// lineText holds the text associated with the command for line-based operations.
type command struct {
	cmd      int
	lineNum  int
	lineText string
}

// printLine displays the given text on a specified line of the LCD if the line index is within bounds.
// Returns an error if displaying the message fails.
// An empty string is replaced with a single space to prevent library panics.
func (lcd *lcdData) printLine(line int, text string) (err error) {
	if line >= 0 && line < numLines {
		if len(text) == 0 {
			text = " " // avoid panic, because the library can't handle empty strings
		}
		err = lcd.dev.ShowMessage(text, lcd.lines[line])
		return err
	}
	return nil
}

// runTicker continuously scrolls the given text on the specified line using a ticker with an interval
// based on scrollSpeed.
func (lcd *lcdData) runTicker(line int, text string) {
	lcd.ticker[line] = time.NewTicker(time.Duration(lcd.scrollSpeed) * time.Millisecond)
	textWithPadding := text + "     "
	for range lcd.ticker[line].C {
		lcd.cmdChan <- command{
			cmd:      cmdPrintln,
			lineNum:  line,
			lineText: textWithPadding,
		}
		textWithPadding = textWithPadding[1:] + textWithPadding[:1]
	}
}

// printAndScrollLine displays text on a specified line or scrolls it if the text exceeds the character limit.
func (lcd *lcdData) printAndScrollLine(line int, text string) {
	line = line % numLines
	if lcd.ticker[line] != nil {
		lcd.ticker[line].Stop()
		lcd.ticker[line] = nil
	}
	if len(text) <= numChars {
		lcd.cmdChan <- command{
			cmd:      cmdPrintln,
			lineNum:  line,
			lineText: text,
		}
	} else {
		go lcd.runTicker(line, text)
	}
}

// commandHandler processes commands received via the cmdChan channel and performs actions on the LCD.
func (lcd *lcdData) commandHandler() {
	var err error
	for {
		err = nil
		c := <-lcd.cmdChan
		switch c.cmd {
		case cmdClear:
			err = lcd.dev.Clear()
			time.Sleep(100 * time.Millisecond)
		case cmdBacklightOn:
			err = lcd.dev.BacklightOn()
		case cmdBacklightOff:
			err = lcd.dev.BacklightOff()
		case cmdPrintln:
			err = lcd.printLine(c.lineNum, c.lineText)
		default:
			panic("unhandled default case")
		}
		if err != nil {
			lgLcd.Error(err.Error())
			lcd.retryDevice()
		}
	}
}

// Backlight controls the LCD backlight. Pass true to enable the backlight or false to disable it.
func (lcd *lcdData) Backlight(on bool) {
	if on {
		lcd.cmdChan <- command{
			cmd: cmdBacklightOn,
		}
	} else {
		lcd.cmdChan <- command{
			cmd: cmdBacklightOff,
		}
	}
}

// ClearLine clears the content on the specified line of the LCD by sending an empty string to it.
func (lcd *lcdData) ClearLine(line int) {
	// dummy function, not really needed for lcdData
	lcd.cmdChan <- command{
		cmd:      cmdPrintln,
		lineNum:  line,
		lineText: "",
	}
}

// Clear sends a command to clear the LCD and reset its content.
func (lcd *lcdData) Clear() {
	lcd.cmdChan <- command{
		cmd: cmdClear,
	}
}

// Close releases resources associated with the LCD, stops active tickers, and closes the I2C bus connection.
func (lcd *lcdData) Close() {
	if lcd.i2cBus != nil {
		for i := 0; i < numLines; i++ {
			if lcd.ticker[i] != nil {
				lcd.ticker[i].Stop()
				lcd.ticker[i] = nil
			}
		}
		time.Sleep(2 * time.Second)
		_ = lcd.i2cBus.Close()
	}
}

// PrintLine writes text to a specific line on the LCD or scrolls it if the scroll flag is true and text exceeds the limit.
func (lcd *lcdData) PrintLine(line int, text string, scroll bool) {
	if line < 0 || line >= numLines {
		lgLcd.Error("LCD display row is ouf of bounds: ", line)
		return
	}
	if scroll {
		lcd.printAndScrollLine(line, text)
	} else {
		if lcd.ticker[line] != nil {
			lcd.ticker[line].Stop()
			lcd.ticker[line] = nil
		}
		lcd.cmdChan <- command{
			cmd:      cmdPrintln,
			lineNum:  line,
			lineText: text,
		}
	}
}

// GetCharsPerLine returns the maximum number of characters that can be displayed per line on the LCD.
func (lcd *lcdData) GetCharsPerLine() int {
	return lcd.charsPerLine
}

// GetMinMaxRowNum returns the minimum and maximum row numbers available on the LCD.
func (lcd *lcdData) GetMinMaxRowNum() (int, int) {
	return 0, numLines - 1
}

// retryDevice attempts to reinitialize the LCD device and its I2C connection after a failure, incrementing retryCount.
func (lcd *lcdData) retryDevice() {
	lgLcd.Info("Start of retryDevice(): ", lcd.retryCount)
	var err error
	lcd.i2cBus, err = i2c.NewI2C(0x27, 1)
	if err != nil {
		lgLcd.Error(err.Error())
	}
	time.Sleep(3 * time.Second)

	lcd.dev, err = device.NewLcd(lcd.i2cBus, device.LCD_20x4)
	if err != nil {
		lgLcd.Error(err.Error())
	}
	time.Sleep(time.Duration(lcd.initDelay) * time.Second)
	lcd.retryCount++
	lcd.Clear()
	lcd.Backlight(true)
	lgLcd.Info("End of retryDevice(): %d", lcd.retryCount)
}

// New initializes and returns a new Display instance with the given scroll header option, speed, and
// initial delay.
// Returns an error if the initialization process fails.
func New(scrollHeader bool, speed int, initDelay int) (_ Display, err error) {
	lgLcd.Debug("LCD initializing...")
	_ = logger.ChangePackageLogLevel("i2c", logger.WarnLevel)
	lcd := lcdData{scrollSpeed: speed, charsPerLine: numChars, cmdChan: make(chan command)}
	err = nil

	lcd.retryCount = 0
	lcd.initDelay = initDelay
	lcd.lines[0] = device.SHOW_LINE_1 | device.SHOW_BLANK_PADDING
	if !scrollHeader {
		lcd.lines[0] |= device.SHOW_ELIPSE_IF_NOT_FIT
	}
	lcd.lines[1] = device.SHOW_LINE_2 | device.SHOW_BLANK_PADDING
	lcd.lines[2] = device.SHOW_LINE_3 | device.SHOW_BLANK_PADDING
	lcd.lines[3] = device.SHOW_LINE_4 | device.SHOW_BLANK_PADDING

	lcd.i2cBus, err = i2c.NewI2C(0x27, 1)
	if err != nil {
		lgLcd.Error(err.Error())
		return &lcd, err
	}
	time.Sleep(3 * time.Second)

	lcd.dev, err = device.NewLcd(lcd.i2cBus, device.LCD_20x4)
	if err != nil {
		lgLcd.Error(err.Error())
		return &lcd, err
	}
	// time.Sleep(time.Duration(lcd.initDelay) * time.Second)

	go lcd.commandHandler()

	lcd.Clear()
	lcd.Backlight(true)
	return &lcd, err
}
