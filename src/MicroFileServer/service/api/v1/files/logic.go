package files

import (
	"context"
	"errors"
	"net/http"

	"github.com/MicroFileServer/pkg/models/file"
	op_err "github.com/MicroFileServer/pkg/repositories/errors"
	"github.com/MicroFileServer/pkg/statuscode"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var (
	ErrFileNotFound			= errors.New("File not found")
	ErrFailedUploadFile		= errors.New("Fail to upload file")
	FileIDIsNotValid		= errors.New("File id is not valid")
	FailedToDownloadFile	= errors.New("Failed to download file")
	ErrFailedToDeleteFile	= errors.New("Failed to delete file")
	ErrFailedToGetFile		= errors.New("Failed to get file")
	ErrFailedGetFiles		= errors.New("Failed to get files")
	ErrYouAreNotOwner		= errors.New("You are not file owner")
)

type Server struct {
	Repo 		Repository
	Logger		log.Logger
}

func New(
	Repo	Repository,
	logger	log.Logger,
) *Server {
	return &Server{
		Repo: Repo,
		Logger: log.With(logger, "service", "files"),
	}
}

func (s *Server) UploadFile(
	ctx		context.Context,
	req		*UploadFileReq,
) (*file.File, error) {
	logger := log.With(s.Logger, "method", "UploadFile")
	f, err := s.Repo.UploadFile(
		ctx,
		req.FileName,
		req.RawFile,
		file.Metadata{
			FileSender: req.FileSender,
			FileDescription: req.FileDesc,
		},
	)
	if err == op_err.ErrDocumentNotFound {
		level.Info(logger).Log("upload file not found: err", err)
		return nil, statuscode.WrapStatusError(
			ErrFileNotFound,
			http.StatusNotFound,
		)
	} else if err != nil {
		level.Error(logger).Log("err", err)
		return nil, statuscode.WrapStatusError(
			ErrFailedUploadFile,
			http.StatusInternalServerError,
		)
	}

	return f, nil
}

func (s *Server) DownloadFile(
	ctx		context.Context,
	FileID	string,
) (*DownloadFileResp, error) {
	logger := log.With(s.Logger, "method", "DownloadFile")
	rawFile, err := s.Repo.DownloadFile(
		ctx,
		FileID,
	)
	if err == op_err.ErrNotValidID {
		return nil, statuscode.WrapStatusError(
			FileIDIsNotValid,
			http.StatusBadRequest,
		)
	} else if err == op_err.ErrDocumentNotFound {
		return nil, statuscode.WrapStatusError(
			ErrFileNotFound,
			http.StatusNotFound,
		)
	} else if err != nil {
		level.Error(logger).Log("err", err)
		return nil, statuscode.WrapStatusError(
			FailedToDownloadFile,
			http.StatusInternalServerError,
		)
	}

	file, err := s.Repo.GetFile(
		ctx,
		FileID,
	)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, statuscode.WrapStatusError(
			FailedToDownloadFile,
			http.StatusInternalServerError,
		)
	}

	return &DownloadFileResp{
		File: file,
		RawFile: rawFile,
	}, nil
}

func (s *Server) DeleteFile(
	ctx 	context.Context,
	FileID	string,
) error {
	logger := log.With(s.Logger, "method", "DeleteFile")

	if err := s.Repo.DeleteFile(
		ctx,
		FileID,
	); err == op_err.ErrDocumentNotFound {
		return statuscode.WrapStatusError(
			ErrFileNotFound,
			http.StatusNotFound,
		)
	} else if err != nil {
		level.Error(logger).Log("err", err)
		return statuscode.WrapStatusError(
			ErrFailedToDeleteFile,
			http.StatusInternalServerError,
		)
	}

	return nil
}

func (s *Server) GetFile(
	ctx		context.Context,
	FileID	string,
) (*file.File, error) {
	logger := log.With(s.Logger, "method", "GetFile")

	file, err := s.Repo.GetFile(
		ctx,
		FileID,
	)
	if err == op_err.ErrNotValidID {
		return nil, statuscode.WrapStatusError(
			FileIDIsNotValid,
			http.StatusBadRequest,
		)
	} else if err == op_err.ErrDocumentNotFound {
		return nil, statuscode.WrapStatusError(
			ErrFileNotFound,
			http.StatusNotFound,
		)
	} else if err != nil {
		level.Error(logger).Log("err", err)
		return nil, statuscode.WrapStatusError(
			ErrFailedToGetFile,
			http.StatusInternalServerError,
		)
	}

	return file, nil
}

func (s *Server) GetFiles(
	ctx		context.Context,
	Query	GetFilesQuery,
) ([]*file.File, error) {
	logger := log.With(s.Logger, "method", "GetFiles")

	builder := s.Repo.GetFilesBuilder()
	if Query.UserID != "" {
		builder = builder.SetFileSender(Query.UserID)
	}
	if Query.SortedBy == "name" {
		builder = builder.SetNameSort()
	} else if Query.SortedBy == "data" {
		builder = builder.SetDateSort()
	}

	files, err := s.Repo.GetFiles(
		ctx,
		builder,
	)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, statuscode.WrapStatusError(
			ErrFailedGetFiles,
			http.StatusInternalServerError,
		)
	}
	return files, nil
}

func (s *Server) IsOwner(
	ctx		context.Context,
	FileID	string,
	UserID	string,
) error {
	file, err := s.GetFile(
		ctx,
		FileID,
	)
	if err != nil {
		return err
	}

	if file.Metadata.FileSender != UserID {
		return statuscode.WrapStatusError(
			ErrYouAreNotOwner,
			http.StatusForbidden,
		)
	}

	return nil
}