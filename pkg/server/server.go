package server

import (
	"fmt"
	"github.com/ystepanoff/gowest"
	"net"
	"net/http"
)

func Start(port int) {
	http.HandleFunc("/", wsHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
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
