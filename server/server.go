package server

import (
	"log"
	"net"
	"syscall"

	"github.com/thetinygoat/kqueue-event-loop/eventloop"
)

// Socket wraps the file descriptor and provides useful methods
// we can also implement reader and writer interfaces
type Socket struct {
	fd int
}

type Server struct {
	socket *Socket
}

func NewServer(host string, port int) (*Server, error) {
	// syscall.Socket creates a new socket and returns a file descriptor
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return nil, err
	}

	// bind socket to the address
	socket := &Socket{fd: fd}
	addr := syscall.SockaddrInet4{
		Port: port,
	}
	copy(addr.Addr[:], net.ParseIP(host))
	syscall.Bind(socket.fd, &addr)

	// syscall.Listen marks that the socket will be used for accepting new connections
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
