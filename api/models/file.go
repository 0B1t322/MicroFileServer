package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type File struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Length		 int64				`json:"length"`
	ChunkSize	 int				`json:"chunkSize"`
	UploadDate	 primitive.DateTime `json:"uploadDate"`
	FileName	 string            	`json:"filename" `
	Metadata	 Metadata			`json:"metadata"`
}

type Metadata struct {
	FileSender			string		`json:"fileSender"`
	FileDescription		string		`json:"fileDescription"`
}