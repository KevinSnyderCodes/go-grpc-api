package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	hellov1 "github.com/KevinSnyderCodes/go-grpc-api/gen/hello/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/KevinSnyderCodes/go-grpc-api/internal/env"
	"github.com/KevinSnyderCodes/go-grpc-api/internal/openapi"
	"github.com/KevinSnyderCodes/go-grpc-api/internal/server"
)

const title = "api"

var (
	fOpenAPIFilesGlob = flag.String("openapi-files-glob", env.MustGetStringOrDefault("OPENAPI_FILES_GLOB", "api/**/*.swagger.json"), "Glob for OpenAPI files.")
	fPort             = flag.Uint("port", env.MustGetUintOrDefault("PORT", 8080), "Port to listen on.")
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	flag.Parse()

	if *fOpenAPIFilesGlob == "" {
		return fmt.Errorf("must provide openapi files glob")
	}

	// Merge OpenAPI files
	spec, err := openapi.MergeGlob(title, *fOpenAPIFilesGlob)
	if err != nil {
		return fmt.Errorf("error merging openapi files: %w", err)
	}

	ctx := context.Background()

	addr := fmt.Sprintf(":%d", *fPort)

	// Create listener
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("error listening on %s: %w", addr, err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	hellov1.RegisterHelloServiceServer(grpcServer, &HelloServiceServer{})

	// Register grpc-gateway handlers
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if err := hellov1.RegisterHelloServiceHandlerFromEndpoint(ctx, gwmux, addr, opts); err != nil {
		return fmt.Errorf("error registering hello service handler: %w", err)
	}

	// Create and serve HTTP server
	httpServer := http.Server{
		Addr:    addr,
		Handler: server.NewHandler(grpcServer, gwmux, spec),
	}
	fmt.Printf("Listening on %s...\n", addr)
	if err := httpServer.Serve(lis); err != nil {
		return fmt.Errorf("error serving: %w", err)
	}

	return nil
}

// HelloServiceServer implements the HelloService service.
type HelloServiceServer struct {
	hellov1.UnimplementedHelloServiceServer
}

// Greet implements the Greet method of the HelloService service.
func (o *HelloServiceServer) Greet(ctx context.Context, req *hellov1.GreetRequest) (*hellov1.GreetResponse, error) {
	name := req.GetName()

	greeting := "Hello!"
	if name != "" {
		greeting = fmt.Sprintf("Hello, %s!", name)
	}

	resp := hellov1.GreetResponse{
		Greeting:  greeting,
		CreatedAt: timestamppb.Now(),
	}

	return &resp, nil
}
