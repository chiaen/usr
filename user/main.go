package main

import (
	"log"
	"net"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/chiaen/usr/api/user"
	"github.com/chiaen/usr/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

const (
	defaultPort = "80"
)

var (
	exposing = kingpin.Flag("port", "exposing port of server").Default(defaultPort).Short('p').String()
)

func main() {
	l, err := net.Listen("tcp", *exposing)
	if err != nil {
		log.Fatalf("listen to port %s failed: %v", *exposing, err)
	}
	server := grpc.NewServer(grpc_middleware.WithUnaryServerChain(auth.UnaryTokenVerifier))
	impl, err := newUserService()
	if err != nil {
		log.Fatalf("cannot init api server : %v", err)
	}
	user.RegisterUserServer(server, impl)
	if err := server.Serve(l); err != nil {
		log.Printf("server execution err: %v", err)
		os.Exit(1)
	}
}
