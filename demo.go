package main

import (
	"fmt"
	"context"
	//"encoding/json"
	//"go.mongodb.org/mongo-driver/bson"
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
	res, err := coll.InsertOne(context.TODO(), map[string]interface{}{"name":"tang","age":21})
	if err != nil {
		panic(err)
	}

	fmt.Println("id:",res.InsertedID)
*/
/*
	//插入多条
	res, err := coll.InsertMany(context.TODO(), []interface{}{map[string]interface{}{"name":"liu111","age":20},map[string]interface{}{"name":"ma222","age":25}})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.InsertedIDs)
*/

/*
	//查询单条，如果是find就单条，select就多条
	//filter := bson.M{"age":bson.M{"$gte":33},"name":bson.M{"$eq":"zhang"}}
	filter := map[string]interface{}{"age":map[string]interface{}{"$gte":3},"name":map[string]interface{}{"$eq":"tang"}}

	var result map[string]interface{}
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("没有结果")
			return
		} else {
			panic(err)
		}
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
	findOptions.SetSort(map[string]interface{}{"age": 1})
	findOptions.SetProjection(map[string]interface{}{"name":1,"age":1})

	filter := map[string]interface{}{"age":map[string]interface{}{"$gte":1}}
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
	id, _ := primitive.ObjectIDFromHex("6585296a7670cd5b687cbcea")
	filter := map[string]interface{}{"_id":id}
	update := map[string]interface{}{"$set":map[string]interface{}{"addr":"shanghai"}}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.ModifiedCount)
*/

/*
	//修改多条，如果条件不含_id
	update := map[string]interface{}{"$set":map[string]interface{}{"email":"a@b.c"}}
	filter := map[string]interface{}{"age":map[string]interface{}{"$gte":1}}
	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.ModifiedCount)
*/

/*
	//删除单条，如果条件含_id或limit=1
	id, _ := primitive.ObjectIDFromHex("6585296a7670cd5b687cbce9")
	filter := map[string]interface{}{"_id":id}
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.DeletedCount)
*/

/*
	//删除多条，如果条件不含_id或limit大于1
	filter := map[string]interface{}{"age":map[string]interface{}{"$gt":30}}
	result, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.DeletedCount)
*/

/*
	//统计条数
	filter := map[string]interface{}{"name":map[string]interface{}{"$in":[]string{"mahaha","tang"}}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Println(count)
*/
/*
	//sum查询
    filter := map[string]interface{}{"age":map[string]interface{}{"$gt":10}}
    cursor, _ := coll.Aggregate(context.Background(), []map[string]interface{}{
        {"$match": filter},
        {"$group": map[string]interface{}{"_id": nil, "total": map[string]interface{}{"$sum": "$age"}}},
    })

    if cursor.Next(context.TODO()) {
    	var result map[string]interface{}
		cursor.Decode(&result)
		fmt.Println(result["total"])
    } else {
    	fmt.Println(0)
    }
*/

//*
    //常见聚合count\sum\max\min\avg
    filter := map[string]interface{}{"age":map[string]interface{}{"$gt":1}}
    cursor, _ := coll.Aggregate(context.Background(), []map[string]interface{}{
        {"$match": filter},//where
        {"$limit":10},
        {"$skip":0},
        {"$group": map[string]interface{}{
        	"_id":"$addr",
        	"num": map[string]interface{}{"$sum": 1},
        	"total": map[string]interface{}{"$sum": "$age"},
        	"avg":map[string]interface{}{"$avg": "$age"},
        	"min":map[string]interface{}{"$min": "$age"},
        	"max":map[string]interface{}{"$max": "$age"},
        }},
        {"$sort": map[string]interface{}{"total" : -1 }},//聚合结果排序
        {"$match": map[string]interface{}{"total":map[string]interface{}{"$gt":100}}},//having
    })

    var result map[string]interface{}
    for cursor.Next(context.TODO()) {
		cursor.Decode(&result)
		fmt.Println(result)
    }
//*//
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
