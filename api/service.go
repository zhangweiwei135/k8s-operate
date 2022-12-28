package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sOps/model"
	"net/http"
)

//获取svc列表信息
func Services(clusterName,namespace string) *[]model.ServiceInfo {
	clientSet, _ := InitClientSet(clusterName)
	svcs, err := clientSet.CoreV1().Services(namespace).List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取svc列表失败 %v \n",err)
		return nil
	}

	var svcInfo []model.ServiceInfo
	var svc model.ServiceInfo

	for _, v := range svcs.Items {
		svc.Name = v.Name
		svc.Labels = v.Labels
		svc.Selector = v.Spec.Selector
		svc.SvcType = string(v.Spec.Type)
		svc.Namespace =  v.Namespace
		var ips []string
		ips = append(ips,v.Spec.ClusterIP)
		if len(v.Spec.ExternalIPs) > 0 {
			ips = append(ips,v.Spec.ExternalIPs[0])
		}
		svc.IP = ips
		svcInfo = append(svcInfo,svc)
	}
	return &svcInfo
}


//根据集群获取svc列表
func SvcList(c *gin.Context)  {
	clusterName := c.Request.FormValue("cluster_name")
	namespaceRes := c.Request.FormValue("namespace")
	svc := Services(clusterName,namespaceRes)

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": svc,
		"msg": "success",
	})
}



//删除svc
func DeleteSvc(c *gin.Context)  {
	clusterName := c.Query("cluster_name")
	name := c.Query("name")
	namespace := c.Query("namespace")
	client, _ := InitClientSet(clusterName)

	err := client.CoreV1().Services(namespace).Delete(context.TODO(),name,metav1.DeleteOptions{})
	if err != nil {
		fmt.Sprintln(err)
		c.JSON(http.StatusBadRequest,gin.H{
			"code": 400,
			"msg": fmt.Sprintln(err),
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"msg": "delete svc success!",
	})

}




//获取ing列表信息
func IngressList(clusterName,namespace string) *[]model.IngressInfo {
	clientSet, _ := InitClientSet(clusterName)
	ings, err := clientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(),metav1.ListOptions{})

	if err != nil {
		fmt.Printf("获取ingress列表失败 %v \n",err)
		return nil
	}

	var rInfo []model.IngressInfo
	var r model.IngressInfo

	for _, v := range ings.Items {
		r.Name = v.Name
		r.Labels = v.Labels
		var ip []string
		if len(v.Status.LoadBalancer.Ingress) >0 {
			for _,x := range v.Status.LoadBalancer.Ingress {
				ip = append(ip,x.IP)
			}
		}
		r.Vip = ip
		r.IngType = v.Spec.IngressClassName
		r.Namespace =  v.Namespace
		r.BackendSvc = v.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name
		r.UrlName =  v.Spec.Rules[0].Host + v.Spec.Rules[0].HTTP.Paths[0].Path
		rInfo = append(rInfo,r)
	}
	return &rInfo
}


//根据集群获取ing列表
func IngList(c *gin.Context)  {
	clusterName := c.Request.FormValue("cluster_name")
	namespaceRes := c.Request.FormValue("namespace")
	r := IngressList(clusterName,namespaceRes)

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": r,
		"msg": "success",
	})
}



//删除svc
func DeleteIng(c *gin.Context)  {
	clusterName := c.Query("cluster_name")
	name := c.Query("name")
	namespace := c.Query("namespace")
	client, _ := InitClientSet(clusterName)

	err := client.NetworkingV1().Ingresses(namespace).Delete(context.TODO(),name,metav1.DeleteOptions{})
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
