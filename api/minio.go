package api

import (
	"fmt"
	"github.com/minio/minio-go"
	"io"
	"k8sOps/model"
	"log"
)

//上传文件
func UploadFile(bucketName,objectName string,reader io.Reader,objectSize int64) bool {
	minioClient := InitMinioClient()
	//判断存储桶是否存在
	b, _ := minioClient.BucketExists(bucketName)
	if !b {
		minioClient.MakeBucket(bucketName,"")
	}

	n, err := minioClient.PutObject(bucketName,objectName,reader,objectSize,minio.PutObjectOptions{})
	if err != nil {
		log.Println(err)
		return false
	}

	fmt.Println("successfully uploaded bytes: ",n)
	return true
}


//下载文件
func DownloadFile(objectName string) *minio.Object {
	minioClient := InitMinioClient()
	object, err := minioClient.GetObject(model.KubeConfigBucket,objectName,minio.GetObjectOptions{})
	if err != nil {
		log.Println(err)
		return nil
	}

	//fmt.Printf("%v",b)
	return object

}
