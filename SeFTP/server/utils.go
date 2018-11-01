package main

import (
	"encoding/hex"
	"flag"
	"golang.org/x/crypto/sha3"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
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

//SHA3FileHash is a function to get file's SHA3 hash.
func SHA3FileHash(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	} //Get is a function to handle GET request.

	defer file.Close()

	hash := sha3.New256()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}
