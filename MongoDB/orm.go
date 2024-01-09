package MongoDB

import "os"
import "fmt"
import "context"
import "strings"
import "go.mongodb.org/mongo-driver/mongo"
import "go.mongodb.org/mongo-driver/bson/primitive"

type Orm struct {
	coll *mongo.Collection

	selectFields map[string]interface{}
	selectConds []string
	selectParams []interface{}
	selectPage int64
	selectLimit int64
	selectOrder map[string]interface{}
	selectGroup []string
	selectHaving string

	debug string
}

func (O *Orm) Init(dbtag string, table string) *Orm {
	dbname        := Dbname(dbtag)
	O.coll         = Open(dbtag).Database(dbname).Collection(table)
	O.debug        = os.Getenv("debug")
	O.selectFields = map[string]interface{}{}
	O.selectOrder  = map[string]interface{}{}

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
	return 1
}

func (O *Orm) Update(data map[string]interface{}) int64 {
	return 1
}

func (O *Orm) Field(fields string) *Orm {
	_fields := strings.Split(fields, ",")
	for _, field := range _fields {
		O.selectFields[field] = 1
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
				O.selectConds  = append(O.selectConds, "_id=?")
				O.selectParams = append(O.selectParams, field)
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

									placeholders := []string{}
									for _,v := range criteria {
										placeholders   = append(placeholders, "?")
										O.selectParams = append(O.selectParams, v)
									}

									O.selectConds  = append(O.selectConds, field+" "+opr+"("+strings.Join(placeholders, ",")+")")

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
	return O
}

func (O *Orm) Having(field string, opr string, criteria int) *Orm {
	return O
}

func (O *Orm) Order(field string, sort string) *Orm {
	sort = strings.ToTitle(sort)
	if sort!="DESC" && sort!="ASC" {
		panic("排序类型只能是asc或desc")
	}

	if sort=="ASC" {
		O.selectOrder[field] = 1
	} else {
		O.selectOrder[field] = -1
	}
fmt.Println(O.selectOrder)
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
	return nil
}

func (O *Orm) Find() map[string]interface{} {
	return nil
}

func (O *Orm) Value(field string) string {
	return ""
}

func (O *Orm) Values(field string) []string {
	return nil
}

func (O *Orm) Columns(fields ...string) map[string]string {
	return nil
}

func (O *Orm) Sum(field string) int64 {
	return 1
}

func (O *Orm) Count() int64 {
	return 1
}

func (O *Orm) Exist(primary string) bool {
	return true
}

