package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"
	"timezone/pb/pb"

	"github.com/bradfitz/latlong"
	"github.com/ip2location/ip2location-go/v9"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	pb.UnimplementedTimeZoneServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) TimeZoneDetails(ctx context.Context, req *pb.TimeZoneRequest) (*pb.TimeZoneResponse, error) {
	// p, _ := peer.FromContext(ctx)
	// ipaddress := p.Addr.String()
	// fmt.Println(ipaddress)
	db, err := ip2location.OpenDB("service/IP2LOCATION-LITE-DB11.BIN")
	if err != nil {
		return &pb.TimeZoneResponse{}, err
	}
	var client_ip []string
	var forwarded []string
	ipaddress := req.Ipaddress
	if ipaddress == "" {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			client_ip = md.Get("client-ip")
			forwarded = md.Get("x-forwarded-for")
		}
		if len(forwarded) != 0 {
			ipaddress = strings.Split(forwarded[0], ",")[0]
		} else {
			ipaddress = client_ip[0]
		}
	}
	fmt.Println(ipaddress)
	results, err := db.Get_all(ipaddress)
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}
	db.Close()
	currentTime := req.GetTime()
	fmt.Println(currentTime)
	timezone := latlong.LookupZoneName(float64(results.Latitude), float64(results.Longitude))
	if err := setTimezone(timezone); err != nil {
		return &pb.TimeZoneResponse{}, err // most likely timezone not loaded in Docker OS
	}
	var t time.Time
	if currentTime == "" {
		t = getTime(time.Now())
	} else {
		t, err := time.Parse(time.RFC3339, currentTime)
		fmt.Println(t)
		if err != nil {
			return &pb.TimeZoneResponse{}, err
		}
		t = getTime(t)
	}
	m, _ := t.MarshalText()
	a := string(m)
	fmt.Println(a)
	var utc string
	if strings.Contains(a, "+") {
		u := strings.Split(a, "+")
		utc = fmt.Sprintf("UTC+%v", u[len(u)-1])
	} else {
		u := strings.Split(a, "-")
		utc = fmt.Sprintf("UTC-%v", u[len(u)-1])
	}
	return &pb.TimeZoneResponse{UtcOffset: utc, ZoneName: timezone, TimeInThatZone: t.Format(time.RFC3339), Country: results.Country_long, Latitude: float64(results.Latitude), Longitude: float64(results.Longitude)}, nil
}

var loc *time.Location

func setTimezone(tz string) error {
	location, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}
	loc = location
	return nil
}

func getTime(t time.Time) time.Time {
	return t.In(loc)
}
