package main

import (
	"fmt"
	"context"
	//"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	client := Open()

	defer Close(client)

	coll := client.Database("test").Collection("test")
/*
	//插入单条
	res, err := coll.InsertOne(context.TODO(), bson.M{"name":"liu","age":20})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.InsertedID)
*/
/*
	//插入多条
	res, err := coll.InsertMany(context.TODO(), []interface{}{bson.M{"name":"liu","age":20},bson.M{"name":"ma","age":25}})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.InsertedIDs)
*/

/*
	//查询单条，如果是find就单条，select就多条
	filter := bson.M{"age":bson.M{"$gte":33},"name":bson.M{"$eq":"zhang"}}

	var result map[string]interface{}
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		panic(err)
	}

	for k, v := range result {
		fmt.Println(k,"=",v)
	}
*/

/*
	//多条查询，分页、排序
	pageSize := 2
	pageNum := 2
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(pageSize * (pageNum - 1)))
	findOptions.SetSort(bson.M{"age": 1})

	filter := bson.M{"age":bson.M{"$gte":1}}
	cursor, err := coll.Find(context.TODO(), filter, findOptions)
	if err != nil {
		panic(err)
	}

	var list []map[string]interface{}
	if err = cursor.All(context.TODO(), &list); err != nil {
		panic(err)
	}

	for _, v := range list {
		cursor.Decode(&v)
		fmt.Println(v)
	}
*/

/*
	//修改单条，如果条件是_id
	id, _ := primitive.ObjectIDFromHex("6582b1952290aece3b63803b")
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.M{"addr":"shanghai"}}}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.ModifiedCount)
*/

/*
	//修改多条，如果条件不含_id
	update := bson.D{{"$set", bson.M{"email":"a@b.c"}}}
	filter := bson.M{"age":bson.M{"$gte":1}}
	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.ModifiedCount)
*/

/*
	//删除单条，如果条件含_id或limit=1
	id, _ := primitive.ObjectIDFromHex("6582b1952290aece3b63803b")
	filter := bson.D{{"_id", id}}
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.DeletedCount)
*/

/*
	//删除多条，如果条件不含_id或limit大于1
	filter := bson.M{"age":bson.M{"$gt":30}}
	result, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.DeletedCount)
*/

/*
	//统计条数
	//filter := bson.M{"age":bson.M{"$gt":1}}
	filter := bson.M{"name":bson.M{"$in":bson.A{"lee","jun"}}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Println(count)
*/
/*
	//sum查询
    filter := bson.M{"age":bson.M{"$gt":100}}
    cursor, _ := coll.Aggregate(context.Background(), []bson.M{
        {"$match": filter},
        {"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$age"}}},
    })

    if cursor.Next(context.TODO()) {
    	var result bson.M
		cursor.Decode(&result)
		fmt.Println(result["total"])
    } else {
    	fmt.Println(0)
    }
*/

    //常见聚合count\sum\max\min\avg
    filter := bson.M{"age":bson.M{"$gt":1}}
    cursor, _ := coll.Aggregate(context.Background(), []bson.M{
        {"$match": filter},//where
        {"$limit":10},
        {"$skip":0},
        {"$group": bson.M{"_id":"$addr",
        	"num": bson.M{"$sum": 1},
        	"total": bson.M{"$sum": "$age"},
        	"avg":bson.M{"$avg": "$age"},
        	"min":bson.M{"$min": "$age"},
        	"max":bson.M{"$max": "$age"},
        }},
        {"$sort": bson.M{"total" : -1 }},//聚合结果排序
        {"$match": bson.M{"total":bson.M{"$gt":100}}},//having
    })

    var result bson.M
    for cursor.Next(context.TODO()) {
		cursor.Decode(&result)
		fmt.Println(result)
    }
}

func Open() *mongo.Client {
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func Close(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
