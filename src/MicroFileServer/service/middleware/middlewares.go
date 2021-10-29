package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MicroFileServer/pkg/contextvalue/subcontext"
	"github.com/MicroFileServer/pkg/statuscode"
	"github.com/go-kit/kit/endpoint"
	log "github.com/sirupsen/logrus"
)

type ReqWithID interface{
	GetID() string
}


type IsOwnerChecker interface {
	IsOwner(
		ctx			context.Context,
		ID			string,
		UserID		string,
	) error
}
// req shoild implements ReqWithID
// 
// And in ctx should be subcontext to get userID
// 
// Check if user a owner of object with current ID
func CheckUserIsOwner(
	check 	IsOwnerChecker,
) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(
			ctx context.Context, 
			request interface{},
		) (response interface{}, err error) {
			log.Debug("CheckUserIsOwner")
			switch req := request.(type) {
			case ReqWithID:
				log.Debug("ReqWithID")
				UserID, err := subcontext.GetSubFromContext(ctx)
				if err != nil {
					return nil, statuscode.WrapStatusError(
						fmt.Errorf("Failed to check owner"),
						http.StatusInternalServerError,
					)
				}

				if err := check.IsOwner(ctx,req.GetID(), UserID); err != nil {
					return nil, err
				}
			}
			return next(ctx, request)
		}
	}
}

// One of middlewares should be without error
// 
// It's be good if all this middlewares return similat error because if all of them return error, will return only last failed middleware
func MergeMiddlewaresIntoOr(
	ms ...endpoint.Middleware,
) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(
			ctx context.Context, 
			request interface{},
		) (interface{}, error) {
			var err error
			var success bool = false

			for _, m := range ms {
				_, err = m(endpoint.Nop)(ctx, request)
				if err == nil {
					success = true
					break
				}
			}

			if success {
				return next(ctx, request)
			}

			return nil, err
		}
	}
}

// req should implement method 
// 	SetUserID(userid string)
// 	GetUserID() string
func ValidateAndSetUserID() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(
			ctx		context.Context,
			request		interface{},
		) (response interface{}, err error) {
			type UserIDReq interface{
				SetUserID(userid string)
				GetUserID() string
			}

			log.Info("ValidateAndSetUserID middleware")
			userid, err := subcontext.GetSubFromContext(ctx)
			if err != nil {
				return nil, statuscode.WrapStatusError(
					fmt.Errorf("failed to chek userid"),
					http.StatusInternalServerError,
				)
			}

			switch req := request.(type){
			case UserIDReq:
				reqUserID := req.GetUserID()

				if reqUserID == "" {
					req.SetUserID(userid)
				} else if reqUserID != userid {
					return nil, statuscode.WrapStatusError(
						fmt.Errorf("You are not have permission to do that"),
						http.StatusForbidden,
					)
				}
			}
			return next(ctx, request)
		}
	}
}