package handler

import (
	"database/sql"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var NothingFoundErr = errors.New("nothing found")

func WrapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return status.Errorf(codes.NotFound, "Nothing found")
	}

	if errors.Is(err, NothingFoundErr) {
		return status.Errorf(codes.NotFound, "Nothing found")
	}

	return status.Errorf(codes.Internal, "%v", err)
}
