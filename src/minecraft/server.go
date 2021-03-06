package minecraft

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/hyren/nyrah/minecraft/protocol"
	"net/hyren/nyrah/minecraft/protocol/packet"
)

var NoHandlerException = errors.New("No packet handler has been specified")

type Server struct {
	host string
	port int

	listener net.Listener
	handler  Handler
}

type Handler func(*protocol.Connection, packet.Holder) error

func NewServer(host string, port int, handler Handler) *Server {
	return &Server{host: host, port: port, handler: handler}
}

func (server *Server) SetHandler(handler Handler) {
	server.handler = handler
}

func (server *Server) ListenAndServe() (err error) {
	if server.handler == nil {
		return NoHandlerException
	}

	server.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", server.host, server.port))
	if err != nil {
		return
	}

	for {
		client, err := server.listener.Accept()
		if err != nil {
			log.Println("Error occurred while accepting a connection: " + err.Error())
			continue
		}

		go server.handleConnection(protocol.NewConnection(client))
	}
}

func (server *Server) handleConnection(conn *protocol.Connection) {
	for {
		holder, err := conn.Next()

		if holder == nil {
			conn.Close()
			break
		}

		if err != nil {
			if err == protocol.UnknownPacketType {
				continue
			}

			// Not sure what to do about this one..
			log.Println("Error occurred while reading packet: " + err.Error())
			conn.Close()
			break
		}

		if err = server.handler(conn, holder); err != nil {
			log.Println("Error occurred while handling packet: " + err.Error())
			conn.Close()
			break
		}

		if conn.Stop {
			break
		}
	}
}
