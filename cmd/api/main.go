package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/timestamppb"

	connect "connectrpc.com/connect"

	hellov1 "github.com/KevinSnyderCodes/go-grpc-api/gen/hello/v1"
	"github.com/KevinSnyderCodes/go-grpc-api/gen/hello/v1/hellov1connect"

	"github.com/KevinSnyderCodes/go-grpc-api/internal/env"
)

var (
	fPort = flag.Uint("port", env.MustGetUintOrDefault("PORT", 8080), "Port to listen on.")
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	path, handler := hellov1connect.NewHelloServiceHandler(&HelloServiceServer{})
	mux.Handle(path, handler)

	addr := fmt.Sprintf(":%d", *fPort)
	fmt.Printf("Listening on %s...\n", addr)

	http.ListenAndServe(
		addr,
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

// HelloServiceServer implements the HelloService service.
type HelloServiceServer struct {
	hellov1connect.UnimplementedHelloServiceHandler
}

// Greet implements the Greet method of the HelloService service.
func (o *HelloServiceServer) Greet(ctx context.Context, req *connect.Request[hellov1.GreetRequest]) (*connect.Response[hellov1.GreetResponse], error) {
	name := req.Msg.GetName()

	greeting := "Hello!"
	if name != "" {
		greeting = fmt.Sprintf("Hello, %s!", name)
	}

	resp := hellov1.GreetResponse{
		Greeting:  greeting,
		CreatedAt: timestamppb.Now(),
	}

	return connect.NewResponse(&resp), nil
}
