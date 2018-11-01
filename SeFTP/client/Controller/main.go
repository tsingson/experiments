package Controller

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
	"io"
	"log"
	"net"
)

//TCPController is an interface to control a TCP Dial.
type TCPController struct {
	ServerAddr string
	Conn       net.Conn
	Smux       *smux.Session
	Stream     *smux.Stream
	Passwd     [32]byte
}

//EstabConn is a function to establish a connection with the server.
func (tcpCon *TCPController) EstabConn() {
	conn, err := net.Dial("tcp", tcpCon.ServerAddr)
	if err != nil {
		log.Println(err)
	}
	tcpCon.Conn = conn
	log.Println("Conn Established.")

	session, err := smux.Client(conn, nil)
	if err != nil {
		log.Println(err)
	}
	tcpCon.Smux = session
	log.Println("Smux Established.")

	// Open a new stream
	stream, err := session.OpenStream()
	if err != nil {
		panic(err)
	}
	tcpCon.Stream = stream
	log.Println("Stream Established.")
}

//CloseConn is a function to close TCP connection.
func (tcpCon *TCPController) CloseConn() {
	tcpCon.Stream.Close()
	tcpCon.Smux.Close()
	tcpCon.Conn.Close()
	log.Println("Dial closed.")
}

//SendByte is a function to send byte though TCP socket.
func (tcpCon *TCPController) SendByte(data []byte) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err.Error())
	}

	encByte := GCMEncrypter(data, tcpCon.Passwd, nonce)
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(encByte)))
	finalPac := append(append(nonce, bs...), encByte...)
	tcpCon.Stream.Write(finalPac)
}

//GetByte is a function to get byte though TCP socket.
func (tcpCon *TCPController) GetByte(exbuf []byte) ([]byte, []byte, error) {
	//log.Println("ExBUF: ", exbuf)
	tcpbuf := make([]byte, 65550)
	rLen := 0
	n, rErr := tcpCon.Stream.Read(tcpbuf)
	tcpbuf = tcpbuf[:n]
	buf := append(exbuf, tcpbuf...)
	rLen += n
	rLen += len(exbuf)
	//log.Println("Package Length Received: ", n)

	if rErr == nil {
		lth := buf[12:14]
		//log.Println(lth)
		length := binary.LittleEndian.Uint16(lth)
		//log.Println("Package Length Defined: ", length)
		for {
			if rLen < int(length)+14 {
				subbuf := make([]byte, 65550)
				n, rErr := tcpCon.Stream.Read(subbuf)
				subbuf = subbuf[:n]
				if rErr == nil {
					//log.Println("Package Length Received: ", n)
					buf = append(buf, subbuf...)
					rLen += n
				}
				continue
			} else if rLen > int(length)+14 {
				//log.Println("RECV EXCESSIVE PACKAGE")
				break
			} else {
				//log.Println("RECV PACKAGE COMPLETE")
				break
			}
		}
		//log.Println("BUF: ", buf)
		nonce, buf := buf[:12], buf[14:]
		data, buf := buf[:length], buf[length:]
		decData, err := GCMDecrypter(data, tcpCon.Passwd, nonce)
		if len(buf) == 0 {
			return decData, nil, err
		}
		return decData, buf, err
	}
	return nil, nil, rErr
}

//SendText is a function to send text though TCP socket.
func (tcpCon *TCPController) SendText(text string) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err.Error())
	}

	encByte := GCMEncrypter([]byte(text), tcpCon.Passwd, nonce)
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(encByte)))
	finalPac := append(append(nonce, bs...), encByte...)
	tcpCon.Stream.Write(finalPac)
}

//GetText is a function to get text though TCP socket.
func (tcpCon *TCPController) GetText() (string, error) {
	buf := make([]byte, 4096)
	_, rErr := tcpCon.Stream.Read(buf)

	if rErr == nil {
		nonce, buf := buf[:12], buf[12:]
		lth, buf := buf[:2], buf[2:]
		length := binary.LittleEndian.Uint16(lth)
		data, _ := buf[:length], buf[length:]
		decData, err := GCMDecrypter(data, tcpCon.Passwd, nonce)
		return string(decData), err
	}
	return "", rErr
}

