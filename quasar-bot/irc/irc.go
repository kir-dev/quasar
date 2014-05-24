package irc

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

// Represents the connection to a single IRC server
type IRCConnection struct {
	socket net.Conn
    // the server address, including the port
	server string

	password string
	nick     string
	user     string
    // 0: visible, 8: invisible
    // (http://tools.ietf.org/html/rfc2812&section-3.1.3)
	mode     int
	realname string

    // the nickname may change if the original one is in use
	currentNick string

	write chan string
	err   chan error

    // to signal when the reader and writer routines are stopped
	readerStopped, writerStopped chan struct{}

    // timeout for the read operation on the socket
	timeout time.Duration
}

// to be used as a goroutine
func (con *IRCConnection) readLoop() {
	br := bufio.NewReader(con.socket)
	for {
		// set the timeout, since servers have to send a ping message at
        // regular intervals
		con.socket.SetDeadline(time.Now().Add(con.timeout))
        // read a message ending with \r\n (including it)
		line, err := br.ReadString('\n')
		if err != nil {
			con.err <- err
			break
		}
		// remove crlf
		line = line[:len(line)-2]
		fmt.Println(line)
	}
	con.reader_stopped <- struct{}{}
}

// to be used as a goroutine
func (con *IRCConnection) writeLoop() {
	for {
		msg := <-con.write
        // end the loop if the channel is closed
		if msg != "" {
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

// main loop; to be called after connecting
func (con *IRCConnection) Loop() {
	for {
		err := <-con.err
		if err != nil {
			log.Println(err)
			break
		}
	}
	con.Disconnect("")
}

func (con *IRCConnection) Register() {
	if con.password != "" {
		con.write <- fmt.Sprintf("PASS %s\r\n", con.password)
	}
	con.write <- fmt.Sprintf("NICK %s\r\n", con.nick)
	con.write <- fmt.Sprintf("USER %s %s * :%s\r\n", con.user, con.mode, con.realname)
}

func (con *IRCConnection) Connect() error {
	if con.server == "" {
		return errors.New("empty server")
	}
	if con.nick == "" {
		return errors.New("empty nick")
	}
	if con.user == "" {
		con.user = con.nick
	}
	if con.realname == "" {
		con.realname = con.nick
	}
	if con.timeout == 0 {
        // most servers send a PING message every 3-4 minutes
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

func (con *IRCConnection) Disconnect(quit_message string) {
	if quit_message != "" {
		con.socket.Write([]byte(fmt.Sprintf("QUIT :%s\r\n", quit_message)))
	}
	else {
		con.socket.Write([]byte("QUIT\r\n"))
	}
	close(con.write)
	con.socket.Close()
    // wait for the reader and writer routines to stop
	<-con.reader_stopped
	<-con.writer_stopped
	close(con.err)
}
