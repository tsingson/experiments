package main

import (
	"fmt"
	"sort"

	"github.com/karrick/godirwalk"
	"github.com/tsingson/fastweb/fasthttputils"
)

func main() {
	osDirename, _ := fasthttputils.GetCurrentPath()

	children, err := godirwalk.ReadDirnames(osDirename, nil)
	if err != nil {
		//	return nil, errors.Wrap(err, "cannot get list of directory children")
	}
	sort.Strings(children)
	for _, child := range children {
		fmt.Println(child)
	}

}
