# OnSQL
用SQL查询MongoDB

## 快速上手
```
package main

import "fmt"
import "github.com/boyxp/OnSQL/MongoDB"

func main() {
	//注册MongoDB数据库信息(tag标签、数据库名称、dsn)
	MongoDB.Register("demo", "test", "mongodb://localhost:27017")

	//指定tag标签和集合名，用tag方便数据库改名
	shop := MongoDB.Model{"demo.shop"}

	//插入记录
	shop.Insert(map[string]interface{}{"name":"可口可乐","price":100,"detail":"...","category":"饮料"})
	shop.Insert(map[string]interface{}{"name":"小红帽","price":200,"detail":"...","category":"服装"})
	shop.Insert(map[string]interface{}{"name":"雪碧","price":300,"category":"饮料"})

	//读取记录
	list := shop.Field("name,price").
				Where("price", ">", 100).
				Order("price", "desc").
				Page(1).Limit(10).Select()
	fmt.Println(list)

	//聚合查询
	aggs := shop.Field("category,count(*) total,sum(price) as sum_price,avg(price) as avg_price,min(price) as min_price,max(price) as max_price").
			Group("category").Select()
	for k, v := range aggs {
		fmt.Println(k, v)
	}
}
```

更多示例见[单元测试](https://github.com/boyxp/OnSQL/blob/main/MongoDB/orm_test.go)
