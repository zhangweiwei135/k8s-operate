package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	v13 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8sOps/model"
	"log"
	"net/http"
	"strconv"
)

//获取deploy列表信息
func Deploys(clusterName,namespace string) []model.DeploymentInfo {
	clientSet, _ := InitClientSet(clusterName)
	deploys, err := clientSet.AppsV1().Deployments(namespace).List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取deploy列表失败 %v \n",err)
		return nil
	}

	var deployInfo []model.DeploymentInfo
	var deploy model.DeploymentInfo

	for _, v := range deploys.Items {
		deploy.Name = v.Name
		deploy.Labels = v.Labels
		deploy.Selector = v.Spec.Selector.String()
		deploy.Expected = fmt.Sprintf("%v/%v",v.Status.ReadyReplicas,*v.Spec.Replicas)
		deploy.Replicas = *v.Spec.Replicas
		deploy.Namespace =  v.Namespace
		deployInfo = append(deployInfo,deploy)
	}
	return deployInfo
}


//根据集群获取deploy列表
func DeployList(c *gin.Context)  {
	clusterName := c.Request.FormValue("name")
	namespaceRes := c.Request.FormValue("namespace")
	d := Deploys(clusterName,namespaceRes)

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": d,
		"msg": "success",
	})
}

//修改deploy副本数量
func UpdateDeployReplicas(c *gin.Context)  {
	name := c.Query("name")
	namespace := c.Query("namespace")
	clusterName := c.Query("cluster_name")
	replicasNum := c.Query("replicas")

	replicas,_ := strconv.Atoi(replicasNum)
	var r int32

	r = int32(replicas)

	client, _ := InitClientSet(clusterName)
	deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(),name,metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return
	}

	deployment.Spec.Replicas = &r
	deployment, err = client.AppsV1().Deployments(namespace).Update(context.TODO(),deployment,metav1.UpdateOptions{})

	if err != nil {
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"msg": "update success",
	})
}

//通过yaml文件创建deployment
func CreateDeployFromYaml(clusterName,ns string,b []byte) error {
	client, _ := InitClientSet(clusterName)
	deploy := &v1.Deployment{}
	deployJson, err := yaml.ToJSON(b)
	if err != nil {
		return err
	}

	err = json.Unmarshal(deployJson,deploy)
	if err != nil {
		return err
	}

	if ns == "" {
		ns = "default"
	}

	if _,err = client.AppsV1().Deployments(ns).Create(context.TODO(),deploy,metav1.CreateOptions{}); err !=nil {
		return err
	}

	return nil
}

//通过yaml文件创建sts
func CreateStsFromYaml(clusterName,ns string,b []byte) error  {
	client, _ := InitClientSet(clusterName)
	sts := &v1.StatefulSet{}
	stsJson, err := yaml.ToJSON(b)
	if err != nil {
		return err
	}

	err = json.Unmarshal(stsJson,sts)
	if err != nil {
		return err
	}

	if ns == "" {
		ns = "default"
	}

	if _,err = client.AppsV1().StatefulSets(ns).Create(context.TODO(),sts,metav1.CreateOptions{}); err !=nil {
		return err
	}

	return nil
}

//通过yaml文件创建ds
func CreateDsFromYaml(clusterName,ns string,b []byte) error  {
	client, _ := InitClientSet(clusterName)
	ds := &v1.DaemonSet{}
	dsJson, err := yaml.ToJSON(b)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dsJson,ds)
	if err != nil {
		return err
	}

	if ns == "" {
		ns = "default"
	}

	if _,err = client.AppsV1().DaemonSets(ns).Create(context.TODO(),ds,metav1.CreateOptions{}); err !=nil {
		return err
	}

	return nil
}

