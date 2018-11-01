package main

import (
	"./Controller"
	"io"
	//"encoding/hex"
	"log"
	"os"
	"strconv"
	"strings"
	//"bufio"
	"encoding/binary"
	"github.com/xtaci/smux"
	"time"
)

//SeFTPConfig is a config predefined for convenience.
var SeFTPConfig = Config{}

func handleCommand(seftpCon Controller.TCPController, stream *smux.Stream, plainCommand string) {
	command := strings.Fields(plainCommand)
	switch command[0] {
	case "GET":
		if _, err := os.Stat(string(command[1])); !os.IsNotExist(err) {
			subPort, err := GetOpenPort()
			if !checkerr(err) {
				return
			}
			if (len(command) <= 2) || (command[2] == "TCP") {
				subFtpCon := Controller.TCPController{ServerAddr: SeFTPConfig.ServerAddr + ":" + strconv.Itoa(subPort), Passwd: SeFTPConfig.Passwd}
				subFtpCon.EstabListener()
				defer subFtpCon.CloseListener()

				seftpCon.SendText(stream, "PASV TCP "+strconv.Itoa(subPort))
				for {
					// Get net.TCPConn object
					subconn, err := subFtpCon.Listener.Accept()
					if !checkerr(err) {
						continue
					}
					// Setup server side of smux
					session, err := smux.Server(subconn, nil)
					if !checkerr(err) {
						continue
					}

					// Accept a stream
					substream, err := session.AcceptStream()
					if !checkerr(err) {
						continue
					}
					plainEcho, err := subFtpCon.GetText(substream)
					if !checkerr(err) {
						continue
					}
					if plainEcho == "FILE SIZE" {
						f, err := os.Open(string(command[1]))
						if !checkerr(err) {
							continue
						}
						defer f.Close()
						fileInfo, err := f.Stat()
						if !checkerr(err) {
							continue
						}
						fileSize := int(fileInfo.Size())
						subFtpCon.SendText(substream, "SIZE "+strconv.Itoa(fileSize))
						//result, err := subFtpCon.GetText(conn)
						//checkerr(err)
						//if result == "READY" {
						//	log.Println("CLIENT READY")
						sendSize := 0
						result, err := subFtpCon.GetText(substream)
						if !checkerr(err) {
							continue
						}
						if result == "READY" {
							log.Println("CLIENT READY")
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
								subFtpCon.SendByte(substream, data)
								sendSize += n
								time.Sleep(time.Microsecond)
							}
						}
						log.Println("FILE READ COMPLETE")
						result, err = subFtpCon.GetText(substream)
						if !checkerr(err) {
							break
						}
						if result == "HALT" {
							log.Println("TRANSFER COMPLETE")
							break
						} else {
							log.Println("TRANSFER FAILED: ", result)
						}
					} else {
						subFtpCon.SendText(substream, "UNKNOWN COMMAND")
					}
				}
				log.Println("CLOSE SUBCONN")
				return
			} else if command[2] == "UDP" {
				subFtpCon := Controller.KCPController{ServerAddr: ":" + strconv.Itoa(subPort), Passwd: SeFTPConfig.Passwd}
				subFtpCon.EstabListener()
				defer subFtpCon.CloseListener()

				seftpCon.SendText(stream, "PASV UDP "+strconv.Itoa(subPort))
				for {
					subconn, err := subFtpCon.Listener.AcceptKCP()
					if !checkerr(err) {
						continue
					}
					subconn.SetStreamMode(true)
					subconn.SetWriteDelay(true)
					subconn.SetNoDelay(0, 40, 2, 1)
					subconn.SetWindowSize(1024, 1024)
					subconn.SetMtu(1350)
					// Setup server side of smux
					session, err := smux.Server(subconn, nil)
					if !checkerr(err) {
						continue
					}

					// Accept a stream
					substream, err := session.AcceptStream()
					if !checkerr(err) {
						continue
					}
					plainEcho, err := subFtpCon.GetText(substream)
					if !checkerr(err) {
						continue
					}
					if plainEcho == "FILE SIZE" {
						f, err := os.Open(string(command[1]))
						if !checkerr(err) {
							continue
						}
						defer f.Close()
						fileInfo, err := f.Stat()
						if !checkerr(err) {
							continue
						}
						fileSize := int(fileInfo.Size())
						subFtpCon.SendText(substream, "SIZE "+strconv.Itoa(fileSize))
						//result, err := subFtpCon.GetText(conn)
						//checkerr(err)
						//if result == "READY" {
						//	log.Println("CLIENT READY")
						sendSize := 0
						result, err := subFtpCon.GetText(substream)
						if !checkerr(err) {
							continue
						}
						if result == "READY" {
							log.Println("CLIENT READY")
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
								subFtpCon.SendByte(substream, data)
								sendSize += n
								time.Sleep(time.Microsecond)
							}
						}
						log.Println("FILE READ COMPLETE")
						result, err = subFtpCon.GetText(substream)
						if !checkerr(err) {
							break
						}
						if result == "HALT" {
							log.Println("TRANSFER COMPLETE")
							break
						} else {
							log.Println("TRANSFER FAILED: ", result)
						}
					} else {
						subFtpCon.SendText(substream, "UNKNOWN COMMAND")
					}
				}
				log.Println("CLOSE SUBCONN")
				return
			}
		} else {
			seftpCon.SendText(stream, "FILE NOT EXIST")
		}
		seftpCon.SendText(stream, "")
	case "POST":
		subPort, err := GetOpenPort()
		if !checkerr(err) {
			return
		}
		filePath := string(command[1])
		if (len(command) <= 2) || (command[2] == "TCP") {
			subFtpCon := Controller.TCPController{ServerAddr: SeFTPConfig.ServerAddr + ":" + strconv.Itoa(subPort), Passwd: SeFTPConfig.Passwd}
			subFtpCon.EstabListener()
			defer subFtpCon.CloseListener()

			seftpCon.SendText(stream, "PASV TCP "+strconv.Itoa(subPort))
			for {
				// Get net.TCPConn object
				subconn, err := subFtpCon.Listener.Accept()
				if !checkerr(err) {
					continue
				}
				// Setup server side of smux
				session, err := smux.Server(subconn, nil)
				if !checkerr(err) {
					continue
				}

				// Accept a stream
				substream, err := session.AcceptStream()
				if !checkerr(err) {
					continue
				}
				log.Println("ACCEPT SUBSTREAM")
				plainEcho, err := subFtpCon.GetText(substream)
				if !checkerr(err) {
					continue
				}
				echo := strings.Fields(plainEcho)
				log.Println("ECHO: ", plainEcho)
				if (echo[0] != "SIZE") || (len(echo) != 2) {
					continue
				}
				fileSize, err := strconv.Atoi(echo[1])
				if !checkerr(err) {
					continue
				}
				f, err := os.Create(strings.Fields(filePath)[0])
				if !checkerr(err) {
					continue
				}
				defer f.Close()
				recvSize := 0
				subFtpCon.SendText(substream, "READY")

				var exbuf []byte
				var buf []byte
				for recvSize+len(exbuf) < fileSize {
					buf, exbuf, err = subFtpCon.GetByte(exbuf, substream)
					checkerr(err)
					recvSize += len(buf)
					//log.Println("RECV BYTE LENGTH: ", len(buf))
					f.Write(buf)
				}
				if recvSize < fileSize {
					lth := exbuf[12:14]
					//log.Println(lth)
					length := binary.LittleEndian.Uint16(lth)
					nonce, exbuf := exbuf[:12], exbuf[14:]
					data, _ := exbuf[:length], exbuf[length:]
					decData, err := Controller.GCMDecrypter(data, SeFTPConfig.Passwd, nonce)
					checkerr(err)
					f.Write(decData)
				}
				log.Println("FILE RECEIVED")
				subFtpCon.SendText(substream, "HALT")
				return
			}
		} else if command[2] == "UDP" {
			subFtpCon := Controller.KCPController{ServerAddr: ":" + strconv.Itoa(subPort), Passwd: SeFTPConfig.Passwd}
			subFtpCon.EstabListener()
			defer subFtpCon.CloseListener()

			seftpCon.SendText(stream, "PASV UDP "+strconv.Itoa(subPort))
			for {
				// Get net.TCPConn object
				subconn, err := subFtpCon.Listener.Accept()
				if !checkerr(err) {
					continue
				}
				// Setup server side of smux
				session, err := smux.Server(subconn, nil)
				if !checkerr(err) {
					continue
				}

				// Accept a stream
				substream, err := session.AcceptStream()
				if !checkerr(err) {
					continue
				}
				log.Println("ACCEPT SUBSTREAM")
				plainEcho, err := subFtpCon.GetText(substream)
				if !checkerr(err) {
					continue
				}
				echo := strings.Fields(plainEcho)
				log.Println("ECHO: ", plainEcho)
				if (echo[0] != "SIZE") || (len(echo) != 2) {
					continue
				}
				fileSize, err := strconv.Atoi(echo[1])
				if !checkerr(err) {
					continue
				}
				f, err := os.Create(strings.Fields(filePath)[0])
				if !checkerr(err) {
					continue
				}
				defer f.Close()
				recvSize := 0
				subFtpCon.SendText(substream, "READY")

				var exbuf []byte
				var buf []byte
				for recvSize+len(exbuf) < fileSize {
					buf, exbuf, err = subFtpCon.GetByte(exbuf, substream)
					checkerr(err)
					recvSize += len(buf)
					//log.Println("RECV BYTE LENGTH: ", len(buf))
					f.Write(buf)
				}
				if recvSize < fileSize {
					lth := exbuf[12:14]
					//log.Println(lth)
					length := binary.LittleEndian.Uint16(lth)
					nonce, exbuf := exbuf[:12], exbuf[14:]
					data, _ := exbuf[:length], exbuf[length:]
					decData, err := Controller.GCMDecrypter(data, SeFTPConfig.Passwd, nonce)
					checkerr(err)
					f.Write(decData)
				}
				log.Println("FILE RECEIVED")
				subFtpCon.SendText(substream, "HALT")
				return
			}
		}

	case "CD":
		newPath := command[1]
		err := os.Chdir(newPath)
		if !checkerr(err) {
			seftpCon.SendText(stream, "DIR CHANGE FAILED")
		} else {
			seftpCon.SendText(stream, "DIR CHANGED")
		}
	case "LS":
		var list []string
		if len(command) > 1 {
			path := command[1]
			list = Ls(path)
		} else {
			list = Ls("")
		}
		seftpCon.SendText(stream, strings.Join(list, " | "))
	case "RM":
		if len(command) > 1 {
			err := os.Remove(command[1])
			if !checkerr(err) {
				seftpCon.SendText(stream, "RM FAILED")
			} else {
				seftpCon.SendText(stream, "RM SUCCEEDED")
			}
		} else {
			log.Println("NO SPECIFIC FILE")
		}
	case "SHA3SUM":
		if len(command) > 1 {
			sum, err := SHA3FileHash(command[1])
			checkerr(err)
			seftpCon.SendText(stream, sum)
		} else {
			seftpCon.SendText(stream, "No specific file")
		}
	default:
		seftpCon.SendText(stream, "UNKNOWN COMMAND")
	}
}

