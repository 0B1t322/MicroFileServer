package files

import (
	"github.com/MicroFileServer/pkg/repositories/agregate"
	"github.com/MicroFileServer/pkg/repositories/deleter"
	"github.com/MicroFileServer/pkg/repositories/getter"
	"github.com/MicroFileServer/pkg/repositories/saver"
	"github.com/MicroFileServer/pkg/repositories/updater"
)

type FileRepositorier interface{
	getter.Getter
	deleter.Deleter
	updater.Updater
	agregate.Agregater
	saver.Saver
}