package file

import (
	"time"

	"github.com/MicroFileServer/pkg/models/basemodel"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type BaseFile struct {
	Length		 int64				`json:"length" bson:"length"`
	ChunkSize	 int				`json:"chunkSize" bson:"chunkSize"`
	UploadDate	 time.Time 			`json:"uploadDate" bson:"uploadDate"`
	FileName	 string            	`json:"filename" bson:"filename"`
	Metadata	 Metadata			`json:"metadata" bson:"metadata"`
}

type Metadata struct {
	FileSender			string		`json:"fileSender" bson:"fileSender"`
	FileDescription		string		`json:"fileDescription" bson:"fileDescription"`
}

type File struct {
	ID			string				`json:"id" bson:"_id"`
	*BaseFile						`json:",inline" bson:",inline"`
}

type FileMongoDB struct {
	basemodel.BaseModel				`bson:",inline"`
	*gridfs.File					`bson:",inline"`
}

func (f *FileMongoDB) CollectionName() string {
	return "fs.file"
}