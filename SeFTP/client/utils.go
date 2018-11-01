package main

import (
	"./Controller"
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"golang.org/x/crypto/sha3"
	"gopkg.in/cheggaaa/pb.v2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//Config is the config struct for SeFTP.
type Config struct {
	ServerAddr string
	Passwd     [32]byte
	ServerPort int
}

//Parse is a function to parse flag config to Config struct.
func (config *Config) Parse() {
	serverAddr := flag.String("s", "127.0.0.1", "Server IP Address")
	serverPort := flag.Int("p", 9080, "Server Port")
	plainPasswd := flag.String("k", "WELCOMETOTHEGRID", "Password")
	flag.Parse()

	passwd := GetSHA3Hash(*plainPasswd)

	config.ServerAddr = *serverAddr
	config.ServerPort = *serverPort
	config.Passwd = passwd
}

//GetSHA3Hash is a function to get SHA3 hash of a string.
func GetSHA3Hash(text string) [32]byte {
	return sha3.Sum256([]byte(text))
}

//checkerr is a function to check if there is error.
func checkerr(e error) bool {
	if e != nil {
		log.Println(e)
		return false
	}
	return true
}

//GetOpenPort is a function to get an open TCP port.
func GetOpenPort() (int, error) {
	laddr := net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
	listener, err := net.ListenTCP("tcp4", &laddr)
	if err == nil {
		addr := listener.Addr()
		listener.Close()
		return addr.(*net.TCPAddr).Port, nil
	}
	return 0, err
}

//IsUpper is a function to determine if string uses uppercase.
func IsUpper(str string) bool {
	return str == strings.ToUpper(str)
}

//Ls is a function to handle LS request.
func Ls(path string) []string {
	if path == "" {
		path = "./"
	}
	files, err := ioutil.ReadDir(path)
	checkerr(err)

	var list []string

	for _, f := range files {
		list = append(list, f.Name())
	}
	return list
}

//GET is a function to handle GET request.
func GET(subftpInt interface{}) {
	if subftpCon, ok := subftpInt.(Controller.TCPController); ok {
		subftpCon.EstabConn()
		defer subftpCon.CloseConn()
		subftpCon.SendText("FILE SIZE")
		plainCommand, err := subftpCon.GetText()
		if !checkerr(err) {
			return
		}
		command := strings.Fields(plainCommand)
		switch command[0] {
		case "SIZE":
			fileSize, err := strconv.Atoi(command[1])
			if !checkerr(err) {
				return
			}
			log.Println("FILE SIZE: ", fileSize)
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter File Name: ")
			fileName, _ := reader.ReadString('\n')
			f, err := os.Create(strings.Fields(fileName)[0])
			if !checkerr(err) {
				return
			}
			defer f.Close()
			recvSize := 0
			subftpCon.SendText("READY")

			progressBar := pb.ProgressBarTemplate(`{{bar . | green}} {{speed . | blue }}`).Start(fileSize)

			var exbuf []byte
			var buf []byte
			for recvSize+len(exbuf) < fileSize {
				buf, exbuf, err = subftpCon.GetByte(exbuf)
				checkerr(err)
				recvSize += len(buf)
				progressBar.Add(len(buf))
				//log.Println("RECV BYTE LENGTH: ", len(buf))
				f.Write(buf)
			}
			if recvSize < fileSize {
				lth := exbuf[12:14]
				//log.Println(lth)
				length := binary.LittleEndian.Uint16(lth)
				nonce, exbuf := exbuf[:12], exbuf[14:]
				data, exbuf := exbuf[:length], exbuf[length:]
				decData, err := Controller.GCMDecrypter(data, SeFTPConfig.Passwd, nonce)
				checkerr(err)
				progressBar.Add(len(exbuf))
				f.Write(decData)
			}
			progressBar.Finish()
			log.Println("FILE RECEIVED")
			subftpCon.SendText("HALT")
			return
		}
	} else if subftpCon, ok := subftpInt.(Controller.KCPController); ok {
		subftpCon.EstabConn()
		defer subftpCon.CloseConn()
		subftpCon.SendText("FILE SIZE")
		plainCommand, err := subftpCon.GetText()
		if !checkerr(err) {
			return
		}
		command := strings.Fields(plainCommand)
		switch command[0] {
		case "SIZE":
			fileSize, err := strconv.Atoi(command[1])
			if !checkerr(err) {
				return
			}
			log.Println("FILE SIZE: ", fileSize)
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter File Name: ")
			fileName, _ := reader.ReadString('\n')
			f, err := os.Create(strings.Fields(fileName)[0])
			if !checkerr(err) {
				return
			}
			defer f.Close()
			recvSize := 0
			subftpCon.SendText("READY")

			progressBar := pb.ProgressBarTemplate(`{{bar . | green}} {{speed . | blue }}`).Start(fileSize)

			var exbuf []byte
			var buf []byte
			for recvSize+len(exbuf) < fileSize {
				buf, exbuf, err = subftpCon.GetByte(exbuf)
				checkerr(err)
				recvSize += len(buf)
				progressBar.Add(len(buf))
				//log.Println("RECV BYTE LENGTH: ", len(buf))
				f.Write(buf)
			}
			if recvSize < fileSize {
				lth := exbuf[12:14]
				//log.Println(lth)
				length := binary.LittleEndian.Uint16(lth)
				nonce, exbuf := exbuf[:12], exbuf[14:]
				data, exbuf := exbuf[:length], exbuf[length:]
				decData, err := Controller.GCMDecrypter(data, SeFTPConfig.Passwd, nonce)
				checkerr(err)
				progressBar.Add(len(exbuf))
				f.Write(decData)
			}
			progressBar.Finish()
			log.Println("FILE RECEIVED")
			subftpCon.SendText("HALT")
			return
		}
	}
}

//POST is a function to handle POST request.
func POST(subftpInt interface{}) {
	fmt.Print("Enter File Name: ")
	reader := bufio.NewReader(os.Stdin)
	fileName, _ := reader.ReadString('\n')
	f, err := os.Open(strings.Fields(fileName)[0])
	if !checkerr(err) {
		return
	}
	defer f.Close()
	fileInfo, err := f.Stat()
	if !checkerr(err) {
		return
	}
	fileSize := int(fileInfo.Size())
	if subftpCon, ok := subftpInt.(Controller.TCPController); ok {
		log.Println("TCP Controller")
		subftpCon.EstabConn()
		defer subftpCon.CloseConn()

		log.Println("File size: ", fileSize)

		subftpCon.SendText("SIZE " + strconv.Itoa(fileSize))

		sendSize := 0
		result, err := subftpCon.GetText()
		if !checkerr(err) {
			return
		}
		if result == "READY" {
			log.Println("Server ready")
			progressBar := pb.ProgressBarTemplate(`{{bar . | green}} {{speed . | blue }}`).Start(fileSize)
			for sendSize < fileSize {
				data := make([]byte, 60000)
				n, err := f.Read(data)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Println(err)
					return
				}
				data = data[:n]
				//log.Println("Data:", string(data))
				subftpCon.SendByte(data)
				sendSize += n
				progressBar.Add(n)
				time.Sleep(time.Microsecond)
			}
			progressBar.Finish()
			log.Println("FILE READ COMPLETE")
			result, err = subftpCon.GetText()
			if !checkerr(err) {
				return
			}
			if result == "HALT" {
				log.Println("TRANSFER COMPLETE")
				return
			}
			log.Println("TRANSFER FAILED: ", result)
		}
	} else if subftpCon, ok := subftpInt.(Controller.KCPController); ok {
		log.Println("KCP Controller")
		subftpCon.EstabConn()
		defer subftpCon.CloseConn()

		log.Println("File size: ", fileSize)

		subftpCon.SendText("SIZE " + strconv.Itoa(fileSize))

		sendSize := 0
		result, err := subftpCon.GetText()
		if !checkerr(err) {
			return
		}
		if result == "READY" {
			log.Println("Server ready")
			progressBar := pb.ProgressBarTemplate(`{{bar . | green}} {{speed . | blue }}`).Start(fileSize)
			for sendSize < fileSize {
				data := make([]byte, 60000)
				n, err := f.Read(data)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Println(err)
					return
				}
				data = data[:n]
				//log.Println("Data:", string(data))
				subftpCon.SendByte(data)
				sendSize += n
				progressBar.Add(n)
				time.Sleep(time.Microsecond)
			}
			progressBar.Finish()
			log.Println("FILE READ COMPLETE")
			result, err = subftpCon.GetText()
			if !checkerr(err) {
				return
			}
			if result == "HALT" {
				log.Println("TRANSFER COMPLETE")
				return
			}
			log.Println("TRANSFER FAILED: ", result)
		}
	}
}

//SHA3FileHash is a function to get file's SHA3 hash.
func SHA3FileHash(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := sha3.New256()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}
