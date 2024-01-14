package MongoDB

import "testing"
import "os"

//注册数据库连接
func init() {
	os.Setenv("debug", "yes")
	Register("mongodb", "test", "mongodb://localhost:27017")
}

//插入
func TestInsert(t *testing.T) {
	O := Model{"goods"}
	O.Insert(map[string]any{"name":"可口可乐","price":100,"detail":"...","category":"饮料"})
	O.Insert(map[string]any{"name":"小红帽","price":200,"detail":"...","category":"服装"})
	O.Insert(map[string]any{"name":"雪碧","price":300,"category":"饮料"})
	O.Insert(map[string]any{"name":"高跟鞋","price":400,"detail":"...","category":"服装"})
	O.Insert(map[string]any{"name":"芬达","price":500,"detail":"...","category":"饮料"})
	O.Insert(map[string]any{"name":"海魂衫","price":600,"category":"服装"})
	O.Insert(map[string]any{"name":"和其正","price":700,"detail":"...","category":"饮料"})
	O.Insert(map[string]any{"name":"领带","price":800,"detail":"...","category":"服装"})
	O.Insert(map[string]any{"name":"美年达","price":900,"category":"饮料"})
	O.Insert(map[string]any{"name":"呢子大衣","price":200,"detail":"...","category":"服装"})
}

//主键条件查询
func TestSelectPrimary(t *testing.T) {
	O   := Model{"goods"}
	id  := O.Insert(map[string]any{"name":"主键查询","price":200,"detail":"...","category":"服装"})
	row := O.Where(id).Find()
	_, ok := row["name"]
	if ok {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//主键条件检查记录是否存在
func TestExist(t *testing.T) {
	O := Model{"goods"}
	ok := O.Exist("_test_")
	if !ok {
		t.Log("yes")
	} else {
		t.Fail()
	}
}

//指定返回字段
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

//=条件查询
func TestSelectEq(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("price", 400).Find()
	if row!=nil {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//大于等于条件查询
func TestSelectGtEq(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("price", ">=", 700).Find()
	if row!=nil {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//多条件查询
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

//in条件查询
func TestSelectIn(t *testing.T) {
	O := Model{"goods"}
	rows := O.Where("name", "in", []string{"小红帽","可口可乐"}).Select()
	if len(rows)>0 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//null过滤条件查询
func TestSelectNull(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("detail", "is", "null").Find()
	if row!=nil {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//not null过滤条件查询
func TestSelectNotNull(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("detail", "is not", "null").Find()
	if row!=nil {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//between区间条件查询
func TestSelectBetween(t *testing.T) {
	O := Model{"goods"}
	rows := O.Where("price", "BETWEEN", []string{"300","800"}).Select()
	if len(rows)>1 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//复杂语句参数代入查询
func TestSelectExp(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("detail is not null AND category=? AND price>?", "服装",100).Find()
	if len(row)>1 {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//like条件搜索查询
func TestSelectLike(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("name", "like", "%帽").Find()
	if len(row)>1 {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//字段排序
func TestSelectOrder(t *testing.T) {
	O := Model{"goods"}
	rows := O.Field("category,price").Order("category","asc").Order("price","desc").Select()
	if len(rows)>1 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//group+having查询
func TestSelectGroup(t *testing.T) {
	O := Model{"goods"}
	rows := O.Field("count(*) as num,category,price").
			Where("price", ">", 50).
			Group("category","price").
			Having("num",">",1).
			Select()

	if len(rows)>0 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//直接取单记录指定字段值
func TestSelectValue(t *testing.T) {
	O := Model{"goods"}
	name := O.Field("name").Value("name")
	if name!=nil {
		t.Log(name)
	} else {
		t.Fail()
	}
}

//取多记录指定字段切片
func TestSelectValues(t *testing.T) {
	O := Model{"goods"}
	names := O.Field("name").Values("name")

	if len(names)>1 {
		t.Log(names)
	} else {
		t.Fail()
	}
}

//取K=>V字段记录map
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

//取最大最小值
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

//查询条件复用
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

//更新操作，可选更新条数
func TestUpdate(t *testing.T) {
	O := Model{"goods"}
	af := O.Where("name", "可口可乐").Update(map[string]any{"name":"可可口口","price":123})
	if af > 0 {
		t.Log(af)
	} else {
		t.Fail()
	}
}

//删除操作，可选删除条数
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
