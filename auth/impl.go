package main

import (
	"context"
	"github.com/chiaen/usr/utils/crypto"

	authapi "github.com/chiaen/usr/api/auth"
	"github.com/chiaen/usr/utils/uuid"
	"github.com/gocraft/dbr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUnimplemented = status.Error(codes.Unimplemented, "not implement yet")
)

type userRecord struct {
	UserID   string `db:"uid"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type serviceImpl struct {
	dbconn *dbr.Connection
}

func newAuthService() (authapi.AuthenticationServer, error) {
	conn, err := dbr.Open("mysql", "root:@(mysql:3306)/orb", nil)
	if err != nil {
		return nil, err
	}
	return &serviceImpl{conn}, nil
}

func (s *serviceImpl) SignupNewUser(ctx context.Context, req *authapi.PasswordRequest) (*authapi.TokenResponse, error) {
	var email, password string
	if cred := req.GetPassword(); cred == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty credential ")
	} else if email = cred.GetEmail(); email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty email")
	} else if password = cred.GetPassword(); password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty password")
	}

	uid, err := uuid.New()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "uid generation error: %v", err)
	}

	pwHash, err := crypto.HashAndSalt(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password hash error: %v", err)
	}

	sess := s.dbconn.NewSession(nil)
	if _, err := sess.InsertInto(TableUserRecords).
		Columns(TableUserRecordsColums...).
		Record(&userRecord{
			UserID:   uid.String(),
			Email:    email,
			Password: pwHash,
		}).Exec(); err != nil {
		return nil, status.Errorf(codes.Internal, "db insertion error: %v", err)
	}

	// TODO: issue token

	return nil, errUnimplemented
}

func (s *serviceImpl) SignInWithPassword(ctx context.Context, req *authapi.PasswordRequest) (*authapi.TokenResponse, error) {
	var email, password string
	if cred := req.GetPassword(); cred == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty credential ")
	} else if email = cred.GetEmail(); email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty email")
	} else if password = cred.GetPassword(); password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty password")
	}

	var user userRecord
	sess := s.dbconn.NewSession(nil)
	_, err := sess.Select("*").
		From(TableUserRecords).
		Where(dbr.Eq(ColEmail, email)).
		Load(&user)

	if err == dbr.ErrNotFound || !crypto.Compare(user.Password, password) {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect login parameters")
	}
	// TODO: issue token
	return nil, errUnimplemented
}

func (s *serviceImpl) UpdatePassword(ctx context.Context, req *authapi.PasswordRequest) (*authapi.TokenResponse, error) {
	var email, password string
	if cred := req.GetPassword(); cred == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty credential ")
	} else if email = cred.GetEmail(); email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty email")
	} else if password = cred.GetPassword(); password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty password")
	}

	var user userRecord
	sess := s.dbconn.NewSession(nil)
	_, err := sess.Select("*").
		From(TableUserRecords).
		Where(dbr.Eq(ColEmail, email)).
		Load(&user)

	if err == dbr.ErrNotFound {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if crypto.Compare(user.Password, password) {
		return nil, status.Errorf(codes.InvalidArgument, "new password should not be the same as old password")
	}

	pwHash, err := crypto.HashAndSalt(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password hash error: %v", err)
	}

	if _, err := sess.InsertInto(TableUserRecords).
		Columns(TableUserRecordsColums...).
		Record(&userRecord{
			UserID:   user.UserID,
			Email:    user.Email,
			Password: pwHash,
		}).Exec(); err != nil {
		return nil, status.Errorf(codes.Internal, "db insertion error: %v", err)
	}
	// TODO: issue token
	return nil, errUnimplemented
}
