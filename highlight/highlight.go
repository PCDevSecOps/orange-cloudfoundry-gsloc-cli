// Code from github.com/andreazorzetto/yh/highlight made by Andrea Zorzetto but using messages

package highlight

import (
	"bufio"
	"fmt"
	msg "github.com/ArthurHlt/messages"
	"io"
	"strings"
)

func Highlight(r io.Reader) (string, error) {
	// Service vars
	foundChompingIndicator := false
	indentationSpacesBeforeComment := 0

	// Warm-up the engine
	scanner := bufio.NewScanner(r)

	var buf strings.Builder

	// Get the juice
	for scanner.Scan() {
		if scanner.Text() == "EOF" {
			break
		}

		// Check for errors during Stdin read
		if err := scanner.Err(); err != nil {
			return "", err
		}

		// Drink the juice
		l := yamlLine{raw: scanner.Text()}

		if foundChompingIndicator && (l.indentationSpaces() > indentationSpacesBeforeComment) {
			// Found multiline comment or configmap, not treated as YAML at all
			buf.WriteString(multiline(l))

		} else if l.isKeyValue() {
			// This is a valid YAML key: value line. Key and value are returned in l

			if l.isComment() {
				// This line is a comment
				buf.WriteString(comment(l))
			} else if l.valueIsNumberOrIP() {
				// The value is a number or an IP address x.x.x.x
				buf.WriteString(keyNumberOrIP(l))

			} else if l.valueIsBoolean() {
				// The value is boolean true or false
				buf.WriteString(keyBool(l))

			} else {
				// The is a normal key/value line
				buf.WriteString(keyValue(l))
			}

			if l.valueContainsChompingIndicator() {
				// This line contains a chomping indicator, sign of a possible multiline text coming next

				// Setting flag for next execution
				foundChompingIndicator = true

				// Saving current number of indentation spaces
				indentationSpacesBeforeComment = l.indentationSpaces()

			} else {
				// Resetting multiline flag
				foundChompingIndicator = false
			}

		} else if !l.isEmptyLine() {
			// This is not a YAML key: value line and is not empty

			if l.isUrl() {
				// Value is a URL
				buf.WriteString(url(l))
			} else if l.isComment() {
				// This line is a comment
				buf.WriteString(comment(l))
			} else if l.isElementOfList() {
				// This line is an element of a list
				buf.WriteString(listElement(l))
			} else {
				// This line is not valid YAML
				buf.WriteString(invalidLine(l))
			}

			foundChompingIndicator = false

		} else if l.isEmptyLine() {
			// This is an empty line
			fmt.Println(l.raw)
		}

	}

	return buf.String(), nil
}

func keyValue(l yamlLine) string {
	return fmt.Sprintf("%v: %v\n", msg.BrightRed(l.key), msg.Yellow(l.value))
}

func keyNumberOrIP(l yamlLine) string {
	return fmt.Sprintf("%v: %v\n", msg.BrightRed(l.key), msg.Blue(l.value))
}

func keyBool(l yamlLine) string {
	return fmt.Sprintf("%v: %v\n", msg.BrightRed(l.key), msg.Blue(l.value))
}

func comment(l yamlLine) string {
	return fmt.Sprintf("%v %v\n", msg.Gray(13, l.key), msg.Gray(13, l.value))
}

func listElement(l yamlLine) string {
	return fmt.Sprintf("%v\n", msg.Yellow(l.raw))
}

func invalidLine(l yamlLine) string {
	return fmt.Sprintf("%v\n", msg.Black(l.raw).BgBrightRed())
}

func multiline(l yamlLine) string {
	return fmt.Sprintf("%v\n", msg.Gray(20-1, l.raw))
}
func url(l yamlLine) string {
	return fmt.Sprintf("%v\n", msg.Yellow(l.raw))
}
