package client

import (
	"goChat/chat"
	"goChat/config"
	"goChat/usersession"
	"log"

	"github.com/gorilla/websocket"
)

var clients []*Client

type Client struct {
	conn *websocket.Conn
	send chan *chat.Message

	roomId string
	user   *usersession.UserSession
}

func New(conn *websocket.Conn, roomId string, u *usersession.UserSession) {
	c := &Client{
		conn:   conn,
		send:   make(chan *chat.Message, config.Socket().MsgChanelSize),
		roomId: roomId,
		user:   u,
	}

	clients = append(clients, c)

	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) Close() {
	for i, client := range clients {
		if client == c {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}

		close(c.send)
		c.conn.Close()
		log.Printf("close connection. addr: $s", c.conn.RemoteAddr())
	}
}

func (c *Client) read() (*chat.Message, error) {
	var msg *chat.Message

	if err := c.conn.ReadJSON(&msg); err != nil {
		return nil, err
	}

	msg.User = c.user

	log.Println("read from websocket: ", msg)
	return msg, nil
}

func broadcast(m *chat.Message) {
	for _, client := range clients {
		client.send <- m
	}
}

func (c *Client) readLoop() {
	for {
		m, err := c.read()
		if err != nil {
			log.Println("read message error: ", err)
			break
		}

		m.Create()
		broadcast(m)
	}
	c.Close()
}

func (c *Client) write(m *chat.Message) error {
	log.Println("write to websocket: ", m)

	return c.conn.WriteJSON(m)
}

func (c *Client) writeLoop() {
	for msg := range c.send {
		if c.roomId == msg.RoomID.Hex() {
			c.write(msg)
		}
	}
}
