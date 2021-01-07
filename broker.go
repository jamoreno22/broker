package main

import (
	"context"
	"log"
	"net"

	lab3 "github.com/jamoreno22/broker/pkg/proto"
	"google.golang.org/grpc"
)

type brokerServer struct {
	lab3.UnimplementedBrokerServer
}

func main() {

	// create a listener on TCP port 8000
	lis, err := net.Listen("tcp", "10.10.28.20:8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	bs := brokerServer{}                             // create a gRPC server object
	grpcBrokerServer := grpc.NewServer()             // attach the Ping service to the server
	lab3.RegisterBrokerServer(grpcBrokerServer, &bs) // start the server

	log.Println("BrokerServer running ...")
	if err := grpcBrokerServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

//DNSIsAvailable server side
func (b *brokerServer) DNSIsAvailable(ctx context.Context, msg *lab3.Message) (*lab3.DNSState, error) {
	dnsIps := []string{"10.10.28,17:8000", "10.10.28,18:8000", "10.10.28,19:8000"}
	state := lab3.DNSState{Dns1: true, Dns2: true, Dns3: true}
	for i, val := range dnsIps {
		// connection to each dns server, if fail server not available
		_, err := grpc.Dial(val, grpc.WithInsecure())
		if err != nil {
			switch i {
			case 1:
				state.Dns1 = false

			case 2:
				state.Dns2 = false

			case 3:
				state.Dns3 = false
			}
		}
	}
	return &state, nil
}

//GetIP server side
func (b *brokerServer) GetIP(ctx context.Context, cmd *lab3.Command) (*lab3.PageInfo, error) {
	var pageInfo lab3.PageInfo
	//Call to GetPageInfo in DNSServer as a DNS client
	dnsIps := []string{"10.10.28.17:8000", "10.10.28.18:8000", "10.10.28.19:8000"}
	for _, val := range dnsIps {
		conn1, err1 := grpc.Dial(val, grpc.WithInsecure())
		if err1 != nil {
			continue
		}
		dnsc := lab3.NewDNSClient(conn1)
		pageInfo, _ = &dnsc.GetIP(ctx, cmd)
		break
	}
	return &pageInfo, nil
}
