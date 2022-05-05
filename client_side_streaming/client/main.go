package main

import (
	"context"
	"log"
	"fmt"

	pb "client_side_streaming/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Port struct {
	client pb.PortServiceClient
}
 
func NewPort(conn grpc.ClientConnInterface) Port {
	return Port{
		client: pb.NewPortServiceClient(conn),
	}
}

func (p Port) Create(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
 
	// I am hardcoding this here but you should not!
	requests := []*pb.CreateRequest{
		{Id: "id-1", Name: "name-1"},
		{Id: "id-2", Name: "name-2"},
		{Id: "id-3", Name: "name-3"},
		{Id: "id-4", Name: "name-4"},
	}
 
	stream, err := p.client.Create(ctx)
	if err != nil {
		return fmt.Errorf("create stream: %w", err)
	}
 
	for _, request := range requests {
		if err := stream.Send(request); err != nil {
			return fmt.Errorf("send stream: %w", err)
		}
	}
 
	response, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("close and receive: %w", err)
	}
 
	log.Printf("%+v\n", response)
 
	return nil
}

func main() {
	log.Println("Client")
    conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
 
	portClient := NewPort(conn)
	if err := portClient.Create(context.Background()); err != nil {
		log.Fatalln(err)
	}
}
