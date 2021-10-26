package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/MicroFileServer/pkg/repositories/files"
	mgm "github.com/kamva/mgm/v3"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type Repositories struct {
	File	files.FileRepositorier
}

type Config struct {
	DBURI string
}

func New(cfg *Config) (*Repositories, error) {
	client, err := mongo.NewClient(
		options.Client().
			ApplyURI(cfg.DBURI).
			SetMaxConnIdleTime(0).
			SetLocalThreshold(10*time.Millisecond),
	)
	if err != nil {
		return nil, errors.New("Error on created client")
	}
	defer client.Disconnect(context.Background())
	

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		return nil, errors.New("Error on connection")
	}

	ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("Error on ping: %v", err)
	}
	constr, err := connstring.Parse(cfg.DBURI)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI: err: %v", err)
	}
	
	dbName := constr.Database
	
	if err := mgm.SetDefaultConfig(
		nil,
		dbName,
		options.Client().ApplyURI(cfg.DBURI),
	); err != nil {
		return nil, err
	}

	
	return &Repositories{
		File: files.New(),
	}, 
	nil
}