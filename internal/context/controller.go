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

func (ctrl *Controller) RegisterRoutes(router *gin.Engine) {

	router.OPTIONS("/", ctrl.ginOptions)


	router.POST("/", func(c *gin.Context) {
		var queries = c.Request.URL.Query()
		if queries.Get("list") != "" {
			// @api {POST} /?list={list} Query objects by input object prefix
			// @apiGroup 1.Object
			// @apiVersion 1.0.0
			// @apiParam {String} bucket The oss bucket
			// @apiParam {String} accesskey The access key
			// @apiParam {String} secret The access secret
			// @apiParam {String} list The object prefix
			// @apiParam {Blob} file The upload file
			// @apiSuccess {Number} code The return code (most sync with HTTP status code)
			// @apiSuccess {String} msg The return msg (blank when no error)
			// @apiSuccess {Array} result The retrieved objects
			ctrl.ginListObject(c)
		}else {
			// @api {POST} / Upload object
			// @apiGroup 1.Object
			// @apiVersion 1.0.0
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
			ctrl.ginUploadObject(c)
		}
	})

	// @api {GET} /:bucket/:object
	// @apiGroup 1.Object
	// @apiVersion 1.0.0

	// @apiParam {String} bucket The oss bucket
	// @apiParam {String} object The object key
	// @apiParam {Blob} The data
	router.GET(":bucket/*object", ctrl.ginPublicGetResource)

	// @api {DELETE} / Delete object by object key
	// @apiGroup 1.Object
	// @apiVersion 1.0.0

	// @apiParam {String} bucket The oss bucket
	// @apiParam {String} object The object key
	// @apiParam {String} accesskey The access key
	// @apiParam {String} secret The access secret
	// @apiSuccess {Number} code The return code (most sync with HTTP status code)
	// @apiSuccess {String} msg The return msg (blank when no error)
	router.DELETE("/", ctrl.ginDeleteObject)

	/* TODO conflicted with :bucket/*object closing
	if config.Config.OutputMode == consts.OutputModeReDirect {
		gin.StaticFS(consts.StaticFSRouteName, http.Dir(config.Config.StorePath))
	}*/
}
