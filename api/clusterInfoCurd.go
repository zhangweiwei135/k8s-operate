package api

import (
	"context"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sOps/model"
	"log"
	"time"

	"net/http"
)

//判断集群是否存在
func IsExist(name string) bool  {
	UserSlice := K8sSel(name)
	if len(UserSlice) != 0 {
		return true
	}
	return false
}

//获取集群
func ClusterList(c *gin.Context)  {
	k := K8sList()

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": k,
		"msg": "seccuss",
	})
}

//添加集群
func ClusterAdd(c *gin.Context)  {
	name := c.Request.FormValue("name")
	token := c.Request.FormValue("config")
	fileName := c.Request.FormValue("filename")


	k := model.ClusterInfo{
		Name: name,
		Token: token,
		KubeConfigPath: fileName,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	Bool := IsExist(name)
	if Bool {
		//集群是否存在
		model.ClusterState["state"] = 2
		model.ClusterState["text"] = "此集群已经存在!"
	} else {
		K8sAdd(&k)
		model.ClusterState["state"] = 1
		model.ClusterState["text"] = "集群新增成功！"
	}

	if name == "" {
		model.ClusterState["state"] = 2
		model.ClusterState["text"] = "集群为空，请重新输入"
	}

	if model.ClusterState["state"] == 1 {
		c.JSON(http.StatusOK,gin.H{
			"code": 200,
			"data": "test",
			"msg": model.ClusterState,
		})
	} else {
		c.JSON(http.StatusBadRequest,gin.H{
			"code": 400,
			"data": "test",
			"msg": model.ClusterState,
		})
	}
}

//删除集群
func ClusterDel(c *gin.Context)  {
	name := c.Request.FormValue("name")
	k := K8sSel(name)
	if len(k) == 0 {
		c.JSON(http.StatusNotFound,gin.H{
			"msg": "集群不存在",
		})
		return
	}

	res := K8sDel(name)
	if res {
		c.JSON(http.StatusOK,gin.H{
			"code": 200,
			"msg": "success",
		})
	} else {
		c.JSON(http.StatusBadRequest,gin.H{
			"code": 400,
			"msg": "集群删除失败",
		})
	}

}

//编辑集群
func ClusterEdit(c *gin.Context)  {
	name := c.Request.FormValue("name")
	token := c.Request.FormValue("config")
	oldName := c.Request.FormValue("oldname")

	k := &model.ClusterInfo{
		Name: name,
		Token: token,
	}

	selK := K8sSel(oldName)
	if len(selK) == 0 {
		c.JSON(http.StatusBadRequest,gin.H{
			"code": 400,
			"msg": "未获取到集群信息或者输入的名为空",
		})
		return
	}

	selk2 := K8sSel(name)
	if len(selk2) != 0 {
		c.JSON(http.StatusBadRequest,gin.H{
			"code": 400,
			"msg": "集群名称已存在",
		})
		return
	}

	res := K8sUpdate(oldName,k)
	if res {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg": "集群信息修改成功",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg": "信息修改失败",
		})
	}
}


//查询命名空间
func NamespaceSearch(c *gin.Context)  {
	cluster_name := c.Query("cluster_name")

	client, _ := InitClientSet(cluster_name)
	ns, err := client.CoreV1().Namespaces().List(context.TODO(),metav1.ListOptions{})

	if err != nil {
		log.Println(err)
		return
	}

	var nsList []string
	for _, v := range ns.Items {
		nsList = append(nsList,v.Name)
	}

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": nsList,
		"msg": "search ns ok",
	})
}
