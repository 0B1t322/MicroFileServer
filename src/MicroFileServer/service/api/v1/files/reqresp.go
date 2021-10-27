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
	"github.com/MicroFileServer/pkg/models/file"
	"github.com/MicroFileServer/pkg/statuscode"
	"github.com/MicroFileServer/pkg/urlvalue/encode"
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

	req := &UploadFileReq{
		RawFile: rawFile,
		FileDesc: desc,
		FileName: handler.Filename,
	}

	return req, nil
}

type UploadFileResp struct {
	*file.File	`json:",inline"`
}

func (resp *UploadFileResp) StatusCode() int {
	return http.StatusCreated
}

func (resp *UploadFileResp) Encode(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(resp)
}

func (resp *UploadFileResp) Headers(ctx context.Context, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

type DeleteFileReq struct {
	FileID	string
}

func (d *DeleteFileReq) GetID() string {
	return d.FileID
}

func HTTPDecodeDeleteFileReq(
	ctx		context.Context,
	r		*http.Request,
) (interface{}, error) {
	vars := mux.Vars(r)
	req := &DeleteFileReq{
		FileID: vars["id"],
	}

	return req, nil
}

type DeleteFileResp struct {
}

func (resp *DeleteFileResp) StatusCode() int {
	return http.StatusOK
}

func (resp *DeleteFileResp) Encode(w http.ResponseWriter) error {
	return nil
}

func (resp *DeleteFileResp) Headers(ctx context.Context, w http.ResponseWriter) {
	return
}

type GetFileReq struct {
	FileID	string
}

func HTTPDecodeGetFileReq(
	ctx		context.Context,
	r		*http.Request,
) (interface{}, error) {
	vars := mux.Vars(r)
	req := &GetFileReq{
		FileID: vars["id"],
	}

	return req, nil
}

type GetFileResp struct {
	*file.File	`json:",inline"`
}

func (resp *GetFileResp) StatusCode() int {
	return http.StatusOK
}

func (resp *GetFileResp) Encode(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(resp)
}

func (resp *GetFileResp) Headers(ctx context.Context, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

type GetFilesReq struct {
	Query	GetFilesQuery
}

func (g *GetFilesReq) SetUserID(userid string) {
	g.Query.UserID = userid
}

type GetFilesQuery struct {
	UserID		string	`query:"user,string"`
	SortedBy	string	`query:"sorted_by,string"`
}

func HTTPDecodeGetFilesReq(
	ctx		context.Context,
	r		*http.Request,
) (interface{}, error) {
	req := &GetFilesReq{}

	if err := encode.UrlQueryUnmarshall(
		&req.Query,
		r.URL.Query(),
	); err != nil {
		return nil, statuscode.WrapStatusError(
			fmt.Errorf("Failed to decode request"),
			http.StatusBadRequest,
		)
	}

	return req, nil
}

type GetFilesResp struct {
	Files	[]*file.File	`json:"files"`
}

func (resp *GetFilesResp) StatusCode() int {
	return http.StatusOK
}

func (resp *GetFilesResp) Encode(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(resp)
}

func (resp *GetFilesResp) Headers(ctx context.Context, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}