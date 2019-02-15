package main

import (
	"context"
	"github.com/chiaen/usr/api/user"
	"log"
	"net/http"

	"github.com/alecthomas/kingpin"
	"github.com/chiaen/usr/api/auth"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	kingpin.Parse()
	if err := sandbox(); err != nil {
		log.Fatal(err)
	}
}

func sandbox() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jsonPb := &runtime.JSONPb{}
	jsonPb.EnumsAsInts = true
	jsonPb.EmitDefaults = true
	opts := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption("application/json", jsonPb),
		runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
	}

	gw := runtime.NewServeMux(opts...)

	o := []grpc.DialOption{grpc.WithInsecure()}
	if err := auth.RegisterAuthenticationHandlerFromEndpoint(ctx, gw, "auth:80", o); err != nil {
		return err
	}
	if err := user.RegisterUserHandlerFromEndpoint(ctx, gw, "user:80", o); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", gw)
	return http.ListenAndServe(":8080", mux)
}
