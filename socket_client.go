package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"log"
)

const channelBufferSize = 1024

var currentId int = 0

type SocketClient struct {
	id             int
	ws             *websocket.Conn
	server         *SocketServer
	messageChannel chan SocketMessage
	doneChannel    chan bool
}

// Create new socket client.
func NewClient(socket *websocket.Conn, server *SocketServer) *SocketClient {

	if socket == nil {
		panic("socket cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	currentId++
	ch := make(chan SocketMessage, channelBufferSize)
	doneCh := make(chan bool)

	return &SocketClient{currentId, socket, server, ch, doneCh}
}

func (c *SocketClient) Write(response SocketMessage) {
	select {
	case c.messageChannel <- response:
	default:
		c.server.Delete(c)
		err := fmt.Errorf("client %d is disconnected.", c.id)
		c.server.Err(err)
	}
}

func (c *SocketClient) Done() {
	c.doneChannel <- true
}

// listen for read/write
func (c *SocketClient) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// listen for requests to write to web client
func (c *SocketClient) listenWrite() {
	for {
		select {
		// send msg to client
		case msg := <-c.messageChannel:
			websocket.JSON.Send(c.ws, msg)

		// done
		case <-c.doneChannel:
			c.server.Delete(c)
			c.doneChannel <- true // for listenRead
			return
		}
	}
}

// listen for requests to read data from web client
func (c *SocketClient) listenRead() {
	log.Println("Listening for read request from web client")
	for {
		select {
		case <-c.doneChannel:
			c.server.Delete(c)
			c.doneChannel <- true // for listenWrite
			return

		default:
			var msg SocketMessage
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneChannel <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				c.server.HandleIncomingRequest(msg)
			}
		}
	}
}
