package main

import (
	"fmt"
	"k8sOps/route"
)


func main() {
	route.InitRouter()

	//api.Deploys("cluster1","default")

	//api.InitClientSet("cluster1")

	//api.DownloadFile("kubeconfig","kubeConfig.yaml")

	//api.CreateDeployFromYaml("cluster3","")

	//name := "cluster3"
	//api.Deploys(name,"")

	//clientSet := api.InitClientSet(name)
	//pods, err := clientSet.CoreV1().Pods("default").List(context.TODO(),metav1.ListOptions{})
	//if err != nil {
	//	fmt.Printf("获取失败 %v \n",err)
	//	return
	//}
	//
	//for _, v := range pods.Items {
	//	fmt.Printf("name: %v\t status: %v\n",v.Name,v.Status.Phase)
	//}

	fmt.Println()
}
