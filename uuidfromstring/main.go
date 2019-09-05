package main

import (
	"fmt"

	"github.com/tsingson/uuid"
)

func main() {
	u, _ := uuid.NewV4()
	uuidStr := u.String()
	fmt.Println(" uuid string:   ", uuidStr)
	uu, _ := uuid.FromString(uuidStr)
	fmt.Println("uuid from string :  ", uu.String())

	fmt.Println("-------- uuid string:   ", "f84bf4ab-f965-49ff-afb2-ffc1474533f3")
	ui, _ := uuid.FromString("f84bf4ab-f965-49ff-afb2-ffc1474533f3")
	fmt.Println("---------------uuid from string :  ", ui.String())
}
