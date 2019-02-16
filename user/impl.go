package main

import (
	"context"
	"github.com/chiaen/usr/auth"
	"github.com/chiaen/usr/utils/dbr"

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
	uid, ok := auth.UserIDfromValidAccessToken(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	sess := s.dbconn.NewSession(nil)

	var email string
	if err := sess.Select(ColEmail).From(TableUserRecords).Where(dbr.Eq(ColUID, uid)).LoadOneContext(ctx, &email); err != nil {
		return nil, status.Errorf(codes.Internal, "query error %v", err)
	}
	var interests []string
	if _, err :=sess.Select(ColName).
		From(TableUserHasInterest).
		Join(TableInterests, dbrutils.EqExpr(TableUserHasInterest, ColInterestID, TableInterests, ColID)).
		Where(dbr.Eq(ColUserID, uid)).
		LoadContext(ctx, &interests); err != nil && err != dbr.ErrNotFound {
		return nil, status.Errorf(codes.Internal, "query error %v", err)
	}
	return &userapi.ProfileResponse{
		Email: email,
		Interest: interests,
	}, nil
}

func (s *serviceImpl) AddInterest(ctx context.Context, req *userapi.AddInterestRequest) (*userapi.AddInterestResponse, error) {
	uid, ok := auth.UserIDfromValidAccessToken(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	var interest string
	var published bool
	if interest = req.GetInterest(); interest == "" {
		return nil, status.Errorf(codes.InvalidArgument, "interest should not be empty")
	}
	sess := s.dbconn.NewSession(nil)
	tx, err := sess.Begin()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database transaction init failed: %v", err)
	}
	defer tx.RollbackUnlessCommitted()

	result, err := tx.InsertInto(TableInterests).
		Pair(ColName, interest).
		ExecContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database insert error: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	_, err = tx.InsertInto(TableUserHasInterest).
		Pair(ColUserID, uid).
		Pair(ColInterestID, id).
		Pair(ColPublished, published).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database insert error: %v", err)
	}
	err = tx.Commit()
	return &userapi.AddInterestResponse{}, err
}

func (s *serviceImpl) ListUserByInterest(ctx context.Context, req *userapi.ListUserByInterestRequest) (*userapi.ListUserByInterestResponse, error) {
	var interest string
	if interest = req.GetInterest(); interest == "" {
		return nil, status.Errorf(codes.InvalidArgument, "interest should not be empty")
	}
	sess := s.dbconn.NewSession(nil)


	var emails []string
	if _, err := sess.Select(ColEmail).
		From(TableUserHasInterest).
		LeftJoin(TableInterests, dbrutils.EqExpr(TableUserHasInterest, ColInterestID, TableInterests, ColID)).
		LeftJoin(TableUserRecords, dbrutils.EqExpr(TableUserHasInterest, ColUserID, TableUserRecords, ColUID)).
		Where(dbr.Eq(ColPublished, true)).
		Where(dbr.Eq(ColName, interest)).
		LoadContext(ctx, &emails); err == dbr.ErrNotFound {
			return &userapi.ListUserByInterestResponse{}, nil
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database query error: %v", err)
	}
	return &userapi.ListUserByInterestResponse{
		Email: emails,
	}, errUnimplemented
}
