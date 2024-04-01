package server

import (
	"bufio"
	"encoding/json"
	"flotta-home/mindbond/websocket-server/pkg/client"
	"fmt"
	"github.com/ystepanoff/gowest"
	"net"
	"net/http"
)

type Server struct {
	Port        int
	AuthClient  client.AuthServiceClient
	ChatClient  client.ChatServiceClient
	RWBuffers   map[string]*bufio.ReadWriter
	RemoteAddrs map[int64]string
}

type inputMessage struct {
	UserId  int64           `json:"userId"`
	Token   string          `json:"token"`
	Request string          `json:"request"`
	Data    json.RawMessage `json:"data"`
}

type sendRequest struct {
	ContactId int64  `json:"contactId"`
	Message   string `json:"message"`
}

func (s *Server) Start() {
	s.RWBuffers = make(map[string]*bufio.ReadWriter)
	s.RemoteAddrs = make(map[int64]string)
	http.HandleFunc("/", s.wsHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil); err != nil {
		panic(err)
	}
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, bufrw, err := gowest.GetConnection(w, r)
	fmt.Println("New Connection", conn.RemoteAddr().String())
	s.RWBuffers[conn.RemoteAddr().String()] = bufrw
	if err != nil {
		panic(err)
	}
	defer s.closeConnection(conn)
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
			fmt.Printf("failed to unmarshal: %v", err)
			continue
		}
		validationResponse, err := s.AuthClient.Validate(data.Token)
		if err != nil {
			fmt.Printf("token validation failed: %v", err)
			continue
		}
		if validationResponse.Status != http.StatusOK {
			fmt.Printf("token validation failed with status %v: %v", validationResponse.Status, validationResponse.Error)
			continue
		}
		switch req := data.Request; req {
		case "init":
			fmt.Println("Received init request from user", data.UserId)
			s.RemoteAddrs[data.UserId] = conn.RemoteAddr().String()
		case "sendMessage":
			fmt.Println("Received send message request")
			sendReq := sendRequest{}
			if err = json.Unmarshal(data.Data, &sendReq); err != nil {
				fmt.Printf("failed to unmarshal send request: %v", err)
				continue
			}
			fmt.Println(sendReq.Message, sendReq.ContactId)
		}

		responseMessage := fmt.Sprintf("You sent me %s!", message)
		if err := gowest.WriteString(bufrw, []byte(responseMessage)); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Server) closeConnection(conn net.Conn) {
	s.RWBuffers[conn.RemoteAddr().String()] = nil
	if err := conn.Close(); err != nil {
		panic(err)
	}
}
