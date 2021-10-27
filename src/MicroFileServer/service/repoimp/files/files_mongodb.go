package files

import (
	"bytes"
	"context"

	"github.com/MicroFileServer/pkg/models/file"
	op_err "github.com/MicroFileServer/pkg/repositories/errors"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/MicroFileServer/pkg/repositories/files"
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

	if _, err := stream.Write(rawFile); err != nil {
		stream.Close()
		return nil, err
	}
	fileId := stream.FileID
	stream.Close()
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

func (f *FilesMongoDBImp) DeleteFile(
	ctx		context.Context,
	FileId	string,
) error {
	fileId, err := primitive.ObjectIDFromHex(FileId)
	if err != nil {
		return op_err.ErrNotValidID
	}

	bucket, err := f.Repo.NewBucket()
	if err != nil {
		return err
	}

	if err := bucket.Delete(fileId); err == mongo.ErrNoDocuments {
		return op_err.ErrDocumentNotFound
	} else if err != nil {
		return err
	}

	return nil
}

type getFilesQueryBuilder struct {
	filter	bson.M
	sort	bson.D
}

func (g *getFilesQueryBuilder) SetFileSender(userid string) GetFilesBuilder {
	g.filter["fileSender"] = userid
	return g
}

func (g *getFilesQueryBuilder) SetDateSort() GetFilesBuilder {
	g.sort = append(g.sort, bson.E{Key: "uploadDate", Value: 1})
	return g
}

func (g *getFilesQueryBuilder) SetNameSort() GetFilesBuilder {
	g.sort = append(g.sort, bson.E{Key: "metadata.fileSender", Value: 1})
	return g
}

func (f *FilesMongoDBImp) GetFilesBuilder() GetFilesBuilder {
	return &getFilesQueryBuilder{
		filter: bson.M{},
		sort: bson.D{},
	}
}

func (f *FilesMongoDBImp) GetFiles(
	ctx		context.Context,
	Builder	GetFilesBuilder,
) ([]*file.File, error) {
	builder := Builder.(*getFilesQueryBuilder)

	opt := options.Find()
	{
		if len(builder.sort) != 0 {
			opt = opt.SetSort(builder.sort)
		}
	}

	files := []*file.File{}
	{
		if err := f.Repo.GetAllFiltered(
			ctx,
			builder.filter,
			func(c *mongo.Cursor) error {
				if c.RemainingBatchLength() == 0 {
					return op_err.ErrDocumentNotFound
				}

				return c.All(ctx, &files)
			},
			opt,
		); err == op_err.ErrDocumentNotFound {
			// Pass
		} else if err != nil {
			return nil, err
		}
	}

	return files, nil
}