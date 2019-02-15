package main

import (
  "context"
  "log"
  "net"

  "google.golang.org/grpc"
  pb "github.com/nickcen/concord_grpc/concord"
  "github.com/go-redis/redis"
)

const (
  port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{
  client *redis.Client
}

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
  log.Printf("Received Get Request: %v", in.Key)

  val, err := s.client.Get(in.Key).Result()
  if err != nil {
    panic(err)
  }

  return &pb.GetReply{Ret: true, Error: "", Value: []byte("Hello World " + val)}, nil
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
  log.Printf("Received Set Request: %v", in.Key)

  err := s.client.Set(in.Key, string(in.Value), 0).Err()
  if err != nil {
    panic(err)
  }

  return &pb.SetReply{Ret: true, Error: ""}, nil
}

func main() {
  lis, err := net.Listen("tcp", port)
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }
  s := grpc.NewServer()

  client := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
  })

  pb.RegisterConcordServer(s, &server{client: client})
  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}
