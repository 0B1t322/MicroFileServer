package files

import (
	e "github.com/MicroFileServer/pkg/err"
	file "github.com/MicroFileServer/pkg/models/file"
)

func init() {
	_ = e.Message{}
	_ = file.File{}
}

// UploadFile
// 
// @Tags files
// 
// @Summary upload file
// 
// @Description upload file to service
// 
// @Security ApiKeyAuth
// 
// @Router /files/upload [post]
// 
// @Param uploadingForm formData file true "file that need to upload"
// 
// @Param fileDescription formData string false "file description"
// 
// @Accept multipart/form-data
// 
// @Produce json
// 
// @Success 200 {object} files.UploadFileResp
// 
// @Failure 401 {object} e.Message
// 
// @Failure 404 {object} e.Message "if file not found after upload"
// 
// @Failure 500 {object} e.Message
func (Server) HTTPUploadFile() {}

// DownloadFile
// 
// @Tags files
// 
// @Summary download file
// 
// @Description download file from service
// 
// @Router /download/{id} [get]
// 
// @Param id path string true "id of the file"
// 
// @Produce */*,jpeg,png,gif,video/*,audio/*,image/*,application/pdf,application/msword,application/vnd.ms-excel
// 
// @Success 200 {file} file
// 
// @Failure 400 {object} e.Message "if file id is not valid"
// 
// @Failure 404 {object} e.Message "if file not found after upload"
// 
// @Failure 500 {object} e.Message
func (Server) HTTPDownloadFile() {}

// DeleteFile
// 
// @Tags files
// 
// @Summary delete file
// 
// @Description delete file from service
// 
// @Description if you not admin you can only delete files that you upload
// 
// @Security ApiKeyAuth
// 
// @Router /files/{id} [delete]
// 
// @Param id path string true "id of the file"
// 
// @Success 200
// 
// @Failure 400 {object} e.Message "if file id is not valid"
// 
// @Failure 404 {object} e.Message "if file not found after upload"
// 
// @Failure 401 {object} e.Message
// 
// @Failure 403 {object} e.Message "if it's not your file"
// 
// @Failure 500 {object} e.Message
func (Server) HTTPDeleteFile() {}

// GetFile
// 
// @Tags files
// 
// @Summary get file info
// 
// @Description get info about file
// 
// @Security ApiKeyAuth
// 
// @Router /files/{id} [get]
// 
// @Param id path string true "id of the file"
// 
// @Produce json
// 
// @Success 200 {object} files.GetFileResp
// 
// @Failure 400 {object} e.Message "if file id is not valid"
// 
// @Failure 404 {object} e.Message "if file not found after upload"
// 
// @Failure 401 {object} e.Message
// 
// @Failure 500 {object} e.Message
func (Server) HTTPGetFile() {}

// GetFiles
// 
// @Tags files
// 
// @Summary get files
// 
// @Description return files info
// 
// @Description if you are not admin you can get info only about you files
// 
// @Security ApiKeyAuth
// 
// @Param user query string false "id of the user which files you want get"
// 
// @Param sorted_by query string false "sort by ascendig; can be name or date "
// 
// @Router /files [get]
// 
// @Produce json
// 
// @Success 200 {array} file.File
// 
// @Failure 401 {object} e.Message
// 
// @Failure 500 {object} e.Message
func (Server) HTTPGetFiles() {}