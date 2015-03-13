package goconfig

import (
	"bytes"
	"os"
	"strings"
)

// 在等号边加上空格更好看
var PrettyFormat = true

func (c *ConfigFile) saveFile(filename string) (err error) {
	if len(filename) == 0 {
		filename = c.fileName
	}

	var f *os.File
	if f, err = os.Create(filename); err != nil {
		return err
	}

	equalSign := DEFAULT_SEPARATOR
	if PrettyFormat {
		equalSign = " " + DEFAULT_SEPARATOR + " "
	}

	buf := bytes.NewBuffer(nil)
	for _, section := range c.sectionList {

		if comments := c.GetSectionComments(section); len(comments) > 0 {
			if _, err = buf.WriteString(comments + breakLine); err != nil {
				return err
			}
		}

		if section != DEFAULT_SECTION {
			if _, err = buf.WriteString("[" + section + "]" + breakLine); err != nil {
				return err
			}
		}

		for _, key := range c.keyList[section] {
			if key != " " {

				if comments := c.GetKeyComments(section, key); len(comments) > 0 {
					if _, err = buf.WriteString(comments + breakLine); err != nil {
						return nil
					}
				}

				keyName := key
				if strings.Contains(keyName, DEFAULT_SEPARATOR) {
					if strings.Contains(keyName, "`") {
						if strings.Contains(keyName, `"`) {
							keyName = `"""` + keyName + `"""`
						} else {
							keyName = `"` + keyName + `"`
						}
					} else {
						keyName = "`" + keyName + "`"
					}
				}

				value := c.data[section][key]
				if strings.Contains(value, "`") {
					if strings.Contains(value, `"`) {
						value = `"""` + value + `"""`
					} else {
						value = `"` + value + `"`
					}
				}

				if _, err = buf.WriteString(keyName + equalSign + value + breakLine); err != nil {
					return err
				}
			}
		}
		// 在节点后加上一个空行
		if _, err = buf.WriteString(breakLine); err != nil {
			return err
		}
	}
	if _, err = buf.WriteTo(f); err != nil {
		return err
	}
	return f.Close()
}
