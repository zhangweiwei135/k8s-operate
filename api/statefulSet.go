package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sOps/model"
	"log"
	"net/http"
	"strconv"
)

//获取sts列表信息
func StatefulSets(clusterName,namespace string) *[]model.StatefulSetInfo {
	clientSet, _ := InitClientSet(clusterName)
	stss, err := clientSet.AppsV1().StatefulSets(namespace).List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取sts列表失败 %v \n",err)
		return nil
	}

	var stsInfo []model.StatefulSetInfo
	var sts model.StatefulSetInfo

	for _, v := range stss.Items {
		sts.Name = v.Name
		sts.Labels = v.Labels
		sts.Selector = v.Spec.Selector.String()
		sts.Expected = fmt.Sprintf("%v/%v",v.Status.ReadyReplicas,*v.Spec.Replicas)
		sts.Replicas = *v.Spec.Replicas
		sts.Namespace =  v.Namespace
		stsInfo = append(stsInfo,sts)
	}
	return &stsInfo
}


//根据集群获取sts列表
func StsList(c *gin.Context)  {
	clusterName := c.Request.FormValue("cluster_name")
	namespaceRes := c.Request.FormValue("namespace")
	s := StatefulSets(clusterName,namespaceRes)

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": s,
		"msg": "success",
	})
}

//修改sts副本数
func UpdateStsReplicas(c *gin.Context)  {
	name := c.Query("name")
	namespace := c.Query("namespace")
	clusterName := c.Query("cluster_name")
	replicasNum := c.Query("replicas")

	replicas,_ := strconv.Atoi(replicasNum)
	var r int32

	r = int32(replicas)

	client, _ := InitClientSet(clusterName)
	sts, err := client.AppsV1().StatefulSets(namespace).Get(context.TODO(),name,metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return
	}

	sts.Spec.Replicas = &r
	sts, err = client.AppsV1().StatefulSets(namespace).Update(context.TODO(),sts,metav1.UpdateOptions{})

	if err != nil {
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"msg": "update success",
	})
}

//删除sts
func DeleteSts(c *gin.Context)  {
	clusterName := c.Query("cluster_name")
	name := c.Query("name")
	namespace := c.Query("namespace")
	client, _ := InitClientSet(clusterName)

	err := client.AppsV1().StatefulSets(namespace).Delete(context.TODO(),name,metav1.DeleteOptions{})
	if err != nil {
		fmt.Sprintln(err)
		c.JSON(http.StatusBadRequest,gin.H{
			"code": 400,
			"msg": fmt.Sprintln(err),
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"msg": "delete sts success!",
	})

}
