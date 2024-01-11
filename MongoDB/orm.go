package MongoDB

import "os"
import "log"
import "context"
import "strings"
import "strconv"
import "go.mongodb.org/mongo-driver/mongo"
import "go.mongodb.org/mongo-driver/bson/primitive"
import "go.mongodb.org/mongo-driver/mongo/options"
import "go.mongodb.org/mongo-driver/bson"

type Orm struct {
	coll *mongo.Collection

	selectFields map[string]interface{}
	selectConds []string
	selectParams []interface{}
	selectPage int64
	selectLimit int64
	selectOrder bson.D
	selectGroupId map[string]interface{}
	selectGroup map[string]interface{}
	selectHaving map[string]interface{}

	debug string
}

func (O *Orm) Init(dbtag string, table string) *Orm {
	dbname        := Dbname(dbtag)
	O.coll         = Open(dbtag).Database(dbname).Collection(table)
	O.debug        = os.Getenv("debug")
	O.selectFields = map[string]interface{}{}
	O.selectOrder  = bson.D{}
	O.selectGroupId= map[string]interface{}{}
	O.selectGroup  = map[string]interface{}{}
	O.selectHaving = map[string]interface{}{}
	O.selectPage   = 1
	O.selectLimit  = 20

	return O
}

func (O *Orm) Insert(data map[string]interface{}) string {
	res, err := O.coll.InsertOne(context.TODO(), data)
	if err != nil {
		panic(err)
	}

	return res.InsertedID.(primitive.ObjectID).Hex()
}

func (O *Orm) Delete() int64 {
	filter := O.filter()
	result, err := O.coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	return result.DeletedCount
}

func (O *Orm) Update(data map[string]interface{}) int64 {
	filter := O.filter()
	update := map[string]interface{}{"$set":data}
	result, err := O.coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	return result.ModifiedCount
}

func (O *Orm) Field(fields string) *Orm {
	_fields := strings.Split(fields, ",")
	for k:=0;k<len(_fields);k++ {
		field := strings.TrimSpace(_fields[k])

		left_idx := strings.Index(field, "(")
		//聚合或函数字段
		if left_idx>-1 {
			right_idx := strings.Index(field, ")")
			if right_idx==-1 {
				k++
				field = field+"^"+strings.TrimSpace(_fields[k])
				right_idx = strings.Index(field, ")")
			}

			aggs      := strings.ToLower(field[0:left_idx])
			alias     := strings.TrimSpace(field[right_idx+1:])
			aggs_field:= strings.TrimSpace(field[left_idx+1:right_idx])

			if len(alias)==0 {
				alias = aggs+"_"+strconv.Itoa(k)
			} else {
				space := strings.LastIndex(alias, " ")
				if space>-1 {
					alias = alias[space+1:]
				}
			}

			switch aggs {
				case "count" : O.selectGroup[alias] = map[string]interface{}{"$sum": 1}
				case "sum"   : fallthrough
				case "avg"   : fallthrough
				case "max"   : fallthrough
				case "min"   : O.selectGroup[alias] = map[string]interface{}{"$"+aggs: "$"+aggs_field}
				case "date_format":
								split_idx := strings.Index(aggs_field, "^")
								if split_idx==-1 {
									panic("date_format必须有参数")
								}
								aggs_param := strings.Trim(aggs_field[split_idx+1:],"'")
								aggs_field  = aggs_field[0:split_idx]

								O.selectGroupId[alias] = map[string]interface{}{"$dateToString": map[string]interface{}{ "format": aggs_param, "date": "$"+aggs_field}}

				default      : panic("不支持的方法："+aggs)
			}

			O.selectFields[alias] = 1

		//普通字段
		} else {
			O.selectFields[field] = 1
		}
	}

	return O
}

