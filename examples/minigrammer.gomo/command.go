package main

import (
	"github.com/wowplus2/golang/examples/minigrammer.gomo/manager"
	"os"
	"github.com/pkg/errors"
	"path/filepath"
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

var  (
	initCommand = manager.Command{
		Name:	"init",
		Usage:	"init	: gomo init",
		Run:	initFunc,
	}
	newCommand = manager.Command{
		Name:	"new",
		Usage:	"new	: gomo new 'content'",
		Run:	newFunc,
	}
	listCommand = manager.Command{
		Name:	"list",
		Usage:	"list	: gomo list",
		Run:	listFunc,
	}
	delCommand = manager.Command{
		Name:	"del",
		Usage:	"del	: gomo del [number]",
		Run:	delFunc,
	}
)

func isExistMemoFile(memoFile string) bool {
	if _, err := os.Stat(memoFile); err != nil {
		return false
	}

	return true
}

func getAllMemos(memoFile string) ([]Memo, error) {
	var memos []Memo

	if !isExistMemoFile(memoFile) {
		return nil, errors.New("You must run 'init' command first")
	}

	fileContent, err := ioutil.ReadFile(memoFile)
	if err != nil {
		return nil, errors.New("Error read the file data")
	}

	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &memos); err != nil {
			return nil, errors.New("Error unmarshaling the memo text")
		}
	}

	return memos, nil
}

func setWriteMemo(memoFile string, memos []Memo) error {
	file, err := os.OpenFile(memoFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return errors.New("Error open the memo file")
	}
	defer file.Close()

	memoListJson, err := json.MarshalIndent(memos, "", "	")
	if err != nil {
		return errors.New("Error mashaling the memo content")
	}

	_, err = file.Write(memoListJson)
	if err != nil {
		return errors.New("Error write the memo content to memo file")
	}

	return nil
}


// initFunc initializes the memo file with filename of manager.
// if there is memo file already, it return error else return nil if success.
func initFunc(memoFile string, args []string) error {
	if isExistMemoFile(memoFile) {
		return errors.New("Memo file already exists...")
	}

	dir := filepath.Dir(memoFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0766)
	}

	_, err := os.Create(memoFile)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error create the memo file...")
	}

	fmt.Println("Memo file was created successfully")

	return nil
}

// newFunc creates new memo and store that to file.
func newFunc(memoFile string, args []string) error {
	if len(args) != 1 {
		return errors.New("Invalid arguments: new command must take only one argument")
	}

	memos, err := getAllMemos(memoFile)
	if err != nil {
		return err
	}

	newMemo := Memo{
		Content:	args[0],
		CreatedAt:	time.Now(),
	}

	memos = append(memos, newMemo)
	err = setWriteMemo(memoFile, memos)
	if err != nil {
		return errors.New("Error write the memo to file")
	}

	fmt.Printf("Memo was created: [%s]\n", newMemo.Content)

	return nil
}

// listFunc prints the all memos on console.
func listFunc(memoFile string, args []string) error {
	memos, err := getAllMemos(memoFile)
	if err != nil {
		return err
	}

	PrintMemos(memos)

	return nil
}

// delFunc deletes a memo.
func delFunc(memoFile string, args []string) error {
	nArgs := len(args)
	if nArgs > 1 {
		return errors.New("Invalid arguments: delete command could take one argument at most")
	}

	memos, err := getAllMemos(memoFile)
	if err != nil {
		return err
	}

	var memoNo int

	if nArgs == 1 {
		if memoNo, err = strconv.Atoi(args[0]); err != nil {
			return errors.New("Integer is allowed only")
		}
	} else {
		PrintMemos(memos)

		fmt.Print("Enter memo number to delete")
		fmt.Scanln(&memoNo)
	}

	if memoNo < 1 || memoNo > len(memos) {
		return errors.New("Invalid memo number")
	}

	memos = append(memos[:memoNo-1], memos[memoNo:]...)
	err = setWriteMemo(memoFile, memos)
	if err != nil {
		return errors.New("Error write the memos to file")
	}

	fmt.Printf("Memo #%d was deleted...\n", memoNo)

	return nil
}