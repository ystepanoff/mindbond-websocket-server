proto:
	git clone git@flotta-home:mindbond/proto.git
	protoc proto/chat.proto --go_out=plugins=grpc:./pkg/
	rm -rf proto/

websocket-server:
	go run cmd/main.go
