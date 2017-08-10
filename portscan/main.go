package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

var (
	host       string
	start_port int
	end_port   int
)

func main() {
	pathOperation()
	getMacAddress()
	user_input()

}

func check_port(host string, start_port, end_port int) {

	for i := start_port; i <= end_port; i++ {
		fmt.Println(i)
		qualified_host := fmt.Sprintf("%s%s%d", host, ":", i)
		conn, err := net.DialTimeout("tcp", qualified_host, 1*time.Second) // Got the timeout code from: https://stackoverflow.com/questions/37294052/golang-why-net-dialtimeout-get-timeout-half-of-the-time
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		status, err := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(status)
	}
}

func pathOperation() {

	execDirAbsPath, _ := os.Getwd()
	log.Println("执行程序所在目录的绝对路径　　　　　　　:", execDirAbsPath)

	execFileRelativePath, _ := exec.LookPath(os.Args[0])
	log.Println("执行程序与命令执行目录的相对路径　　　　:", execFileRelativePath)

	execDirRelativePath, _ := path.Split(execFileRelativePath)
	log.Println("执行程序所在目录与命令执行目录的相对路径:", execDirRelativePath)

	execFileAbsPath, _ := filepath.Abs(execFileRelativePath)
	log.Println("执行程序的绝对路径　　　　　　　　　　　:", execFileAbsPath)

	execDirAbsPath, _ = filepath.Abs(execDirRelativePath)
	log.Println("执行程序所在目录的绝对路径　　　　　　　:", execDirAbsPath)

	os.Chdir(execDirRelativePath) //进入目录
	enteredDirAbsPath, _ := os.Getwd()
	log.Println("所进入目录的绝对路径　　　　　　　　　　:", enteredDirAbsPath)
}

func user_input() {
	fmt.Println("Host> ")
	fmt.Scan(&host)
	fmt.Println("Starting Port (i.e. 80)> ")
	fmt.Scan(&start_port)
	fmt.Println("End Port (i.e. 8080)> ")
	fmt.Scan(&end_port)
	fmt.Println("Running scan... ")
	check_port(host, start_port, end_port)
}

func getMacAddress() {

	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Poor soul, here is what you got: " + err.Error())
	}
	for _, inter := range interfaces {
		fmt.Println(inter.Name, inter.HardwareAddr)
	}

}