func (O *Orm) Where(conds ...interface{}) *Orm {
	args_len := len(conds)
	if args_len < 1 {
		panic("查询参数不应为空")
	}

	field, ok := conds[0].(string)
	if !ok {
			maps, ok := conds[0].(map[string]interface{})
			if ok {
				for _field,_criteria := range maps {
						_set, ok := _criteria.([]interface{})
						if ok {
							_tmp := []interface{}{_field}
							_tmp  = append(_tmp, _set...)
							O.Where(_tmp...)
							continue
						}

						O.Where(_field, _criteria)
				}
				return O
			} else {
				panic("第一个参数应为string类型或map[string]interface{}类型")
			}
	}

	if placeholder_count := strings.Count(field, "?"); placeholder_count > 0 {
		if placeholder_count != args_len -1 {
			panic("查询占位符和参数数量不匹配")
		}

		for k,v := range conds {
			if k==0 {
				continue
			}

			O.selectParams = append(O.selectParams, v)
		}

		O.selectConds  = append(O.selectConds, field)
		return O
	}


	switch args_len {
		case 1 :
				_id, _ := primitive.ObjectIDFromHex(field)
				O.selectConds  = append(O.selectConds, "_id=?")
				O.selectParams = append(O.selectParams, _id)
		case 2 :
				O.selectConds  = append(O.selectConds, field+"=?")
				O.selectParams = append(O.selectParams, conds[1])
		case 3 :
				opr, ok := conds[1].(string)
				if !ok {
					panic("运算符应为string类型")
				}

				opr = strings.ToTitle(opr)
				switch opr {
					case ">"      : fallthrough
					case ">="     : fallthrough
					case "<"      : fallthrough
					case "<="     :
									switch conds[2].(type) {
										case int, int8, int16, int32, int64  :
										case uint, uint8,uint16,uint32,uint64:
										case float32,float64: 
										default : panic(field+"参数应为数值类型")
									}

									O.selectConds  = append(O.selectConds, field+" "+opr+" ?")
									O.selectParams = append(O.selectParams, conds[2])

					case "="      : fallthrough
					case "!="     :
									switch conds[2].(type) {
										case int, int8, int16, int32, int64:
										case uint, uint8,uint16,uint32,uint64:
										case float32,float64:
										case string:
										default : panic(field+"参数应为数值或字符串类型")
									}

									O.selectConds  = append(O.selectConds, field+" "+opr+" ?")
									O.selectParams = append(O.selectParams, conds[2])

					case "LIKE"   :
									criteria, ok := conds[2].(string)
									if !ok {
										panic(field+"查询条件应为string类型")
									}

									O.selectConds  = append(O.selectConds, field+" "+opr+" ?")
									O.selectParams = append(O.selectParams, criteria)

					case "IN"     : fallthrough
					case "NOT IN" :
									criteria, ok := conds[2].([]string)
									if !ok {
										panic(field+"查询条件应为[]string类型")
									}

									if len(criteria)==0 {
										criteria = append(criteria, "_in_query_placeholder_")
									}

									O.selectConds  = append(O.selectConds, field+" "+opr+" ?")
									O.selectParams = append(O.selectParams, criteria)

					case "IS"     : fallthrough
					case "IS NOT" :
									criteria, ok := conds[2].(string)
									if !ok {
										panic(field+"查询条件应为string类型")
									}

									criteria = strings.ToTitle(criteria)
									if criteria!="NULL" {
										panic(field+"查询条件只能为null")
									}

									O.selectConds  = append(O.selectConds, field+" "+opr+" "+criteria)

					case "BETWEEN":
									criteria, ok := conds[2].([]string)
									if !ok {
										panic(field+"查询条件应为[]string类型")
									}

									if len(criteria)!=2 {
										panic(field+"查询条件应为[]string类型,且必须2个元素")
									}

									O.selectConds  = append(O.selectConds, field+" >= ? AND "+field+" <= ? ")
									for _,v := range criteria {
										O.selectParams = append(O.selectParams, v)
									}

					default        :
									panic(field+"不支持的操作类型:"+opr)
				}
		default : panic("查询参数不应超过3个")
	}

	return O
}

func (O *Orm) Group(fields ...string) *Orm {
	if len(O.selectGroup)==0 {
		panic("非聚合查询不支持此操作，如需聚合请先通过Field()设置聚合字段")
	}

	for _, field := range fields {
		_, ok := O.selectGroupId[field]
		if !ok {
			O.selectGroupId[field] = "$"+field
		}
	}

	return O
}

