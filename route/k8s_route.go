
package route

import (
"github.com/gin-gonic/gin"
	"k8sOps/api"
	"k8sOps/common"

	"net/http"
)

//解决跨域问题
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}

func  InitRouter()  {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(Cors())
	router.NoRoute(common.NotFound)

	cluster := router.Group("kubernetes/clusterInfo")
	{
		cluster.GET("/list",api.ClusterList)
		cluster.POST("/add",api.ClusterAdd)
		cluster.DELETE("/del",api.ClusterDel)
		cluster.PUT("/update",api.ClusterEdit)
		cluster.POST("/upload",api.AddUploads)
		cluster.GET("/nslist",api.NamespaceSearch)
	}

	workload := router.Group("kubernetes/workload")
	{
		workload.GET("/deployment/list",api.DeployList)
		workload.GET("/pod/list",api.PodsList)
		workload.POST("/deployment/replicas/update",api.UpdateDeployReplicas)
		workload.DELETE("/deployment/del",api.DeleteDeployment)
		workload.POST("/code/edit",api.ResYaml)  //部署yaml
		workload.GET("/pod/log",api.GetPodLogs)  //部署yaml
	}

	pod := router.Group("pod")
	{
		pod.GET("/terminal",api.PodConnect)
	}

	sts := router.Group("kubernetes/workload/statefulSet")
	{
		sts.GET("/list",api.StsList)
		sts.PUT("/replicas/update",api.UpdateStsReplicas)
		sts.DELETE("/del",api.DeleteSts)
	}

	ds := router.Group("kubernetes/workload/ds")
	{
		ds.GET("/list",api.DsList)
		ds.DELETE("/del",api.DeleteDs)
	}

	svc := router.Group("kubernetes/workload/svc")
	{
		svc.GET("/list",api.SvcList)
		svc.DELETE("/del",api.DeleteSvc)
	}

	ing := router.Group("kubernetes/workload/ingress")
	{
		ing.GET("/list",api.IngList)
		ing.DELETE("/del",api.DeleteIng)
	}

	config := router.Group("kubernetes/workload/cm")
	{
		config.GET("/list",api.CmList)
		config.DELETE("/del",api.DeleteCM)
	}

	secret := router.Group("kubernetes/workload/secret")
	{
		secret.GET("/list",api.SecretList)
		secret.DELETE("/del",api.DeleteSecret)
	}


	router.Run(":8083")

}

