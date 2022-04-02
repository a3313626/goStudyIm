package cache

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var (
	RedisClient *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

func init() {
	//从本地读取环境
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println("初始化redis配置文件失败", err)
	}

	//读取配置信息文件内容
	LoadRedis(file)
	//redis链接
	Redis()

}

func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}

//链接redis
func Redis() {
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		DB:   int(db),
	})

	_, err := client.Ping().Result()
	if err != nil {
		logrus.Info(err)
		panic(err)
	}

	RedisClient = client

}
