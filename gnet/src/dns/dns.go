package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type dnsHeader struct {
	Id                                 uint16
	Bits                               uint16
	Qdcount, Ancount, Nscount, Arcount uint16
}

func (header *dnsHeader) SetFlag(QR uint16, OperationCode uint16, AuthoritativeAnswer uint16, Truncation uint16, RecursionDesired uint16, RecursionAvailable uint16, ResponseCode uint16) {
	header.Bits = QR<<15 + OperationCode<<11 + AuthoritativeAnswer<<10 + Truncation<<9 + RecursionDesired<<8 + RecursionAvailable<<7 + ResponseCode
}

type dnsQuery struct {
	QuestionType  uint16
	QuestionClass uint16
}

func ParseDomainName(domain string) []byte {
	var (
		buffer   bytes.Buffer
		segments []string = strings.Split(domain, ".")
	)
	for _, seg := range segments {
		binary.Write(&buffer, binary.BigEndian, byte(len(seg)))
		binary.Write(&buffer, binary.BigEndian, []byte(seg))
	}
	binary.Write(&buffer, binary.BigEndian, byte(0x00))

	return buffer.Bytes()
}

func GetShift(buf []byte, idx int) string {
	var str string
	if buf[idx] == 0x00 {
		return ""
	}
	for buf[idx] != 0 {
		if buf[idx] == 0xc0 {
			str += GetShift(buf, int(buf[idx+1]))
			break
		} else {
			str += string(buf[idx+1:idx+1+int(buf[idx])]) + "."
			idx += int(buf[idx]) + 1
		}
	}

	return str
}

func SendAndRecvMsg(dnsServer, domain string) ([]byte, int, time.Duration) {
	requestHeader := dnsHeader{
		Id:      0x0010,
		Qdcount: 1,
		Ancount: 0,
		Nscount: 0,
		Arcount: 0,
	}
	requestHeader.SetFlag(0, 0, 0, 0, 1, 0, 0)

	requestQuery := dnsQuery{
		QuestionType:  1,
		QuestionClass: 1,
	}

	var (
		conn   net.Conn
		err    error
		buffer bytes.Buffer
	)

	if conn, err = net.Dial("udp", dnsServer); err != nil {
		fmt.Println(err.Error())
		return make([]byte, 0), 0, 0
	}
	defer conn.Close()

	binary.Write(&buffer, binary.BigEndian, requestHeader)
	binary.Write(&buffer, binary.BigEndian, ParseDomainName(domain))
	binary.Write(&buffer, binary.BigEndian, requestQuery)

	buf := make([]byte, 1024)
	t1 := time.Now()
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		fmt.Println(err.Error())
		return make([]byte, 0), 0, 0
	}
	length, err := conn.Read(buf)
	t := time.Now().Sub(t1)
	return buf, length, t
}

func ParseMsg(msg []byte, length int) (cname, ip []string) {
	reponseHeader := dnsHeader{
		Id:      uint16(msg[0])<<8 + uint16(msg[1]),
		Bits:    uint16(msg[2])<<8 + uint16(msg[3]),
		Qdcount: uint16(msg[4])<<8 + uint16(msg[5]),
		Ancount: uint16(msg[6])<<8 + uint16(msg[7]),
		Arcount: uint16(msg[8])<<8 + uint16(msg[9]),
		Nscount: uint16(msg[10])<<8 + uint16(msg[11]),
	}
	idx := 12

	var domain string
	for i := uint16(0); i < reponseHeader.Qdcount; i++ {
		for msg[idx] != 0 {
			domain += string(msg[idx+1:idx+1+int(msg[idx])]) + "."
			idx += int(msg[idx]) + 1
		}
		idx++
	}
	domain = domain[0 : len(domain)-1]
	idx += 4

	for i := uint16(0); i < reponseHeader.Ancount; i++ {

		TYPE := msg[idx+3]
		idx += 11
		length := int(msg[idx])
		var str string
		if TYPE == 5 {
			idx += 1
			for l := 0; l+2 < length; {
				l += int(msg[idx]) + 1
				str += string(msg[idx+1:idx+1+int(msg[idx])]) + "."
				idx += int(msg[idx]) + 1
			}
			str += GetShift(msg, idx)
			cname = append(cname, str)
			if msg[idx] == 0x00 {
				idx++
			} else {
				idx += 2
			}
		} else if TYPE == 1 {
			for l := 0; l < length; l++ {
				idx++
				str += strconv.Itoa(int(msg[idx])) + "."
			}
			idx++
			ip = append(ip, str[:len(str)-1])
		}
	}
	return cname, ip
}

func LookUpHost(dnsServer, domain string) (time.Duration, []string, []string) {

	buf, length, t := SendAndRecvMsg(dnsServer, domain)

	if length == 0 {
		return t, make([]string, 0), make([]string, 0)
	}

	cname, ip := ParseMsg(buf, length)

	return t, cname, ip
}
