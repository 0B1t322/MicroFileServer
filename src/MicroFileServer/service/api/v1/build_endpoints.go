package v1

import (
	log "github.com/sirupsen/logrus"
	"context"

	"github.com/MicroFileServer/pkg/contextvalue/rolecontext"
	"github.com/MicroFileServer/pkg/contextvalue/subcontext"
	"github.com/MicroFileServer/service/api/v1/files"
	"github.com/MicroFileServer/service/middleware"
	"github.com/MicroFileServer/service/middleware/mgsess"
	"github.com/go-kit/kit/endpoint"
)

func (a *Api) endpoints() ApiEndpoints {
	return ApiEndpoints{
		Files: files.MakeEndpoints(a.FileService),
	}
}

func (a *Api) buildEndpoints() ApiEndpoints {
	endpoints := a.endpoints()

	endpoints.Files.DeleteFile = endpoint.Chain(
		a.Auth.AuthMiddleware(),
		mgsess.PutMongoSessIntoCtx(),
		middleware.MergeMiddlewaresIntoOr(
			a.Auth.IsAdmin(),
			middleware.CheckUserIsOwner(
				a.FileService,
			),
		),
	)(endpoints.Files.DeleteFile)

	endpoints.Files.GetFile = endpoint.Chain(
		a.Auth.AuthMiddleware(),
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Files.GetFile)

	endpoints.Files.UploadFile = endpoint.Chain(
		a.Auth.AuthMiddleware(),
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Files.UploadFile)

	endpoints.Files.GetFiles = endpoint.Chain(
		a.Auth.AuthMiddleware(),
		mgsess.PutMongoSessIntoCtx(),
		middleware.MergeMiddlewaresIntoOr(
			a.Auth.IsAdmin(),
			middleware.SetUserID(),
		),
	)(endpoints.Files.GetFiles)


	endpoints.Files.DownloadFile = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Files.DownloadFile)
	return endpoints
}

func (a *Api) _buildEndpoints() ApiEndpoints {
	endpoints := a.endpoints()

	endpoints.Files.DeleteFile = endpoint.Chain(
		testAuth(),
		mgsess.PutMongoSessIntoCtx(),
		middleware.MergeMiddlewaresIntoOr(
			a.Auth.IsAdmin(),
			middleware.CheckUserIsOwner(
				a.FileService,
			),
		),
	)(endpoints.Files.DeleteFile)

	endpoints.Files.GetFile = endpoint.Chain(
		testAuth(),
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Files.GetFile)
	
	endpoints.Files.UploadFile = endpoint.Chain(
		testAuth(),
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Files.UploadFile)

	endpoints.Files.GetFiles = endpoint.Chain(
		testAuth(),
		mgsess.PutMongoSessIntoCtx(),
		middleware.MergeMiddlewaresIntoOr(
			a.Auth.IsAdmin(),
			middleware.SetUserID(),
		),
	)(endpoints.Files.GetFiles)

	return endpoints
}

func testAuth() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(
			ctx		context.Context,
			request		interface{},
		) (response interface{}, err error) {
			log.Debug("test auth")
			newCtx := subcontext.New(
				rolecontext.New(
					ctx,
					"mfs.admin",
				),
				"mock_user_id",
			)
			return next(newCtx, request)
		}
	}
}