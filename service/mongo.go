package service

import (
	"chat/conf"
	"chat/model/ws"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

func InsertMsg(database, id string, content string, read uint, expire int64) error {
	//插入到mongoDB中
	collection := conf.MongoDBClient.Database(database).Collection(id)
	comment := ws.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}

	_, err := collection.InsertOne(context.TODO(), comment)
	return err
}

func FindMany(database string, sendId string, uid string, funcTime int64, pageSize int) (results []ws.Result, err error) {
	var resultMe []ws.Trainer
	var resultSend []ws.Trainer

	sendIdCollection := conf.MongoDBClient.Database(database).Collection(sendId)
	uidCollection := conf.MongoDBClient.Database(database).Collection(uid)

	//context.Background(), bson.D{}
	sendIdTimeCurcor, err := sendIdCollection.Find(context.TODO(), bson.D{}, options.Find().SetSort(bson.D{{"startTime", -1}}), options.Find().SetLimit(int64(pageSize)))

	uidTimeCurcor, err := uidCollection.Find(context.TODO(), bson.D{}, options.Find().SetSort(bson.D{{"startTime", -1}}), options.Find().SetLimit(int64(pageSize)))
	err = sendIdTimeCurcor.All(context.TODO(), &resultSend)
	err = uidTimeCurcor.All(context.TODO(), &resultMe)

	results, _ = AppendAndSort(resultMe, resultSend)

	return

}

func AppendAndSort(resultMe, resultSend []ws.Trainer) (results []ws.Result, err error) {
	for _, r := range resultMe {
		sendSort := SendSortMsg{ //构造返回的msg
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{ //构造返回所有的内容
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "me",
		}
		results = append(results, result)
	}

	for _, r := range resultSend {
		sendSort := SendSortMsg{ //构造返回的msg
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{ //构造返回所有的内容
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "you",
		}

		results = append(results, result)
	}

	return
}
