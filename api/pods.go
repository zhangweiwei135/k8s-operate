package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sOps/model"
	"net/http"
)

//通过工作负载获取关联的pod
func PodsGetWithWorkload(clusterName,workloadName,namespace string) []model.PodsInfo {
	clientSet,_ := InitClientSet(clusterName)
	lebal := "app="+workloadName
	pods, err := clientSet.CoreV1().Pods(namespace).List(context.TODO(),metav1.ListOptions{
		LabelSelector: lebal,
	})

	if err != nil {
		fmt.Printf("获取pods失败: %v!\n",err)
		return nil
	}

	var podList []model.PodsInfo
	var p model.PodsInfo

	for _, v := range pods.Items {
		p.Name =  v.Name
		p.HostIP = v.Status.HostIP
		p.PodIP =  v.Status.PodIP
		p.Status = string(v.Status.Phase)
		p.PodContainers = v.Spec.Containers
		p.Namespace = v.Namespace
		fmt.Println(v.Name,v.Status.Phase,v.Status.HostIP,v.Status.PodIP,v.Spec.Containers)
		podList = append(podList,p)
	}
	return podList
}

//返回pod组件
func PodsList(c *gin.Context)  {
	clusterName := c.Request.FormValue("name")
	namespaceRes := c.Request.FormValue("namespace")
	workloadName := c.Request.FormValue("workloadname")
	p := PodsGetWithWorkload(clusterName,workloadName,namespaceRes)

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": p,
		"msg": "success",
	})
}




