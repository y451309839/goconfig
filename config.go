package goconfig

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

const (
	ErrSectionNotFound = iota + 1
	ErrKeyNotFound
	ErrBlankSectionName
	ErrCouldNotParse
)

const (
	DEFAULT_COMMENT   = "#"
	DEFAULT_SEPARATOR = "="
	DEFAULT_ARRAYSEP  = ","
	DEFAULT_SECTION   = "DEFAULT"
)

var breakLine = "\n"

func init() {
	if runtime.GOOS == "windows" {
		breakLine = "\r\n"
	}
}

type ConfigFile struct {
	configName  string                       //配名
	data        map[string]map[string]string //节点 -> 键 : 值
	sectionList []string                     //节点
	keyList     map[string][]string          //节点 -> 键
}

func NewConfigFile(name string) *ConfigFile {
	c := new(ConfigFile)
	c.configName = name
	c.data = make(map[string]map[string]string)
	c.keyList = make(map[string][]string)
	return c
}

func (c *ConfigFile) HasSection(section string) bool {
	_, ok := c.data[section]
	return ok
}

func (c *ConfigFile) GetSection(section string) map[string]string {
	if sections, ok := c.data[section]; ok {
		return sections
	}
	return nil
}

func (c *ConfigFile) SetValue(section, key, value string) bool {
	if len(section) == 0 {
		section = DEFAULT_SECTION
	}
	if len(key) == 0 {
		return false
	}

	if _, ok := c.data[section]; !ok {
		c.data[section] = make(map[string]string)
		c.sectionList = append(c.sectionList, section)
	}

	_, ok := c.data[section][key]
	c.data[section][key] = value
	if !ok {
		c.keyList[section] = append(c.keyList[section], key)
	}
	return !ok
}

func (c *ConfigFile) DeleteKey(section, key string) bool {
	if len(section) == 0 {
		section = DEFAULT_SECTION
	}
	if _, ok := c.data[section][key]; !ok {
		return false
	}
	delete(c.data[section], key)
	i := 0
	for _, keyName := range c.keyList[section] {
		if key == keyName {
			break
		}
		i++
	}
	c.keyList[section] = append(c.keyList[section][:i], c.keyList[section][i+1:]...)
	return true
}

func (c *ConfigFile) GetValue(section, key string) (string, error) {
	if len(section) == 0 {
		section = DEFAULT_SECTION
	}
	if _, ok := c.data[section]; !ok {
		return "", getError{ErrKeyNotFound, section}
	}
	value, ok := c.data[section][key]
	if !ok {
		return "", getError{ErrKeyNotFound, key}
	}
	//暂不支持变量
	return value, nil
}

func (c *ConfigFile) String(section, key string) (string, error) {
	return c.GetValue(section, key)
}

func (c *ConfigFile) Bool(section, key string) (bool, error) {
	value, err := c.GetValue(section, key)
	if err != nil {
		return false, nil
	}
	return strconv.ParseBool(value)
}

func (c *ConfigFile) Float64(section, key string) (float64, error) {
	value, err := c.GetValue(section, key)
	if err != nil {
		return 0.0, err
	}
	return strconv.ParseFloat(value, 64)
}

func (c *ConfigFile) Int(section, key string) (int, error) {
	value, err := c.GetValue(section, key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

func (c *ConfigFile) Int64(section, key string) (int64, error) {
	value, err := c.GetValue(section, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(value, 10, 64)
}

func (c *ConfigFile) Array(section, key, delim string) []string {
	val, err := c.GetValue(section, key)
	if err != nil || len(val) == 0 {
		return []string{}
	}

	vals := strings.Split(val, delim)
	for i := range vals {
		vals[i] = strings.TrimSpace(vals[i])
	}
	return vals
}

func (c *ConfigFile) MustValue(section, key string, defaultVal ...string) string {
	val, err := c.GetValue(section, key)
	if len(defaultVal) > 0 && (err != nil || len(val) == 0) {
		return defaultVal[0]
	}
	return val
}

// readError 格式化错误信息
type getError struct {
	Reason int
	Name   string
}

// Error implements Error interface.
func (err getError) Error() string {
	switch err.Reason {
	case ErrSectionNotFound:
		return fmt.Sprintf("section '%s' not found", err.Name)
	case ErrKeyNotFound:
		return fmt.Sprintf("key '%s' not found", err.Name)
	}
	return "invalid get error"
}
