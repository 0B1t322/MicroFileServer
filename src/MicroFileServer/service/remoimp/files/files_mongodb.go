package files

import (
	"bytes"
	"context"

	"github.com/MicroFileServer/pkg/models/file"
	op_err "github.com/MicroFileServer/pkg/repositories/errors"
	"github.com/MicroFileServer/pkg/repositories/files"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FilesMongoDBImp struct {
	Repo	files.FileRepositorier
}

func (f *FilesMongoDBImp) UploadFile(
	ctx			context.Context,
	fileName	string,
	rawFile		[]byte,
	Metadata	file.Metadata,
) (*file.File, error) {
	bucket, err := f.Repo.NewBucket()
	if err != nil {
		return nil, err
	}

	stream, err := bucket.OpenUploadStream(
		fileName,
		options.GridFSUpload().SetMetadata(
			bson.M{
				"fileSender": Metadata.FileSender,
				"fileDescription": Metadata.FileDescription,
			},
		),
	)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	if _, err := stream.Write(rawFile); err != nil {
		return nil, err
	}

	var file file.File
	{
		if err := f.Repo.GetOne(
			ctx,
			bson.M{
				"_id": stream.FileID,
			},
			func(sr *mongo.SingleResult) error {
				return sr.Decode(&file)
			},
		); err == mongo.ErrNoDocuments {
			return nil, op_err.ErrDocumentNotFound
		} else if err != nil {
			return nil, err
		}
	}

	return &file, nil
}

func (f *FilesMongoDBImp) DownloadFile(
	ctx		context.Context,
	FileID	string,
) ([]byte, error) {
	fileId, err := primitive.ObjectIDFromHex(FileID)
	if err != nil {
		return nil, op_err.ErrNotValidID
	}

	bucket, err := f.Repo.NewBucket()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	_, err = bucket.DownloadToStream(fileId, buf)
	if err == mongo.ErrNoDocuments {
		return nil, op_err.ErrDocumentNotFound
	} else if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (f *FilesMongoDBImp) GetFile(
	ctx		context.Context,
	FileID	string,
) (*file.File, error) {
	fileId, err := primitive.ObjectIDFromHex(FileID)
	if err != nil {
		return nil, op_err.ErrNotValidID
	}

	return f.getFile(
		ctx,
		fileId,
	)
}

func (f *FilesMongoDBImp) getFile(
	ctx		context.Context,
	fileId	primitive.ObjectID,
) (*file.File, error) {
	var file file.File
	{
		if err := f.Repo.GetOne(
			ctx,
			bson.M{
				"_id": fileId,
			},
			func(sr *mongo.SingleResult) error {
				return sr.Decode(&file)
			},
		); err == mongo.ErrNoDocuments {
			return nil, op_err.ErrDocumentNotFound
		} else if err != nil {
			return nil, err
		}
	}

	return &file, nil
}