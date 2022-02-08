package server

import (
	"log"
	"net"
	"syscall"

	"github.com/thetinygoat/kqueue-event-loop/eventloop"
)

type Socket struct {
	fd int
}

type Server struct {
	socket *Socket
}

func NewServer(host string, port int) (*Server, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return nil, err
	}
	socket := &Socket{fd: fd}
	addr := syscall.SockaddrInet4{
		Port: port,
	}
	copy(addr.Addr[:], net.ParseIP(host))
	syscall.Bind(socket.fd, &addr)
	err = syscall.Listen(socket.fd, syscall.SOMAXCONN)

	if err != nil {
		return nil, err
	}

	return &Server{socket: socket}, nil
}

func (s *Server) Listen() {
	loop, err := eventloop.NewEventLoop(s.socket.fd)
	if err != nil {
		log.Fatal(err)
	}
	loop.Start()
}

func (s *Server) Close() error {
	return syscall.Close(s.socket.fd)
}

func (s *Socket) Fd() int {
	return s.fd
}
