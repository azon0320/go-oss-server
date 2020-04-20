package context

import (
	"fmt"
	"github.com/dormao/go-oss-server/internal/context/config"
	"github.com/dormao/go-oss-server/internal/context/consts"
	"github.com/dormao/go-oss-server/internal/context/database"
	"github.com/dormao/go-oss-server/internal/context/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
)

func (ctrl *Controller) ginDeleteObject(c *gin.Context){
	var ak = c.PostForm("accesskey")
	var secret = c.PostForm("secret")
	var bkName = c.PostForm("bucket")
	var objectName = strings.Trim(c.PostForm("object"), "/")
	if len(bkName) == 0 {
		bkName = config.Config.Bucket
	}
	err := ctrl.DataProvider.RemoveObject(objectName, database.RetrieveOption{
		AccessKey: ak, AccessSecret: secret, Bucket: &bkName,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.BaseResponse{
			Code: http.StatusBadRequest, Msg: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &models.BaseResponse{
		Code: http.StatusOK, Msg: "",
	})
}

func (ctrl *Controller) ginListObject(c *gin.Context){
	var ak = c.PostForm("accesskey")
	var secret = c.PostForm("secret")
	var bkName = c.PostForm("bucket")
	if len(bkName) == 0 {
		bkName = config.Config.Bucket
	}
	var objectPref = c.Request.URL.Query().Get("list")
	var result, err = ctrl.DataProvider.ListObject(objectPref, database.RetrieveOption{
		Bucket:       &bkName,
		AccessKey:    ak,
		AccessSecret: secret,
	})
	if err != nil {
		logrus.Errorf("error while listing objects in bucket (%s): %s", bkName, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.BaseResponse{
			Code: http.StatusBadRequest, Msg: err.Error(),
		})
		return
	}
	var resultList = result.([]string)
	c.JSON(http.StatusOK, &models.BaseResponse{
		Code: http.StatusOK, Msg: "", Result: gin.H{
			"count": len(resultList),
			"result": &resultList,
		},
	})
}

func (ctrl *Controller) ginUploadObject(c *gin.Context) {

	var queries = c.Request.URL.Query()
	if queries.Get("list") != "" {
		ctrl.ginListObject(c)
		return
	}

	type UploadResult struct {
		Object string `json:"object"`
		URL    string `json:"url"`
	}
	parseExtension := func(filename string) string {
		nodes := strings.Split(filename, ".")
		var ext = ""
		if len(nodes) > 0 {
			ext = nodes[len(nodes)-1]
		}
		return ext
	}
	var ak = c.PostForm("accesskey")
	var secret = c.PostForm("secret")
	var objectName = strings.Trim(c.PostForm("object"), "/")
	var bkName = c.PostForm("bucket")
	if len(bkName) == 0 {
		bkName = config.Config.Bucket
	}
	if objectName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.BaseResponse{Code: http.StatusBadRequest, Msg: "object name can not be empty"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.BaseResponse{Code: http.StatusBadRequest, Msg: err.Error()})
		return
	}
	rand, _ := uuid.NewRandom()
	var ext = parseExtension(file.Filename)
	if ext != "" {
		ext = "." + ext
	}
	var fullFileName = fmt.Sprintf("%s%s", rand.String(), ext)

	err = ctrl.DataProvider.PutObject(objectName, fullFileName, database.SetUpOption{
		RetrieveOption: database.RetrieveOption{
			Bucket:       &bkName,
			AccessKey:    ak,
			AccessSecret: secret,
		},
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.BaseResponse{
			Code: http.StatusBadRequest, Msg: err.Error(),
		})
		return
	}

	err = c.SaveUploadedFile(file, ctrl.FileProvider.Path(fullFileName))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.BaseResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
		return
	}

	ctrl.DataProvider.Save()
	c.JSON(http.StatusOK, &models.BaseResponse{
		Code: http.StatusOK, Msg: "",
		Result: &UploadResult{
			Object: objectName, URL: ctrl.GetFullUrl(objectName),
		},
	})
}

func (ctrl *Controller) ginPublicGetResource(c *gin.Context) {
	var bkName = c.Param("bucket")
	var objectName = strings.Trim(c.Param("object"), "/")
	var path string
	if objectName == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	path, err := ctrl.DataProvider.GetObject(objectName, database.RetrieveOption{
		Bucket: &bkName,
	})
	if err != nil {
		ctrl.Logger.Infof("object key (%s) not found", objectName)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if config.Config.OutputMode == consts.OutputModeReDirect {
		// TODO disabled now, reason: controller.go:59
		ctrl.Logger.Warnf("Redirect now disabled! please use ( outputmode:'direct' ) in config file!!")
		c.AbortWithStatus(http.StatusInternalServerError)
		//c.Redirect(http.StatusTemporaryRedirect, ctrl.GetStatic(path))
	} else if config.Config.OutputMode == consts.OutputModeDirect {
		var fPath = ctrl.FileProvider.Path(path)
		f, err := os.Open(fPath)
		if err != nil {
			ctrl.Logger.Errorf("error open file: %s", err)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		defer f.Close()
		_, err = io.Copy(c.Writer, f)
		if err != nil {
			// TODO will occur error when user cancel binary download
			ctrl.Logger.Errorf("error while echo file bytes: %s", err)
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
}

func (ctrl *Controller) ginOptions(c *gin.Context) {
	c.Status(http.StatusOK)
}
