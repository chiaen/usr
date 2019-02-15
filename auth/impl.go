package main

import (
	"context"

	authapi "github.com/chiaen/usr/api/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUnimplemented = status.Error(codes.Unimplemented, "not implement yet")
)

type serviceImpl struct {
}

func newAuthService() (authapi.AuthenticationServer, error) {
	return &serviceImpl{}, nil
}

func (s *serviceImpl) SignupNewUser(ctx context.Context, req *authapi.PasswordRequest) (*authapi.TokenResponse, error) {
	return nil, errUnimplemented
}

func (s *serviceImpl) SignInWithPassword(ctx context.Context, req *authapi.PasswordRequest) (*authapi.TokenResponse, error) {
	return nil, errUnimplemented
}

func (s *serviceImpl) UpdatePassword(ctx context.Context, req *authapi.PasswordRequest) (*authapi.TokenResponse, error) {
	return nil, errUnimplemented
}
