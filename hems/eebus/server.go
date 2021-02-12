package eebus

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func NewServer(addr string, cert tls.Certificate) (*http.Server, error) {
	s := &http.Server{
		Addr:    addr,
		Handler: &Handler{},
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	go func() {
		if err := s.ListenAndServeTLS("", ""); err != nil {
			log.Fatal(err)
		}
	}()

	return s, nil
}

type Handler struct{}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return false },
	}

	// upgrade
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("client connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Fatal(err)
	}
	// listen indefinitely for new messages coming
}
