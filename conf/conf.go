package conf

import (
	"chat/model"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string

	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string

	MongoDBName     string
	MongoDBAddr     string
	MongoDBPwd      string
	MongoDBPort     string
	MongoDBUserName string
	MongoDBClient   *mongo.Client
)

func Init() {
	//从本地读取环境
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println("初始化配置文件失败", err)
	}

	LoadServer(file)
	LoadMySQL(file)
	LoadMongoDB(file)

	MongoDB()

	path := DbUser + ":" + DbPassword + "@tcp(" + DbHost + ":" + DbPort + ")/" + DbName + "?charset=utf8mb4&parseTime=true"
	model.Database(path)

}

func MongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://" + MongoDBUserName + ":" + MongoDBPwd + "@" + MongoDBAddr + ":" + MongoDBPort)

	var err error
	MongoDBClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Info(err)
		panic(err)
	}

	logrus.Info("MongoDB connect Successfully")

}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMySQL(file *ini.File) {
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassword = file.Section("mysql").Key("DbPassword").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadMongoDB(file *ini.File) {
	MongoDBName = file.Section("MongoDB").Key("MongoDBName").String()
	MongoDBAddr = file.Section("MongoDB").Key("MongoDBAddr").String()
	MongoDBPwd = file.Section("MongoDB").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("MongoDB").Key("MongoDBPort").String()
	MongoDBUserName = file.Section("MongoDB").Key("MongoDBUserName").String()
}
