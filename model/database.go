package model

import v1 "k8s.io/api/core/v1"

type ClusterInfo struct {
	ID int `json:"id" gorm:"primaryKey,autoIncrement,index"`
	KubeConfigPath string `json:"kube_config_path"`
	Token string `json:"token"` //token信息
	Name string `json:"name" gorm:"type:not null"`
	CreateTime string `json:"create_time"`
}


//数据库连接信息
type DBInfo struct {
	Host string
	User string
	Password string
	Db string
	Port int
}




//用于临时存储集群信息的Map
var ClusterState = make(map[string]interface{})

//返回给前端的deployment信息
type DeploymentInfo struct {
	Name string `json:"name"`
	Labels map[string]string `json:"labels"`
	Selector string `json:"selector"`
	Expected string `json:"expected"`
	Replicas int32 `json:"replicas"`
	Namespace string `json:"namespace"`
}


//返回给前端的sts信息
type StatefulSetInfo struct {
	Name string `json:"name"`
	Labels map[string]string `json:"labels"`
	Selector string `json:"selector"`
	Expected string `json:"expected"`
	Replicas int32 `json:"replicas"`
	Namespace string `json:"namespace"`
}

//返回给前端的ds信息
type DaemonSetInfo struct {
	Name string `json:"name"`
	Labels map[string]string `json:"labels"`
	Selector string `json:"selector"`
	HostName string `json:"host_name"`
	Namespace string `json:"namespace"`
}

//返回给前端的svc信息
type ServiceInfo struct {
	Name string `json:"name"`
	Labels map[string]string `json:"labels"`
	Selector map[string]string `json:"selector"`
	IP []string `json:"ip"`
	Namespace string `json:"namespace"`
	SvcType string `json:"svc_type"`
}


//返回给前端的ing信息
type IngressInfo struct {
	Name string `json:"name"`
	Labels map[string]string `json:"labels"`
	Vip []string `json:"vip"`
	BackendSvc string `json:"backend_svc"`
	UrlName string `json:"url_name"`
	Namespace string `json:"namespace"`
	IngType *string `json:"ing_type"`
}

//返回给前端的cm信息
type ConfigMapInfo struct {
	Name string `json:"name"`
	Labels map[string]string `json:"labels"`
	Namespace string `json:"namespace"`
}

//返回给前端的secret信息
type SecretInfo struct {
	Name string `json:"name"`
	Labels map[string]string `json:"labels"`
	Namespace string `json:"namespace"`
	SecretType string `json:"secret_type"`
}

//返回给前端的pod信息
type PodsInfo struct {
	Name string `json:"name"`
	Status string `json:"status"`
	HostIP string `json:"host_ip"`
	PodContainers []v1.Container `json:"pod_containers"`
	PodIP string `json:"pod_ip"`
	Namespace string `json:"namespace"`
}


//minio连接信息
type MinioClient struct {
	Endpoint string `json:"endpoint"`
	Username string `json:"username"`
	Password string `json:"password"`
	UseSSl bool
}

