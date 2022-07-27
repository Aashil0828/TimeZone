package main

import (
	"context"
	"fmt"
	"log"
	"timezone/pb/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("0.0.0.0:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot dial to the server : %v", err)
	}
	client := pb.NewTimeZoneServiceClient(conn)
	req := &pb.TimeZoneRequest{
		Latitude:  43.8287724,
		Longitude: -79.542908,
	}
	res, err := client.TimeZoneDetails(context.Background(), req)
	if err != nil {
		log.Fatalf("cannot recieve response from the server: %v", err)
	}
	fmt.Printf("response :\n utcOffSet: %v\n Zonename: %v\n Time In That Zone: %v\n", res.UtcOffset, res.ZoneName, res.TimeInThatZone)
}
