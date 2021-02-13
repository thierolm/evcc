package ship

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	CmiTypeEnd      byte = 3
	CmiCloseTimeout      = 100 * time.Millisecond
)

type CmiCloseMsg struct {
	ConnectionClose []ConnectionClose `json:"connectionClose"`
}

const (
	ConnectionCloseReasonUnspecific        = "unspecific"
	ConnectionCloseReasonRemovedConnection = "removedConnection"

	CmiClosePhaseAnnounce = "announce"
	CmiClosePhaseConfirm  = "confirm"
)

type ConnectionClose struct {
	Phase   string `json:"phase"`
	MaxTime int    `json:"maxTime,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

func (c *Connection) close() error {
	// always send READY
	errC := make(chan error)
	go func(errC chan<- error) {
		msg := CmiCloseMsg{
			ConnectionClose: []ConnectionClose{
				{
					Phase:   CmiClosePhaseAnnounce,
					MaxTime: int(CmiCloseTimeout / time.Millisecond),
				},
			},
		}

		if err := c.writeJSON(CmiTypeEnd, msg); err != nil {
			errC <- fmt.Errorf("close send failed: %w", err)
		}
	}(errC)

	readC := make(chan CmiCloseMsg, 1)
	closeC := make(chan struct{})
	defer close(closeC)

	// read loop
	go func(readC chan<- CmiCloseMsg, closeC chan struct{}, errC chan error) {
		var msg CmiCloseMsg
		for {
			select {
			case <-closeC:
				return
			default:
				typ, err := c.readJSON(&msg)
				if err == nil {
					if typ == CmiTypeEnd {
						readC <- msg
						continue
					} else {
						err = fmt.Errorf("close: invalid type: %0x", typ)
					}
				}

				errC <- fmt.Errorf("close: %w", err)
			}
		}
	}(readC, closeC, errC)

	timer := time.NewTimer(CmiCloseTimeout)
	for {
		select {
		case msg := <-readC:
			log.Printf("close recv: %+v", msg)

			switch msg.ConnectionClose[0].Phase {
			case CmiClosePhaseConfirm:
				return nil
			default:
				return errors.New("invalid close response")
			}
		case err := <-errC:
			return err
		case <-timer.C:
			return errors.New("close timeout")
		}
	}
}
