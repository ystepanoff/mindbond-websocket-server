package server

import (
	"flotta-home/mindbond/websocket-server/pkg/client"
	"fmt"
	"github.com/ystepanoff/gowest"
	"net"
	"net/http"
)

type Server struct {
	Port       int
	AuthClient client.AuthServiceClient
}

func (s *Server) Start() {
	http.HandleFunc("/", s.wsHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil); err != nil {
		panic(err)
	}
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, bufrw, err := gowest.GetConnection(w, r)
	fmt.Println("New Connection", conn.RemoteAddr().String())
	if err != nil {
		panic(err)
	}
	defer closeConnection(conn)
	for {
		msg, err := gowest.Read(bufrw)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed by client", conn.RemoteAddr().String())
				break
			}
			fmt.Println(err)
			continue
		}
		message := string(msg)
		fmt.Println(message)
		validationResponse, err := s.AuthClient.Validate(message)
		fmt.Println(validationResponse, err)
		responseMessage := fmt.Sprintf("You sent me %s!", message)
		if err := gowest.WriteString(bufrw, []byte(responseMessage)); err != nil {
			fmt.Println(err)
		}
	}
}

func closeConnection(conn net.Conn) {
	if err := conn.Close(); err != nil {
		panic(err)
	}
}
