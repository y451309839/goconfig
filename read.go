package goconfig

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func (c *ConfigFile) read(read io.Reader) (err error) {
	buf := bufio.NewReader(read)

	// UTF8去除BOM
	// http://en.wikipedia.org/wiki/Byte_order_mark#Representations_of_byte_order_marks_by_encoding
	mask, err := buf.Peek(3)
	if err == nil && len(mask) >= 3 &&
		mask[0] == 239 && mask[1] == 187 && mask[2] == 191 {
		buf.Read(mask)
	}

	section := DEFAULT_SECTION
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		lineLength := len(line)
		if err != nil {
			if err != io.EOF {
				return err
			}

			if lineLength == 0 {
				break
			}
		}

		switch {
		case lineLength == 0:
			continue
		case line[0] == '#' || line[0] == ';':
			continue
		case section[0] == '[' && section[len(section)-1] == ']':
			section = strings.TrimSpace(line[1 : lineLength-1])
			c.SetValue(section, "", "")
		case section == "":
			return readError{ErrBlankSectionName, line}
		default:
			var (
				i        int
				key      string
				keyQuote string
				value    string
				valQuote string
			)

			if line[0] == '"' {
				if lineLength > 6 && line[0:3] == `"""` {
					keyQuote = `"""`
				} else {
					keyQuote = `"`
				}
			} else if line[0] == '`' {
				keyQuote = "`"
			}
			if keyQuote != "" {
				qLen := len(keyQuote)
				pos := strings.Index(line[qLen:], keyQuote)
				if pos == -1 {
					return readError{ErrCouldNotParse, line}
				}
				pos = pos + qLen
				i = strings.IndexAny(line[pos:], DEFAULT_SEPARATOR)
				if i <= 0 {
					return readError{ErrCouldNotParse, line}
				}
				i = i + pos
				key = line[qLen:pos] //引号保留空格
			} else {
				i = strings.IndexAny(line, DEFAULT_SEPARATOR)
				if i <= 0 {
					return readError{ErrCouldNotParse, line}
				}
				key = strings.TrimSpace(line[0:i])
			}

			//字符串支持引号包围
			lineRight := strings.TrimSpace(line[i+1:])
			lineRightLen := len(lineRight)
			firstChar := ""
			if lineRightLen >= 2 {
				firstChar = lineRight[0:1]
			}
			if firstChar == "`" {
				valQuote = "`"
			} else if lineRightLen >= 6 && lineRight[0:3] == `"""` {
				valQuote = `"""`
			}

			if valQuote != "" {
				qLen := len(valQuote)
				pos := strings.LastIndex(lineRight[qLen:], valQuote)
				if pos == -1 {
					return readError{ErrCouldNotParse, line}
				}
				pos = pos + qLen
				value = lineRight[qLen:pos]
			} else {
				value = strings.TrimSpace(lineRight[0:])
			}

			c.SetValue(section, key, value)
		}
	}
	return nil
}

func (c *ConfigFile) LoadFile(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.read(f)
}

// readError 格式化错误信息
type readError struct {
	Reason  int
	Content string // Line content
}

// Error implement Error interface.
func (err readError) Error() string {
	switch err.Reason {
	case ErrBlankSectionName:
		return "empty section name not allowed"
	case ErrCouldNotParse:
		return fmt.Sprintf("could not parse line: %s", string(err.Content))
	}
	return "invalid read error"
}
