proto:
	git clone git@flotta-home:mindbond/proto.git
	protoc proto/auth.proto --go_out=plugins=grpc:./pkg/
	protoc proto/chat.proto --go_out=plugins=grpc:./pkg/
	rm -rf proto/

websocket-server:
	go run cmd/main.go
