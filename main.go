package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/abdukahhor/streamer/handlers/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	var wg sync.WaitGroup
	c, err := grpc.Dial(":9090", grpc.WithInsecure())
	defer c.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	cl := pb.NewStreamerClient(c)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go getURL(cl, &wg)
	}
	wg.Wait()
}

func getURL(cl pb.StreamerClient, wg *sync.WaitGroup) {
	defer wg.Done()
	str, err := cl.GetRandomDataStream(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Printf("failed to start timer: %v", err)
	}

	for {
		msg, err := str.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("stream read failed: %v", err)
		}

		fmt.Println(msg)
	}
}
