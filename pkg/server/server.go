package server

import (
	"encoding/json"
	"flotta-home/mindbond/websocket-server/pkg/client"
	"fmt"
	"github.com/ystepanoff/gowest"
	"net"
	"net/http"
)

type Server struct {
	Port       int
	AuthClient client.AuthServiceClient
	ChatClient client.ChatServiceClient
	Pool       map[int]chan interface{}
}

type inputMessage struct {
	UserId  int             `json:"userId"`
	Token   string          `json:"token"`
	Request string          `json:"request"`
	Data    json.RawMessage `json:"data"`
}

type sendRequest struct {
	ContactId int    `json:"contactId"`
	Message   string `json:"message"`
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
		data := inputMessage{}
		if err = json.Unmarshal(msg, &data); err != nil {
			fmt.Errorf("failed to unmarshal: %v", err)
			continue
		}
		validationResponse, err := s.AuthClient.Validate(data.Token)
		if validationResponse.Status != http.StatusOK {
			fmt.Errorf("token validation failed with status %v: %v", validationResponse.Status, validationResponse.Error)
		}
		switch req := data.Request; req {
		case "init":
			fmt.Println("Received init request")
		case "sendMessage":
			fmt.Println("Received send message request")
		}

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
