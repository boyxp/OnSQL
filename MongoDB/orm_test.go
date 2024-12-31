package MongoDB

import "testing"
import "os"

//Initialization
func init() {
	os.Setenv("debug", "yes")
	Register("mongodb", "test", "mongodb://localhost:27017")
}

//Insert Records
func TestInsert(t *testing.T) {
	O := Model{"goods"}
	O.Insert(map[string]any{"name":"可口可乐","price":100,"detail":"...","category":"饮料","create_time":"2024-01-01 01:01:01"})
	O.Insert(map[string]any{"name":"小红帽","price":200,"detail":"...","category":"服装","create_time":"2024-01-03 01:01:01"})
	O.Insert(map[string]any{"name":"雪碧","price":300,"category":"饮料","create_time":"2024-01-05 01:01:01"})
	O.Insert(map[string]any{"name":"高跟鞋","price":400,"detail":"...","category":"服装","create_time":"2024-01-10 01:01:01"})
	O.Insert(map[string]any{"name":"芬达","price":500,"detail":"...","category":"饮料","create_time":"2024-01-12 01:01:01"})
	O.Insert(map[string]any{"name":"海魂衫","price":600,"category":"服装","create_time":"2024-01-13 01:01:01"})
	O.Insert(map[string]any{"name":"和其正","price":700,"detail":"...","category":"饮料","create_time":"2024-01-14 01:01:01"})
	O.Insert(map[string]any{"name":"领带","price":800,"detail":"...","category":"服装","create_time":"2024-01-15 01:01:01"})
	O.Insert(map[string]any{"name":"美年达","price":900,"category":"饮料","create_time":"2024-01-16 01:01:01"})
	O.Insert(map[string]any{"name":"呢子大衣","price":200,"detail":"...","category":"服装","create_time":"2024-01-17 01:01:01"})
}

//Primary Key Query
func TestSelectPrimary(t *testing.T) {
	O   := Model{"goods"}
	id  := O.Insert(map[string]any{"name":"主键查询","price":200,"detail":"...","category":"服装"})
	row1 := O.Where(id).Find()
	if row1!=nil {
		t.Log(row1)
	} else {
		t.Fail()
	}

	row2 := O.Where("_id", id).Find()
	if row2!=nil {
		t.Log(row2)
	} else {
		t.Fail()
	}

	row3 := O.Where("_id", "in", []string{id}).Find()
	if row3!=nil {
		t.Log(row3)
	} else {
		t.Fail()
	}
}

//Existence Check
func TestExist(t *testing.T) {
	O := Model{"goods"}
	ok := O.Exist("_test_")
	if !ok {
		t.Log("yes")
	} else {
		t.Fail()
	}
}

//Field Selection
func TestSelectField(t *testing.T) {
	O := Model{"goods"}
	row := O.Field("name,category").Find()
	_, ok := row["name"]

	if ok {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//Equal Condition
func TestSelectEq(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("price", 400).Find()
	if row!=nil {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//Greater or Equal Condition
func TestSelectGtEq(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("price", ">=", 700).Find()
	if row!=nil {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//Multiple Conditions
func TestSelectMulti(t *testing.T) {
	O := Model{"goods"}
	rows := O.Where(map[string]any{
		"category":"服装",
		"name":[]any{"is not", "null"},
		"price":[]any{">=", 200},
	}).Select()

	if len(rows)>0 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//In Condition
func TestSelectIn(t *testing.T) {
	O := Model{"goods"}
	rows := O.Where("name", "in", []string{"小红帽","可口可乐"}).Select()
	if len(rows)>0 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//Null Conditions
func TestSelectNull(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("detail", "is", "null").Find()
	if row!=nil {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//Not Null Conditions
func TestSelectNotNull(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("detail", "is not", "null").Find()
	if row!=nil {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//Between Condition
func TestSelectBetween(t *testing.T) {
	O := Model{"goods"}
	rows := O.Where("create_time", "BETWEEN", []string{"2024-01-01 01:01:01","2024-01-03 01:01:01"}).Select()
	if len(rows)>1 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//Complex Expressions
func TestSelectExp(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("detail is not null AND category=? AND price>?", "服装",100).Find()
	if len(row)>1 {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//Like Condition
func TestSelectLike(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("name", "like", "%帽").Find()
	if len(row)>1 {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//Ordering
func TestSelectOrder(t *testing.T) {
	O := Model{"goods"}
	rows := O.Field("category,price").Order("category","asc").Order("price","desc").Select()
	if len(rows)>1 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//Grouping and Aggregation
func TestSelectGroup(t *testing.T) {
	O := Model{"goods"}
	rows := O.Field("count(*) as num,category,price,date_format(create_time, '%Y_%m_%d_%H') as d").
			Where("price", ">", 50).
			Where("create_time", "is not", "null").
			Group("category","price","d").
			Select()

	if len(rows)>0 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//Value Selection
func TestSelectValue(t *testing.T) {
	O := Model{"goods"}
	name := O.Field("name").Value("name")
	if name!=nil {
		t.Log(name)
	} else {
		t.Fail()
	}
}

func TestSelectValues(t *testing.T) {
	O := Model{"goods"}
	names := O.Field("name").Values("name")

	if len(names)>1 {
		t.Log(names)
	} else {
		t.Fail()
	}
}

func TestSelectColumns(t *testing.T) {
	O := Model{"goods"}
	//names := O.Columns("name")
	names := O.Columns("category", "name")

	if len(names)>1 {
		t.Log(names)
	} else {
		t.Fail()
	}
}

//Max and Min Values
func TestSelectMaxMin(t *testing.T) {
	O := Model{"goods"}
	max_id := O.Field("MAX(price) as max_price").Select()
	if max_id!=nil {
		t.Log(max_id)
	} else {
		t.Fail()
	}

	min_id := O.Field("MIN(price) as min_price").Select()
	if min_id!=nil {
		t.Log(min_id)
	} else {
		t.Fail()
	}
}

//Query Reuse
func TestSelectQueryReuse(t *testing.T) {
	O := Model{"goods"}

	query := O.Where("detail", "is not", "null")

	rows := query.Select()
	if len(rows)>1 {
		t.Log(rows)
	} else {
		t.Fail()
	}

	row := query.Find()
	if len(row)>1 {
		t.Log(row)
	} else {
		t.Fail()
	}

	name := query.Value("name")
	if name!=nil {
		t.Log(name)
	} else {
		t.Fail()
	}

	count := query.Count()
	if count>1 {
		t.Log(count)
	} else {
		t.Fail()
	}

	sum := query.Sum("price")
	if sum>100 {
		t.Log(sum)
	} else {
		t.Fail()
	}
}

//Update Records
func TestUpdate(t *testing.T) {
	O := Model{"goods"}
	af := O.Where("name", "可口可乐").Update(map[string]any{"name":"可可口口","price":123})
	if af > 0 {
		t.Log(af)
	} else {
		t.Fail()
	}
}

//Delete Records
func TestDelete(t *testing.T) {
	O := Model{"goods"}
	id := O.Field("_id").Find()
	af := O.Where(id["_id"]).Delete()
	if af > 0 {
		t.Log(af)
	} else {
		t.Fail()
	}
}
