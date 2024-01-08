package MongoDB

import "os"
import "go.mongodb.org/mongo-driver/mongo"

type Orm struct {
	coll *mongo.Collection
	debug string
}

func (O *Orm) Init(dbtag string, table string) *Orm {
	dbname := Dbname(dbtag)
	O.coll  = Open(dbtag).Database(dbname).Collection(table)
	O.debug = os.Getenv("debug")

	return O
}

func (O *Orm) Insert(data map[string]interface{}) string {
	return ""
}

func (O *Orm) Delete() int64 {
	return 1
}

func (O *Orm) Update(data map[string]interface{}) int64 {
	return 1
}

func (O *Orm) Field(fields string) *Orm {
	return O
}

func (O *Orm) Where(conds ...interface{}) *Orm {
	return O
}

func (O *Orm) Group(fields ...string) *Orm {
	return O
}

func (O *Orm) Having(field string, opr string, criteria int) *Orm {
	return O
}

func (O *Orm) Order(field string, sort string) *Orm {
	return O
}

func (O *Orm) Page(page int) *Orm {
	return O
}

func (O *Orm) Limit(limit int) *Orm {
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

