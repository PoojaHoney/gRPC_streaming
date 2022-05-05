package main

import (
	pb "bidirectional_streaming/proto"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type container struct {
	users []*pb.User
}

func main() {
	conn, err := grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Error connecting: %v \n", err)
	}

	c := pb.NewUsersClient(conn)
	bulkUsers(c)
}

func (c container) initUsers() []*pb.User {
	c.users = append(c.users, c.getUser("1", "Carl", "Phill", 23))
	c.users = append(c.users, c.getUser("2", "Marisol", "Richardson", 29))
	c.users = append(c.users, c.getUser("3", "Mia", "Phill", 27))
	c.users = append(c.users, c.getUser("4", "Tomas", "Smith", 25))
	c.users = append(c.users, c.getUser("5", "Zian", "Heat", 28))

	return c.users
}

func bulkUsers(c pb.UsersClient) {
	stream, err := c.CreateUser(context.Background())
	if err != nil {
		log.Fatalf("Error when getting stream object: %v", err)
		return
	}

	requests := container{}.initUsers()

	done := make(chan bool)

	// Use a go routine to send request messages to the server
	go func() {
		for _, req := range requests {
			stream.Send(req)
			time.Sleep(100 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Use a go routine to receive response messages from the server
	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error when receiving response: %v", err)
			}
			fmt.Println("Server Response: ", res)
		}
		done <- true
	}()

	<-done
	log.Printf("finished")
}

func (c container) getUser(id, name, lastName string, age int32) *pb.User {
	return &pb.User{
		Id:       id,
		Name:     name,
		LastName: lastName,
		Age:      age,
	}
}