func (O *Orm) Having(field string, opr string, criteria int) *Orm {
	if len(O.selectGroup)==0 {
		panic("非聚合查询不支持此操作")
	}

	_, ok := O.selectGroup[field]
	if !ok {
		panic(field+"：having条件字段必须是聚合别名")
	}

	var _opr string
	switch opr {
		case "="  : _opr = "$eq"
		case ">"  : _opr = "$gt"
		case ">=" : _opr = "$gte"
		case "<"  : _opr = "$lt"
		case "<=" : _opr = "$lte"
		default   : panic("having操作符仅支持=、>、>=、<、<=")
	}

	O.selectHaving[field] = map[string]interface{}{_opr:criteria}

	return O
}

func (O *Orm) Order(field string, sort string) *Orm {
	sort = strings.ToTitle(sort)
	if sort!="DESC" && sort!="ASC" {
		panic("排序类型只能是asc或desc")
	}

	if len(O.selectGroup)>0 {
		_, ok := O.selectGroup[field]
		if !ok {
			panic("聚合查询排序只能是聚合字段")
		}
	}

	if sort=="ASC" {
		O.selectOrder = append(O.selectOrder, bson.E{field, 1})
	} else {
		O.selectOrder = append(O.selectOrder, bson.E{field, -1})
	}

	return O
}

func (O *Orm) Page(page int64) *Orm {
	if page < 1 {
		panic("页码不应小于1")
	}

	O.selectPage = page

	return O
}

func (O *Orm) Limit(limit int64) *Orm {
	if limit < 1 {
		panic("每页条数不应小于1")
	}

	O.selectLimit = limit

	return O
}

func (O *Orm) Select() []map[string]interface{} {
	filter := O.filter()
	result := []map[string]interface{}{}

	var cursor *mongo.Cursor

	//聚合查询
	if len(O.selectGroup)>0 {
		aggs := []map[string]interface{}{}

		aggs = append(aggs, map[string]interface{}{"$match": filter})

		O.selectGroup["_id"] = O.selectGroupId
		aggs = append(aggs, map[string]interface{}{"$group":O.selectGroup})

		if len(O.selectOrder)>0 {
			aggs = append(aggs, map[string]interface{}{"$sort":O.selectOrder})
		}

		if len(O.selectHaving)>0 {
			aggs = append(aggs, map[string]interface{}{"$match":O.selectHaving})
		}

		aggs = append(aggs, map[string]interface{}{"$skip": int64(O.selectLimit * (O.selectPage - 1))})
		aggs = append(aggs, map[string]interface{}{"$limit": O.selectLimit})

		if O.debug=="yes" {
			log.Println(aggs)
		}

	    //常见聚合count\sum\max\min\avg
	    var err error
	    cursor, err = O.coll.Aggregate(context.Background(), aggs)
		if err != nil {
			panic(err)
		}

	    //普通查询
	} else {

		findOptions := options.Find()
		findOptions.SetLimit(O.selectLimit)
		findOptions.SetSkip(int64(O.selectLimit * (O.selectPage - 1)))

		if len(O.selectOrder)>0 {
			findOptions.SetSort(O.selectOrder)
		}

		if len(O.selectFields)>0 {
			findOptions.SetProjection(O.selectFields)
		}

		var err error
		cursor, err = O.coll.Find(context.TODO(), filter, findOptions)
		if err != nil {
			panic(err)
		}
	}

	//遍历取出结果
	var list []map[string]interface{}
	if err := cursor.All(context.TODO(), &list); err != nil {
		panic(err)
	}

	for _, v := range list {
		cursor.Decode(&v)
		result = append(result, v)
	}

	return result
}

func (O *Orm) Find() map[string]interface{} {
	if len(O.selectGroup)>0 {
		panic("聚合查询不支持此操作")
	}

	filter := O.filter()

	findOptions := options.FindOne()

	if len(O.selectOrder)>0 {
		findOptions.SetSort(O.selectOrder)
	}

	if len(O.selectFields)>0 {
		findOptions.SetProjection(O.selectFields)
	}

	var result map[string]interface{}
	err := O.coll.FindOne(context.TODO(), filter, findOptions).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			panic(err)
		}
	}

	return result
}

func (O *Orm) Value(field string) interface{} {
	O.selectFields = map[string]interface{}{field:1}

	doc := O.Find()

	value, ok := doc[field]
	if ok {
		return value
	}

	return nil
}

