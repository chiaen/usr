package main

import (
	"context"

	userapi "github.com/chiaen/usr/api/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUnimplemented = status.Error(codes.Unimplemented, "not implement yet")
)

type serviceImpl struct {
}

func newUserService() (userapi.UserServer, error) {
	return &serviceImpl{}, nil
}

func (s *serviceImpl) GetProfile(ctx context.Context, req *userapi.GetProfileRequest) (*userapi.ProfileResponse, error) {
	return nil, errUnimplemented
}

func (s *serviceImpl) AddInterest(ctx context.Context, req *userapi.AddInterestRequest) (*userapi.AddInterestResponse, error) {
	return nil, errUnimplemented
}

func (s *serviceImpl) ListUserByInterest(ctx context.Context, req *userapi.ListUserByInterestRequest) (*userapi.ListUserByInterestResponse, error) {
	return nil, errUnimplemented
}
