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
	Listener   net.Listener
	Passwd     [32]byte
}

//EstabListener is a function to establish a listener.
func (tcpCon *TCPController) EstabListener() {
	ln, err := net.Listen("tcp", tcpCon.ServerAddr)
	if err != nil {
		log.Println(err)
	}
	tcpCon.Listener = ln
	log.Println("Listener Established.")
}

//CloseListener is a function to close TCP listener.
func (tcpCon *TCPController) CloseListener() {
	tcpCon.Listener.Close()
	log.Println("Listener closed.")
}

//SendByte is a function to send byte though TCP socket.
func (tcpCon *TCPController) SendByte(stream *smux.Stream, data []byte) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err.Error())
	}

	encByte := GCMEncrypter(data, tcpCon.Passwd, nonce)
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(encByte)))
	finalPac := append(append(nonce, bs...), encByte...)
	stream.Write(finalPac)
}

//GetByte is a function to get byte though TCP socket.
func (tcpCon *TCPController) GetByte(exbuf []byte, stream *smux.Stream) ([]byte, []byte, error) {
	//log.Println("ExBUF: ", exbuf)
	tcpbuf := make([]byte, 65550)
	rLen := 0
	n, rErr := stream.Read(tcpbuf)
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
				n, rErr := stream.Read(subbuf)
				subbuf = subbuf[:n]
				if rErr == nil {
					//log.Println("Package Length Received: ", n)
					buf = append(buf, subbuf...)
					rLen += n
				}
				continue
			} else if rLen > int(length)+14 {
				log.Println("RECV EXCESSIVE PACKAGE")
				break
			} else {
				log.Println("RECV PACKAGE COMPLETE")
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
func (tcpCon *TCPController) SendText(stream *smux.Stream, text string) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err.Error())
	}

	encByte := GCMEncrypter([]byte(text), tcpCon.Passwd, nonce)
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(encByte)))
	finalPac := append(append(nonce, bs...), encByte...)
	stream.Write(finalPac)
}

//GetText is a function to get text though TCP socket.
func (tcpCon *TCPController) GetText(stream *smux.Stream) (string, error) {
	buf := make([]byte, 4096)
	_, rErr := stream.Read(buf)

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
	Listener   *kcp.Listener
	Passwd     [32]byte
}

//EstabListener is a function to establish a listener.
func (kcpCon *KCPController) EstabListener() {
	ln, err := kcp.ListenWithOptions(kcpCon.ServerAddr, nil, 10, 3)
	if err != nil {
		log.Println(err)
	}
	kcpCon.Listener = ln
	log.Println("Listener Established.")
}

//CloseListener is a function to close KCP connection.
func (kcpCon *KCPController) CloseListener() {
	kcpCon.Listener.Close()
	log.Println("Listener closed.")
}

//SendByte is a function to send byte though KCP socket.
func (kcpCon *KCPController) SendByte(stream *smux.Stream, data []byte) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err.Error())
	}

	encByte := GCMEncrypter(data, kcpCon.Passwd, nonce)
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(encByte)))
	finalPac := append(append(nonce, bs...), encByte...)
	stream.Write(finalPac)
}

//GetByte is a function to get byte though KCP socket.
func (kcpCon *KCPController) GetByte(exbuf []byte, stream *smux.Stream) ([]byte, []byte, error) {
	//log.Println("ExBUF: ", exbuf)
	tcpbuf := make([]byte, 65550)
	rLen := 0
	n, rErr := stream.Read(tcpbuf)
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
				n, rErr := stream.Read(subbuf)
				subbuf = subbuf[:n]
				if rErr == nil {
					//log.Println("Package Length Received: ", n)
					buf = append(buf, subbuf...)
					rLen += n
				}
				continue
			} else if rLen > int(length)+14 {
				log.Println("RECV EXCESSIVE PACKAGE")
				break
			} else {
				log.Println("RECV PACKAGE COMPLETE")
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
func (kcpCon *KCPController) SendText(stream *smux.Stream, text string) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err.Error())
	}

	encByte := GCMEncrypter([]byte(text), kcpCon.Passwd, nonce)
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(encByte)))
	finalPac := append(append(nonce, bs...), encByte...)
	stream.Write(finalPac)
}

//GetText is a function to get text though KCP socket.
func (kcpCon *KCPController) GetText(stream *smux.Stream) (string, error) {
	buf := make([]byte, 4096)
	_, rErr := stream.Read(buf)

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
