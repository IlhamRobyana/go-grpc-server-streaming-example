package main

import (
	"context"
	"io"
	"log"
	"math/rand"

	pb "github.com/pramonow/go-grpc-server-streaming-example/src/proto"
	"github.com/pramonow/go-grpc-server-streaming-example/utils"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	rand.Seed(time.Now().Unix())

	// dail server
	conn, err := grpc.Dial("localhost:50005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	client := pb.NewStreamServiceClient(conn)
	defer utils.TimeTrack(time.Now(), "grpc streaming")
	in := &pb.Request{Id: 1}
	stream, err := client.FetchResponse(context.Background(), in)
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	// ctx := stream.Context()
	done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true //close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			log.Printf("Resp received: %+v", resp)
		}
	}()

	<-done
	log.Printf("finished")
}
