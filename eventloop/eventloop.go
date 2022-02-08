package eventloop

import (
	"errors"
	"syscall"
)

type EventLoop struct {
	kqueueFd int
	sockFd   int
}

func NewEventLoop(sockFd int) (*EventLoop, error) {
	kqueueFd, err := syscall.Kqueue()
	if err != nil {
		return nil, err
	}

	loop := &EventLoop{kqueueFd: kqueueFd, sockFd: sockFd}

	socketEvent := syscall.Kevent_t{
		Ident:  uint64(loop.sockFd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
	}

	r, err := syscall.Kevent(loop.kqueueFd, []syscall.Kevent_t{socketEvent}, nil, nil)

	if err != nil {
		return nil, err
	}

	if r == -1 {
		return nil, errors.New("failed to register socket with kqueue")
	}

	return loop, nil
}

func (e *EventLoop) Start() {
	for {
		events := make([]syscall.Kevent_t, 1)
		numEvents, err := syscall.Kevent(e.kqueueFd, nil, events, nil)

		if err != nil {
			continue
		}
		for i := 0; i < numEvents; i++ {
			event := events[i]
			eventFd := int(event.Ident)

			// we reached eof so close connection
			if event.Flags&syscall.EV_EOF != 0 {
				syscall.Close(eventFd)
			} else if eventFd == e.sockFd {
				// we received an event on the socket,
				// which means a new connection arrived, so we accept the new connection and add it to the kqueue

				// first we create a new fd for the new connection
				sockFd, _, err := syscall.Accept(eventFd)
				if err != nil {
					continue
				}

				// create an event to register with the kqueue
				sockEvent := syscall.Kevent_t{
					Ident:  uint64(sockFd),
					Filter: syscall.EVFILT_READ,
					Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
				}

				r, err := syscall.Kevent(e.kqueueFd, []syscall.Kevent_t{sockEvent}, nil, nil)

				if err != nil || r == -1 {
					continue
				}
			} else if event.Filter&syscall.EVFILT_READ != 0 {
				buf := make([]byte, 1024)

				n, err := syscall.Read(eventFd, buf)
				if err != nil {
					continue
				}
				syscall.Write(eventFd, buf[:n])
			}
		}
	}
}
