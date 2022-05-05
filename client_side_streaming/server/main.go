package main
 
import (
	"log"
	"net"
	"io"
 
	gogrpc "google.golang.org/grpc"
 
	pb "client_side_streaming/proto"
)
 
type Port struct{
	pb.UnimplementedPortServiceServer
}
 
func (p Port) Create(stream pb.PortService_CreateServer) error {
	var total int32
 
	for {
		port, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.CreateResponse{
				Total: total,
			})
		}
		if err != nil {
			return err
		}
 
		total++
		log.Printf("%+v\n", port)
	}
}

func main() {
	log.Println("server")
 
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
 
	server := gogrpc.NewServer()
	protoServer := Port{}
 
	pb.RegisterPortServiceServer(server, protoServer)
 
	log.Fatalln(server.Serve(listener))
}