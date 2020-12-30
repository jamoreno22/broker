package main

import (
	"context"
	"log"
	"net"

	lab3 "github.com/jamoreno22/broker/pkg/proto"
	"google.golang.org/grpc"
)

type brokerServer struct {
	lab3.UnimplementedDNSServer
}

func main() {

	// create a listener on TCP port 8000
	lis, err := net.Listen("tcp", "10.10.28.20:8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	bs := brokerServer{}                               // create a gRPC server object
	grpcBrokerServer := grpc.NewServer()               // attach the Ping service to the server
	lab3.RegisterDNSServer(grpcBrokerServer, &bs) // start the server

	log.Println("BrokerServer running ...")
	if err := grpcBrokerServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

//DNSIsAvailable server side
func (b *brokerServer) DNSIsAvailable(ctx context.Context, msg *lab3.Message) (*lab3.DNSState, error) {

	return &lab3.DNSState{Dns1:true, Dns2:true, Dns3:true}, nil
}

//GetIP server side
func (b *brokerServer) GetIP(ctx context.Context, cmd *lab3.Command) (*lab3.PageInfo, error) {
	return &lab3.PageInfo{Ip:"10.10.10.0"}, nil
}