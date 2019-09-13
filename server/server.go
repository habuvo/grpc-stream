package main

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/habuvo/grpc-stream/hello"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (_ *server) Exchange(stream hello.Hello_ExchangeServer) error {
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

func (_ *server) Command(_ *empty.Empty, stream hello.Hello_CommandServer) error {
	tick := time.Tick(time.Second * 1)
	for range tick {
		err := stream.Send(&hello.HelloResponse{})
		if err == io.EOF {
			log.Println("Received EOF")
			return nil
		}
		if err != nil {
			log.Println("Received error:", err)
			return err
		}
		log.Println("Send hello response")
	}
	return nil
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
