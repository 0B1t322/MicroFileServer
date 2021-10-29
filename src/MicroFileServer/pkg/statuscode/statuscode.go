package statuscode

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatusCode struct {
	Err			error
	Status		int
	GRPCCode	codes.Code
}

func (s *StatusCode) GRPCStatus() *status.Status {
	if s.GRPCCode == codes.Unknown {
		switch s.Status {
		case http.StatusNotFound:
			return status.New(
				codes.NotFound,
				s.Err.Error(),
			)
		case http.StatusBadRequest:
			return status.New(
				codes.InvalidArgument,
				s.Err.Error(),
			)
		case http.StatusRequestTimeout:
			return status.New(
				codes.DeadlineExceeded,
				s.Err.Error(),
			)
		case http.StatusServiceUnavailable:
			return status.New(
				codes.Canceled,
				s.Err.Error(),
			)
		case http.StatusForbidden:
			return status.New(
				codes.PermissionDenied,
				s.Err.Error(),
			)
		case http.StatusUnauthorized:
			return status.New(
				codes.Unauthenticated,
				s.Err.Error(),
			)
		}
	}

	return status.New(
		s.GRPCCode,
		s.Err.Error(),
	)
}

func (s *StatusCode) Error() string {
	return fmt.Sprintf("%v: %v", s.Status, s.Err)
}

func (s *StatusCode) SetGRPCStatus(
	code	codes.Code,
) *StatusCode {
	s.GRPCCode = code
	return s
}

func WrapStatusError(
	err 	error,
	status 	int,
) *StatusCode {
	return &StatusCode{
		Err: err,
		Status: status,
		GRPCCode: codes.Unknown,
	}
}

func GetStatus(
	err error,
) (status int, ok bool) {
	if err == nil {
		return http.StatusOK, true
	}
	
	StatusCode, ok := err.(*StatusCode)
	if !ok {
		return 0, ok
	}

	return StatusCode.Status, ok
}

func GetError(
	err error,
) error {
	StatusCode, ok := err.(*StatusCode)
	if !ok {
		return nil
	}

	return StatusCode.Err
}

func Is(
	err		error,
	target 	error,
) bool {
	if StatusCode, ok := err.(*StatusCode); ok {
		return errors.Is(StatusCode.Err, target)
	} else {
		return false
	}
}

// If error is not StatusCode just return wrapped error
func WrapErrorOnStatus(
	err			error,
	Message		string,
) error {
	statusErr, ok := err.(*StatusCode)
	if !ok {
		return errors.Wrap(err, Message)
	}

	statusErr.Err = errors.Wrap(statusErr.Err, Message)

	return statusErr
}