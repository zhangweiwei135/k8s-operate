package model


var KubeConfigBucket = "kubeconfig"

var YamlBucket = "yaml"

var MinioInfo = MinioClient{
	Endpoint: "139.9.58.140:9000",
	Username: "admin",
	Password: "admin123",
	UseSSl: false,
}

//实例化数据库信息
var DbInfo  = DBInfo{
	Host: "124.221.177.224",
	Port: 31671,
	Password: "ZTEidc123",
	User: "root",
	Db: "kubernetes",
}
