package MongoDB

import "os"
import "log"
import "sync"
import "context"
import "go.mongodb.org/mongo-driver/mongo"
import "go.mongodb.org/mongo-driver/mongo/options"

var cache sync.Map

func Register(tag string, dbname string, dsn string) {
	if tag=="mongodb" && (dbname=="" || dsn=="") {
		log.Printf("\033[1;31;40m%s\033[0m\n",".env配置文件不存在或database.dbname和database.dsn未设置")
		os.Exit(1)
	}

	cache.Store("dsn."+tag, dsn)
	cache.Store("dbname."+tag, dbname)
}

func Dsn(tag string) string {
	value, ok := cache.Load("dsn."+tag);
    if !ok {
        panic("dsn配置不存在:"+tag)
    }

    dsn, _ := value.(string)

	return dsn
}

func Dbname(tag string) string {
	value, ok := cache.Load("dbname."+tag)
	if !ok {
		panic("dbname配置不存在:"+tag)
	}

	dbname, _ := value.(string)

	return dbname
}

func NewOrm(table string, tag ...string) *Orm {
	var dbtag string
	if len(tag)==0 {
		dbtag = "mongodb"
	} else {
		dbtag = tag[0]
	}

	return (&Orm{}).Init(dbtag, table)
}

func Open(tag string) *mongo.Client {
	dsn := Dsn(tag)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn).SetMinPoolSize(5).SetMaxPoolSize(100))
	if err != nil {
		panic(err)
	}

	return client
}

func Close(client *mongo.Client) bool {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	return true
}
