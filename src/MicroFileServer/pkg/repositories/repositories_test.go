package repositories_test

import "github.com/MicroFileServer/pkg/repositories"

var Repositories *repositories.Repositories

func init() {
	repositories.New(&repositories.Config{DBURI: "mongodb://"})
}