package main

import (
	"golinksmart/models/terminal"
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/fsnotify/fsnotify"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"time"
)

var (
	excelFileName string = "/Users/qinshen/git/linksmart-alpha/HD-740-2016-09-29MAC.xlsx"
)

func main() {
	Ximport(excelFileName)

	t, _ := timeFormat("01/02/19 19:04")
	p := fmt.Println
	p(t.Format("2006-01-02T15:04:05Z07:00"))

}

func Ximport(excelFileName string) error {
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println("file open error ")
		os.Exit(1)
	}
	//var terminalArray = map[int]terminal.TerminalSn{}
	for _, sheet := range xlFile.Sheets {
		for j, row := range sheet.Rows {

			var terminal terminal.TerminalSn

			terminal.Sn, _ = row.Cells[0].String()
			terminal.Model, _ = row.Cells[1].String()
			terminal.TerminalOwner, _ = row.Cells[2].String()
			terminal.ReleaseDate, _ = row.Cells[3].String()
			//terminalArray[i] = terminal
			fmt.Println("%d", j, terminal)
		}

	}

	return nil
	//return terminalArray, nil
}

func timeFormat(importTime string) (time.Time, error) {
	withNanos := "01/02/06 15:04"
	// importTime := "01/02/19 19:04"
	t, err := time.Parse(withNanos, importTime)
	return t, err
}

func XlsxWatch(pathname string) error {

	if len(pathname) == 0 {
		return errors.New("Path error")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(pathname)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	return nil
}
