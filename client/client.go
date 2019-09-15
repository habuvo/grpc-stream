package main

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/habuvo/grpc-stream/hello"
	"google.golang.org/grpc/keepalive"
	"log"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

func main() {
	conn, err := grpc.Dial("localhost:1234",
		grpc.WithBackoffConfig(grpc.BackoffConfig{MaxDelay: 2 * time.Second}),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    time.Second * 10,
			Timeout: time.Second * 5,
		}),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't connect")
	}
	client := hello.NewHelloClient(conn)
	//send
	go func() {
		for {
			stream, err := client.Exchange(context.Background())
			if err != nil {
				log.Println("error get stream", err)
				time.Sleep(time.Second)
			} else {
				tick := time.Tick(time.Second * 1)
				for range tick {
					err := stream.Send(&hello.HelloRequest{})
					if err != nil {
						log.Println("Error sending:", err)
						break
					} else {
						log.Println("message send")
					}
				}
			}
		}
	}()
	//receive
	go func() {
		for {
			stream, err := client.Command(context.Background(), &empty.Empty{})
			if err != nil {
				log.Println("error get stream", err)
				time.Sleep(time.Second)
			} else {
				tick := time.Tick(time.Second * 1)
				for range tick {
					_, err := stream.Recv()
					if err != nil {
						log.Println("Error sending:", err)
						break
					} else {
						log.Println("message receive")
					}
				}
			}
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
