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
	//查询单条
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
	//修改单条
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
	//修改多条
	update := bson.D{{"$set", bson.M{"email":"a@b.c"}}}
	filter := bson.M{"age":bson.M{"$gte":1}}
	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.ModifiedCount)
*/

/*
	//删除单条
	id, _ := primitive.ObjectIDFromHex("6582b1952290aece3b63803b")
	filter := bson.D{{"_id", id}}
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.DeletedCount)
*/

/*
	//删除多条
	filter := bson.M{"age":bson.M{"$gt":30}}
	result, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.DeletedCount)
*/


	//统计条数
	filter := bson.M{"age":bson.M{"$gt":1}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Println(count)
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
