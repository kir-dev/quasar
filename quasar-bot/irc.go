package irc

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

// server and nick must be set; password, user, mode and realname are optional
type IRCConnection struct {
	socket net.Conn
	server string

	password string
	nick     string
	user     string
	mode     int
	realname string

	current_nick string

	write chan string
	err   chan error

	reader_stopped, writer_stopped chan struct{}

	timeout time.Duration
}

// To be used as a goroutine
func (con *IRCConnection) readLoop() {
	br := bufio.NewReader(con.socket)
	for {
		// Most servers send a PING message every 3-4 minutes
		con.socket.SetDeadline(time.Now().Add(con.timeout))
		line, err := br.ReadString('\n')
		if err != nil {
			con.err <- err
			break
		}
		// Remove crlf
		line = line[:len(line)-2]
		fmt.Println(line)
	}
	con.reader_stopped <- struct{}{}
}

// To be used as a goroutine
func (con *IRCConnection) writeLoop() {
	for {
		msg := <-con.write
		if len(msg) == 0 {
			break
		}
		_, err := con.socket.Write([]byte(msg))
		if err != nil {
			con.err <- err
			break
		}
	}
	con.writer_stopped <- struct{}{}
}

// Main loop; to be called after connecting
func (con *IRCConnection) Loop() {
	for {
		err := <-con.err
		if err != nil {
			log.Println(err)
			break
		}
	}
	con.Disconnect()
}

func (con *IRCConnection) Register() {
	if len(con.password) > 0 {
		con.write <- fmt.Sprintf("PASS %s\r\n", con.password)
	}
	con.write <- fmt.Sprintf("NICK %s\r\n", con.nick)
	con.write <- fmt.Sprintf("USER %s %s * :%s\r\n", con.user, con.mode, con.realname)
}

func (con *IRCConnection) Connect() error {
	if len(con.server) == 0 {
		return errors.New("empty server")
	}
	if len(con.nick) == 0 {
		return errors.New("empty nick")
	}
	if len(con.user) == 0 {
		con.user = con.nick
	}
	if len(con.realname) == 0 {
		con.realname = con.nick
	}
	if con.timeout == 0 {
		con.timeout = 5 * time.Minute
	}

	var err error
	con.socket, err = net.Dial("tcp", con.server)
	if err != nil {
		return err
	}

	con.write = make(chan string)
	con.err = make(chan error)
	con.reader_stopped = make(chan struct{})
	con.writer_stopped = make(chan struct{})
	go con.writeLoop()
	go con.readLoop()

	con.Register()

	return nil
}

func (con *IRCConnection) Disconnect() {
	con.socket.Write([]byte("QUIT\r\n"))
	close(con.write)
	con.socket.Close()
	<-con.reader_stopped
	<-con.writer_stopped
	close(con.err)
}