//通过yaml文件创建svc
func CreateSvcFromYaml(clusterName,ns string,b []byte) error  {
	client, _ := InitClientSet(clusterName)
	svc := &v12.Service{}
	svcJson, err := yaml.ToJSON(b)
	if err != nil {
		return err
	}

	err = json.Unmarshal(svcJson,svc)
	if err != nil {
		return err
	}

	if ns == "" {
		ns = "default"
	}

	if _,err = client.CoreV1().Services(ns).Create(context.TODO(),svc,metav1.CreateOptions{}); err !=nil {
		return err
	}

	return nil
}


//通过yaml文件创建ingress
func CreateIngressFromYaml(clusterName,ns string,b []byte) error  {
	client, _ := InitClientSet(clusterName)
	r := &v13.Ingress{}
	rJson, err := yaml.ToJSON(b)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rJson,r)
	if err != nil {
		return err
	}

	if ns == "" {
		ns = "default"
	}

	if _,err = client.NetworkingV1().Ingresses(ns).Create(context.TODO(),r,metav1.CreateOptions{}); err !=nil {
		return err
	}

	return nil
}

//通过yaml文件创建cm
func CreateCMFromYaml(clusterName,ns string,b []byte) error  {
	client, _ := InitClientSet(clusterName)
	r := &v12.ConfigMap{}
	rJson, err := yaml.ToJSON(b)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rJson,r)
	if err != nil {
		return err
	}

	if ns == "" {
		ns = "default"
	}

	if _,err = client.CoreV1().ConfigMaps(ns).Create(context.TODO(),r,metav1.CreateOptions{}); err !=nil {
		return err
	}

	return nil
}

//通过yaml文件创建secret
func CreateSecretFromYaml(clusterName,ns string,b []byte) error  {
	client, _ := InitClientSet(clusterName)
	r := &v12.Secret{}
	rJson, err := yaml.ToJSON(b)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rJson,r)
	if err != nil {
		return err
	}

	if ns == "" {
		ns = "default"
	}

	if _,err = client.CoreV1().Secrets(ns).Create(context.TODO(),r,metav1.CreateOptions{}); err !=nil {
		return err
	}

	return nil
}



//获取前端yaml文件
func ResYaml(c *gin.Context)  {

	code := c.Query("code")
	clusterName := c.Query("cluster_name")
	ns := c.Query("namespace")
	resourceKind := c.Query("resource_kind")

	b := []byte(code)

	if ns == "" {
		ns = "default"
	}

	var err error

	switch resourceKind {
	case "deployment":
		err = CreateDeployFromYaml(clusterName,ns,b)
	case "statefulSet":
		err = CreateStsFromYaml(clusterName,ns,b)
	case "daemonSet":
		err = CreateDsFromYaml(clusterName,ns,b)
	case "service":
		err = CreateSvcFromYaml(clusterName,ns,b)
	case "ingress":
		err = CreateIngressFromYaml(clusterName,ns,b)
	case "configMap":
		err = CreateCMFromYaml(clusterName,ns,b)
	case "secret":
		err = CreateSecretFromYaml(clusterName,ns,b)
	default:
		err = errors.New("请输入正确的资源类型")
	}


	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"code": 200,
			"msg": resourceKind + "deploy fail" + fmt.Sprintln(err),
		})
	} else {
		c.JSON(http.StatusOK,gin.H{
			"code": 200,
			"msg": resourceKind + " deploy success",
		})
	}


}

//删除deploy
func DeleteDeployment(c *gin.Context)  {
	clusterName := c.Query("cluster_name")
	name := c.Query("name")
	namespace := c.Query("namespace")
	client, _ := InitClientSet(clusterName)

	err := client.AppsV1().Deployments(namespace).Delete(context.TODO(),name,metav1.DeleteOptions{})
	if err != nil {
		fmt.Sprintln(err)
		c.JSON(http.StatusBadRequest,gin.H{
			"code": 400,
			"msg": fmt.Sprintln(err),
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"code": 200,
		"msg": "delete deployment success!",
	})

}




