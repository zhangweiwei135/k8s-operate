package api

import (
	"fmt"
	"github.com/minio/minio-go"
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8sOps/model"
	"log"
)

//func KubeConfig(name string)  *string  {
//	c := K8sSel(name)
//	var temp string
//	for _,v := range c {
//		temp = path.Join("./" + v.KubeConfigPath)
//	}
//
//
//	return &temp
//}

func InitClientSet(clusterName string) (*kubernetes.Clientset,*rest.Config)  {
	c := K8sSel(clusterName)

	var objectName string
	for _, v := range c {
		objectName = v.KubeConfigPath
	}

	object := DownloadFile(objectName)

	defer object.Close()
	b, err := io.ReadAll(object)
	if err != nil {
		log.Println(err)
		return nil,nil
	}
	clientConfig, err := clientcmd.NewClientConfigFromBytes(b)
	if err != nil {
		fmt.Printf("初始化kubeconfig配置文件失败：v%\n",err)
		return nil,nil
	}
	config, err := clientConfig.ClientConfig()
	clientSet, _ := kubernetes.NewForConfig(config)

	//pods, err := clientSet.CoreV1().Pods("").List(context.TODO(),metav1.ListOptions{})
	//
	//for _,v := range pods.Items {
	//	fmt.Println(v.Name)
	//}

	return clientSet,config
}

func InitMinioClient()  *minio.Client {

	// 初使化 minio client对象。
	minioClient, err := minio.New(model.MinioInfo.Endpoint,model.MinioInfo.Username,model.MinioInfo.Password,model.MinioInfo.UseSSl)
	if err != nil {
		log.Println(err)
		return nil
	}

	return minioClient
}
