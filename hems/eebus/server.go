package eebus

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/andig/evcc/hems/eebus/ship"
	"github.com/gorilla/websocket"
)

func NewServer(addr string, cert tls.Certificate) (*http.Server, error) {
	s := &http.Server{
		Addr:    addr,
		Handler: &Handler{},
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.NoClientCert,
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

	// ship
	sc := ship.New(conn)
	sc.Log = log.New(&writer{os.Stdout, "2006/01/02 15:04:05 "}, "[server] ", 0)

	if err := sc.Connect(); err != nil {
		sc.Log.Println(err)
		return
	}

	// go func() {
	// 	b := []byte{0, 0}
	// 	if err := conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
	// 		log.Printf("hello: %v", err)
	// 	}
	// }()

	// p := make([]byte, 1024)
	// for {
	// 	typ, r, err := conn.NextReader()
	// 	if err != nil {
	// 		log.Println("ws nextreader:", err)
	// 		return
	// 	}

	// 	log.Println("ws nextreader:", typ)

	// 	n, err := r.Read(p)
	// 	if err != nil {
	// 		log.Println("ws read:", err)
	// 		return
	// 	}

	// 	var req map[string]interface{}
	// 	if err = json.Unmarshal(p[:n], &req); err == nil {
	// 		log.Printf("ws json: %+v", req)

	// 		if err := conn.WriteMessage(websocket.BinaryMessage, []byte(`{"connectionHello":{"phase":"ready"}}`)); err != nil {
	// 			log.Printf("resp: %v", err)
	// 		}
	// 	} else {
	// 		log.Printf("ws read %d: %0x %s", n, p[:n], string(p[:n]))
	// 	}
	// }

	// listen indefinitely for new messages coming
}
