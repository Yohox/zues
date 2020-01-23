package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"unsafe"
)

type ClientRequestMessage struct {
	conn net.Conn
	Len uint32
	ShakeStep uint8 //握手情况0-3
	Body []byte
}

type ClientResponseMessage struct {
	Len uint32
	Body []byte
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}


func (p *ClientRequestMessage) readHead() error {
	bytes := make([]byte, 5)
	_, err := p.conn.Read(bytes)
	if err != nil {
		return err
	}
	p.Len = binary.LittleEndian.Uint32(bytes[0:4])
	p.ShakeStep = bytes[4]
	return nil
}

func (p *ClientRequestMessage) readBody() error {
	bytes := make([]byte, p.Len)
	_, err := p.conn.Read(bytes)
	if err != nil {
		return err
	}
	p.Body = bytes
	return nil
}

func(p *ClientRequestMessage) judgeHead() error {
	bytes := make([]byte, 5)
	_, err := p.conn.Read(bytes)
	if err != nil {
		return err
	}
	str := string(bytes)
	if str != "head|" {
		return errors.New(fmt.Sprintf("判定消息头失败！%s", str))
	} else {
		return nil
	}
}

func ReadClientRequestMessage(conn net.Conn) (*ClientRequestMessage, error) {
	clientRequestMessage := &ClientRequestMessage{conn: conn}
	if err := clientRequestMessage.judgeHead(); err != nil {
		return nil, err
	}
	if err := clientRequestMessage.readHead(); err != nil {
		return nil, err
	}
	if err := clientRequestMessage.readBody(); err != nil {
		return nil, err
	}
	return clientRequestMessage, nil
}


func GetClientResponseMessage(clientResponseMessage *ClientResponseMessage) []byte {
	Len := unsafe.Sizeof(*clientResponseMessage)
	bytes := &SliceMock{
		addr: uintptr(unsafe.Pointer(clientResponseMessage)),
		len:  int(Len),
		cap:  int(Len),
	}
	return *(*[]byte)(unsafe.Pointer(bytes))
}