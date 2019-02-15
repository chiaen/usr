package main

import (
	"context"

	userapi "github.com/chiaen/usr/api/user"
	"github.com/gocraft/dbr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUnimplemented = status.Error(codes.Unimplemented, "not implement yet")
)

type userInterest struct {
	UserID       string `db:"uid"`
	Email        string `db:"email"`
	InterestName string `db:"interest_name"`
}

type userInterests []userInterest

type serviceImpl struct {
	dbconn *dbr.Connection
}

func newUserService() (userapi.UserServer, error) {
	conn, err := dbr.Open("mysql", "root:@(mysql:3306)/orb", nil)
	if err != nil {
		return nil, err
	}
	return &serviceImpl{conn}, nil
}

func (s *serviceImpl) GetProfile(ctx context.Context, req *userapi.GetProfileRequest) (*userapi.ProfileResponse, error) {
	//TODO: extract user form context
	uid := "1234567"

	return nil, errUnimplemented
}

func (s *serviceImpl) AddInterest(ctx context.Context, req *userapi.AddInterestRequest) (*userapi.AddInterestResponse, error) {
	return nil, errUnimplemented
}

func (s *serviceImpl) ListUserByInterest(ctx context.Context, req *userapi.ListUserByInterestRequest) (*userapi.ListUserByInterestResponse, error) {
	return nil, errUnimplemented
}
