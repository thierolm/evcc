package ship

import (
	"errors"
	"fmt"
)

const (
	CmiTypeControl byte = 1
)

const (
	ProtocolHandshakeFormatJSON = "JSON-UTF8"

	ProtocolHandshakeTypeAnnounceMax = "announceMax"
	ProtocolHandshakeTypeSelect      = "select"

	SubProtocol = "ship"
)

type ProtocolHandshake struct {
	HandshakeType string   `json:"handshakeType"`
	Version       Version  `json:"version"`
	Formats       []string `json:"formats"`
}

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

const (
	CmiProtocolHandshakeErrorUnexpectedMessage = 2
)

type CmiProtocolHandshakeError struct {
	Error int `json:"error"`
}

func (c *Connection) handshakeReceiveSelect() (CmiHandshakeMsg, error) {
	var resp CmiHandshakeMsg
	typ, err := c.readJSON(&resp)

	if err == nil && typ != CmiTypeControl {
		err = fmt.Errorf("handshake: invalid type: %0x", typ)
	}

	if err == nil {
		if len(resp.ProtocolHandshake) != 1 {
			return resp, errors.New("handshake: invalid length")
		}

		handshake := resp.ProtocolHandshake[0]

		if handshake.HandshakeType != ProtocolHandshakeTypeSelect ||
			len(handshake.Formats) != 1 ||
			handshake.Formats[0] != ProtocolHandshakeFormatJSON {
			msg := CmiProtocolHandshakeError{
				Error: CmiProtocolHandshakeErrorUnexpectedMessage,
			}

			_ = c.writeJSON(CmiTypeControl, msg)
			err = errors.New("handshake: invalid response")
		} else {
			// send selection back to server
			err = c.writeJSON(CmiTypeControl, resp)
		}
	}

	return resp, err
}

func (c *Connection) clientProtocolHandshake() error {
	req := CmiHandshakeMsg{
		ProtocolHandshake: []ProtocolHandshake{
			{
				HandshakeType: ProtocolHandshakeTypeAnnounceMax,
				Version:       Version{Major: 1, Minor: 0},
				Formats:       []string{ProtocolHandshakeFormatJSON},
			},
		},
	}
	err := c.writeJSON(CmiTypeControl, req)

	if err == nil {
		_, err = c.handshakeReceiveSelect()
	}

	return err
}

func (c *Connection) serverProtocolHandshake() error {
	var req CmiHandshakeMsg
	typ, err := c.readJSON(&req)

	if err == nil && typ != CmiTypeControl {
		err = fmt.Errorf("handshake: invalid type: %0x", typ)
	}

	if err == nil {
		if len(req.ProtocolHandshake) != 1 {
			return errors.New("handshake: invalid length")
		}

		handshake := req.ProtocolHandshake[0]

		if handshake.HandshakeType != ProtocolHandshakeTypeAnnounceMax ||
			len(handshake.Formats) != 1 ||
			handshake.Formats[0] != ProtocolHandshakeFormatJSON {
			msg := CmiProtocolHandshakeError{
				Error: CmiProtocolHandshakeErrorUnexpectedMessage,
			}

			_ = c.writeJSON(CmiTypeControl, msg)
			err = errors.New("handshake: invalid response")
		} else {
			// send selection back to server
			req.ProtocolHandshake[0].HandshakeType = ProtocolHandshakeTypeSelect
			err = c.writeJSON(CmiTypeControl, req)
		}
	}

	if err == nil {
		_, err = c.handshakeReceiveSelect()
	}

	return err
}
