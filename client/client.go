package main

import (
	"context"
	"flag"
	proto "github.com/SeizenPass/rusprofilegrpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

type app struct {
	client   proto.SearchServiceClient
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	host := flag.String("host", "localhost:8081", "host")
	grpcServerEndpoint := flag.String("grpc-server-endpoint", "localhost:8080", "gRPC server endpoint")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	connGrpc, err := grpc.Dial(*grpcServerEndpoint, grpc.WithInsecure())
	if err != nil {
		errorLog.Printf("could not connect: %v", err)
	}
	defer connGrpc.Close()

	client := proto.NewSearchServiceClient(connGrpc)

	app := &app{
		client:   client,
		errorLog: errorLog,
		infoLog:  infoLog,
	}
	app.infoLog.Printf("Starting server on %s", *host)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	regErr := proto.RegisterSearchServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if regErr != nil {
		app.errorLog.Fatal(err)
	}
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	httpErr := http.ListenAndServe(*host, mux)
	if httpErr != nil {
		app.errorLog.Fatal(err)
	}
}
