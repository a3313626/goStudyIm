package model

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Database(connstring string) {
	db, err := gorm.Open("mysql", connstring)

	if err != nil {
		panic("Mysql 数据库连接错误")
	}
	//fmt.Println("数据库连接成功")

	db.LogMode(true)
	if gin.Mode() == "release" {
		db.LogMode(false)
	}

	db.SingularTable(true)                       //表名不加s
	db.DB().SetMaxIdleConns(20)                  //设置连接池
	db.DB().SetMaxOpenConns(100)                 //最大连接数
	db.DB().SetConnMaxLifetime(time.Second * 30) //30秒超时
	DB = db

	migration()

}
