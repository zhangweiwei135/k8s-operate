
package common

import (
_ "database/sql"
"fmt"
"gorm.io/driver/mysql"
"gorm.io/gorm"
"gorm.io/gorm/schema"
	"k8sOps/model"

	"time"
)

func InitDb() *gorm.DB {


	constr := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		model.DbInfo.User, model.DbInfo.Password, model.DbInfo.Host, model.DbInfo.Port, model.DbInfo.Db)
	//database, err := sqlx.Open("mysql",constr)
	db, err := gorm.Open(mysql.Open(constr),&gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		fmt.Errorf("数据库连接失败： %v\n", err)
		return nil
	}

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(20)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)

	return db
}

