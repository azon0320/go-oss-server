package main

import (
	"github.com/dormao/go-oss-server/internal/context"
	"github.com/dormao/go-oss-server/internal/context/config"
	"github.com/dormao/go-oss-server/internal/context/database"
	"github.com/dormao/go-oss-server/internal/filesystem"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	var logger = logrus.New()
	logger.SetOutput(os.Stdout)
	var conf = config.Config
	var dataProvider database.Provider
	if conf.Provider.Type == "yaml" {
		dataProvider = &database.YamlProvider{Bucket: conf.Bucket}
	} else if conf.Provider.Type == "postgres" {
		dataProvider = &database.PostgresProvider{
			Bucket:      conf.Bucket,
			ResourceURI: conf.Provider.PostgresURI,
			DB:          nil,
		}
	}
	var ctrl = &context.Controller{
		DataProvider: dataProvider,
		FileProvider: filesystem.CreateFileSystem(conf.StorePath, 0755),
		Logger:       logger,
	}
	err := ctrl.Init()
	if err != nil {
		ctrl.Logger.Errorf("error while init controller: %s", err)
		return
	}
	var corsHosts = cors.New(cors.Config{
		AllowOrigins:     strings.Split(conf.CORSHosts, ","),
		AllowMethods:     []string{"GET", "PATCH", "POST", "DELETE", "OPTION", "PUT"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	})
	var engine = gin.Default()
	engine.Use(corsHosts)
	ctrl.RegisterRoutes(engine)
	engine.Run(conf.ServeAddress)
}
