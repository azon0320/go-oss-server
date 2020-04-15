package main

import (
	"github.com/dormao/go-oss-server/internal/context"
	"github.com/dormao/go-oss-server/internal/context/config"
	"github.com/dormao/go-oss-server/internal/context/database"
	"github.com/dormao/go-oss-server/internal/filesystem"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	var logger = logrus.New()
	logger.SetOutput(os.Stdout)
	var conf = config.Config
	var dataProvider = database.YamlProvider{Bucket: conf.Bucket}
	var ctrl = &context.Controller{
		DataProvider: &dataProvider,
		FileProvider: filesystem.CreateFileSystem(conf.StorePath, 0755),
		Logger:       logger,
	}
	ctrl.Init()
	var engine = gin.Default()
	ctrl.RegisterRoutes(engine)
	engine.Run(conf.ServeAddress)
}
