// Package commonmark provides functionality to convert CommonMark syntax to
// HTML.
package commonmark

import (
	"bufio"
	"bytes"
)

// ToHTMLBytes converts text formatted in CommonMark into the corresponding
// HTML.
//
// The input must be encoded as UTF-8.
//
// Line breaks in the output will be single '\n' bytes, regardless of line
// endings in the input (which can be CR, LF or CRLF).
//
// Note that the output might contain unsafe tags (e.g. <script>); if you are
// accepting untrusted user input, you must run the output through a sanitizer
// before sending it to a browser.
func ToHTMLBytes(markdown []byte) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(markdown))
	scanner.Split(scanLines)
	var html []byte
	for scanner.Scan() {
		line := scanner.Bytes()
		line = tabsToSpaces(line)

		html = append(html, line...)
		html = append(html, '\n')
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return html, nil
}

// scanLines is a split function for bufio.Scanner that splits on CR, LF or
// CRLF pairs. We need this because bufio.ScanLines itself only does CR and
// CRLF.
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	for i := 0; i < len(data); i++ {
		if data[i] == '\r' {
			if i+1 < len(data) && data[i+1] == '\n' {
				return i + 2, data[0:i], nil
			} else {
				return i + 1, data[0:i], nil
			}
		} else if data[i] == '\n' {
			return i + 1, data[0:i], nil
		}
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

// tabsToSpaces returns a slice (possibly the same one) in which tabs have been
// replaced by up to 4 spaces, depending on the tab stop position.
//
// It does not modify the input slice; a copy is made if needed.
func tabsToSpaces(line []byte) []byte {
	const tabStop = 4

	var tabCount int
	for remaining := line; len(remaining) > 0; {
		i := bytes.IndexByte(remaining, '\n')
		if i == -1 {
			break
		}
		remaining = remaining[i+1:]
		tabCount++
	}
	if tabCount == 0 {
		return line
	}

	output := make([]byte, 0, len(line)+4*tabCount)
	for i := 0; i < len(line); i++ {
		if line[i] == '\t' {

			spaces := bytes.Repeat([]byte{' '}, tabStop-i%tabStop)
			output = append(output, spaces...)
		} else {
			output = append(output, line[i])
		}
	}
	return output
}