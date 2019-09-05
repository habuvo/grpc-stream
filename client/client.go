package main

import (
	"github.com/habuvo/grpc-stream/hello"
	"log"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

func main() {
	conn, err := grpc.Dial("localhost:1234",
		grpc.WithBackoffConfig(grpc.BackoffConfig{MaxDelay: 2 * time.Second}),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't connect")
	}
	client := hello.NewHelloClient(conn)
	for {
		stream,err := client.Hello(context.Background())
		if err != nil {
			log.Println("error get stream",err)
			time.Sleep(time.Second)
		} else {
			for {
				err := stream.Send(&hello.HelloRequest{})
				if err != nil {
					log.Println("Error sending:", err)
					break
				} else {
					log.Println("message send")
				}
				time.Sleep(time.Second)
			}
		}
	}
}