package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"timezone/pb/pb"

	"github.com/bradfitz/latlong"
	"github.com/oschwald/geoip2-golang"
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
	var client_ip []string
	var forwarded []string
	var ipaddress string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		client_ip = md.Get("client-ip")
		forwarded = md.Get("x-forwarded-for")
	}
	if len(forwarded) != 0 {
		ipaddress = forwarded[0]
	} else {
		ipaddress = client_ip[0]
	}
	db, err := geoip2.Open("GeoIP2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ip := net.ParseIP(ipaddress)
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(record.City.Names)
	latitude := req.GetLatitude()
	fmt.Println(latitude)
	longitude := req.GetLongitude()
	fmt.Println(longitude)
	currentTime := req.GetTime()
	timezone := latlong.LookupZoneName(latitude, longitude)
	if err := setTimezone(timezone); err != nil {
		return &pb.TimeZoneResponse{}, err // most likely timezone not loaded in Docker OS
	}
	var t time.Time
	if currentTime == "" {
		t = getTime(time.Now())
	} else {
		t, err := time.Parse("2022-07-25 15:29:45.7725535 +0530 IST", currentTime)
		if err != nil {
			return &pb.TimeZoneResponse{}, err
		}
		t = getTime(t)
	}
	m, _ := t.MarshalText()
	a := string(m)
	var utc string
	if strings.Contains(a, "+") {
		u := strings.Split(a, "+")
		utc = fmt.Sprintf("UTC+%v", u[len(u)-1])
	} else {
		u := strings.Split(a, "-")
		utc = fmt.Sprintf("UTC-%v", u[len(u)-1])
	}
	return &pb.TimeZoneResponse{UtcOffset: utc, ZoneName: timezone, TimeInThatZone: t.Format(time.RFC3339)}, nil
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
