package ship

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const cmiReadWriteTimeout = 10 * time.Second

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type Connection struct {
	conn *websocket.Conn
	Log  Logger
}

func New(conn *websocket.Conn) (c *Connection) {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) log() Logger {
	return c.Log
}

// Connect performs the connection handshake
func (c *Connection) Connect() error {
	c.log().Println("connect")

	err := c.handshake()
	if err == nil {
		err = c.hello()
	}
	if err == nil {
		err = c.protocolHandshake()
	}

	// close connection if handshake or hello fails
	if err != nil {
		_ = c.conn.Close()
	}

	return err
}

func (c *Connection) writeBinary(msg []byte) error {
	c.log().Println("send:", string(msg))

	err := c.conn.SetWriteDeadline(time.Now().Add(cmiReadWriteTimeout))
	if err == nil {
		c.conn.WriteMessage(websocket.BinaryMessage, msg)
	}
	return err
}

func (c *Connection) writeJSON(jsonMsg interface{}) error {
	msg, err := json.Marshal(jsonMsg)
	if err != nil {
		return err
	}

	q := []byte(strconv.Quote(string(msg)))

	return c.writeBinary(q)
}

func (c *Connection) readBinary() ([]byte, error) {
	err := c.conn.SetReadDeadline(time.Now().Add(cmiReadWriteTimeout))
	if err != nil {
		return nil, err
	}

	typ, msg, err := c.conn.ReadMessage()

	if err == nil {
		c.log().Println("recv:", string(msg))

		if typ != websocket.BinaryMessage {
			err = fmt.Errorf("invalid message type: %d", typ)
		}
	}

	return msg, err
}

func (c *Connection) readJSON(jsonMsg interface{}) error {
	msg, err := c.readBinary()
	if err == nil {
		var q string
		q, err = strconv.Unquote(string(msg))

		if err == nil {
			qq := []byte(q)
			err = json.Unmarshal(qq, &jsonMsg)
		}
	}

	return err
}

// Close closes the service connection
func (c *Connection) Close() error {
	err := c.close()
	_ = c.conn.Close()
	return err
}
