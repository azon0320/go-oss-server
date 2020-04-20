package context

import (
	"fmt"
	"github.com/dormao/go-oss-server/internal/context/config"
	"github.com/dormao/go-oss-server/internal/context/consts"
	"github.com/dormao/go-oss-server/internal/context/database"
	"github.com/dormao/go-oss-server/internal/filesystem"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
)

type Controller struct {
	DataProvider database.Provider
	FileProvider *filesystem.FileSystem
	Logger       *logrus.Logger
}

func (ctrl *Controller) GetFullUrl(object string) string {
	return fmt.Sprintf("%s/%s/%s",
		config.Config.BaseURL,
		config.Config.Bucket,
		strings.Trim(object, ctrl.FileProvider.FileSeparator()))
}
func (ctrl *Controller) GetStatic(filename string) string {
	return fmt.Sprintf("%s/%s/%s", config.Config.BaseURL, consts.StaticFSRouteName, strings.Trim(filename, ctrl.FileProvider.FileSeparator()))
}
func (ctrl *Controller) Init() error {
	return ctrl.DataProvider.Init()
}

func (ctrl *Controller) RegisterRoutes(gin *gin.Engine) {
	// @api {POST} /:bucket/:object
	// @apiGroup 1.Upload
	// @apiVersion 1.0.0

	gin.OPTIONS("/", ctrl.ginOptions)

	// @apiParam {String} bucket The oss bucket
	// @apiParam {String} object The object key
	// @apiParam {String} accesskey The access key
	// @apiParam {String} secret The access secret
	// @apiParam {Blob} file The upload file
	// @apiSuccess {Number} code The return code (most sync with HTTP status code)
	// @apiSuccess {String} msg The return msg (blank when no error)
	// @apiSuccess {Object} result The upload result
	// @apiSuccess {String} result.object The upload object key
	// @apiSuccess {String} result.url The public url of the uploaded file
	gin.POST("*object", ctrl.ginUploadObject)

	// @api {POST} /:bucket/:object
	// @apiGroup 1.Upload
	// @apiVersion 1.0.0

	// @apiParam {String} bucket The oss bucket
	// @apiParam {String} object The object key
	// @apiParam {Blob} The data
	gin.GET(":bucket/*object", ctrl.ginPublicGetResource)

	/* TODO conflicted with :bucket/*object closing
	if config.Config.OutputMode == consts.OutputModeReDirect {
		gin.StaticFS(consts.StaticFSRouteName, http.Dir(config.Config.StorePath))
	}*/
}
