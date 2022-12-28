package api

import (
	"fmt"
	"gorm.io/gorm/clause"
	"k8sOps/common"
	"k8sOps/model"
)

//添加集群
func K8sAdd(k *model.ClusterInfo)  {
	db := common.InitDb()

	if k.Name == "" || (k.KubeConfigPath == "" && k.Token == "") {
		fmt.Println("集群名称，认证信息不能为空")
		return
	}

	err :=db.Create(k).Error
	if err != nil {
		fmt.Printf("添加新群信息失败:v% \n",err)
		return
	}
	fmt.Println("添加集群成功！",k.Name)
}

//查询集群
func K8sSel(name string) []*model.ClusterInfo {
	db := common.InitDb()

	var k []*model.ClusterInfo
	err :=db.Where("name=?",name).Find(&k).Error
	if err != nil {
		fmt.Printf("查询集群失败：%v\n",err)
		return nil
	}

	return k
}

func K8sList() []*model.ClusterInfo  {
	var k []*model.ClusterInfo

	db := common.InitDb()
	err := db.Find(&k).Error

	if err != nil {
		fmt.Printf("集群查询失败: %v\n",err)
		return nil
	}

	fmt.Println("select secc: ",k)
	return k
}


//修改集群
func K8sUpdate(name string, k *model.ClusterInfo) bool  {
	db := common.InitDb()

	if k.Name == "" {
		fmt.Println("修改集群名不能为空")
		return false
	}

	err := db.Model(model.ClusterInfo{}).Clauses(clause.Returning{}).Where("name = ?",name).
		Updates(k).Error

	if err != nil {
		fmt.Printf("修改集群失败: %v \n",err)
		return false
	}

	fmt.Println("update succ: ",k)
	return true
}

//删除集群
func K8sDel(name string) bool  {
    db :=common.InitDb()

	err := db.Exec("delete from cluster_info where name=?",name).Error
	if err != nil {
		fmt.Printf("集群删除失败：%v\n",err)
		return false
	}
	fmt.Println("delete succ: ",name)
	return true
}



