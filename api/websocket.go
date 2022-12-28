package api

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
)

// http升级websocket协议的配置
var WsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func PodConnect(c *gin.Context)  {
	//获取容器相关参数
	ns := c.Query("namespace")
	podName := c.Query("name")
	containerName := c.Query("cname")
	clusterName := c.Query("cluster_name")


	wsClient, err := WsUpgrader.Upgrade(c.Writer,c.Request,nil)
	shellClient := NewWsShellClient(wsClient)
	if err != nil {
		log.Println(err)
		return
	}


	defer wsClient.Close()

	//var restConf *rest.Config
	//restConf,err = clientcmd.BuildConfigFromFlags("",*KubeConfig(clusterName));
	clientset, restConf := InitClientSet(clusterName)


	sshReq := clientset.CoreV1().RESTClient().Post().
		Resource("pods").Name(podName).Namespace(ns).SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command: []string{"/bin/bash"},
			Stdin: true,
			Stdout: true,
			Stderr: true,
			TTY: true,
	},scheme.ParameterCodec)

	fmt.Println(sshReq.URL().Path)

	executor, err := remotecommand.NewSPDYExecutor(restConf,"POST",sshReq.URL())
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("connectingi to pod...")

	err = executor.StreamWithContext(context.TODO(),remotecommand.StreamOptions{
		Stdin: shellClient,
		Stdout: shellClient,
		Stderr: shellClient,
		Tty: true,
	})

	if err != nil {
		log.Println(err)
		return
	}
}



type WsShellClient struct {
	client *websocket.Conn
}

func NewWsShellClient(client *websocket.Conn) *WsShellClient {
	return &WsShellClient{client: client}
}

func (this *WsShellClient) Write(p []byte) (n int, err error) {
	err = this.client.WriteMessage(websocket.TextMessage,
		p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
func (this *WsShellClient) Read(p []byte) (n int, err error) {

	_, b, err := this.client.ReadMessage()

	if err != nil {
		return 0, err
	}
	return copy(p, string(b)), nil
}


//获取pod日志
func GetPodLogs(c *gin.Context)  {
	cluster_name := c.Query("cluster_name")
	podName := c.Query("name")
	ns := c.Query("namespace")
	cname := c.Query("cname")



	client,_ := InitClientSet(cluster_name)

	wsClient, err := WsUpgrader.Upgrade(c.Writer,c.Request,nil)
	//shellClient := NewWsShellClient(wsClient)
	if err != nil {
		log.Println(err)
		return
	}


	defer wsClient.Close()


	req := client.CoreV1().Pods(ns).GetLogs(podName,&v1.PodLogOptions{
		Container: cname,
		//follow为true代表的是流式获取，否则只返回单次日志
		Follow: true,
	})

	stream ,err := req.Stream(context.TODO())
	defer stream.Close()

	reader := bufio.NewReader(stream)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		wsClient.WriteMessage(websocket.TextMessage,[]byte(line))
	}

	//fmt.Println(req.URL().Path)
	//
	//executor, err := remotecommand.NewSPDYExecutor(restConf,"GET",req.URL())
	//if err != nil {
	//
	//	log.Println(err)
	//	return
	//}
	//
	//log.Println("connectingi to pod...")
	//
	//err = executor.StreamWithContext(context.TODO(),remotecommand.StreamOptions{
	//	Stdin: shellClient,
	//	Stderr: shellClient,
	//	Stdout: os.Stdout,
	//	Tty: true,
	//})
	//
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	////gin会给每个请求都起一个协程，不设超时时间就会阻塞在read，导致每刷新一次多一个协程
	//cc, err := context.WithTimeout(context.TODO(),time.Minute*5)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	////通过request的流获取到ioReader
	//reader, _ := req.Stream(cc)
	//
	//defer reader.Close()
	//for {
	//	b := make([]byte,1024)
	//	n, err := reader.Read(b)
	//	if err != nil && err == io.EOF {
	//		break
	//	}
	//
	//	if n >0 {
	//		fmt.Println(b[0:n])
	//		c.Writer.Write(b[0:n])
	//		c.Writer.(http.Flusher).Flush()
	//	}
	//}
	//return
}