func (O *Orm) Values(field string) []interface{} {
	if len(O.selectGroup)>0 {
		panic("聚合查询不支持此操作")
	}

	O.selectFields = map[string]interface{}{field:1}

	list := O.Select()

	result := []interface{}{}

	for _, v := range list {
		value, ok := v[field]
		if ok {
			result = append(result, value)
		}
	}

	return result
}

func (O *Orm) Columns(fields ...string) map[string]interface{} {
	if len(O.selectGroup)>0 {
		panic("聚合查询不支持此操作")
	}

	var key   string
	var value string

	if len(fields)==0 {
		panic("参数不可为空")

	} else if(len(fields)==1) {
		key   = "_id"
		value = fields[0]

		O.selectFields = map[string]interface{}{value:1}

	} else {
		key   = fields[1]
		value = fields[0]

		O.selectFields = map[string]interface{}{value:1,key:1}
	}

	list := O.Select()

	var result = map[string]interface{}{}
	if len(list)>0 {
		for _, v := range list {
			result[v[key].(string)] = v[value]
		}
	}

	return result
}

func (O *Orm) Sum(field string) int32 {
	if len(O.selectGroup)>0 {
		panic("聚合查询不支持此操作")
	}

    filter    := O.filter()
    cursor, _ := O.coll.Aggregate(context.Background(), []map[string]interface{}{
        {"$match": filter},
        {"$group": map[string]interface{}{"_id": nil, "sum": map[string]interface{}{"$sum": "$"+field}}},
    })

    if cursor.Next(context.TODO()) {
    	var result map[string]interface{}
		cursor.Decode(&result)
		return result["sum"].(int32)
    }

    return 0
}

func (O *Orm) Count() int32 {
	filter := O.filter()

	if len(O.selectGroupId)>0 {
		aggs := []map[string]interface{}{}

		aggs = append(aggs, map[string]interface{}{"$match": filter})

		O.selectGroup["_id"] = O.selectGroupId
		aggs = append(aggs, map[string]interface{}{"$group":O.selectGroup})

		if len(O.selectHaving)>0 {
			aggs = append(aggs, map[string]interface{}{"$match":O.selectHaving})
		}

		aggs = append(aggs, map[string]interface{}{"$count":"total"})

		if O.debug=="yes" {
			log.Println(aggs)
		}

	    //常见聚合count\sum\max\min\avg
	    var err error
	    cursor, err := O.coll.Aggregate(context.Background(), aggs)
		if err != nil {
			panic(err)
		}

		if cursor.Next(context.TODO()) {
    		var result map[string]interface{}
			cursor.Decode(&result)
			return result["total"].(int32)
    	}

		return 0

	} else {
		count, err := O.coll.CountDocuments(context.TODO(), filter)
		if err != nil {
			panic(err)
		}

		return int32(count)
	}
}

func (O *Orm) Exist(id string) bool {
	if len(O.selectGroup)>0 {
		panic("聚合查询不支持此操作")
	}

	var result map[string]interface{}

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := map[string]interface{}{"_id":_id}
	err    := O.coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		} else {
			panic(err)
		}
	}

	return true
}

func (O *Orm) filter() map[string]interface{} {
	sql := strings.Join(O.selectConds, " AND ")

	var filter map[string]interface{}
	if len(sql)>0 {
		scheme := (&Parser{}).Parse(sql)
		if O.debug=="yes" {
			log.Println(scheme)
		}

		filter = O.bind(scheme)
		if O.debug=="yes" {
			log.Println(filter)
		}
	} else {
		filter = map[string]interface{}{}
	}

	return filter
}

func (O *Orm) bind(filter map[string]interface{}) map[string]interface{} {
	for k,v := range filter {
		switch k {
			case "$and": fallthrough
			case "$or":
						var cond []map[string]interface{} = []map[string]interface{}{}

						for _,e := range v.([]map[string]interface{}) {
							cond = append(cond, O.bind(e))
						}

						return map[string]interface{}{k:cond}
			default:
					set := v.(map[string]interface{})
					value := set["value"]
					opr   := set["opr"].(string)
					idx   := set["placeholder"].(int)
					if value=="?" {
						param := O.selectParams[idx]
						return map[string]interface{}{
							k:map[string]interface{}{opr:param},
						}
					} else {
						return map[string]interface{}{
							k:map[string]interface{}{opr:value},
						}
					}
		}
	}

	return nil
}

