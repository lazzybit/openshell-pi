package main

import (
	"encoding/json"
	"flag"
	"github.com/lazzybit/openshell-pi-interceptor/internal/interceptor"
	pb "github.com/lazzybit/openshell-pi-interceptor/internal/pb/gatewayinterceptorv1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	listen := flag.String("listen", ":50051", "gRPC listen address")
	config := flag.String("templates", "", "JSON file mapping approved names to inline templates")
	flag.Parse()
	impl := &interceptor.Server{}
	if *config != "" {
		data, err := os.ReadFile(*config)
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(data, &impl.Templates); err != nil {
			log.Fatal(err)
		}
	}
	listener, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer(grpc.MaxRecvMsgSize(1048576))
	pb.RegisterGatewayInterceptorServer(server, impl)
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() { <-stop; server.GracefulStop() }()
	log.Printf("openshell-pi-interceptor listening on %s", *listen)
	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
