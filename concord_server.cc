/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

#include <iostream>
#include <memory>
#include <string>
#include <hiredis/hiredis.h>

#include <grpcpp/grpcpp.h>

#include "concord.grpc.pb.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;
using concord::GetRequest;
using concord::GetReply;
using concord::SetRequest;
using concord::SetReply;
using concord::DeleteRequest;
using concord::DeleteReply;
using concord::InitRequest;
using concord::InitReply;
using concord::Concord;

// Logic and data behind the server's behavior.
class ConcordServiceImpl final : public Concord::Service {
  Status Get(ServerContext* context, const GetRequest* request,
    GetReply* reply) override {

    redisContext *c = redisConnect("127.0.0.1", 6379);
    if (c == NULL || c->err) {
      if (c) {
        printf("Error: %s\n", c->errstr);
        // handle error
      } else {
        printf("Can't allocate redis context\n");
      }
    }
    redisReply *pRedisReply = (redisReply*)redisCommand(c, "GET %s %s", request->key());
    std::cout << pRedisReply->str << std::endl;
    freeReplyObject(pRedisReply); 

    reply->set_value(pRedisReply->str);

    return Status::OK;
  }

  Status Set(ServerContext* context, const SetRequest* request,
    SetReply* reply) override {
    std::cout << "received Set request " << request->key() << ":" << request->value() << std::endl;

    redisContext *c = redisConnect("127.0.0.1", 6379);
    if (c == NULL || c->err) {
      if (c) {
        printf("Error: %s\n", c->errstr);
        // handle error
      } else {
        printf("Can't allocate redis context\n");
      }
    }
    std::string k = request->key();
    std::string v = request->value();
    redisReply *pRedisReply = (redisReply*)redisCommand(c, "SET %s %s", "test", "bar");
    std::cout << pRedisReply->str << std::endl;
    freeReplyObject(pRedisReply); 
    return Status::OK;
  }

  Status Delete(ServerContext* context, const DeleteRequest* request,
    DeleteReply* reply) override {

    return Status::OK;
  }

  Status Init(ServerContext* context, const InitRequest* request,
    InitReply* reply) override {

    return Status::OK;
  }
};

void RunServer() {
  std::string server_address("0.0.0.0:50051");
  ConcordServiceImpl service;

  ServerBuilder builder;
  // Listen on the given address without any authentication mechanism.
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  // Register "service" as the instance through which we'll communicate with
  // clients. In this case it corresponds to an *synchronous* service.
  builder.RegisterService(&service);
  // Finally assemble the server.
  std::unique_ptr<Server> server(builder.BuildAndStart());
  std::cout << "Server listening on " << server_address << std::endl;

  // Wait for the server to shutdown. Note that some other thread must be
  // responsible for shutting down the server for this call to ever return.
  server->Wait();
}

int main(int argc, char** argv) {
  RunServer();

  return 0;
}