//KCPController is an interface to control a KCP Dial.
type KCPController struct {
	ServerAddr string
	Conn       *kcp.UDPSession
	Smux       *smux.Session
	Stream     *smux.Stream
	Passwd     [32]byte
}

//EstabConn is a function to establish a connection with the server.
func (kcpCon *KCPController) EstabConn() {
	conn, err := kcp.DialWithOptions(kcpCon.ServerAddr, nil, 10, 3)
	if err != nil {
		log.Println(err)
	}
	conn.SetStreamMode(true)
	conn.SetWriteDelay(true)
	conn.SetNoDelay(0, 40, 2, 1)
	conn.SetWindowSize(1024, 1024)
	conn.SetMtu(1350)
	kcpCon.Conn = conn
	log.Println("Listener Established.")
	session, err := smux.Client(conn, nil)
	if err != nil {
		log.Println(err)
	}
	kcpCon.Smux = session
	log.Println("Smux Established.")

	// Open a new stream
	stream, err := session.OpenStream()
	if err != nil {
		panic(err)
	}
	kcpCon.Stream = stream
	log.Println("Stream Established.")
}

//CloseConn is a function to close KCP connection.
func (kcpCon *KCPController) CloseConn() {
	kcpCon.Stream.Close()
	kcpCon.Smux.Close()
	kcpCon.Conn.Close()
	log.Println("Dial closed.")
}

//SendByte is a function to send byte though KCP socket.
func (kcpCon *KCPController) SendByte(data []byte) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err.Error())
	}

	encByte := GCMEncrypter(data, kcpCon.Passwd, nonce)
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(encByte)))
	finalPac := append(append(nonce, bs...), encByte...)
	kcpCon.Stream.Write(finalPac)
}

//GetByte is a function to get byte though KCP socket.
func (kcpCon *KCPController) GetByte(exbuf []byte) ([]byte, []byte, error) {
	//log.Println("ExBUF: ", exbuf)
	tcpbuf := make([]byte, 65550)
	rLen := 0
	n, rErr := kcpCon.Stream.Read(tcpbuf)
	tcpbuf = tcpbuf[:n]
	buf := append(exbuf, tcpbuf...)
	rLen += n
	rLen += len(exbuf)
	//log.Println("Package Length Received: ", n)

	if rErr == nil {
		lth := buf[12:14]
		//log.Println(lth)
		length := binary.LittleEndian.Uint16(lth)
		//log.Println("Package Length Defined: ", length)
		for {
			if rLen < int(length)+14 {
				subbuf := make([]byte, 65550)
				n, rErr := kcpCon.Stream.Read(subbuf)
				subbuf = subbuf[:n]
				if rErr == nil {
					//log.Println("Package Length Received: ", n)
					buf = append(buf, subbuf...)
					rLen += n
				}
				continue
			} else if rLen > int(length)+14 {
				//log.Println("RECV EXCESSIVE PACKAGE")
				break
			} else {
				//log.Println("RECV PACKAGE COMPLETE")
				break
			}
		}
		//log.Println("BUF: ", buf)
		nonce, buf := buf[:12], buf[14:]
		data, buf := buf[:length], buf[length:]
		decData, err := GCMDecrypter(data, kcpCon.Passwd, nonce)
		if len(buf) == 0 {
			return decData, nil, err
		}
		return decData, buf, err
	}
	return nil, nil, rErr
}

//SendText is a function to send text though KCP socket.
func (kcpCon *KCPController) SendText(text string) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err.Error())
	}

	encByte := GCMEncrypter([]byte(text), kcpCon.Passwd, nonce)
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(encByte)))
	finalPac := append(append(nonce, bs...), encByte...)
	kcpCon.Stream.Write(finalPac)
}

//GetText is a function to get text though KCP socket.
func (kcpCon *KCPController) GetText() (string, error) {
	buf := make([]byte, 4096)
	_, rErr := kcpCon.Stream.Read(buf)

	if rErr == nil {
		nonce, buf := buf[:12], buf[12:]
		lth, buf := buf[:2], buf[2:]
		length := binary.LittleEndian.Uint16(lth)
		data, _ := buf[:length], buf[length:]
		decData, err := GCMDecrypter(data, kcpCon.Passwd, nonce)
		return string(decData), err
	}
	return "", rErr
}
