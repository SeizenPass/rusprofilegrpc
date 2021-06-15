package main

import (
	"flag"
	proto "github.com/SeizenPass/rusprofilegrpc/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type app struct {
	proto.UnimplementedSearchServiceServer
	errorLog *log.Logger
	infoLog  *log.Logger
	service SearchServiceInterface
}

func main() {
	host := flag.String("host", "localhost:8080", "host")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	service := &SearchServiceImpl{}

	app := &app{
		UnimplementedSearchServiceServer: proto.UnimplementedSearchServiceServer{},
		errorLog:                         errorLog,
		infoLog:                          infoLog,
		service:                          service,
	}

	l, err := net.Listen("tcp", *host)
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}

	s := grpc.NewServer()
	proto.RegisterSearchServiceServer(s, app)

	app.infoLog.Printf("Starting server on %s", *host)
	if err := s.Serve(l); err != nil {
		app.errorLog.Fatalf("failed to serve:%v", err)
	}
}
