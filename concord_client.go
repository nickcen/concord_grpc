package main

import (
  "context"
  "log"
  "time"

  "google.golang.org/grpc"
  pb "github.com/nickcen/concord_grpc/concord"
)

const (
  address     = "localhost:50051"
)

func main() {
  // Set up a connection to the server.
  conn, err := grpc.Dial(address, grpc.WithInsecure())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer conn.Close()
  c := pb.NewConcordClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()

  s_r, err := c.Set(ctx, &pb.SetRequest{Key: "test", Value: []byte("hello world")})
  if err != nil {
    log.Fatalf("could not set: %v", err)
  }
  _ = s_r

  g_r, err := c.Get(ctx, &pb.GetRequest{Key: "test"})
  if err != nil {
    log.Fatalf("could not get: %v", err)
  }
  log.Printf("Greeting: %s", g_r.Value)
}
