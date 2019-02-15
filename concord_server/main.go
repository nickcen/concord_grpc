package main

import (
  "context"
  "log"
  "net"

  "google.golang.org/grpc"
  pb "github.com/nick/concord_grpc/concord/concord"
)

const (
  port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
  log.Printf("Received: %v", in.Key)
  return &pb.GetReply{Body: "Hello World" + in.Key}, nil
}

func main() {
  lis, err := net.Listen("tcp", port)
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }
  s := grpc.NewServer()
  pb.RegisterConcordServer(s, &server{})
  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}
