package file

import (
	"time"

	"github.com/MicroFileServer/pkg/models/basemodel"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type BaseFile struct {
	Length		 int64				`json:"length"`
	ChunkSize	 int				`json:"chunkSize"`
	UploadDate	 time.Time 			`json:"uploadDate"`
	FileName	 string            	`json:"filename" `
	Metadata	 Metadata			`json:"metadata"`
}

type Metadata struct {
	FileSender			string		`json:"fileSender"`
	FileDescription		string		`json:"fileDescription"`
}

type File struct {
	ID			string				`json:"id"`
	*BaseFile						`json:",inline"`
}

type FileMongoDB struct {
	basemodel.BaseModel				`bson:",inline"`
	*gridfs.File					`bson:",inline"`
}

func (f *FileMongoDB) CollectionName() string {
	return "fs.file"
}