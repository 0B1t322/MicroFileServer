package v1

import (
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
			middleware.ValidateAndSetUserID(),
		),
	)(endpoints.Files.GetFiles)


	endpoints.Files.DownloadFile = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Files.DownloadFile)
	return endpoints
}

func (a *Api) buildAMQPEndpoints() ApiEndpoints {
	endpoints := a.endpoints()

	endpoints.Files.DeleteFile = endpoint.Chain(
		mgsess.PutMongoSessIntoCtx(),
	)(endpoints.Files.DeleteFile)
	
	return endpoints
}