# OnSQL

Using SQL Queries with MongoDB

![](https://img.shields.io/npm/l/vue.svg)
![Test](https://github.com/boyxp/OnSQL/actions/workflows/go.yml/badge.svg)


OnSQL allows you to use SQL-like queries to interact with MongoDB. Below is a guide on how to get started with OnSQL.

## Quick Start

1. **Install OnSQL**:
   To use OnSQL, you need to get the package:

   ```sh
   go get github.com/boyxp/OnSQL
   ```

2. **Create a Basic Application**:
   Here’s an example of how to create a basic application using OnSQL:

   ```go
   package main

   import (
       "fmt"
       "github.com/boyxp/OnSQL/MongoDB"
   )

   func main() {
       // Register MongoDB database information (tag, database name, DSN)
       MongoDB.Register("demo", "test", "mongodb://localhost:27017")

       // Specify the tag and collection name
       shop := MongoDB.Model{"demo.shop"}

       // Insert records
       shop.Insert(map[string]interface{}{
           "name":     "Coca Cola",
           "price":    100,
           "detail":   "...",
           "category": "Beverage",
       })
       shop.Insert(map[string]interface{}{
           "name":     "Red Riding Hood",
           "price":    200,
           "detail":   "...",
           "category": "Clothing",
       })
       shop.Insert(map[string]interface{}{
           "name":     "Sprite",
           "price":    300,
           "category": "Beverage",
       })

       // Read records
       list := shop.Field("name,price").
           Where("price", ">", 100).
           Order("price", "desc").
           Page(1).Limit(10).Select()
       fmt.Println(list)

       // Aggregate query
       aggs := shop.Field("category,count(*) total,sum(price) as sum_price,avg(price) as avg_price,min(price) as min_price,max(price) as max_price").
           Group("category").Select()
       for k, v := range aggs {
           fmt.Println(k, v)
       }
   }
   ```

### Explanation:

- **Register MongoDB**:
  ```go
  MongoDB.Register("demo", "test", "mongodb://localhost:27017")
  ```
  This registers the MongoDB database with a tag (`demo`), database name (`test`), and the connection string (DSN).

- **Model Definition**:
  ```go
  shop := MongoDB.Model{"demo.shop"}
  ```
  This defines a model for the `shop` collection in the `demo` database.

- **Insert Records**:
  ```go
  shop.Insert(map[string]interface{}{"name":"Coca Cola","price":100,"detail":"...","category":"Beverage"})
  shop.Insert(map[string]interface{}{"name":"Red Riding Hood","price":200,"detail":"...","category":"Clothing"})
  shop.Insert(map[string]interface{}{"name":"Sprite","price":300,"category":"Beverage"})
  ```
  These lines insert several records into the `shop` collection.

- **Read Records**:
  ```go
  list := shop.Field("name,price").
      Where("price", ">", 100).
      Order("price", "desc").
      Page(1).Limit(10).Select()
  fmt.Println(list)
  ```
  This queries the `shop` collection, selecting records where the price is greater than 100, ordering them by price in descending order, and limiting the results to 10 per page.

- **Aggregate Query**:
  ```go
  aggs := shop.Field("category,count(*) total,sum(price) as sum_price,avg(price) as avg_price,min(price) as min_price,max(price) as max_price").
      Group("category").Select()
  for k, v := range aggs {
      fmt.Println(k, v)
  }
  ```
  This performs an aggregate query to group records by category and calculate various statistics (total count, sum, average, min, and max price).

## Additional Resources

For more examples and detailed usage, refer to the [OnSQL Unit Tests](https://github.com/boyxp/OnSQL/blob/main/MongoDB/orm_test.go).

This should help you get started with OnSQL for MongoDB. If you have any specific questions or need further assistance, feel free to ask!

