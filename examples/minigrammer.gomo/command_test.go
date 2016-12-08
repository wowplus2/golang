package main

import (
	"testing"
	"github.com/wowplus2/golang/examples/minigrammer.gomo/manager"
	"path/filepath"
	"os"
	"io/ioutil"
	"strings"
)

var (
	homeDir		string
	memoFile	string
)

func setupTestCase(t *testing.T) {
	t.Log("Setup testcase")

	homeDir = manager.GetHomeDir()
	memoFile = filepath.Join(homeDir, ".gomo-tmp", "memo.json")

	os.Remove(memoFile)
}

func teardownTestCase(t *testing.T) {
	t.Log("Teardown testcase")

	homeDir = manager.GetHomeDir()
	memoFile = filepath.Join(homeDir, ".gomo-tmp", "memo.json")

	os.Remove(memoFile)
}

func TestInitFunc(t *testing.T) {
	setupTestCase(t)
	defer teardownTestCase(t)

	args := []string{}

	if err := initFunc(memoFile, args); err != nil {
		t.Errorf("Error occur when run init command: %v", err)
	}

	if _, err := os.Stat(memoFile); os.IsNotExist(err) {
		t.Error("Memo file is not created correctly")
	}
}

func TestNewFunc(t *testing.T) {
	setupTestCase(t)
	defer teardownTestCase(t)

	var args []string
	initFunc(memoFile, args)

	args = []string{"first", "second"}
	if err := newFunc(memoFile, args); err != nil {
		t.Error("Except error")
	}

	args = []string{"first"}
	if err := newFunc(memoFile, args); err != nil {
		t.Error("Except no error")
	}

	fileContents, _ := ioutil.ReadFile(memoFile)
	if !strings.Contains(string(fileContents), "first") {
		t.Error("Except memo file contains 'first'")
	}
}

func TestDelete(t *testing.T) {
	setupTestCase(t)
	defer teardownTestCase(t)

	var args []string

	initFunc(memoFile, args)
	args = []string{"first"}

	newFunc(memoFile, args)

	args = []string{"a"}
	if err := delFunc(memoFile, args); err != nil {
		t.Error("Except error")
	}

	args = []string{"2"}
	if err := delFunc(memoFile, args); err != nil {
		t.Error("Except error")
	}

	args = []string{"1"}
	if err := delFunc(memoFile, args); err != nil {
		t.Error("Except no error")
	}

	fileContents, _ := ioutil.ReadFile(memoFile)

	if strings.Contains(string(fileContents), "first") {
		t.Error("Except memo file should not contains 'first'")
	}
}
