package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"timezone/pb/pb"
	"timezone/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

func main() {
	server := service.NewServer()
	// grpcServer := grpc.NewServer()
	// pb.RegisterTimeZoneServiceServer(grpcServer, server)
	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Fatalf("cannot start server : %v", err)
	}
	log.Print("server started")
	// err = grpcServer.Serve(listener)
	// if err != nil {
	// 	log.Fatalf("cannot start server : %v", err)
	// }
	// conn, err := grpc.Dial("0.0.0.0:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("cannot dial to the server : %v", err)
	// }
	// client := pb.NewTimeZoneServiceClient(conn)
	// req := &pb.TimeZoneRequest{
	// 	Latitude:  43.8287724,
	// 	Longitude: -79.542908,
	// }
	// res, err := client.TimeZoneDetails(context.Background(), req)
	// if err != nil {
	// 	log.Fatalf("cannot recieve response from the server: %v", err)
	// }
	// fmt.Printf("response :\n utcOffSet: %v\n Zonename: %v\n Time In That Zone: %v\n", res.UtcOffset, res.ZoneName, res.TimeInThatZone)
	mux := runtime.NewServeMux(runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		return metadata.Pairs("client-ip", req.RemoteAddr, "client-ip-header", req.Header.Get("x-forwarded-for"))
	}))
	ctx := context.Background()
	err = pb.RegisterTimeZoneServiceHandlerServer(ctx, mux, server)
	if err != nil {
		log.Fatalf("cannot start server : %v", err)
	}
	http.Serve(listener, mux)
}
