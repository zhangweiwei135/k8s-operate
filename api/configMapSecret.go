package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sOps/model"
	"net/http"
)

//获取cm列表信息
func ConfigMaps(clusterName,namespace string) *[]model.ConfigMapInfo {
	clientSet, _ := InitClientSet(clusterName)
	rs, err := clientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(),metav1.ListOptions{})

	if err != nil {
		fmt.Printf("获取ingress列表失败 %v \n",err)
		return nil
	}

	var rInfo []model.ConfigMapInfo
	var r model.ConfigMapInfo

	for _, v := range rs.Items {
		r.Name = v.Name
		r.Labels = v.Labels
		r.Namespace =  v.Namespace
		rInfo = append(rInfo,r)
	}
	return &rInfo
}


//根据集群获取cm列表
func CmList(c *gin.Context)  {
	clusterName := c.Request.FormValue("cluster_name")
	namespaceRes := c.Request.FormValue("namespace")
	r := ConfigMaps(clusterName,namespaceRes)

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": r,
		"msg": "success",
	})
}



//删除CM
func DeleteCM(c *gin.Context)  {
	clusterName := c.Query("cluster_name")
	name := c.Query("name")
	namespace := c.Query("namespace")
	client, _ := InitClientSet(clusterName)

	err := client.CoreV1().ConfigMaps(namespace).Delete(context.TODO(),name,metav1.DeleteOptions{})
	if err != nil {
		fmt.Sprintln(err)
		c.JSON(http.StatusBadRequest,gin.H{
			"code": 400,
			"msg": fmt.Sprintln(err),
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"msg": "delete ingress success!",
	})

}


//获取secret列表信息
func Secrets(clusterName,namespace string) *[]model.SecretInfo {
	clientSet, _ := InitClientSet(clusterName)
	rs, err := clientSet.CoreV1().Secrets(namespace).List(context.TODO(),metav1.ListOptions{})

	if err != nil {
		fmt.Printf("获取列表失败 %v \n",err)
		return nil
	}

	var rInfo []model.SecretInfo
	var r model.SecretInfo

	for _, v := range rs.Items {
		r.Name = v.Name
		r.Labels = v.Labels
		r.Namespace =  v.Namespace
		r.SecretType = string(v.Type)
		rInfo = append(rInfo,r)
	}
	return &rInfo
}


//根据集群获取secret列表
func SecretList(c *gin.Context)  {
	clusterName := c.Request.FormValue("cluster_name")
	namespaceRes := c.Request.FormValue("namespace")
	r := Secrets(clusterName,namespaceRes)

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": r,
		"msg": "success",
	})
}



//删除secret
func DeleteSecret(c *gin.Context)  {
	clusterName := c.Query("cluster_name")
	name := c.Query("name")
	namespace := c.Query("namespace")
	client, _ := InitClientSet(clusterName)

	err := client.CoreV1().Secrets(namespace).Delete(context.TODO(),name,metav1.DeleteOptions{})
	if err != nil {
		fmt.Sprintln(err)
		c.JSON(http.StatusBadRequest,gin.H{
			"code": 400,
			"msg": fmt.Sprintln(err),
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"msg": "delete  success!",
	})

}
