package main

import (
  "context"
  "log"
  "net"

  "google.golang.org/grpc"
  pb "github.com/nickcen/concord_grpc/msgs"
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
  log.Printf("Received Get Request: [%v]", in.Key)

  val, err := s.client.Get(in.Key).Result()
  if err != nil {
    return &pb.GetReply{Ret: true, Error: "", Value: []byte(nil)}, nil
    // panic(err)
  }

  return &pb.GetReply{Ret: true, Error: "", Value: []byte(val)}, nil
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
  log.Printf("Received Set Request: [%v] - [%v]", in.Key, string(in.Value))

  err := s.client.Set(in.Key, string(in.Value), 0).Err()
  if err != nil {
    panic(err)
  }

  return &pb.SetReply{Ret: true, Error: ""}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
  log.Printf("Received Delete Request: [%v]", in.Key)

  err := s.client.Del(in.Key).Err()
  if err != nil {
    panic(err)
  }

  return &pb.DeleteReply{Ret: true, Error: ""}, nil
}

func (s *server) Init(ctx context.Context, in *pb.InitRequest) (*pb.InitReply, error) {
  log.Print("Received Init Request")

  err := s.client.FlushAll().Err()
  if err != nil {
    panic(err)
  }

  return &pb.InitReply{}, nil
}

func main() {
  log.Print("listen on "+port)
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
