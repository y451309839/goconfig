package goconfig

import (
	"fmt"
	"testing"
)

func Test_newConfig(t *testing.T) {
	cfg := NewConfigFile("test")
	cfg.LoadFile("test.ini")
	fmt.Printf("cfg.data is: %s\n", cfg.data)
}

func Test_String(t *testing.T) {
	cfg := NewConfigFile("test")
	cfg.LoadFile("test.ini")
	username, _ := cfg.String("", "username")
	password, _ := cfg.String("DEFAULT", "password")
	testdata, _ := cfg.String("CONFIGS", "testdata")
	fmt.Printf("cfg.String(\"\", \"username\") result is: %s\n", username)
	fmt.Printf("cfg.String(\"DEFAULT\", \"password\") result is: %s\n", password)
	fmt.Printf("cfg.String(\"CONFIGS\", \"testdata\") result is: %s\n", testdata)
}
