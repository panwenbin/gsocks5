package structs

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/panwenbin/gsocks5/consts"
	"io"
)

type NBytes struct {
	Len   byte
	Bytes []byte
}

type ConnectRequest struct {
	Ver     byte
	Methods NBytes
}

// Read reads a ConnectRequest
//   +----+----------+----------+
//   |VER | NMETHODS |  METHODS |
//   +----+----------+----------+
//   |  1 |     1    | 1 to 255 |
//   +----+----------+----------+
func (cr *ConnectRequest) Read(connReader *bufio.Reader) error {
	var err error
	// check version
	cr.Ver, err = connReader.ReadByte()
	if err != nil {
		return errors.New("read version error, " + err.Error())
	}
	if cr.Ver != consts.VERSION5 {
		return errors.New(fmt.Sprintf("version %d is not accepted", cr.Ver))
	}

	cr.Methods.Len, err = connReader.ReadByte()
	if err != nil {
		return errors.New("read nMethods error, " + err.Error())
	}
	// 读取指定长度字节好像不对
	methods := make([]byte, int(cr.Methods.Len))
	_, err = io.ReadFull(connReader, methods)
	if err != nil {
		return errors.New("read methods error, " + err.Error())
	}
	cr.Methods.Bytes = methods

	return nil
}

type ConnectResponse struct {
	Ver    byte
	Method byte
}

// Write writes a ConnectResponse
//   +----+--------+
//   |VER | METHOD |
//   +----+--------+
//   |  1 |    1   |
//   +----+--------+
func (cr ConnectResponse) Write(connWriter *bufio.Writer) error {

	err := connWriter.WriteByte(cr.Ver)
	if err != nil {
		return err
	}
	err = connWriter.WriteByte(cr.Method)
	if err != nil {
		return err
	}
	return connWriter.Flush()
}

type Addr struct {
	Atyp     byte
	Ipv4Addr [4]byte
	Ipv6Addr [16]byte
	Domain   NBytes
	Port     [2]byte
}

type CmdRequest struct {
	Ver byte
	Cmd byte
	Rsv byte
	Dst Addr
}

// Read reads a CmdRequest
//   +----+-----+-------+------+----------+----------+
//   |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
//   +----+-----+-------+------+----------+----------+
//   | 1  |  1  | X'00' |  1   | Variable |    2     |
//   +----+-----+-------+------+----------+----------+
func (cr *CmdRequest) Read(connReader *bufio.Reader) error {
	var err error
	cr.Ver, err = connReader.ReadByte()
	if err != nil {
		return errors.New("read version error, " + err.Error())
	}
	if cr.Ver != consts.VERSION5 {
		return errors.New(fmt.Sprintf("version %d is not accepted", cr.Ver))
	}

	cr.Cmd, err = connReader.ReadByte()
	if err != nil {
		return errors.New("read cmd error, " + err.Error())
	}
	cr.Rsv, err = connReader.ReadByte()
	if err != nil {
		return errors.New("read rsv error, " + err.Error())
	}
	switch cr.Cmd {
	case consts.CMD_CONNECT:
		cr.Dst.Atyp, err = connReader.ReadByte()
		if err != nil {
			return errors.New("read cmd error, " + err.Error())
		}
		switch cr.Dst.Atyp {
		case consts.ATYP_IPv4:
			_, err := io.ReadFull(connReader, cr.Dst.Ipv4Addr[:])
			if err != nil {
				return err
			}
		case consts.ATYP_IPv6:
			_, err := io.ReadFull(connReader, cr.Dst.Ipv6Addr[:])
			if err != nil {
				return err
			}
		case consts.ATYP_DOMAIN:
			cr.Dst.Domain.Len, err = connReader.ReadByte()
			if err != nil {
				return err
			}
			domainBytes := make([]byte, cr.Dst.Domain.Len)
			_, err := io.ReadFull(connReader, domainBytes)
			if err != nil {
				return err
			}
			cr.Dst.Domain.Bytes = domainBytes
		}
		_, err := io.ReadFull(connReader, cr.Dst.Port[:])
		if err != nil {
			return err
		}
	case consts.CMD_BIND:
	case consts.CMD_UDP_ASSOCIATE:
	default:
	}

	return nil
}

type CmdResponse struct {
	Ver byte
	Rep byte
	Rsv byte
	Bnd Addr
}

// Write writes a CmdResponse
//   +----+-----+-------+------+----------+----------+
//   |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
//   +----+-----+-------+------+----------+----------+
//   | 1  |  1  | X'00' |  1   | Variable |    2     |
//   +----+-----+-------+------+----------+----------+
func (cr CmdResponse) Write(connWriter *bufio.Writer) error {
	var err error
	err = connWriter.WriteByte(cr.Ver)
	err = connWriter.WriteByte(cr.Rep)
	err = connWriter.WriteByte(cr.Rsv)
	err = connWriter.WriteByte(cr.Bnd.Atyp)
	_, err = connWriter.Write(cr.Bnd.Ipv4Addr[:])
	_, err = connWriter.Write(cr.Bnd.Port[:])
	err = connWriter.Flush()
	return err
}
