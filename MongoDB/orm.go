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
