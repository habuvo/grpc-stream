package main

import (
	"context"
	"github.com/habuvo/grpc-stream/hello"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (_ *server) Hello(stream hello.Hello_HelloServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			log.Println("Received EOF")
			return nil
		}
		if err != nil {
			log.Println("Received error:", err)
			return err
		}
		log.Println("Received hello request.")
	}
}

func (_ *server) HelloOnce(ctx context.Context,in *hello.HelloRequest) (*hello.HelloResponse,error) {
	return &hello.HelloResponse{},nil
}

func main() {
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	hello.RegisterHelloServer(s, &server{})
	println("serving :1234")
	if err := s.Serve(l); err != nil {
		println(err)
	}
}
