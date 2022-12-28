package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sOps/model"
	"net/http"
)

//获取ds列表信息
func DaemonSets(clusterName,namespace string) *[]model.DaemonSetInfo {
	clientSet, _ := InitClientSet(clusterName)
	dss, err := clientSet.AppsV1().DaemonSets(namespace).List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取ds列表失败 %v \n",err)
		return nil
	}

	var dsInfo []model.DaemonSetInfo
	var ds model.DaemonSetInfo

	for _, v := range dss.Items {
		ds.Name = v.Name
		ds.Labels = v.Labels
		ds.Selector = v.Spec.Selector.String()
		ds.HostName = v.Spec.Template.Spec.Hostname
		ds.Namespace =  v.Namespace
		dsInfo = append(dsInfo,ds)
	}
	return &dsInfo
}


//根据集群获取ds列表
func DsList(c *gin.Context)  {
	clusterName := c.Request.FormValue("cluster_name")
	namespaceRes := c.Request.FormValue("namespace")
	d := DaemonSets(clusterName,namespaceRes)

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": d,
		"msg": "success",
	})
}



//删除ds
func DeleteDs(c *gin.Context)  {
	clusterName := c.Query("cluster_name")
	name := c.Query("name")
	namespace := c.Query("namespace")
	client, _ := InitClientSet(clusterName)

	err := client.AppsV1().DaemonSets(namespace).Delete(context.TODO(),name,metav1.DeleteOptions{})
	if err != nil {
		fmt.Sprintln(err)
		c.JSON(http.StatusBadRequest,gin.H{
			"code": 400,
			"msg": fmt.Sprintln(err),
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"msg": "delete ds success!",
	})

}