func handleConnection(seftpCon Controller.TCPController, stream *smux.Stream) {
	log.Println("Handling new connection...")

	// Close connection when this function ends
	defer func() {
		log.Println("Closing connection...")
		stream.Close()
	}()

	for {
		text, rErr := seftpCon.GetText(stream)

		if rErr == nil {
			log.Println("Got Command:", text)
			handleCommand(seftpCon, stream, text)
			continue
		}

		if rErr == io.EOF {
			log.Println("END OF LINE.")

			break
		}
		break
	}
}

func main() {
	SeFTPConfig.Parse()
	seftpCon := Controller.TCPController{ServerAddr: SeFTPConfig.ServerAddr + ":" + strconv.Itoa(SeFTPConfig.ServerPort), Passwd: SeFTPConfig.Passwd}
	seftpCon.EstabListener()

	defer seftpCon.CloseListener()

	for {
		// Get net.TCPConn object
		conn, err := seftpCon.Listener.Accept()
		if !checkerr(err) {
			continue
		}
		// Setup server side of smux
		session, err := smux.Server(conn, nil)
		if !checkerr(err) {
			continue
		}

		// Accept a stream
		stream, err := session.AcceptStream()
		if !checkerr(err) {
			continue
		}

		go handleConnection(seftpCon, stream)
	}
}
