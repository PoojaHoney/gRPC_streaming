package main

import (
	pb "bidirectional_streaming/proto"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUsersServer
}

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Unable to listen on port 9090: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUsersServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}

func (s *server) CreateUser(stream pb.Users_CreateUserServer) error {
	for {
		// Receive the request and possible error from the stream object
		req, err := stream.Recv()

		// If there are no more requests, we return
		if err == io.EOF {
			return nil
		}

		// Handle error from the stream object
		if err != nil {
			log.Fatalf("Error when reading client request stream: %v", err)
		}

		// Get name, last name and user id, form the request
		name, lastName, userID := req.GetName(), req.GetLastName(), req.GetId()
		fmt.Printf("Request: name: %v, last_name: %v, id: %v \n", name, lastName, userID)

		// Initialize the errors and success variables
		errors := []string{}
		success := true

		// Run some validations
		if len(name) <= 3 {
			errors = append(errors, "Name is too short")
			success = false
		}

		if lastName == "Phill" {
			errors = append(errors, "Last Name already taken")
			success = false
		}

		// Build and send response to the client
		res := stream.Send(&pb.CreateUserRes{
			UserId:  userID,
			Success: success,
			Errors:  errors,
		})

		// Handle any possible error, when sending the response
		if res != nil {
			log.Fatalf("Error when response was sent to the client: %v", res)
		}
	}
}
