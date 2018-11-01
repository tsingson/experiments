package main

import (
	"./Controller"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//SeFTPConfig is a config predefined for convenience.
var SeFTPConfig = Config{}

func handleGet(serverCommand []string, clientCommand []string) {
	if serverCommand[0] != "FILE" {
		if (len(clientCommand) <= 2) || (clientCommand[2] == "TCP") {
			subftpCon := Controller.TCPController{ServerAddr: SeFTPConfig.ServerAddr + ":" + serverCommand[2], Passwd: SeFTPConfig.Passwd}
			GET(subftpCon)
		} else if clientCommand[2] == "UDP" {
			subftpCon := Controller.KCPController{ServerAddr: SeFTPConfig.ServerAddr + ":" + serverCommand[2], Passwd: SeFTPConfig.Passwd}
			GET(subftpCon)
		}
	}
}

func handlePost(serverCommand []string, clientCommand []string) {
	if len(clientCommand) >= 2 {
		if (len(clientCommand) <= 2) || (clientCommand[2] == "TCP") {
			subftpCon := Controller.TCPController{ServerAddr: SeFTPConfig.ServerAddr + ":" + serverCommand[2], Passwd: SeFTPConfig.Passwd}
			POST(subftpCon)
		} else if clientCommand[2] == "UDP" {
			subftpCon := Controller.KCPController{ServerAddr: SeFTPConfig.ServerAddr + ":" + serverCommand[2], Passwd: SeFTPConfig.Passwd}
			POST(subftpCon)
		}
	}
}

func processRemoteCommand(plainClientCommand string, seftpCon Controller.TCPController) {
	clientCommand := strings.Fields(plainClientCommand)
	seftpCon.SendText(plainClientCommand)
	plainServerCommand, rErr := seftpCon.GetText()
	checkerr(rErr)
	serverCommand := strings.Fields(plainServerCommand)
	log.Println("Response From Server:", plainServerCommand)
	switch clientCommand[0] {
	case "GET":
		handleGet(serverCommand, clientCommand)
	case "POST":
		handlePost(serverCommand, clientCommand)
	}
}

func processLocalCommand(plainClientCommand string) {
	clientCommand := strings.Fields(plainClientCommand)
	switch clientCommand[0] {
	case "cd":
		newPath := clientCommand[1]
		err := os.Chdir(newPath)
		if !checkerr(err) {
			log.Println("Dir change failed")
		} else {
			log.Println("Dir changed")
		}
	case "ls":
		var list []string
		if len(clientCommand) > 1 {
			path := clientCommand[1]
			list = Ls(path)
		} else {
			list = Ls("")
		}
		log.Println(strings.Join(list, " | "))
	case "rm":
		if len(clientCommand) > 1 {
			err := os.Remove(clientCommand[1])
			checkerr(err)
		} else {
			log.Println("No specific file")
		}
	case "sha3sum":
		if len(clientCommand) > 1 {
			sum, err := SHA3FileHash(clientCommand[1])
			checkerr(err)
			log.Println(sum)
		} else {
			log.Println("No specific file")
		}
	case "exit":
		log.Println("Exit SeFTP")
		os.Exit(0)
	default:
		log.Println("Unknown command")
	}
}

func processCommand(plainClientCommand string, seftpCon Controller.TCPController) {
	if IsUpper(strings.Fields(plainClientCommand)[0]) {
		processRemoteCommand(plainClientCommand, seftpCon)
	} else {
		processLocalCommand(plainClientCommand)
	}
}

func main() {
	SeFTPConfig.Parse()
	seftpCon := Controller.TCPController{ServerAddr: SeFTPConfig.ServerAddr + ":" + strconv.Itoa(SeFTPConfig.ServerPort), Passwd: SeFTPConfig.Passwd}
	seftpCon.EstabConn()

	defer seftpCon.CloseConn()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		//log.Println(text)
		processCommand(text, seftpCon)
	}
}
