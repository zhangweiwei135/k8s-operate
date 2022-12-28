package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8sOps/model"
	"net/http"
)

func AddUploads(c *gin.Context)  {
	file, err := c.FormFile("file")

	if err == nil {
		fmt.Println(file.Filename)
		f, _ := file.Open()
		defer f.Close()
		b := UploadFile(model.KubeConfigBucket,file.Filename,f,-1)
		if b {
			c.JSON(http.StatusOK,gin.H{
				"code": 200,
				"msg": "file upload success",
				"fileName": file.Filename,
			})
		}
		//dst := path.Join("./config",file.Filename)
		//c.SaveUploadedFile(file,dst)

	}
}
