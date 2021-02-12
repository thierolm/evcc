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
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	// upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	p := make([]byte, 1024)
	for {
		typ, r, err := conn.NextReader()
		if err != nil {
			log.Println("ws nextreader:", err)
			return
		}

		log.Println("ws nextreader:", typ)

		n, err := r.Read(p)
		if err != nil {
			log.Println("ws read:", err)
			return
		}

		log.Printf("ws read: %02x %s", p[:n], string(p[:n]))
	}

	// listen indefinitely for new messages coming
}
