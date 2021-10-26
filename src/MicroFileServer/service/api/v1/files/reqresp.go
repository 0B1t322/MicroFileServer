package files

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/MicroFileServer/pkg/contextvalue/reqcontext"
	"github.com/MicroFileServer/pkg/contextvalue/subcontext"
	"github.com/MicroFileServer/pkg/models/file"
	"github.com/MicroFileServer/pkg/statuscode"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type DownloadFileReq struct {
	ID	string
}

func HTTPDecodeDownloadFileReq(
	ctx		context.Context,
	r		*http.Request,
) (interface{}, error) {
	vars := mux.Vars(r)
	req := &DownloadFileReq{
		ID: vars["id"],
	}

	return req, nil
}

type DownloadFileResp struct {
	File	*file.File
	RawFile	[]byte
}

func (resp *DownloadFileResp) StatusCode() int {
	return http.StatusOK
}

func (resp *DownloadFileResp) EncodeCTX(
	ctx context.Context, 
	w http.ResponseWriter,
) error {
	httpReq, err := reqcontext.GetRequestFromContext(ctx)
	if err != nil {
		log.WithFields(
			log.Fields{
				"func": "DownloadFileResp.EncodeCTX",
				"err": err,
			},
		).Error()
		return statuscode.WrapStatusError(
			fmt.Errorf("Failed to encode file"),
			http.StatusInternalServerError,
		)
	}

	http.ServeContent(
		w, 
		httpReq, 
		resp.File.FileName, 
		time.Now(), 
		bytes.NewReader(resp.RawFile),
	)

	return nil
}

func (resp *DownloadFileResp) Headers(
	ctx context.Context, 
	w http.ResponseWriter,
) {
	mime := mimetype.Detect(resp.RawFile)
	if !strings.Contains(mime.String(), "video") && !strings.Contains(mime.String(), "audio") {
		w.Header().Set("Content-Type", mime.String())
		w.Header().Set("Content-Disposition", "attachment; filename=\""+resp.File.FileName+"\"")
	}
}

type UploadFileReq struct {
	RawFile		[]byte
	FileSender	string
	FileDesc	string
	FileName	string
}

func HTTPDecodeUploadFileReq(
	ctx		context.Context,
	r		*http.Request,
) (interface{}, error) {
	data, handler, err := r.FormFile("uploadingForm")
	if err != nil {
		return nil, statuscode.WrapStatusError(
			fmt.Errorf("File is not appropriate"),
			http.StatusBadRequest,
		)
	}

	rawFile, err := ioutil.ReadAll(data)
	if err != nil {
		log.Errorf("Failed to read file data: %v", err)
		return nil, statuscode.WrapStatusError(
			fmt.Errorf("Failed to read file data"),
			http.StatusInternalServerError,
		)
	}

	desc := r.FormValue("fileDescription")
	fileSender, err := subcontext.GetSubFromContext(ctx)
	if err != nil {
		log.Errorf("Don't find sub id")
		return nil, statuscode.WrapStatusError(
			fmt.Errorf("Failed to find sub id"),
			http.StatusInternalServerError,
		)
	}

	req := &UploadFileReq{
		RawFile: rawFile,
		FileSender: fileSender,
		FileDesc: desc,
		FileName: handler.Filename,
	}

	return req, nil
}

type UploadFileResp struct {
	File	*file.File
}

func (resp *UploadFileResp) StatusCode() int {
	return http.StatusCreated
}

func (resp *UploadFileResp) Encode(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(resp.File)
}

func (resp *UploadFileResp) Headers(ctx context.Context, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}