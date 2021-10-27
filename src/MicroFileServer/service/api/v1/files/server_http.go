package files

import (
	"net/http"

	"github.com/MicroFileServer/service/encoder"
	"github.com/MicroFileServer/service/errorencoder"
	"github.com/MicroFileServer/service/middleware"
	serverbefore "github.com/MicroFileServer/service/serverbefore/http"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type Config struct {
	MaxFileSizeMB	int64
}

func NewHTTPServer(
	cfg			*Config,
	enpoints	Endpoints,
	r			*mux.Router,
) {
	// UploadFile
	r.Handle(
		"/files/upload",
		middleware.HTTPSetMaxFileSize(cfg.MaxFileSizeMB)(
			httptransport.NewServer(
				enpoints.UploadFile,
				HTTPDecodeUploadFileReq,
				encoder.EncodeResponce,
				httptransport.ServerBefore(
					serverbefore.TokenFromReq,
				),
				httptransport.ServerErrorEncoder(
					errorencoder.ErrorEncoder,
				),
			),
		),
	).Methods(http.MethodPost)

	// DownloadFile
	r.Handle(
		"/files/download/{id}",
		httptransport.NewServer(
			enpoints.DownloadFile,
			HTTPDecodeDownloadFileReq,
			encoder.EncodeResponce,
			httptransport.ServerBefore(
				serverbefore.TokenFromReq,
				serverbefore.PutRequestInCTX,
			),
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods(http.MethodGet)

	// DeleteFile
	r.Handle(
		"/files/{id}",
		httptransport.NewServer(
			enpoints.DeleteFile,
			HTTPDecodeDeleteFileReq,
			encoder.EncodeResponce,
			httptransport.ServerBefore(
				serverbefore.TokenFromReq,
			),
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods(http.MethodDelete)

	// GetFile
	r.Handle(
		"/files/{id}",
		httptransport.NewServer(
			enpoints.GetFile,
			HTTPDecodeGetFileReq,
			encoder.EncodeResponce,
			httptransport.ServerBefore(
				serverbefore.TokenFromReq,
			),
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods(http.MethodGet)

	// GetFiles
	r.Handle(
		"/files",
		httptransport.NewServer(
			enpoints.GetFiles,
			HTTPDecodeGetFilesReq,
			encoder.EncodeResponce,
			httptransport.ServerBefore(
				serverbefore.TokenFromReq,
			),
			httptransport.ServerErrorEncoder(
				errorencoder.ErrorEncoder,
			),
		),
	).Methods(http.MethodGet)
}