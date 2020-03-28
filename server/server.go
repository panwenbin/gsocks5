package server

import (
	"bufio"
	"github.com/panwenbin/gsocks5/consts"
	"github.com/panwenbin/gsocks5/structs"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"
)

var ReverseProxy string

func init() {
	ReverseProxy = os.Getenv("REVERSE_PROXY")
	if ReverseProxy != "" {
		log.Printf("use reverse proxy %s\n", ReverseProxy)
	}
}

// ListenAndServe
func ListenAndServe(network, address string) {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatalln(err)
	}

	Serve(l)
}

// Serve
func Serve(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go HandleConn(conn)
	}
}

// HandleConn
func HandleConn(conn net.Conn) {
	defer conn.Close()

	connReader := bufio.NewReader(conn)
	connWriter := bufio.NewWriter(conn)

	conReq := &structs.ConnectRequest{}
	err := conReq.Read(connReader)
	if err != nil {
		log.Println(err)
		return
	}

	conRes := structs.ConnectResponse{
		Ver:    consts.VERSION5,
		Method: consts.METHOD_NO_AUTHENTICATION_REQUIRED,
	}
	err = conRes.Write(connWriter)
	if err != nil {
		log.Println(err)
		return
	}

	cmdReq := &structs.CmdRequest{}
	err = cmdReq.Read(connReader)
	if err != nil {
		log.Println(err)
		return
	}

	cmdRes := structs.CmdResponse{}
	cmdRes.Ver = consts.VERSION5
	cmdRes.Rep = consts.REP_SUCCEEDED
	cmdRes.Rsv = consts.RSV
	cmdRes.Bnd.Atyp = consts.ATYP_IPv4

	// bind and udp associate are not supported
	if cmdReq.Cmd != consts.CMD_CONNECT {
		cmdRes.Rep = consts.REP_COMMAND_NOT_SUPPORTED
		_ = cmdRes.Write(connWriter)
		return
	}

	host := string(cmdReq.Dst.Domain.Bytes)
	// reverse proxy
	if ReverseProxy != "" {
		host = ReverseProxy
	}
	port := (int(cmdReq.Dst.Port[0]) << 8) | int(cmdReq.Dst.Port[1])
	remote, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), time.Second*3)
	if err != nil {
		switch err.(type) {
		case *net.OpError:
			err2 := err.(*net.OpError).Err
			switch err2.(type) {
			case *net.DNSError:
				cmdRes.Rep = consts.REP_HOST_UNREACHABLE
			default:
				if err2, ok := err2.(net.Error); ok && err2.Timeout() {
					cmdRes.Rep = consts.REP_HOST_UNREACHABLE
				} else {
					cmdRes.Rep = consts.REP_GENERAL_SOCKS_SERVER_FAILURE
				}
			}
			log.Println(2, reflect.TypeOf(err2), err2)
		default:
			log.Println(1, reflect.TypeOf(err), err)
			cmdRes.Rep = consts.REP_GENERAL_SOCKS_SERVER_FAILURE
		}
		_ = cmdRes.Write(connWriter)
		return
	}

	err = cmdRes.Write(connWriter)
	if err != nil {
		log.Println(err)
		return
	}

	go io.Copy(remote, conn)
	io.Copy(conn, remote)
}
