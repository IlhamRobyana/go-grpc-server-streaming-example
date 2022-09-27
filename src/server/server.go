package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/pramonow/go-grpc-server-streaming-example/src/proto"

	"google.golang.org/grpc"
)

type server struct{}

func (s server) FetchResponse(in *pb.Request, srv pb.StreamService_FetchResponseServer) error {

	log.Printf("fetch response for id : %d", in.Id)

	// var wg sync.WaitGroup
	for i := 0; i < 10000000; i++ {
		// wg.Add(1)
		// go func(count int64) {
		// 	defer wg.Done()
		// time.Sleep(time.Duration(count) * time.Second)
		resp := pb.Response{
			Id:      int64(i),
			Message: fmt.Sprintf("Sending you \"%d\"", i),
			Name:    "John Doe",
			Address: "Bandung",
			Amount:  100,
			Price:   90290000,
		}
		if err := srv.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
		log.Printf("finishing request number : %d", i)
		// }(int64(i))
	}

	// wg.Wait()
	return nil
}

func main() {
	// create listiner
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc server
	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, server{})

	log.Println("start server")
	// and start...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
