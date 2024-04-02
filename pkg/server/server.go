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

func (s *Server) initClient(userId int64, conn net.Conn) {
	s.RemoteAddrs[userId] = conn.RemoteAddr().String()
}

func (s *Server) processMessage(userId int64, token string, data json.RawMessage) (*sendRequest, error) {
	sendReq := sendRequest{}
	if err := json.Unmarshal(data, &sendReq); err != nil {
		return nil, err
	}
	addMessageResponse, err := s.ChatClient.AddMessage(userId, sendReq.ContactId, sendReq.Message, token)
	if err != nil {
		return nil, err
	}
	fmt.Println(sendReq.Message, sendReq.ContactId, addMessageResponse)
	return &sendReq, nil
}

func (s *Server) notifyClient(userId int64, response string) error {
	remoteAddr, ok := s.RemoteAddrs[userId]
	if !ok {
		return fmt.Errorf("User %v not found among existing connections", userId)
	}
	bufrw, ok := s.RWBuffers[remoteAddr]
	if !ok {
		return fmt.Errorf("RW buffer does not exist for %s", remoteAddr)
	}
	if err := gowest.WriteString(bufrw, []byte(response)); err != nil {
		return err
	}
	return nil
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, bufrw, err := gowest.GetConnection(w, r)
	if err != nil {
		panic(err)
	}
	fmt.Println("New Connection", conn.RemoteAddr().String())
	s.RWBuffers[conn.RemoteAddr().String()] = bufrw
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
			fmt.Printf("failed to unmarshal: %d %v", len(msg), err)
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
			s.initClient(data.UserId, conn)
		case "sendMessage":
			fmt.Println("Received send message request")
			sendReq, err := s.processMessage(data.UserId, data.Token, data.Data)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if err := s.notifyClient(data.UserId, "hello"); err != nil {
				fmt.Println(err)
			}
			if err := s.notifyClient(sendReq.ContactId, "World"); err != nil {
				fmt.Println(err)
			}
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
