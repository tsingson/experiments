// +build !plan9

package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/tsingson/gin/fasthttputils"
)

func main() {
	path, _ := fasthttputils.GetCurrentExecDir()
	path = path + "/vod"
	fsNotifyWatcher(path)
	select {}
}
func fsNotifyWatcher(path ...string) {
	watcher, err := fasthttputils.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	//	defer watcher.Close()
	done := make(chan bool)
	go watchWorker(watcher)

	for _, mypath := range path {
		if len(mypath) > 0 {
			err = watcher.AddRecursive(mypath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	<-done
}

func watchWorker(watcher *fasthttputils.RWatcher) {

	for {
		select {
		case event := <-watcher.Events:
			//	log.Println("event:", event)
			if event.Op&fsnotify.Create == fsnotify.Create {
				// call listener
				go fileWorker(event.Name)
			}
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				filenmae := event.Name
				fmt.Println("chmod        ", filenmae)
			}

		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}

}

func fileWorker(filenme string) {
	checksum, filesize, err := getFileInfo(filenme)
	fmt.Println("file neme: ", filenme)
	if err == nil {
		fmt.Println("file checksum: ", checksum)
		fmt.Println("file size: ", filesize)
		fmt.Println(" ")
	}
}
